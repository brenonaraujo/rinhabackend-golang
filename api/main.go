package api

import "github.com/gin-gonic/gin"

var router = gin.Default()

func Run() {
	getRoutes()
	router.Run(":33888")
}

func getRoutes() {
	v1 := router.Group("/v1")
	addCostumerRoutes(v1)
}
