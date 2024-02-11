package service

import (
	"brenonaraujo/rinhabackend-q12024/domain"
	"brenonaraujo/rinhabackend-q12024/infra/database"
	"context"
	"fmt"
)

func DeductBalance(customerId, amount int, description string) (domain.Balance, error) {
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

	var customer *domain.Customer
	customer, err = GetCustomer(customerId)
	if err != nil {
		tx.Rollback(ctx)
		return customerBalance, err
	}

	var currentBalance int
	err = tx.QueryRow(ctx,
		"SELECT valor FROM saldos WHERE cliente_id=$1 FOR UPDATE",
		customerId).Scan(&currentBalance)
	if err != nil {
		tx.Rollback(ctx)
		return customerBalance, err
	}

	newBalance := currentBalance - amount
	if newBalance < -customer.AccountLimit {
		tx.Rollback(ctx)
		return customerBalance, fmt.Errorf("Amount %d would violate customer limit",
			amount)
	}

	_, err = tx.Exec(ctx,
		"UPDATE saldos SET valor=$1 WHERE id=$2", newBalance, customerId)
	if err != nil {
		tx.Rollback(ctx)
		return customerBalance, fmt.Errorf("updating customer balance: %w", err)
	}

	_, err = tx.Exec(ctx,
		"INSERT INTO transacoes (id, cliente_id, valor, tipo, descricao, realizada_em) values(default, $1, $2, 'd', $3, now())",
		customerId, amount, description)
	if err != nil {
		tx.Rollback(ctx)
		return customerBalance, fmt.Errorf("Isert transaction error: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return customerBalance, fmt.Errorf("committing transaction: %w", err)
	}

	return domain.Balance{Limite: customer.AccountLimit, Saldo: newBalance}, nil
}
