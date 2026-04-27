//TASK 3 — JWT Authentication + Refresh Token
package main

import (
	"time"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)
var secret = []byte("secret")

// Helper function to generate both Access and Refresh tokens
func generateToken(username string) (string, string) {
	accessClaims := jwt.MapClaims{
		"user": username,
		"exp":  time.Now().Add(time.Minute * 5).Unix(), // Expires in 5 mins
	}
	refreshClaims := jwt.MapClaims{
		"user": username,
		"exp":  time.Now().Add(time.Hour * 24).Unix(), // Expires in 24 hours
	}

	access := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	refresh := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)

	accessToken, _ := access.SignedString(secret)
	refreshToken, _ := refresh.SignedString(secret)

	return accessToken, refreshToken
}

func main() {
	r := gin.Default()

	// LOGIN ENDPOINT
	r.POST("/login", func(c *gin.Context) {
		// In a real app, you'd verify a password here
		access, refresh := generateToken("admin")
		c.JSON(http.StatusOK, gin.H{
			"access":  access,
			"refresh": refresh,
		})
	})

	// REFRESH ENDPOINT
	r.POST("/refresh", func(c *gin.Context) {
		var body map[string]string
		
		// 1. Check if JSON is valid
		if err := c.ShouldBindJSON(&body); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
			return
		}

		tokenStr := body["refresh"]
		if tokenStr == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Refresh token missing"})
			return
		}

		// 2. Parse and Validate the token
		token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
			return secret, nil
		})

		// 3. Check if token is valid or expired
		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired refresh token"})
			return
		}

		// 4. Safely extract claims and generate new access token
		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			username, ok := claims["user"].(string)
			if !ok {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid token data"})
				return
			}
			
			newAccess, _ := generateToken(username)
			c.JSON(http.StatusOK, gin.H{"access": newAccess})
		}
	})

	r.Run(":8080")
}

