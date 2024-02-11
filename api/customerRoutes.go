package api

import (
	"brenonaraujo/rinhabackend-q12024/domain"
	"brenonaraujo/rinhabackend-q12024/infra/database"
	"brenonaraujo/rinhabackend-q12024/service"
	"context"
	"log"
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

		var result domain.Balance
		if dto.Tipo == "d" {
			result, err = service.DeductBalance(customerId, dto.Valor, dto.Descricao)
			if err != nil {
				c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
				return
			}
		} else if dto.Tipo == "c" {
			result, err = service.AddBalance(customerId, dto.Valor)
			if err != nil {
				c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
				return
			}
		}

		c.JSON(http.StatusOK, result)
	})

	customer.GET("/:id/extrato", func(c *gin.Context) {
		customerId, err := strconv.Atoi(c.Param("id"))
		if err != nil || !customerExists(customerId) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Customer not found"})
			return
		}
		log.Printf("Handling extrato request!")
		var customerStatement service.CustomerStatement
		customerStatement, err = service.GetCustomerStatement(customerId)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, customerStatement)
	})
}

func customerExists(customerId int) bool {
	var exists bool
	err := database.GetDBPool().QueryRow(context.Background(),
		"SELECT EXISTS(SELECT 1 FROM clientes WHERE id=$1)", customerId).Scan(&exists)
	return err == nil && exists
}
