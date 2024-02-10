package service

import (
	"brenonaraujo/rinhabackend-q12024/infra/database"
	"context"
	"fmt"
)

type Saldo struct {
	Total       int    `json:"total"`
	Limite      int    `json:"limite"`
	DataExtrato string `json:"data_extrato"`
}

type UltimasTransacoes []struct {
	Valor       int    `json:"valor"`
	Tipo        string `json:"tipo"`
	Descricao   string `json:"descricao"`
	RealizadaEm string `json:"realizada_em"`
}

type CustomerStattement struct {
	saldo              Saldo
	ultimas_transacoes UltimasTransacoes
}

func GetTransactions(customerId int) (CustomerStattement, error) {
	var customerStattement CustomerStattement

	ctx := context.Background()
	db := database.GetDBPool()

	var saldo Saldo
	err := db.QueryRow(ctx, "SELECT valor, limite, NOW() FROM saldos JOIN clientes ON saldos.cliente_id = clientes.id WHERE cliente_id=$1", customerId).Scan(&saldo.Total, &saldo.Limite, &saldo.DataExtrato)
	if err != nil {
		return customerStattement, fmt.Errorf("querying customer transactions: %w", err)
	}

	rows, err := db.Query(ctx, "SELECT valor, tipo, descricao, realizada_em FROM transacoes WHERE cliente_id=$1 ORDER BY realizada_em DESC LIMIT 10", customerId)
	if err != nil {
		return customerStattement, fmt.Errorf("Failed to fetch transactions: %w", err)
	}
	defer rows.Close()
	var lastTranscations UltimasTransacoes
	for rows.Next() {
		var t struct {
			Valor       int    `json:"valor"`
			Tipo        string `json:"tipo"`
			Descricao   string `json:"descricao"`
			RealizadaEm string `json:"realizada_em"`
		}
		if err := rows.Scan(&t.Valor, &t.Tipo, &t.Descricao, &t.RealizadaEm); err != nil {
			return customerStattement, fmt.Errorf("Failed to read transaction data: %w", err)
		}
		lastTranscations = append(lastTranscations, t)
	}
	return CustomerStattement{saldo: saldo, ultimas_transacoes: lastTranscations}, nil
}
