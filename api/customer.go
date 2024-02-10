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
		var dto CustomerDto
		if err := c.ShouldBindJSON(&dto); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		customerId, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid customer ID"})
			return
		}

		if dto.Tipo != "c" && dto.Tipo != "d" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid transaction type"})
			return
		}

		if dto.Tipo == "d" && !canDeduct(customerId, dto.Valor) {
			c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "Insufficient funds"})
			return
		}

		message := fmt.Sprintf("[Transaction received] id:%v, value:%v, kind:%v, description:%v",
			customerId, dto.Valor, dto.Tipo, dto.Descricao)
		c.JSON(http.StatusOK, gin.H{"message": message})
	})

	customer.GET("/:id/extrato", func(c *gin.Context) {
		c.JSON(http.StatusOK, "extrato")
	})
}

func customerExists(customerId string) bool {
	var exists bool
	err := database.GetDBPool().QueryRow(context.Background(),
		"SELECT EXISTS(SELECT 1 FROM clientes WHERE id=$1)", customerId).Scan(&exists)
	return err == nil && exists
}

func AddBalance(customerId, amount int) {
	// Implement the logic to add balance to the customer's account
}

func DeductBalance(customerId, amount int) {
	// Implement the logic to deduct balance from the customer's account ensuring the balance does not go below the limit
}

func canDeduct(customerId, amount int) bool {
	var currentBalance, limit int
	err := database.GetDBPool().QueryRow(context.Background(),
		"SELECT saldo, limite FROM clientes WHERE id=$1", customerId).Scan(&currentBalance, &limit)
	if err != nil {
		fmt.Println("Error retrieving customer balance and limit:", err)
		return false
	}
	newBalance := currentBalance - amount
	return newBalance >= -limit
}

func GetCustomerBalance(customerId int) (CustomerBalanceDto, error) {
	// Implement the logic to retrieve the customer's current balance and limit
	return CustomerBalanceDto{Limite: 100000, Saldo: -9098}, nil // Placeholder return, replace with actual logic
}
