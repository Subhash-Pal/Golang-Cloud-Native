package main

import (
	"lab3-gin/routes"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	routes.RegisterUserRoutes(r)
	r.Run(":8080")
}