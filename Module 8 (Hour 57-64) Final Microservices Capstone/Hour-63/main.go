package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

var ctx = context.Background()

func main() {
	// Redis client
	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	// Simple in-memory credential (keep it minimal)
	username := "shubh"
	password := "1234"

	router := gin.Default()

	router.POST("/login", func(c *gin.Context) {
		var req map[string]string
		if err := c.BindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
			return
		}

		if req["username"] != username || req["password"] != password {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
			return
		}

		// ✅ Publish event
		event := "user_logged_in:" + username
		err := rdb.Publish(ctx, "user_events", event).Err()
		if err != nil {
			fmt.Println("Redis publish error:", err)
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "login successful",
		})
	})

	fmt.Println("🚀 API running on http://localhost:8080")
	router.Run(":8081")
}