package main

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

var secretKey = []byte("my_secret_key")

// =======================
// Request Struct
// =======================
type LoginRequest struct {
	Username string `json:"username" binding:"required,min=3"`
	Password string `json:"password" binding:"required,min=6"`
}

// =======================
// Standard Response Struct
// =======================
type APIResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Error   interface{} `json:"error,omitempty"`
}

// =======================
// Response Helpers
// =======================
func Success(c *gin.Context, message string, data interface{}) {
	c.JSON(http.StatusOK, APIResponse{
		Success: true,
		Message: message,
		Data:    data,
	})
}

func ErrorResponse(c *gin.Context, status int, message string, err interface{}) {
	c.JSON(status, APIResponse{
		Success: false,
		Message: message,
		Error:   err,
	})
}

// =======================
// JWT
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
// Middleware
// =======================
func authMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		authHeader := c.GetHeader("Authorization")

		if authHeader == "" {
			ErrorResponse(c, http.StatusUnauthorized, "Missing token", nil)
			c.Abort()
			return
		}

		if !strings.HasPrefix(authHeader, "Bearer ") {
			ErrorResponse(c, http.StatusUnauthorized, "Invalid token format", nil)
			c.Abort()
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method")
			}
			return secretKey, nil
		})

		if err != nil || !token.Valid {
			ErrorResponse(c, http.StatusUnauthorized, "Invalid token", nil)
			c.Abort()
			return
		}

		c.Next()
	}
}

// =======================
// Error Middleware
// =======================
func errorMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) > 0 {
			err := c.Errors.Last().Err
			ErrorResponse(c, http.StatusInternalServerError, "Internal server error", err.Error())
		}
	}
}

// =======================
// MAIN
// =======================
func main() {

	router := gin.Default()
	router.Use(errorMiddleware())

	// =======================
	// PUBLIC ROUTES
	// =======================
	router.POST("/login", func(c *gin.Context) {

		var req LoginRequest

		if err := c.ShouldBindJSON(&req); err != nil {
			ErrorResponse(c, http.StatusBadRequest, "Invalid request payload", nil)
			return
		}

		if req.Username != "admin" || req.Password != "password123" {
			ErrorResponse(c, http.StatusUnauthorized, "Invalid credentials", nil)
			return
		}

		token, err := generateToken(req.Username)
		if err != nil {
			c.Error(err)
			return
		}

		Success(c, "Login successful", gin.H{"token": token})
	})

	// =======================
	// VERSIONED API
	// =======================
	api := router.Group("/api")
	v1 := api.Group("/v1")
	v1.Use(authMiddleware())

	// Protected route
	v1.GET("/protected", func(c *gin.Context) {
		Success(c, "Access granted", nil)
	})

	// File upload
	v1.POST("/upload", func(c *gin.Context) {

		file, err := c.FormFile("file")
		if err != nil {
			ErrorResponse(c, http.StatusBadRequest, "File not provided", nil)
			return
		}

		ext := filepath.Ext(file.Filename)
		if ext != ".jpg" && ext != ".png" && ext != ".txt" {
			ErrorResponse(c, http.StatusBadRequest, "Invalid file type", nil)
			return
		}

		if err := os.MkdirAll("./uploads", os.ModePerm); err != nil {
			c.Error(err)
			return
		}

		filePath := "./uploads/" + file.Filename

		if err := c.SaveUploadedFile(file, filePath); err != nil {
			c.Error(errors.New("file save failed"))
			return
		}

		Success(c, "File uploaded", gin.H{
			"path": filePath,
		})
	})

	fmt.Println("Server running on :8080")
	router.Run(":8080")
}

