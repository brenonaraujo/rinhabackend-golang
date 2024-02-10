package api

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func addCostumerRoutes(rg *gin.RouterGroup) {
	costumer := rg.Group("/clientes")
	appEnv := os.Getenv("INSTANCE_ID")

	costumer.POST("/:id/transacoes", func(c *gin.Context) {
		var dto CostumerDto
		if err := c.ShouldBindJSON(&dto); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		costumerId := c.Param("id")
		message := fmt.Sprintf("[%v][Transaction received] id:%v, value:%v, kind:%v, description:%v",
			appEnv, costumerId, dto.Valor, dto.Tipo, dto.Descricao)

		c.JSON(http.StatusOK, gin.H{"message": message})
	})

	costumer.GET("/:id/extrato", func(c *gin.Context) {
		c.JSON(http.StatusOK, "extrato")
	})

}
