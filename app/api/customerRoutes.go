package api

import (
	"brenonaraujo/rinhabackend-q12024/domain"
	"brenonaraujo/rinhabackend-q12024/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func addCustomerRoutes(rg *gin.RouterGroup) {
	customer := rg.Group("/clientes")

	customer.POST("/:id/transacoes", func(c *gin.Context) {
		var transaction TransactionRequest
		if err := c.ShouldBindJSON(&transaction); err != nil {
			c.Status(http.StatusUnprocessableEntity)
			return
		}
		customerId, err := strconv.Atoi(c.Param("id"))
		if !customerExists(customerId) {
			c.Status(http.StatusNotFound)
			return
		}

		var result domain.Balance
		if transaction.Tipo == "d" {
			result, err = service.DeductBalance(customerId, transaction.Valor, transaction.Descricao)
			if err != nil {
				c.Status(http.StatusUnprocessableEntity)
				return
			}
		} else {
			result, err = service.AddBalance(customerId, transaction.Valor, transaction.Descricao)
			if err != nil {
				c.Status(http.StatusUnprocessableEntity)
				return
			}
		}

		c.JSON(http.StatusOK, result)
	})

	customer.GET("/:id/extrato", func(c *gin.Context) {
		customerId, err := strconv.Atoi(c.Param("id"))
		if err != nil || !customerExists(customerId) {
			c.Status(http.StatusNotFound)
			return
		}
		var customerStatement service.CustomerStatement
		customerStatement, err = service.GetCustomerStatement(customerId)
		if err != nil {
			c.Status(http.StatusInternalServerError)
			return
		}

		c.JSON(http.StatusOK, customerStatement)
	})
}

func customerExists(customerId int) bool {
	_, err := service.GetCustomer(customerId)
	if err != nil {
		return false
	}
	return true
}
