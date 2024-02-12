package service

import (
	"brenonaraujo/rinhabackend-q12024/domain"
	"brenonaraujo/rinhabackend-q12024/infra/database"
	"context"
	"fmt"
)

type OperationType string

const (
	Debit  OperationType = "d"
	Credit OperationType = "c"
)

func TransactionProcess(ctx context.Context, customerId, amount int, description string, opType OperationType) (domain.Balance, error) {
	tx, err := database.GetDBPool().Begin(ctx)
	if err != nil {
		return domain.Balance{}, fmt.Errorf("starting transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	customer, err := GetCustomer(customerId)
	if err != nil {
		return domain.Balance{}, err
	}

	var currentBalance int
	err = tx.QueryRow(ctx, "SELECT valor FROM saldos WHERE cliente_id=$1 FOR UPDATE", customerId).Scan(&currentBalance)
	if err != nil {
		return domain.Balance{}, err
	}

	newBalance, err := calculateNewBalance(currentBalance, amount, opType, customer.AccountLimit) // Encapsulate balance calculation and limit check
	if err != nil {
		return domain.Balance{}, err
	}

	_, err = tx.Exec(ctx, "UPDATE saldos SET valor=$1 WHERE cliente_id=$2", newBalance, customerId)
	if err != nil {
		return domain.Balance{}, fmt.Errorf("updating customer balance: %w", err)
	}

	_, err = tx.Exec(ctx, "INSERT INTO transacoes (cliente_id, valor, tipo, descricao, realizada_em) VALUES ($1, $2, $3, $4, NOW())", customerId, amount, opType, description)
	if err != nil {
		return domain.Balance{}, fmt.Errorf("insert transaction error: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return domain.Balance{}, fmt.Errorf("committing transaction: %w", err)
	}

	return domain.Balance{Limite: customer.AccountLimit, Saldo: newBalance}, nil
}

func calculateNewBalance(currentBalance, amount int, opType OperationType, accountLimit int) (int, error) {
	switch opType {
	case Debit:
		if newBalance := currentBalance - amount; newBalance < -accountLimit {
			return 0, fmt.Errorf("debit amount %d would violate customer limit", amount)
		} else {
			return newBalance, nil
		}
	case Credit:
		return currentBalance + amount, nil
	default:
		return 0, fmt.Errorf("unknown operation type: %s", opType)
	}
}
