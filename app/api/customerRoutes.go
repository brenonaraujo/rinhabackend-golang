package api

import (
	"net/http"
	"strconv"

	"brenonaraujo/rinhabackend-q12024/domain/entities"
	"brenonaraujo/rinhabackend-q12024/service"

	"github.com/gin-gonic/gin"
)

func addCustomerRoutes(rg *gin.RouterGroup) {
	customer := rg.Group("/clientes")

	customer.POST("/:id/transacoes", handleCreateTransaction)
	customer.GET("/:id/extrato", handleGetCustomerStatement)
}

func handleCreateTransaction(c *gin.Context) {
	var tx TransactionRequest
	if err := c.ShouldBindJSON(&tx); err != nil {
		c.Status(http.StatusUnprocessableEntity)
		return
	}

	customerID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	if !customerExists(customerID) {
		c.Status(http.StatusNotFound)
		return
	}

	result, err := service.TransactionProcess(c.Request.Context(), customerID, tx.Valor, tx.Descricao, entities.OperationType(tx.Tipo))
	if err != nil {
		c.Status(http.StatusUnprocessableEntity)
		return
	}

	c.JSON(http.StatusOK, result)
}

func handleGetCustomerStatement(c *gin.Context) {
	customerID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	if !customerExists(customerID) {
		c.Status(http.StatusNotFound)
		return
	}

	customerStatement, err := service.GetCustomerStatement(customerID)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, customerStatement)
}

func customerExists(customerID int) bool {
	_, err := service.GetCustomer(customerID)
	return err == nil
}
