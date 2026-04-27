package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	_ "github.com/lib/pq"
)

type Account struct {
	ID      int
	Name    string
	Balance float64
}

func main() {
	db, err := sql.Open("postgres", postgresConnString())
	if err != nil {
		log.Fatalf("open database: %v", err)
	}
	defer db.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		log.Fatalf("ping database: %v", err)
	}

	if err := ensureSchema(ctx, db); err != nil {
		log.Fatalf("ensure schema: %v", err)
	}

	if err := seedAccounts(ctx, db); err != nil {
		log.Fatalf("seed accounts: %v", err)
	}

	fmt.Println("Balances before transfer")
	printAccounts(ctx, db)

	err = transferFunds(ctx, db, 1, 2, 250, true)
	if err != nil {
		fmt.Printf("\nROLLBACK: %v\n", err)
	}

	fmt.Println("\nBalances after failed transfer")
	printAccounts(ctx, db)

	if err := transferFunds(ctx, db, 1, 2, 250, false); err != nil {
		log.Fatalf("second transfer failed unexpectedly: %v", err)
	}

	fmt.Println("\nCOMMIT: transfer completed successfully")
	fmt.Println("\nBalances after successful transfer")
	printAccounts(ctx, db)
}

func ensureSchema(ctx context.Context, db *sql.DB) error {
	query := `
	CREATE TABLE IF NOT EXISTS hour43_accounts (
		id INT PRIMARY KEY,
		name TEXT NOT NULL,
		balance NUMERIC(12,2) NOT NULL
	);`

	_, err := db.ExecContext(ctx, query)
	return err
}

func seedAccounts(ctx context.Context, db *sql.DB) error {
	_, err := db.ExecContext(ctx, `
		INSERT INTO hour43_accounts (id, name, balance)
		VALUES
			(1, 'Alice', 1000.00),
			(2, 'Bob', 500.00)
		ON CONFLICT (id) DO UPDATE
		SET name = EXCLUDED.name,
			balance = EXCLUDED.balance;
	`)
	return err
}

func transferFunds(ctx context.Context, db *sql.DB, fromID, toID int, amount float64, simulateFailure bool) error {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	if _, err = tx.ExecContext(ctx, `UPDATE hour43_accounts SET balance = balance - $1 WHERE id = $2`, amount, fromID); err != nil {
		return fmt.Errorf("debit source account: %w", err)
	}

	if simulateFailure {
		err = fmt.Errorf("simulated failure after debit, transaction must rollback")
		return err
	}

	if _, err = tx.ExecContext(ctx, `UPDATE hour43_accounts SET balance = balance + $1 WHERE id = $2`, amount, toID); err != nil {
		return fmt.Errorf("credit destination account: %w", err)
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("commit transaction: %w", err)
	}

	return nil
}

func printAccounts(ctx context.Context, db *sql.DB) {
	rows, err := db.QueryContext(ctx, `SELECT id, name, balance FROM hour43_accounts ORDER BY id`)
	if err != nil {
		log.Printf("query accounts: %v", err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var account Account
		if err := rows.Scan(&account.ID, &account.Name, &account.Balance); err != nil {
			log.Printf("scan account: %v", err)
			return
		}
		fmt.Printf("- [%d] %s balance = %.2f\n", account.ID, account.Name, account.Balance)
	}
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
