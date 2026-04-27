package broker

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/lib/pq"
)

const (
	ConnStr = "host=127.0.0.1 port=5432 user=postgres password=root dbname=postgres sslmode=disable"
	Channel = "hour45_orders"
)

func OpenDB() (*sql.DB, error) {
	return sql.Open("postgres", ConnStr)
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
