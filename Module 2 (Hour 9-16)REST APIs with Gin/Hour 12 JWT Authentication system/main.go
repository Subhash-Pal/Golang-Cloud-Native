package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var secretKey = []byte("my_secret_key")

// Generate JWT
func generateToken(username string) (string, error) {
	claims := jwt.MapClaims{
		"username": username,
		"exp":      time.Now().Add(time.Minute * 5).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secretKey)
}

// Middleware to validate JWT
func verifyToken(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		tokenString := r.Header.Get("Authorization")

		if tokenString == "" {
			http.Error(w, "Missing token", http.StatusUnauthorized)
			return
		}

		tokenString = tokenString[len("Bearer "):]

		token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
			return secretKey, nil
		})

		if err != nil || !token.Valid {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		next(w, r)
	}
}

// Login endpoint
func loginHandler(w http.ResponseWriter, r *http.Request) {
	token, err := generateToken("admin")

	if err != nil {
		http.Error(w, "Error generating token", http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "TOKEN:\n%s", token)
}

// Protected endpoint
func protectedHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Protected data accessed")
}

func main() {
	http.HandleFunc("/login", loginHandler)
	http.HandleFunc("/protected", verifyToken(protectedHandler))

	fmt.Println("Server running on :8080")
	http.ListenAndServe(":8080", nil)
}

/*
Below is a **clean, lab-manual–style documentation** for your working JWT setup in Go, specifically tuned for **Windows + VS Code + PowerShell**. It follows execution order and validation checkpoints so failures don’t cascade.

---

# 🧪 LAB MANUAL

## Golang JWT Authentication (Windows + VS Code + PowerShell)

---

# 🎯 Objective

Implement a minimal JWT-based authentication system in Go with:

* Token generation (`/login`)
* Protected route (`/protected`)
* Middleware-based validation

---

# 🧱 System Requirements

| Component | Version       |
| --------- | ------------- |
| Go        | ≥ 1.20        |
| VS Code   | Latest        |
| OS        | Windows 10/11 |
| Terminal  | PowerShell    |

---

# 📁 Project Structure

```id="7r2l3n"
go-jwt-lab/
│── go.mod
│── go.sum
│── main.go
```

---

# ⚙️ STEP 1 — Project Initialization

### Command

```powershell id="u2n8av"
mkdir go-jwt-lab
cd go-jwt-lab
go mod init go-jwt-lab
```

### ✅ Validation

```powershell id="zkl7zx"
go mod tidy
```

✔ No errors → proceed

---

# ⚙️ STEP 2 — Install Dependency

### Command

```powershell id="x3h5lw"
go get github.com/golang-jwt/jwt/v5
```

### ✅ Validation

```powershell id="i8m3pc"
go list -m all
```

✔ Must include:

```id="1yck6k"
github.com/golang-jwt/jwt/v5
```

---

# ⚙️ STEP 3 — Create Application Code

### File: `main.go`

```go id="pk71t0"
package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var secretKey = []byte("my_secret_key")

// Generate JWT
func generateToken(username string) (string, error) {
	claims := jwt.MapClaims{
		"username": username,
		"exp":      time.Now().Add(time.Minute * 5).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secretKey)
}

// Middleware to verify token
func verifyToken(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		tokenString := r.Header.Get("Authorization")

		if tokenString == "" {
			http.Error(w, "Missing token", http.StatusUnauthorized)
			return
		}

		// Validate Bearer format
		if len(tokenString) < 7 || tokenString[:7] != "Bearer " {
			http.Error(w, "Invalid token format", http.StatusUnauthorized)
			return
		}

		tokenString = tokenString[7:]

		token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
			return secretKey, nil
		})

		if err != nil || !token.Valid {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		next(w, r)
	}
}

// Login endpoint
func loginHandler(w http.ResponseWriter, r *http.Request) {
	token, err := generateToken("admin")

	if err != nil {
		http.Error(w, "Error generating token", http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "TOKEN:\n%s", token)
}

// Protected endpoint
func protectedHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Protected data accessed")
}

func main() {
	http.HandleFunc("/login", loginHandler)
	http.HandleFunc("/protected", verifyToken(protectedHandler))

	fmt.Println("Server running on :8080")
	http.ListenAndServe(":8080", nil)
}
```

---

# ⚙️ STEP 4 — Run the Server

### Command

```powershell id="p9o1xk"
go run main.go
```

### ✅ Expected Output

```id="vru0tm"
Server running on :8080
```

---

# 🧪 STEP 5 — Testing via PowerShell

---

## 🔹 5.1 Generate Token

```powershell id="d5lq7x"
curl.exe http://localhost:8080/login
```

### Output

```id="5psv4k"
TOKEN:
<JWT_TOKEN>
```

✔ Copy token

---

## 🔹 5.2 Access Protected Route

```powershell id="kz4e8v"
$token="<PASTE_TOKEN_HERE>"

curl.exe -H "Authorization: Bearer $token" http://localhost:8080/protected
```

### ✅ Expected Output

```id="s2uj0v"
Protected data accessed
```

---

## 🔹 5.3 Negative Testing

### No Token

```powershell id="2nyb4m"
curl.exe http://localhost:8080/protected
```

✔ Output:

```id="38mj0s"
Missing token
```

---

### Invalid Token

```powershell id="5gdl6x"
curl.exe -H "Authorization: Bearer invalidtoken" http://localhost:8080/protected
```

✔ Output:

```id="5cmn6c"
Invalid token
```

---

# ⚠️ Windows-Specific Notes

### 1. curl vs curl.exe

| Command    | Behavior                                  |
| ---------- | ----------------------------------------- |
| `curl`     | Alias to PowerShell (`Invoke-WebRequest`) |
| `curl.exe` | Actual curl binary                        |

✔ Always use:

```powershell id="hjv5u6"
curl.exe -H "Authorization: Bearer TOKEN"
```

---

### 2. Alternative (PowerShell Native)

```powershell id="p87w5n"
Invoke-WebRequest -Uri "http://localhost:8080/protected" -Headers @{Authorization="Bearer $token"}
```

---

# 🧠 JWT Flow (Conceptual)

1. Client calls `/login`
2. Server issues JWT
3. Client sends JWT in header
4. Middleware validates:

   * Signature
   * Expiry
5. Access granted/denied

---

# 🔐 Security Notes (Production)

* Replace hardcoded key with environment variable
* Use HTTPS only
* Add:

  * Refresh tokens
  * Token revocation (Redis)
  * Role-based claims
* Validate signing algorithm explicitly

---

# 🚀 Suggested Extensions

You can evolve this lab into:

### 1. Framework Upgrade

* Use Gin for production APIs

### 2. Auth Enhancements

* Login with DB (PostgreSQL)
* Password hashing (bcrypt)

### 3. Token Strategy

* Access + Refresh token flow
* Token rotation

### 4. Architecture

* Middleware chaining
* Clean architecture (handler → service → repo)

---

# ✅ Final Validation Checklist

| Check                 | Status |
| --------------------- | ------ |
| Server starts         | ✔      |
| Token generated       | ✔      |
| Protected route works | ✔      |
| Invalid token blocked | ✔      |

---

# 🎉 Congratulations	!
*/
