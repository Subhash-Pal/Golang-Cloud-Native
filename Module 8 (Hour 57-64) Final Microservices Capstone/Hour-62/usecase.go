package main

import (
	"context"
	"database/sql"
	//"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"github.com/redis/go-redis/v9"
)

var ctx = context.Background()

// ---------- ENV ----------
func envOrDefault(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

func envIntOrDefault(key string, def int) int {
	if v := os.Getenv(key); v != "" {
		if i, err := strconv.Atoi(v); err == nil {
			return i
		}
	}
	return def
}

// ---------- POSTGRES ----------
func postgresConnString() string {
	host := envOrDefault("DB_HOST", "127.0.0.1")
	port := envIntOrDefault("DB_PORT", 5432)
	user := envOrDefault("DB_USER", "postgres")
	password := envOrDefault("DB_PASSWORD", "root")
	dbName := envOrDefault("DB_NAME", "postgres")

	return fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbName,
	)
}

func initPostgres() *sql.DB {
	db, err := sql.Open("postgres", postgresConnString())
	if err != nil {
		log.Fatal(err)
	}
	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}
	fmt.Println("✅ PostgreSQL connected")
	return db
}

// ---------- REDIS ----------
func initRedis() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
}

// ---------- MODEL ----------
type User struct {
	ID       int
	Username string
	Password string
}

// ---------- SETUP ----------
func setupDB(db *sql.DB) {
	query := `
	CREATE TABLE IF NOT EXISTS users (
		id SERIAL PRIMARY KEY,
		username TEXT UNIQUE,
		password TEXT
	);`
	_, err := db.Exec(query)
	if err != nil {
		log.Fatal(err)
	}

	var count int
	db.QueryRow("SELECT COUNT(*) FROM users").Scan(&count)

	if count == 0 {
		db.Exec(`
			INSERT INTO users (username, password)
			VALUES ('shubh', '1234')
		`)
		fmt.Println("✅ Seed user created (shubh / 1234)")
	}
}

// ---------- RATE LIMIT ----------
func allowLogin(rdb *redis.Client, username string) bool {
	key := "login_attempt:" + username

	count, err := rdb.Incr(ctx, key).Result()
	if err != nil {
		log.Println("Redis error:", err)
		return false
	}

	if count == 1 {
		rdb.Expire(ctx, key, time.Minute)
	}

	if count > 5 {
		return false
	}
	return true
}

// ---------- AUTH ----------
func authenticate(db *sql.DB, username, password string) bool {
	var stored string
	err := db.QueryRow(
		"SELECT password FROM users WHERE username=$1",
		username,
	).Scan(&stored)

	if err != nil {
		return false
	}

	return stored == password
}

// ---------- HANDLER ----------
func loginHandler(db *sql.DB, rdb *redis.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req map[string]string
		if err := c.BindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
			return
		}

		username := req["username"]
		password := req["password"]

		// Rate limit check
		if !allowLogin(rdb, username) {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error": "too many attempts, try later",
			})
			return
		}

		// Authenticate
		if !authenticate(db, username, password) {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "invalid credentials",
			})
			return
		}

		// Success
		resp := map[string]interface{}{
			"message": "login successful",
			"user":    username,
			"time":    time.Now(),
		}

		c.JSON(http.StatusOK, resp)
	}
}

// ---------- MAIN ----------
func main() {
	db := initPostgres()
	defer db.Close()

	rdb := initRedis()

	setupDB(db)

	router := gin.Default()
	router.POST("/login", loginHandler(db, rdb))

	fmt.Println("🚀 Server running on http://localhost:8080")
	router.Run(":8080")
}

/*
Good—now that everything is working, here’s a **clean, reproducible run sequence** for your full setup (Redis + PostgreSQL + Gin API + test script).

---

# 🚀 Step-by-Step Run Guide

## ✅ Step 1 — Start Redis (Docker)

```powershell
docker run -d -p 6379:6379 --name redis-local redis:7-alpine
```

Verify:

```powershell
docker ps
```

---

## ✅ Step 2 — Start PostgreSQL

### Option A (Local install)

Ensure PostgreSQL service is running.

### Option B (Docker – recommended for consistency)

```powershell
docker run -d -p 5432:5432 `
-e POSTGRES_PASSWORD=root `
--name postgres-local postgres
```

---

## ✅ Step 3 — Verify PostgreSQL

Login:

```powershell
psql -U postgres -h localhost -p 5432 -d postgres
```

Fix schema (IMPORTANT):

```sql
DROP TABLE IF EXISTS users;

CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    username TEXT UNIQUE,
    password TEXT
);

INSERT INTO users (username, password)
VALUES ('shubh', '1234');
```

Exit:

```sql
\q
```

---

## ✅ Step 4 — Run Go API

Inside your project folder:

```powershell
go mod init login-api   # only first time
go get github.com/gin-gonic/gin
go get github.com/lib/pq
go get github.com/redis/go-redis/v9

go run main.go
```

Expected:

```text
🚀 Server running on http://localhost:8080
```

---

## ✅ Step 5 — Test API (manual)

```powershell
Invoke-RestMethod -Uri "http://localhost:8080/login" `
-Method POST `
-Body (@{username="shubh"; password="1234"} | ConvertTo-Json) `
-ContentType "application/json"
```

---

## ✅ Step 6 — Run Success Script

```powershell
Set-ExecutionPolicy -Scope Process Bypass
.\test-success.ps1
```

---

## ✅ Step 7 — (Optional) Run Failure Script

```powershell
.\test-rate-limit.ps1
```

---

## ⚠️ Important Reset Step (when switching tests)

If failure test was run:

```bash
redis-cli DEL login_attempt:shubh
```

---

# 🧠 Final Architecture You Are Running

```text
PowerShell Script / Client
        ↓
Gin API (Go)
        ↓
Redis (Rate Limiter)
        ↓
PostgreSQL (Auth DB)
```

---

# 🔥 Common Failure Points (quick debug)

| Issue                 | Fix                        |
| --------------------- | -------------------------- |
| `invalid credentials` | Check DB schema/data       |
| `too many attempts`   | Clear Redis key            |
| connection refused    | Check Docker containers    |
| script blocked        | Use ExecutionPolicy Bypass |

---

*/