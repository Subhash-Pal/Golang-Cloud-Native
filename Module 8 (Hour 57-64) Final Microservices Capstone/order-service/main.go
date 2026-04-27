package main

import (
	"github.com/gin-gonic/gin"
)

type OrderRequest struct {
	UserID   string `json:"user_id"`
	Product  string `json:"product"`
	Quantity int    `json:"quantity"`
}

func main() {
	r := gin.Default()

	r.POST("/create", func(c *gin.Context) {

		var req OrderRequest

		// Bind JSON body
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(400, gin.H{
				"error": "Invalid request body",
			})
			return
		}

		// Simulated order creation
		c.JSON(200, gin.H{
			"order_id":  1001,
			"user_id":   req.UserID,
			"product":   req.Product,
			"quantity":  req.Quantity,
			"status":    "created",
		})
	})

	r.Run(":8003")
}