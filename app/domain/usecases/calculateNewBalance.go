package usecases

import (
	domain "brenonaraujo/rinhabackend-q12024/domain/entities"
	"fmt"
)

func CalculateNewBalance(currentBalance, amount int, opType domain.OperationType, accountLimit int) (int, error) {
	switch opType {
	case domain.Debit:
		if newBalance := currentBalance - amount; newBalance < -accountLimit {
			return 0, fmt.Errorf("debit amount %d would violate customer limit", amount)
		} else {
			return newBalance, nil
		}
	case domain.Credit:
		return currentBalance + amount, nil
	default:
		return 0, fmt.Errorf("unknown operation type: %s", opType)
	}
}
