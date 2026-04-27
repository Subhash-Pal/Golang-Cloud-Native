package main

import (
	"net/http"

	"gin-jwt/middleware"
	"gin-jwt/utils"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	// Public route
	router.POST("/login", func(c *gin.Context) {

		// In real app → validate DB user
		token, err := utils.GenerateToken("admin")

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Token generation failed"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"token": token,
		})
	})

	// Protected route
	protected := router.Group("/api")
	protected.Use(middleware.AuthMiddleware())

	protected.GET("/protected", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Access granted",
		})
	})

	router.Run(":8080")
}

/*
Good—moving to **Gin + JWT** is the right step toward production structure. I’ll keep it **clean, modular, but not over-engineered**, so you can extend it later.

---

# 🧪 LAB: Gin-Based JWT Authentication (Production-Oriented, Minimal)

---

# 🎯 Objective

* Use Gin
* Implement:

  * `/login` → generate JWT
  * `/protected` → secured route
  * Middleware → validate JWT

---

# 📁 Project Structure (Simple but Scalable)

```bash
gin-jwt/
│── main.go
│── go.mod
│── middleware/
│     └── auth.go
│── utils/
│     └── jwt.go
```

---

# ⚙️ STEP 1 — Initialize Project

```powershell
mkdir gin-jwt
cd gin-jwt
go mod init gin-jwt
```

---

# ⚙️ STEP 2 — Install Dependencies

```powershell
go get github.com/gin-gonic/gin
go get github.com/golang-jwt/jwt/v5
```

👉 Validate:

```powershell
go mod tidy
```

---

# ⚙️ STEP 3 — JWT Utility

## 📄 `utils/jwt.go`

```go
package utils

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var SecretKey = []byte("my_secret_key")

func GenerateToken(username string) (string, error) {
	claims := jwt.MapClaims{
		"username": username,
		"exp":      time.Now().Add(time.Minute * 10).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(SecretKey)
}
```

---

# ⚙️ STEP 4 — Middleware

## 📄 `middleware/auth.go`

```go
package middleware

import (
	"net/http"
	"strings"

	"gin-jwt/utils"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func AuthMiddleware() gin.HandlerFunc {
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
			return utils.SecretKey, nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		c.Next()
	}
}
```

---

# ⚙️ STEP 5 — Main Application

## 📄 `main.go`

```go
package main

import (
	"net/http"

	"gin-jwt/middleware"
	"gin-jwt/utils"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	// Public route
	router.POST("/login", func(c *gin.Context) {

		// In real app → validate DB user
		token, err := utils.GenerateToken("admin")

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Token generation failed"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"token": token,
		})
	})

	// Protected route
	protected := router.Group("/api")
	protected.Use(middleware.AuthMiddleware())

	protected.GET("/protected", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Access granted",
		})
	})

	router.Run(":8080")
}
```

---

# ▶️ STEP 6 — Run Server

```powershell
go run main.go
```

👉 Expected:

```
[GIN-debug] Listening and serving HTTP on :8080
```

---

# 🧪 STEP 7 — Testing (PowerShell)

---

## 🔹 7.1 Login

```powershell
curl.exe -X POST http://localhost:8080/login
```

👉 Output:

```json
{"token":"<JWT_TOKEN>"}
```

---

## 🔹 7.2 Access Protected Route

```powershell
$token="<PASTE_TOKEN>"

curl.exe -H "Authorization: Bearer $token" http://localhost:8080/api/protected
```

👉 Output:

```json
{"message":"Access granted"}
```






Yes—you can chain **Step 7.1 (login)** and **Step 7.2 (use token)** into a single PowerShell flow. The clean way is: **call `/login` → extract token → reuse it immediately**.

Below are **two reliable approaches**.

---

# ✅ OPTION 1 — Using PowerShell (Recommended)

This avoids parsing issues and works consistently.

```powershell
$response = Invoke-RestMethod -Method POST -Uri "http://localhost:8080/login"

$token = $response.token

Invoke-RestMethod -Uri "http://localhost:8080/api/protected" `
    -Headers @{ Authorization = "Bearer $token" }
```

### ✅ Expected Output

```json
{
  "message": "Access granted"
}
```

👉 Why this works:

* `Invoke-RestMethod` automatically converts JSON → object
* `$response.token` directly extracts the JWT

---

# ✅ OPTION 2 — One-Liner (Compact)

```powershell
$token = (Invoke-RestMethod -Method POST http://localhost:8080/login).token; `
Invoke-RestMethod http://localhost:8080/api/protected -Headers @{Authorization="Bearer $token"}
```

---

# ⚠️ OPTION 3 — Using `curl.exe` (Not Recommended, but possible)

PowerShell cannot easily parse JSON from `curl.exe` without extra tools.

If you still want it:

```powershell
$token = (curl.exe -s -X POST http://localhost:8080/login | ConvertFrom-Json).token
curl.exe -H "Authorization: Bearer $token" http://localhost:8080/api/protected
```

---

# 🔍 Debug Check (Important)

If it fails, validate:

### 1. Token extracted?

```powershell
echo $token
```

### 2. Server running?

```powershell
curl.exe http://localhost:8080/login
```

---

# 🧠 Key Insight

You’ve essentially implemented a **client-side token pipeline**:

```text
/login → JWT → Authorization Header → Protected API
```

This is exactly how:

* frontend apps
* API clients
* microservices

communicate securely.

---


---

## 🔹 7.3 Negative Test

```powershell
curl.exe http://localhost:8080/api/protected
```

👉 Output:

```json
{"error":"Missing token"}
```

---

# 🧠 Production Notes (Important)

### 🔐 1. Move Secret to ENV

```powershell
$env:JWT_SECRET="super_secure"
```

Then read via `os.Getenv`.

---

### 🔐 2. Validate Signing Method (Security Fix)

Replace:

```go
return utils.SecretKey, nil
```

With:

```go
if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
	return nil, fmt.Errorf("unexpected signing method")
}
return utils.SecretKey, nil
```

---

### 🔐 3. Add Claims Struct (Better than MapClaims)

For stricter typing:

```go
type Claims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}
```

---

### 🔐 4. Add Expiry Handling

Already included via:

```go
"exp": time.Now().Add(...)
```

---

# ✅ Final Checklist

| Feature               | Status |
| --------------------- | ------ |
| Gin server            | ✔      |
| JWT generation        | ✔      |
| Middleware protection | ✔      |
| Route grouping        | ✔      |

---

*/
