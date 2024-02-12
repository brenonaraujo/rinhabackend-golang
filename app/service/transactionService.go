package service

import (
	"brenonaraujo/rinhabackend-q12024/domain/entities"
	"brenonaraujo/rinhabackend-q12024/domain/usecases"
	"brenonaraujo/rinhabackend-q12024/infra/database"
	"context"
	"fmt"
)

func TransactionProcess(ctx context.Context, customerId, amount int, description string, opType entities.OperationType) (entities.Balance, error) {
	tx, err := database.GetDBPool().Begin(ctx)
	if err != nil {
		return entities.Balance{}, fmt.Errorf("starting transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	customer, err := GetCustomer(customerId)
	if err != nil {
		return entities.Balance{}, err
	}

	var currentBalance int
	err = tx.QueryRow(ctx, "SELECT valor FROM saldos WHERE cliente_id=$1 FOR UPDATE", customerId).Scan(&currentBalance)
	if err != nil {
		return entities.Balance{}, err
	}

	newBalance, err := usecases.CalculateNewBalance(currentBalance, amount, opType, customer.AccountLimit)
	if err != nil {
		return entities.Balance{}, err
	}

	_, err = tx.Exec(ctx, "UPDATE saldos SET valor=$1 WHERE cliente_id=$2", newBalance, customerId)
	if err != nil {
		return entities.Balance{}, fmt.Errorf("updating customer balance: %w", err)
	}

	_, err = tx.Exec(ctx, "INSERT INTO transacoes (cliente_id, valor, tipo, descricao, realizada_em) VALUES ($1, $2, $3, $4, NOW())", customerId, amount, opType, description)
	if err != nil {
		return entities.Balance{}, fmt.Errorf("insert transaction error: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return entities.Balance{}, fmt.Errorf("committing transaction: %w", err)
	}

	return entities.Balance{Limite: customer.AccountLimit, Saldo: newBalance}, nil
}
