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
// VERSIONED ROUTES
// =======================
func v1Routes(api *gin.RouterGroup) {
	v1 := api.Group("/v1")
	v1.Use(authMiddleware())

	v1.GET("/protected", func(c *gin.Context) {
		user, _ := c.Get("username")
		Success(c, "Access granted", gin.H{"user": user})
	})
}

func v2Routes(api *gin.RouterGroup) {
	v2 := api.Group("/v2")
	v2.Use(authMiddleware())

	v2.GET("/protected", func(c *gin.Context) {
		user, _ := c.Get("username")
		Success(c, "Access granted with new features", gin.H{"user": user, "version": "v2"})
	})
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
	// VERSIONED ROUTES
	// =======================
	api := router.Group("/api")
	v1Routes(api)
	v2Routes(api)

	// =======================
	// START SERVER
	// =======================
	fmt.Println("Server running on :8080")
	router.Run(":8080")
}

/*
After running your Go server (with `go run main.go`), you can test the API endpoints using PowerShell. Below are step-by-step commands to interact with the versioned API (`/api/v1` and `/api/v2`) after starting the server.

---

### **Step 1: Start the Server**
Ensure your server is running in a separate terminal:
```powershell
go run main.go
```
> ✅ You'll see: `Server running on :8080`

---

### **Step 2: Login to Get Tokens**
Log in to get the `access_token` and `refresh_token`:
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

### **Step 3: Access Versioned Endpoints**

#### **Access `/api/v1/protected`**
```powershell
Invoke-RestMethod -Uri "http://localhost:8080/api/v1/protected" -Method Get -Headers @{
    Authorization = "Bearer $accessToken"
}
```

> ✅ **Expected Output**:
```json
{
  "success": true,
  "message": "Access granted",
  "data": {
    "user": "admin"
  }
}
```

---

#### **Access `/api/v2/protected`**
```powershell
Invoke-RestMethod -Uri "http://localhost:8080/api/v2/protected" -Method Get -Headers @{
    Authorization = "Bearer $accessToken"
}
```

> ✅ **Expected Output**:
```json
{
  "success": true,
  "message": "Access granted with new features",
  "data": {
    "user": "admin",
    "version": "v2"
  }
}
```

---

### **Step 4: Refresh Access Token**
If the access token expires, refresh it using the refresh token:
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
> New `access_token` and successful protected endpoint access.

---

### **Step 5: Error Scenarios (Optional)**

#### ❌ Invalid Credentials
```powershell
Invoke-RestMethod -Uri "http://localhost:8080/login" -Method Post -Body (@{
    username = "hacker"
    password = "wrong"
} | ConvertTo-Json) -ContentType "application/json" -ErrorAction Stop
```
> 💡 **Expected**: `401 Unauthorized` with `"Invalid credentials"`

---

#### ❌ Missing Token
```powershell
Invoke-RestMethod -Uri "http://localhost:8080/api/v1/protected" -Method Get -ErrorAction Stop
```
> 💡 **Expected**: `401 Unauthorized` with `"Missing or invalid token"`

---

#### ❌ Invalid Token Type (Using Refresh Token as Access Token)
```powershell
Invoke-RestMethod -Uri "http://localhost:8080/api/v1/protected" -Method Get -Headers @{
    Authorization = "Bearer $refreshToken"
} -ErrorAction Stop
```
> 💡 **Expected**: `401 Unauthorized` with `"Invalid access token"`

---

### **Full Successful Flow Example**
Here’s how you can test the entire flow in one go:

```powershell
# Step 1: Login
$login = Invoke-RestMethod -Uri "http://localhost:8080/login" -Method Post -Body (@{
    username = "admin"
    password = "password123"
} | ConvertTo-Json) -ContentType "application/json"
$accessToken = $login.data.access_token
$refreshToken = $login.data.refresh_token

# Step 2: Access v1 endpoint
Invoke-RestMethod -Uri "http://localhost:8080/api/v1/protected" -Method Get -Headers @{
    Authorization = "Bearer $accessToken"
}

# Step 3: Access v2 endpoint
Invoke-RestMethod -Uri "http://localhost:8080/api/v2/protected" -Method Get -Headers @{
    Authorization = "Bearer $accessToken"
}

# Step 4: Refresh token
$refresh = Invoke-RestMethod -Uri "http://localhost:8080/refresh" -Method Post -Body (@{
    refresh_token = $refreshToken
} | ConvertTo-Json) -ContentType "application/json"
$accessToken = $refresh.data.access_token

# Step 5: Access v1 endpoint again with new token
Invoke-RestMethod -Uri "http://localhost:8080/api/v1/protected" -Method Get -Headers @{
    Authorization = "Bearer $accessToken"
}
```

---

### **Key Notes**
1. **Token Expiry**:
   - Access token expires in **10 minutes** (re-run Step 4 to refresh).
   - Refresh token expires in **1 hour** (re-run Step 2 to get new tokens).

2. **PowerShell Tips**:
   - Use `-ErrorAction Stop` to see error responses clearly.
   - Tokens are stored in PowerShell variables (`$accessToken`, `$refreshToken`).
   - All commands assume the server is running on `localhost:8080`.

3. **Security Note**:
   - This uses a hardcoded secret key (`my_secret_key`) – **never use this in production**. Rotate keys in real systems.

---

By following these steps, you can fully test both `/api/v1` and `/api/v2` endpoints in your Go application!
*/