/*
Below is a **clean lab manual + documentation** for **Hour 15 (API Versioning + Centralized Response Formatting)** built on your existing Gin project.

---

# 🧪 LAB MANUAL

## Hour 15 — API Versioning + Centralized Response System

---

# 🎯 Objective

Enhance your API to:

* Introduce **versioned endpoints** (`/api/v1/...`)
* Enforce a **standard response structure**
* Replace scattered responses with **central helper functions**

Using Gin

---

# 🧠 Concept Overview

---

## 🔹 1. API Versioning

### Before

```text
/api/protected
```

### After

```text
/api/v1/protected
```

👉 Why:

* Backward compatibility
* Safe upgrades
* Multiple client support

---

## 🔹 2. Centralized Response Format

### Standard Contract

```json
{
  "success": true,
  "message": "Operation successful",
  "data": {},
  "error": null
}
```

👉 All APIs must return this format.

---

# 📁 Project Context

You already have:

* JWT auth
* Middleware
* File upload
* Validation

👉 Now we are restructuring responses + routes.

---

# ⚙️ STEP 1 — Add Response Struct

Add in `main.go`:

```go
type APIResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Error   interface{} `json:"error,omitempty"`
}
```

---

# ⚙️ STEP 2 — Add Helper Functions

```go
func Success(c *gin.Context, message string, data interface{}) {
	c.JSON(http.StatusOK, APIResponse{
		Success: true,
		Message: message,
		Data:    data,
	})
}

func ErrorResponse(c *gin.Context, status int, message string, err interface{}) {
	c.JSON(status, APIResponse{
		Success: false,
		Message: message,
		Error:   err,
	})
}
```

---

# ⚙️ STEP 3 — Replace All Responses

---

## ❌ Old Style

```go
c.JSON(http.StatusOK, gin.H{"message": "ok"})
```

---

## ✅ New Standard

```go
Success(c, "Operation successful", data)
```

---

## ❌ Old Error

```go
c.JSON(http.StatusBadRequest, gin.H{"error": "invalid"})
```

---

## ✅ New Error

```go
ErrorResponse(c, http.StatusBadRequest, "Invalid request", nil)
```

---

# ⚙️ STEP 4 — Add API Versioning

---

## Replace:

```go
api := router.Group("/api")
```

---

## With:

```go
api := router.Group("/api")
v1 := api.Group("/v1")
v1.Use(authMiddleware())
```

---

## Move routes into `v1`

```go
v1.GET("/protected", handler)
v1.POST("/upload", handler)
```

---

# ▶️ STEP 5 — Run Server

```powershell
go run main.go
```

---

## ✅ Expected Output

```text
Server running on :8080
```

---

# 🧪 STEP 6 — Testing (Updated Endpoints)

---

## 🔹 6.1 Login

```powershell
$response = Invoke-RestMethod -Method POST `
-Uri "http://localhost:8080/login" `
-ContentType "application/json" `
-Body '{"username":"admin","password":"password123"}'
```

---

## 🔹 6.2 Extract Token

```powershell
$token = $response.data.token
```

---

## 🔹 6.3 Access Protected API (v1)

```powershell
Invoke-RestMethod `
-Uri "http://localhost:8080/api/v1/protected" `
-Headers @{Authorization="Bearer $token"}

```

---

## ✅ Expected Response

```json
{
  "success": true,
  "message": "Access granted"
}
```

---

## 🔹 6.4 Upload File (v1)

```powershell
curl.exe -X POST http://localhost:8080/api/v1/upload `
-H "Authorization: Bearer $token" `
-F "file=@test.txt"
```

---

## ✅ Expected Response

```json
{
  "success": true,
  "message": "File uploaded",
  "data": {
    "path": "./uploads/test.txt"
  }
}
```

---

# ⚠️ STEP 7 — Negative Testing

---

## ❌ Missing Token

```powershell
curl.exe http://localhost:8080/api/v1/protected
```

✔ Output:

```json
{
  "success": false,
  "message": "Missing token"
}
```

---

## ❌ Invalid Version (Route Not Found)

```powershell
curl.exe http://localhost:8080/api/v2/protected
```

✔ Output:

```text
404 page not found
```

---

## ❌ Invalid Payload

```powershell
curl.exe -X POST http://localhost:8080/login `
-H "Content-Type: application/json" `
-d "{}"
```

✔ Output:

```json
{
  "success": false,
  "message": "Invalid request payload"
}
```

---

# 🧠 Architecture Insight

---

## 🔹 Route Hierarchy

```text
/
/login
/api
   /v1
      /protected
      /upload
```

---

## 🔹 Response Flow

```text
Handler → Success()/ErrorResponse() → JSON Output
```

---

## 🔹 Middleware Flow

```text
Request → Gin → Auth Middleware → Handler → Response Helper
```

---

# 🔐 Production Considerations

---

## 1. Versioning Strategy

* Maintain `/v1` stable
* Introduce `/v2` for breaking changes
* Never break existing clients

---

## 2. Response Contract Stability

Once defined:

```json
success, message, data, error
```

👉 **Do not change structure later**

---

## 3. Error Handling

* Avoid exposing internal errors in `error` field
* Log internally instead

---

# 🚀 Learning Outcome

You now have:

| Feature                        | Status |
| ------------------------------ | ------ |
| API Versioning                 | ✔      |
| Unified Response Format        | ✔      |
| Clean Response Helpers         | ✔      |
| Backward-compatible API design | ✔      |

---


---

*/
