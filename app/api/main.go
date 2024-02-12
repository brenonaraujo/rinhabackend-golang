package api

import (
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
)

var router = gin.New()

func Run() {
	port := os.Getenv("HTTP_PORT")
	getRoutes()
	router.Run(fmt.Sprint(port))
}

func getRoutes() {
	v1 := router.Group("")
	addCustomerRoutes(v1)
}