/*
Hands on 
# Step 1: Login
$login = Invoke-RestMethod -Uri "http://localhost:8080/login" -Method Post -Body (@{
    username = "admin"
    password = "password123"
} | ConvertTo-Json) -ContentType "application/json"
$accessToken = $login.data.access_token
$refreshToken = $login.data.refresh_token

# Step 2: Access v1 endpoint
Invoke-RestMethod -Uri "http://localhost:8080/api/v1/protected" -Method Get -Headers @{
    Authorization = "Bearer $accessToken"
}

# Step 3: Access v2 endpoint
Invoke-RestMethod -Uri "http://localhost:8080/api/v2/protected" -Method Get -Headers @{
    Authorization = "Bearer $accessToken"
}

# Step 4: Refresh token
$refresh = Invoke-RestMethod -Uri "http://localhost:8080/refresh" -Method Post -Body (@{
    refresh_token = $refreshToken
} | ConvertTo-Json) -ContentType "application/json"
$accessToken = $refresh.data.access_token

# Step 5: Access v1 endpoint again with new token
Invoke-RestMethod -Uri "http://localhost:8080/api/v1/protected" -Method Get -Headers @{
    Authorization = "Bearer $accessToken"
}
*/

/*
Here’s the **complete set of PowerShell commands** to test your API in one go. Copy and paste this into a **new PowerShell session** after starting your server (`go run main.go`):

```powershell
# Step 1: Login
$login = Invoke-RestMethod -Uri "http://localhost:8080/login" -Method Post -Body (@{
    username = "admin"
    password = "password123"
} | ConvertTo-Json) -ContentType "application/json"
$accessToken = $login.data.access_token
$refreshToken = $login.data.refresh_token

# Step 2: Access v1 endpoint
Invoke-RestMethod -Uri "http://localhost:8080/api/v1/protected" -Method Get -Headers @{
    Authorization = "Bearer $accessToken"
}

# Step 3: Access v2 endpoint
Invoke-RestMethod -Uri "http://localhost:8080/api/v2/protected" -Method Get -Headers @{
    Authorization = "Bearer $accessToken"
}

# Step 4: Refresh token
$refresh = Invoke-RestMethod -Uri "http://localhost:8080/refresh" -Method Post -Body (@{
    refresh_token = $refreshToken
} | ConvertTo-Json) -ContentType "application/json"
$accessToken = $refresh.data.access_token

# Step 5: Access v1 endpoint again with new token
Invoke-RestMethod -Uri "http://localhost:8080/api/v1/protected" -Method Get -Headers @{
    Authorization = "Bearer $accessToken"
}
```

---

### **What Happens When You Run This?**

1. **Login**:
   - Sends a POST request to `/login` with valid credentials.
   - Extracts `access_token` and `refresh_token`.

2. **Access `/api/v1/protected`**:
   - Sends a GET request to `/api/v1/protected` using the `access_token`.
   - Verifies access to the v1 endpoint.

3. **Access `/api/v2/protected`**:
   - Sends a GET request to `/api/v2/protected` using the same `access_token`.
   - Verifies access to the v2 endpoint.

4. **Refresh Token**:
   - Sends a POST request to `/refresh` with the `refresh_token`.
   - Updates the `access_token` with the new one.

5. **Access `/api/v1/protected` Again**:
   - Sends another GET request to `/api/v1/protected` using the refreshed `access_token`.
   - Verifies that the new token works.

---

### **Expected Output**
Each command will return JSON responses like:

#### After Login:
```json
{
  "success": true,
  "message": "Login successful",
  "data": {
    "access_token": "ey...",
    "refresh_token": "ey..."
  }
}
```

#### After Accessing `/api/v1/protected`:
```json
{
  "success": true,
  "message": "Access granted",
  "data": {
    "user": "admin"
  }
}
```

#### After Accessing `/api/v2/protected`:
```json
{
  "success": true,
  "message": "Access granted with new features",
  "data": {
    "user": "admin",
    "version": "v2"
  }
}
```

#### After Refreshing Token:
```json
{
  "success": true,
  "message": "Token refreshed",
  "data": {
    "access_token": "ey..."
  }
}
```

---

### **Key Notes**
1. **Server Must Be Running**:
   - Ensure the Go server is running on `localhost:8080` before executing these commands.

2. **Error Handling**:
   - If any command fails, check the error message for details (e.g., invalid credentials, expired token).

3. **Security Reminder**:
   - The hardcoded secret key (`my_secret_key`) is for development only. Use proper key management in production.

---

This single block of commands tests all major functionalities of your API in sequence!
*/