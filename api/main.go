package api

import (
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
)

var router = gin.Default()

func Run() {
	appEnv := os.Getenv("INSTANCE_PORT")
	getRoutes()
	router.Run(fmt.Sprintf(":%v", appEnv))
}

func getRoutes() {
	v1 := router.Group("/v1")
	addCostumerRoutes(v1)
}
