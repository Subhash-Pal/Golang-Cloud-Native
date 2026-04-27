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
// Request Struct
// =======================
type LoginRequest struct {
	Username string `json:"username" binding:"required,min=3"`
	Password string `json:"password" binding:"required,min=6"`
}

// =======================
// JWT Claims
// =======================
type Claims struct {
	Username string `json:"username"`
	Type     string `json:"type"` // access / refresh
	jwt.RegisteredClaims
}

// =======================
// Standard Response
// =======================
type APIResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Error   interface{} `json:"error,omitempty"`
}

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
// TOKEN GENERATION
// =======================
func generateAccessToken(username string) (string, error) {
	claims := Claims{
		Username: username,
		Type:     "access",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(10 * time.Minute)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secretKey)
}

func generateRefreshToken(username string) (string, error) {
	claims := Claims{
		Username: username,
		Type:     "refresh",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(1 * time.Hour)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secretKey)
}

// =======================
// AUTH MIDDLEWARE
// =======================
func authMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		authHeader := c.GetHeader("Authorization")

		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			ErrorResponse(c, http.StatusUnauthorized, "Missing or invalid token", nil)
			c.Abort()
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(t *jwt.Token) (interface{}, error) {
			return secretKey, nil
		})

		if err != nil || !token.Valid {
			ErrorResponse(c, http.StatusUnauthorized, "Invalid token", nil)
			c.Abort()
			return
		}

		claims, ok := token.Claims.(*Claims)
		if !ok || claims.Type != "access" {
			ErrorResponse(c, http.StatusUnauthorized, "Invalid access token", nil)
			c.Abort()
			return
		}

		c.Set("username", claims.Username)
		c.Next()
	}
}

// =======================
// MAIN
// =======================
func main() {

	router := gin.Default()

	// =======================
	// LOGIN
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

		access, err1 := generateAccessToken(req.Username)
		refresh, err2 := generateRefreshToken(req.Username)

		if err1 != nil || err2 != nil {
			ErrorResponse(c, http.StatusInternalServerError, "Token generation failed", nil)
			return
		}

		Success(c, "Login successful", gin.H{
			"access_token":  access,
			"refresh_token": refresh,
		})
	})

	// =======================
	// REFRESH TOKEN
	// =======================
	router.POST("/refresh", func(c *gin.Context) {

		var body struct {
			RefreshToken string `json:"refresh_token" binding:"required"`
		}

		if err := c.ShouldBindJSON(&body); err != nil {
			ErrorResponse(c, http.StatusBadRequest, "Invalid request payload", nil)
			return
		}

		token, err := jwt.ParseWithClaims(body.RefreshToken, &Claims{}, func(t *jwt.Token) (interface{}, error) {
			return secretKey, nil
		})

		if err != nil || !token.Valid {
			ErrorResponse(c, http.StatusUnauthorized, "Invalid refresh token", nil)
			return
		}

		claims, ok := token.Claims.(*Claims)
		if !ok || claims.Type != "refresh" {
			ErrorResponse(c, http.StatusUnauthorized, "Invalid token type", nil)
			return
		}

		newAccess, err := generateAccessToken(claims.Username)
		if err != nil {
			ErrorResponse(c, http.StatusInternalServerError, "Failed to generate access token", nil)
			return
		}

		Success(c, "Token refreshed", gin.H{
			"access_token": newAccess,
		})
	})

	// =======================
	// VERSIONED PROTECTED API
	// =======================
	api := router.Group("/api")
	v1 := api.Group("/v1")
	v1.Use(authMiddleware())

	v1.GET("/protected", func(c *gin.Context) {
		user, _ := c.Get("username")

		Success(c, "Access granted", gin.H{
			"user": user,
		})
	})

	// =======================
	// START SERVER
	// =======================
	fmt.Println("Server running on :8080")
	router.Run(":8080")
}

