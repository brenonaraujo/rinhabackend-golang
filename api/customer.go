package api

import (
	"brenonaraujo/rinhabackend-q12024/infra/database"
	"context"
	"fmt"
	"net/http"

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

		customerId := c.Param("id")
		if !customerExists(customerId) {
			c.JSON(http.StatusNotFound, "")
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
