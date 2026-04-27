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

type User struct {
	ID        int
	Name      string
	Email     string
	CreatedAt time.Time
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

	userID, err := createUser(ctx, db, "Alice", fmt.Sprintf("alice+%d@example.com", time.Now().Unix()))
	if err != nil {
		log.Fatalf("create user: %v", err)
	}
	fmt.Printf("CREATE: inserted user ID %d\n", userID)

	users, err := listUsers(ctx, db)
	if err != nil {
		log.Fatalf("list users: %v", err)
	}

	fmt.Println("\nREAD: users in table")
	for _, user := range users {
		fmt.Printf("- [%d] %s <%s> at %s\n", user.ID, user.Name, user.Email, user.CreatedAt.Format(time.RFC3339))
	}

	if err := updateUserName(ctx, db, userID, "Alice Wonderland"); err != nil {
		log.Fatalf("update user: %v", err)
	}
	fmt.Printf("\nUPDATE: user %d renamed successfully\n", userID)

	if err := deleteUser(ctx, db, userID); err != nil {
		log.Fatalf("delete user: %v", err)
	}
	fmt.Printf("DELETE: user %d removed successfully\n", userID)
}

func ensureSchema(ctx context.Context, db *sql.DB) error {
	query := `
	CREATE TABLE IF NOT EXISTS hour42_users (
		id SERIAL PRIMARY KEY,
		name TEXT NOT NULL,
		email TEXT UNIQUE NOT NULL,
		created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
	);`

	_, err := db.ExecContext(ctx, query)
	return err
}

func createUser(ctx context.Context, db *sql.DB, name, email string) (int, error) {
	var id int
	err := db.QueryRowContext(
		ctx,
		`INSERT INTO hour42_users (name, email) VALUES ($1, $2) RETURNING id`,
		name,
		email,
	).Scan(&id)
	return id, err
}

func listUsers(ctx context.Context, db *sql.DB) ([]User, error) {
	rows, err := db.QueryContext(ctx, `SELECT id, name, email, created_at FROM hour42_users ORDER BY id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var user User
		if err := rows.Scan(&user.ID, &user.Name, &user.Email, &user.CreatedAt); err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, rows.Err()
}

func updateUserName(ctx context.Context, db *sql.DB, userID int, newName string) error {
	result, err := db.ExecContext(ctx, `UPDATE hour42_users SET name = $1 WHERE id = $2`, newName, userID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return fmt.Errorf("no user found with id %d", userID)
	}

	return nil
}

func deleteUser(ctx context.Context, db *sql.DB, userID int) error {
	result, err := db.ExecContext(ctx, `DELETE FROM hour42_users WHERE id = $1`, userID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return fmt.Errorf("no user found with id %d", userID)
	}

	return nil
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
