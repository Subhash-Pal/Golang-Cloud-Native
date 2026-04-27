package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	_ "github.com/lib/pq"
	"github.com/redis/go-redis/v9"
)

var ctx = context.Background()

// ---------------- ENV HELPERS ----------------
func envOrDefault(key, def string) string {
	val := os.Getenv(key)
	if val == "" {
		return def
	}
	return val
}

func envIntOrDefault(key string, def int) int {
	val := os.Getenv(key)
	if val == "" {
		return def
	}
	i, err := strconv.Atoi(val)
	if err != nil {
		return def
	}
	return i
}

// ---------------- POSTGRES CONFIG ----------------
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

// ---------------- MODEL ----------------
type User struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

// ---------------- REDIS ----------------
func initRedis() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
}

// ---------------- POSTGRES ----------------
func initPostgres() *sql.DB {
	db, err := sql.Open("postgres", postgresConnString())
	if err != nil {
		log.Fatal("Postgres open error:", err)
	}

	if err := db.Ping(); err != nil {
		log.Fatal("Postgres ping failed:", err)
	}

	fmt.Println("✅ Connected to PostgreSQL")
	return db
}

// ---------------- AUTO SETUP ----------------
func setupDB(db *sql.DB) {
	query := `
	CREATE TABLE IF NOT EXISTS users (
		id SERIAL PRIMARY KEY,
		name TEXT,
		email TEXT
	);`
	_, err := db.Exec(query)
	if err != nil {
		log.Fatal("Table creation failed:", err)
	}

	// Seed data if empty
	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM users").Scan(&count)
	if err != nil {
		log.Fatal(err)
	}

	if count == 0 {
		_, err = db.Exec(`
			INSERT INTO users (name, email)
			VALUES 
			('Shubh', 'shubh@example.com'),
			('Alice', 'alice@example.com')
		`)
		if err != nil {
			log.Fatal("Seed insert failed:", err)
		}
		fmt.Println("✅ Seed data inserted")
	}
}

// ---------------- CACHE-ASIDE ----------------
func getUser(rdb *redis.Client, db *sql.DB, userID int) (*User, error) {
	key := fmt.Sprintf("user:%d", userID)

	// 1. Redis lookup
	val, err := rdb.Get(ctx, key).Result()
	if err == nil {
		var user User
		json.Unmarshal([]byte(val), &user)
		fmt.Println("⚡ Cache HIT")
		return &user, nil
	}

	fmt.Println("🐢 Cache MISS → DB")

	// 2. DB lookup
	row := db.QueryRow("SELECT id, name, email FROM users WHERE id=$1", userID)

	var user User
	if err := row.Scan(&user.ID, &user.Name, &user.Email); err != nil {
		return nil, err
	}

	// 3. Store in Redis
	data, _ := json.Marshal(user)
	rdb.Set(ctx, key, data, 2*time.Minute)

	return &user, nil
}

// ---------------- MAIN ----------------
func main() {
	rdb := initRedis()
	db := initPostgres()
	defer db.Close()

	setupDB(db)

	// Call twice to demonstrate cache behavior
	user, err := getUser(rdb, db, 1)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("User:", user)

	user, _ = getUser(rdb, db, 1)
	fmt.Println("User (cached):", user)
}