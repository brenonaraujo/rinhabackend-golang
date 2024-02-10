package service

import (
	"brenonaraujo/rinhabackend-q12024/domain"
	"brenonaraujo/rinhabackend-q12024/infra/database"
	"context"
	"fmt"
)

func AddBalance(customerId, amount int) (domain.Balance, error) {
	ctx := context.Background()
	var customerBalance domain.Balance

	tx, err := database.GetDBPool().Begin(ctx)
	if err != nil {
		return customerBalance, fmt.Errorf("starting transaction: %w", err)
	}
	defer func() {
		if p := recover(); p != nil {
			tx.Rollback(ctx)
			panic(p)
		}
	}()

	// Assuming 'saldos' table's 'cliente_id' is unique, use it directly without JOIN
	var currentBalance, limit int
	err = tx.QueryRow(ctx,
		"SELECT valor, limite FROM saldos JOIN clientes ON saldos.cliente_id = clientes.id WHERE cliente_id=$1 FOR UPDATE", customerId).Scan(&currentBalance, &limit)
	if err != nil {
		tx.Rollback(ctx)
		return customerBalance, fmt.Errorf("querying customer balance and limit: %w", err)
	}

	newBalance := currentBalance + amount

	_, err = tx.Exec(ctx,
		"UPDATE saldos SET valor=$1 WHERE cliente_id=$2", newBalance, customerId)
	if err != nil {
		tx.Rollback(ctx)
		return customerBalance, fmt.Errorf("updating customer balance: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return customerBalance, fmt.Errorf("committing transaction: %w", err)
	}

	customerBalance.Limite = limit
	customerBalance.Saldo = newBalance
	return customerBalance, nil
}
