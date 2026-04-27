package main

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5" //go get -u github.com/golang-jwt/jwt/v5
)

var secretKey = []byte("my_secret_key")

// =======================
// Request Struct (Validation)
// =======================
type LoginRequest struct {
	Username string `json:"username" binding:"required,min=3"`
	Password string `json:"password" binding:"required,min=6"`
}

// =======================
// JWT Token Generation
// =======================
func generateToken(username string) (string, error) {

	claims := jwt.MapClaims{
		"username": username,
		"exp":      time.Now().Add(time.Minute * 10).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secretKey)
}

// =======================
// JWT Middleware
// =======================
func authMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		authHeader := c.GetHeader("Authorization")

		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing token"})
			c.Abort()
			return
		}

		if !strings.HasPrefix(authHeader, "Bearer ") {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token format"})
			c.Abort()
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {

			// Validate signing method
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method")
			}

			return secretKey, nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		c.Next()
	}
}

// =======================
// MAIN
// =======================
func main() {

	router := gin.Default()

	// =======================
	// Login Route
	// =======================
	router.POST("/login", func(c *gin.Context) {

		var req LoginRequest

		// Bind + Validate JSON
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid request payload",
			})
			return
		}

		// Dummy authentication (replace with DB later)
		if req.Username != "admin" || req.Password != "password123" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid credentials",
			})
			return
		}

		token, err := generateToken(req.Username)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Token generation failed",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"token": token,
		})
	})

	// =======================
	// Protected Routes
	// =======================
	api := router.Group("/api")
	api.Use(authMiddleware())

	api.GET("/protected", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Access granted",
		})
	})

	// =======================
	// Start Server
	// =======================
	fmt.Println("Server running on :8080")
	router.Run(":8080")
}

/*

Here is a **complete `main.go`** that integrates everything you’ve built so far:

* Gin server
* JWT generation
* Middleware protection
* Struct binding + validation

👉 This is **single-file (for simplicity)** but still follows production patterns.

---

# 📄 `main.go` (Full Working Code)

```go
package main

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

var secretKey = []byte("my_secret_key")

// =======================
// Request Struct (Validation)
// =======================
type LoginRequest struct {
	Username string `json:"username" binding:"required,min=3"`
	Password string `json:"password" binding:"required,min=6"`
}

// =======================
// JWT Token Generation
// =======================
func generateToken(username string) (string, error) {

	claims := jwt.MapClaims{
		"username": username,
		"exp":      time.Now().Add(time.Minute * 10).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secretKey)
}

// =======================
// JWT Middleware
// =======================
func authMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		authHeader := c.GetHeader("Authorization")

		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing token"})
			c.Abort()
			return
		}

		if !strings.HasPrefix(authHeader, "Bearer ") {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token format"})
			c.Abort()
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {

			// Validate signing method
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method")
			}

			return secretKey, nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		c.Next()
	}
}

// =======================
// MAIN
// =======================
func main() {

	router := gin.Default()

	// =======================
	// Login Route
	// =======================
	router.POST("/login", func(c *gin.Context) {

		var req LoginRequest

		// Bind + Validate JSON
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid request payload",
			})
			return
		}

		// Dummy authentication (replace with DB later)
		if req.Username != "admin" || req.Password != "password123" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid credentials",
			})
			return
		}

		token, err := generateToken(req.Username)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Token generation failed",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"token": token,
		})
	})

	// =======================
	// Protected Routes
	// =======================
	api := router.Group("/api")
	api.Use(authMiddleware())

	api.GET("/protected", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Access granted",
		})
	})

	// =======================
	// Start Server
	// =======================
	fmt.Println("Server running on :8080")
	router.Run(":8080")
}
```

---

# ▶️ Run

```powershell
go run main.go
```

---

# 🧪 Test Quickly (All-in-One Flow)

```powershell
step 1
$response = Invoke-RestMethod -Method POST `
-Uri "http://localhost:8080/login" `
-ContentType "application/json" `
-Body '{"username":"admin","password":"password123"}'

step 2
$token = $response.token

Step 3

Invoke-RestMethod `
-Uri "http://localhost:8080/api/protected" `
-Headers @{Authorization="Bearer $token"}
```

---

# 🧠 What This Covers

* ✔ Struct binding (`ShouldBindJSON`)
* ✔ Validation (`required`, `min`)
* ✔ JWT creation
* ✔ Middleware protection
* ✔ Route grouping

---

# ⚠️ Real Production Gaps (Next Step)

This is **clean but still minimal**. Missing for real production:

* Env-based secret (`os.Getenv`)
* Password hashing (bcrypt)
* DB authentication
* Refresh tokens
* Structured error responses

---

*/

