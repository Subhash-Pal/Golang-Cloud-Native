//TASK 4 — Request Validation + Struct Binding
package main

import (
	"net/http"
	"github.com/gin-gonic/gin"
)

type RegisterInput struct {
	Name  string `json:"name" binding:"required"`
	Email string `json:"email" binding:"required,email"`
	Age   int    `json:"age" binding:"gte=18"`
}

func main() {
	r := gin.Default()

	r.POST("/register", func(c *gin.Context) {
		var input RegisterInput

		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(200, gin.H{"message": "Valid input"})
	})

	r.Run(":8080")
}