/*
Here are the step-by-step PowerShell commands to interact with all HTTP services in your Go application. **Run these in a new PowerShell session** (after starting your server with `go run main.go`):

---

### 1. **Login to Get Tokens** (Valid Credentials)
```powershell
# Send login request
$loginResponse = Invoke-RestMethod -Uri "http://localhost:8080/login" -Method Post -Body (@{
    username = "admin"
    password = "password123"
} | ConvertTo-Json) -ContentType "application/json"

# Extract tokens
$accessToken = $loginResponse.data.access_token
$refreshToken = $loginResponse.data.refresh_token

# View tokens (for debugging)
$accessToken
$refreshToken
```

> ✅ **Expected Output**:  
> `access_token` (JWT string starting with `ey...`) and `refresh_token` (similar JWT)

---

### 2. **Access Protected Endpoint** (Using Access Token)
```powershell
# Call protected API with access token
Invoke-RestMethod -Uri "http://localhost:8080/api/v1/protected" -Method Get -Headers @{
    Authorization = "Bearer $accessToken"
}
```

> ✅ **Expected Output**:  
> `@{success=True; message=Access granted; data=}` with `user=admin` in data

---

### 3. **Refresh Access Token** (Using Refresh Token)
```powershell
# Refresh token request
$refreshResponse = Invoke-RestMethod -Uri "http://localhost:8080/refresh" -Method Post -Body (@{
    refresh_token = $refreshToken
} | ConvertTo-Json) -ContentType "application/json"

# Update access token
$accessToken = $refreshResponse.data.access_token

# Verify new token works
Invoke-RestMethod -Uri "http://localhost:8080/api/v1/protected" -Method Get -Headers @{
    Authorization = "Bearer $accessToken"
}
```

> ✅ **Expected Output**:  
> New `access_token` and successful protected endpoint access

---

### 🔥 Critical Error Cases to Test (Optional)
#### ❌ Invalid Credentials
```powershell
Invoke-RestMethod -Uri "http://localhost:8080/login" -Method Post -Body (@{
    username = "hacker"
    password = "wrong"
} | ConvertTo-Json) -ContentType "application/json" -ErrorAction Stop
```
> 💡 **Expected**: `401 Unauthorized` with `"Invalid credentials"`

#### ❌ Missing Token
```powershell
Invoke-RestMethod -Uri "http://localhost:8080/api/v1/protected" -Method Get -ErrorAction Stop
```
> 💡 **Expected**: `401 Unauthorized` with `"Missing or invalid token"`

#### ❌ Invalid Token Type (Using Refresh Token as Access Token)
```powershell
Invoke-RestMethod -Uri "http://localhost:8080/api/v1/protected" -Method Get -Headers @{
    Authorization = "Bearer $refreshToken"
} -ErrorAction Stop
```
> 💡 **Expected**: `401 Unauthorized` with `"Invalid access token"`

---

### ⚙️ How to Start the Server First
1. Open **PowerShell** in your project directory
2. Start the server:
   ```powershell
   go run main.go
   ```
   > ✅ You'll see: `Server running on :8080`

3. **Keep this window open** and run the API commands in a **new PowerShell window**

---

### 📝 Key Notes
1. **Token Expiry**:
   - Access token expires in **10 minutes** (re-run Step 3 to refresh)
   - Refresh token expires in **1 hour** (re-run Step 1 to get new tokens)
   
2. **PowerShell Tips**:
   - Use `-ErrorAction Stop` to see error responses clearly
   - Tokens are stored in PowerShell variables (`$accessToken`, `$refreshToken`)
   - All commands assume server is running on `localhost:8080`

3. **Security Note**:  
   This uses a hardcoded secret key (`my_secret_key`) – **never use this in production**. Rotate keys in real systems.

---

### 🚀 Full Successful Flow Example
```powershell
# 1. Login
$login = Invoke-RestMethod -Uri "http://localhost:8080/login" -Method Post -Body (@{u="admin";p="password123"}|ConvertTo-Json) -ContentType "application/json"
$at = $login.data.access_token; $rt = $login.data.refresh_token

# 2. Access protected endpoint
Invoke-RestMethod -Uri "http://localhost:8080/api/v1/protected" -Headers @{A="Bearer $at"}

# 3. Refresh token
$refresh = Invoke-RestMethod -Uri "http://localhost:8080/refresh" -Body (@{refresh_token=$rt}|ConvertTo-Json) -ContentType "application/json"
$at = $refresh.data.access_token

# 4. Access endpoint again with new token
Invoke-RestMethod -Uri "http://localhost:8080/api/v1/protected" -Headers @{A="Bearer $at"}
```

> 💡 **Pro Tip**: Shorten variable names (`$at`, `$rt`) for faster typing during testing!
*/