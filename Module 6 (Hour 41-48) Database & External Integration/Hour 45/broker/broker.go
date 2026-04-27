package broker

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/lib/pq"
)

const Channel = "hour45_orders"

func ConnString() string {
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

func OpenDB() (*sql.DB, error) {
	return sql.Open("postgres", ConnString())
}

func EnsureSchema(ctx context.Context, db *sql.DB) error {
	query := `
	CREATE TABLE IF NOT EXISTS broker_messages (
		id SERIAL PRIMARY KEY,
		channel_name TEXT NOT NULL,
		payload TEXT NOT NULL,
		received_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
	);`

	_, err := db.ExecContext(ctx, query)
	return err
}

func PublishMessage(ctx context.Context, db *sql.DB, channel, payload string) error {
	query := fmt.Sprintf("SELECT pg_notify('%s', $1)", channel)
	_, err := db.ExecContext(ctx, query, payload)
	return err
}

func WaitForNotification(listener *pq.Listener, timeout time.Duration) (*pq.Notification, error) {
	select {
	case notification := <-listener.Notify:
		if notification == nil {
			return nil, fmt.Errorf("listener returned no notification")
		}
		return notification, nil
	case <-time.After(timeout):
		return nil, fmt.Errorf("timed out waiting for notification")
	}
}

func StoreMessage(ctx context.Context, db *sql.DB, channel, payload string) error {
	_, err := db.ExecContext(
		ctx,
		"INSERT INTO broker_messages (channel_name, payload) VALUES ($1, $2)",
		channel,
		payload,
	)
	return err
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
