package service

import (
	"brenonaraujo/rinhabackend-q12024/domain"
	"brenonaraujo/rinhabackend-q12024/infra/database"
	"context"
	"fmt"
)

func DeductBalance(customerId, amount int) (domain.Balance, error) {
	ctx := context.Background()
	var costumerBalance domain.Balance

	tx, err := database.GetDBPool().Begin(ctx)
	if err != nil {
		return costumerBalance, fmt.Errorf("starting transaction: %w", err)
	}
	defer func() {
		if p := recover(); p != nil {
			tx.Rollback(ctx)
			panic(p)
		}
	}()

	var currentBalance int
	err = tx.QueryRow(ctx,
		"SELECT valor FROM saldos WHERE id=$1 FOR UPDATE", customerId).Scan(&currentBalance)
	if err != nil {
		tx.Rollback(ctx)
		return costumerBalance, fmt.Errorf("querying customer balance: %w", err)
	}

	// TODO: Set client to a memory cache to do not always request from db.
	var limit int
	err = tx.QueryRow(ctx,
		"SELECT limite FROM clientes WHERE id=$1", customerId).Scan(&limit)
	if err != nil {
		tx.Rollback(ctx)
		return costumerBalance, fmt.Errorf("querying customer limit: %w", err)
	}

	newBalance := currentBalance - amount
	if newBalance < -limit {
		tx.Rollback(ctx)
		return costumerBalance, fmt.Errorf("deduction amount %d would violate customer limit", amount)
	}

	_, err = tx.Exec(ctx,
		"UPDATE saldos SET valor=$1 WHERE id=$2", newBalance, customerId)
	if err != nil {
		tx.Rollback(ctx)
		return costumerBalance, fmt.Errorf("updating customer balance: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return costumerBalance, fmt.Errorf("committing transaction: %w", err)
	}

	return domain.Balance{Limite: limit, Saldo: newBalance}, nil
}
