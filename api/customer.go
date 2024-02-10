package api

import (
	"brenonaraujo/rinhabackend-q12024/infra/database"
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func addCustomerRoutes(rg *gin.RouterGroup) {
	customer := rg.Group("/clientes")

	customer.POST("/:id/transacoes", func(c *gin.Context) {
		var dto TransactionDto
		if err := c.ShouldBindJSON(&dto); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		customerId, err := strconv.Atoi(c.Param("id"))
		if !customerExists(customerId) {
			c.JSON(http.StatusNotFound, gin.H{})
			return
		}

		var result CustomerBalanceDto
		if dto.Tipo == "d" {
			result, err = DeductBalance(customerId, dto.Valor)
			if err != nil {
				c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
				return
			}
		} else if dto.Tipo == "c" {
			// Presumably, you would have a similar function for adding balance
			// e.g., err = AddBalance(customerId, dto.Valor)
		}

		c.JSON(http.StatusOK, result)
	})

	customer.GET("/:id/extrato", func(c *gin.Context) {
		c.JSON(http.StatusOK, "extrato")
	})
}

func customerExists(customerId int) bool {
	var exists bool
	err := database.GetDBPool().QueryRow(context.Background(),
		"SELECT EXISTS(SELECT 1 FROM clientes WHERE id=$1)", customerId).Scan(&exists)
	return err == nil && exists
}

func DeductBalance(customerId, amount int) (CustomerBalanceDto, error) {
	ctx := context.Background()
	var costumerBalance CustomerBalanceDto

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

	var limit int
	err = tx.QueryRow(ctx,
		"SELECT limite FROM clientes WHERE id=$1", customerId).Scan(&limit)
	if err != nil {
		tx.Rollback(ctx)
		return costumerBalance, fmt.Errorf("querying customer limit: %w", err)
	}

	// Calculate new balance and check if deduction is possible
	newBalance := currentBalance - amount
	if newBalance < -limit {
		tx.Rollback(ctx)
		return costumerBalance, fmt.Errorf("deduction amount %d would violate customer limit", amount)
	}

	// Update the customer's balance
	_, err = tx.Exec(ctx,
		"UPDATE saldos SET valor=$1 WHERE id=$2", newBalance, customerId)
	if err != nil {
		tx.Rollback(ctx) // Rollback transaction on error
		return costumerBalance, fmt.Errorf("updating customer balance: %w", err)
	}

	// Commit the transaction
	if err := tx.Commit(ctx); err != nil {
		return costumerBalance, fmt.Errorf("committing transaction: %w", err)
	}

	return CustomerBalanceDto{Limite: limit, Saldo: newBalance}, nil
}

func GetCustomerBalance(customerId int) (CustomerBalanceDto, error) {
	// Implement the logic to retrieve the customer's current balance and limit
	return CustomerBalanceDto{Limite: 100000, Saldo: -9098}, nil // Placeholder return, replace with actual logic
}
