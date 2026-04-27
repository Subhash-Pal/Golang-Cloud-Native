package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"
	"sync"
	"time"

	_ "github.com/lib/pq"
)

func main() {
	db, err := sql.Open("postgres", postgresConnString())
	if err != nil {
		log.Fatalf("open database: %v", err)
	}
	defer db.Close()

	db.SetMaxOpenConns(2)
	db.SetMaxIdleConns(2)
	db.SetConnMaxLifetime(2 * time.Minute)
	db.SetConnMaxIdleTime(30 * time.Second)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		log.Fatalf("ping database: %v", err)
	}

	log.Println("Connected to PostgreSQL")
	log.Println("Running 4 concurrent queries with a pool size of 2 to show connection reuse")

	var wg sync.WaitGroup
	for workerID := 1; workerID <= 4; workerID++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			runWorker(ctx, db, id)
		}(workerID)
	}

	wg.Wait()

	stats := db.Stats()
	fmt.Println()
	fmt.Println("database/sql pool stats")
	fmt.Printf("OpenConnections: %d\n", stats.OpenConnections)
	fmt.Printf("InUse: %d\n", stats.InUse)
	fmt.Printf("Idle: %d\n", stats.Idle)
	fmt.Printf("WaitCount: %d\n", stats.WaitCount)
	fmt.Printf("WaitDuration: %s\n", stats.WaitDuration)
	fmt.Printf("MaxIdleClosed: %d\n", stats.MaxIdleClosed)
	fmt.Printf("MaxLifetimeClosed: %d\n", stats.MaxLifetimeClosed)
}

func runWorker(ctx context.Context, db *sql.DB, workerID int) {
	queryCtx, cancel := context.WithTimeout(ctx, 4*time.Second)
	defer cancel()

	var backendPID int
	var databaseName string
	var startedAt time.Time

	query := `
		SELECT pg_backend_pid(), current_database(), NOW()
		FROM pg_sleep(1);
	`

	if err := db.QueryRowContext(queryCtx, query).Scan(&backendPID, &databaseName, &startedAt); err != nil {
		log.Printf("worker %d failed: %v", workerID, err)
		return
	}

	log.Printf("worker %d used backend PID %d on database %s at %s", workerID, backendPID, databaseName, startedAt.Format(time.RFC3339))
}

func postgresConnString() string {
	host := envOrDefault("DB_HOST", "127.0.0.1")
	port := envIntOrDefault("DB_PORT", 5432)
	user := envOrDefault("DB_USER", "postgres")
	password := envOrDefault("DB_PASSWORD", "root")
	dbName := envOrDefault("DB_NAME", "postgres")

	return fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host,
		port,
		user,
		password,
		dbName,
	)
}

func envOrDefault(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

func envIntOrDefault(key string, fallback int) int {
	if value := os.Getenv(key); value != "" {
		if parsed, err := strconv.Atoi(value); err == nil {
			return parsed
		}
	}
	return fallback
}
