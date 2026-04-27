package main

import (
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	r.GET("/login", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Auth Service is working",
		})
	})

	r.Run(":8001")
}