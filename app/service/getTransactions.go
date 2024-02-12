package service

import (
	"brenonaraujo/rinhabackend-q12024/infra/database"
	"context"
	"fmt"
	"time"
)

type Statement struct {
	Total         int       `json:"total"`
	Limit         int       `json:"limite"`
	StatementDate time.Time `json:"data_extrato"`
}

type LastTransaction struct {
	Value       int       `json:"valor"`
	Type        string    `json:"tipo"`
	Description string    `json:"descricao"`
	CreatedAt   time.Time `json:"realizada_em"`
}

type CustomerStatement struct {
	Statement        Statement         `json:"saldo"`
	LastTransactions []LastTransaction `json:"ultimas_transacoes"`
}

func GetCustomerStatement(customerId int) (CustomerStatement, error) {
	var customerStatement CustomerStatement

	ctx := context.Background()
	db := database.GetDBPool()

	var statement Statement
	err := db.QueryRow(ctx, "SELECT valor, limite, NOW() FROM saldos JOIN clientes ON saldos.cliente_id = clientes.id WHERE cliente_id=$1",
		customerId).Scan(&statement.Total, &statement.Limit, &statement.StatementDate)
	if err != nil {
		return customerStatement, fmt.Errorf("querying customer transactions: %w", err)
	}

	rows, err := db.Query(ctx, "SELECT valor, tipo, descricao, realizada_em FROM transacoes WHERE cliente_id=$1 ORDER BY realizada_em DESC LIMIT 10", customerId)
	if err != nil {
		return customerStatement, fmt.Errorf("Failed to fetch transactions: %w", err)
	}
	defer rows.Close()
	var lastTranscations []LastTransaction
	for rows.Next() {
		var t LastTransaction
		if err := rows.Scan(&t.Value, &t.Type, &t.Description, &t.CreatedAt); err != nil {
			return customerStatement, fmt.Errorf("Failed to read transaction data: %w", err)
		}
		lastTranscations = append(lastTranscations, t)
	}
	customerStatement = CustomerStatement{Statement: statement, LastTransactions: lastTranscations}
	return customerStatement, nil
}
