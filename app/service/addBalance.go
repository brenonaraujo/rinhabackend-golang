package service

import (
	"brenonaraujo/rinhabackend-q12024/domain"
	"brenonaraujo/rinhabackend-q12024/infra/database"
	"context"
	"fmt"
)

func AddBalance(customerId int, amount int, description string) (domain.Balance, error) {
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
		"SELECT valor FROM saldos WHERE cliente_id=$1 FOR UPDATE", customerId).Scan(&currentBalance)
	if err != nil {
		tx.Rollback(ctx)
		return customerBalance, err
	}

	newBalance := currentBalance + amount

	_, err = tx.Exec(ctx,
		"UPDATE saldos SET valor=$1 WHERE cliente_id=$2", newBalance, customerId)
	if err != nil {
		tx.Rollback(ctx)
		return customerBalance, err
	}

	_, err = tx.Exec(ctx,
		"INSERT INTO transacoes (id, cliente_id, valor, tipo, descricao, realizada_em) values(default, $1, $2, 'c', $3, now())",
		customerId, amount, description)
	if err != nil {
		tx.Rollback(ctx)
		return customerBalance, err
	}

	if err := tx.Commit(ctx); err != nil {
		return customerBalance, err
	}

	customerBalance.Limite = customer.AccountLimit
	customerBalance.Saldo = newBalance
	return customerBalance, nil
}
