package main

import "github.com/gin-gonic/gin"

func main() {
	r := gin.Default()

	r.GET("/list", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"products": []string{"Laptop", "Phone", "Tablet"},
		})
	})

	r.Run(":8002")
}