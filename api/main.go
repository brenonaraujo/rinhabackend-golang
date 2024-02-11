package api

import "github.com/gin-gonic/gin"

var router = gin.New()

func Run() {
	getRoutes()
	router.Run(":33888")
}

func getRoutes() {
	v1 := router.Group("")
	addCustomerRoutes(v1)
}
