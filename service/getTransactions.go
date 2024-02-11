package service

import (
	"brenonaraujo/rinhabackend-q12024/infra/database"
	"context"
	"fmt"
	"time"
)

type Saldo struct {
	Total       int       `json:"total"`
	Limite      int       `json:"limite"`
	DataExtrato time.Time `json:"data_extrato"`
}

type UltimasTransacoes []struct {
	Valor       int       `json:"valor"`
	Tipo        string    `json:"tipo"`
	Descricao   string    `json:"descricao"`
	RealizadaEm time.Time `json:"realizada_em"`
}

type CustomerStatement struct {
	Saldo              Saldo             `json:"saldo"`
	Ultimas_transacoes UltimasTransacoes `json:"ultimas_transacoes"`
}

func GetCustomerStatement(customerId int) (CustomerStatement, error) {
	var customerStatement CustomerStatement

	ctx := context.Background()
	db := database.GetDBPool()

	var saldo Saldo
	err := db.QueryRow(ctx, "SELECT valor, limite, NOW() FROM saldos JOIN clientes ON saldos.cliente_id = clientes.id WHERE cliente_id=$1", customerId).Scan(&saldo.Total, &saldo.Limite, &saldo.DataExtrato)
	if err != nil {
		return customerStatement, fmt.Errorf("querying customer transactions: %w", err)
	}

	rows, err := db.Query(ctx, "SELECT valor, tipo, descricao, realizada_em FROM transacoes WHERE cliente_id=$1 ORDER BY realizada_em DESC LIMIT 10", customerId)
	if err != nil {
		return customerStatement, fmt.Errorf("Failed to fetch transactions: %w", err)
	}
	defer rows.Close()
	var lastTranscations UltimasTransacoes
	for rows.Next() {
		var t struct {
			Valor       int       `json:"valor"`
			Tipo        string    `json:"tipo"`
			Descricao   string    `json:"descricao"`
			RealizadaEm time.Time `json:"realizada_em"`
		}
		if err := rows.Scan(&t.Valor, &t.Tipo, &t.Descricao, &t.RealizadaEm); err != nil {
			return customerStatement, fmt.Errorf("Failed to read transaction data: %w", err)
		}
		lastTranscations = append(lastTranscations, t)
	}
	customerStatement = CustomerStatement{Saldo: saldo, Ultimas_transacoes: lastTranscations}
	return customerStatement, nil
}
