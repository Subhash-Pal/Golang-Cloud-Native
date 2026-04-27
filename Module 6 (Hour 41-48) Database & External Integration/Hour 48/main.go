package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	_ "github.com/lib/pq"
	"github.com/redis/go-redis/v9"
)

type Order struct {
	ID           int       `json:"id"`
	CustomerName string    `json:"customer_name"`
	Amount       float64   `json:"amount"`
	Status       string    `json:"status"`
	CreatedAt    time.Time `json:"created_at"`
}

type CreateOrderRequest struct {
	CustomerName string  `json:"customer_name"`
	Amount       float64 `json:"amount"`
	Status       string  `json:"status"`
}

type App struct {
	db  *sql.DB
	rdb *redis.Client
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

	rdb := redis.NewClient(&redis.Options{
		Addr:     envOrDefault("REDIS_ADDR", "127.0.0.1:6379"),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       envIntOrDefault("REDIS_DB", 0),
	})
	defer rdb.Close()

	if err := rdb.Ping(ctx).Err(); err != nil {
		log.Fatalf("ping redis: %v", err)
	}

	if err := ensureSchema(ctx, db); err != nil {
		log.Fatalf("ensure schema: %v", err)
	}

	app := &App{db: db, rdb: rdb}

	mux := http.NewServeMux()
	mux.HandleFunc("POST /orders", app.handleCreateOrder)
	mux.HandleFunc("GET /orders/", app.handleGetOrder)
	mux.HandleFunc("GET /health", app.handleHealth)

	addr := ":8080"
	log.Printf("Order service running on http://127.0.0.1%s", addr)
	log.Fatal(http.ListenAndServe(addr, loggingMiddleware(mux)))
}

func (a *App) handleCreateOrder(w http.ResponseWriter, r *http.Request) {
	var req CreateOrderRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid JSON payload"})
		return
	}

	req.CustomerName = strings.TrimSpace(req.CustomerName)
	req.Status = strings.TrimSpace(req.Status)
	if req.CustomerName == "" || req.Amount <= 0 || req.Status == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "customer_name, amount, and status are required"})
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	order, err := createOrder(ctx, a.db, req)
	if err != nil {
		log.Printf("create order: %v", err)
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to create order"})
		return
	}

	if err := cacheOrder(ctx, a.rdb, order); err != nil {
		log.Printf("cache order after create: %v", err)
	}

	writeJSON(w, http.StatusCreated, order)
}

func (a *App) handleGetOrder(w http.ResponseWriter, r *http.Request) {
	idText := strings.TrimPrefix(r.URL.Path, "/orders/")
	orderID, err := strconv.Atoi(idText)
	if err != nil || orderID < 1 {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid order ID"})
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	order, source, err := getOrder(ctx, a.db, a.rdb, orderID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			writeJSON(w, http.StatusNotFound, map[string]string{"error": "order not found"})
			return
		}

		log.Printf("get order: %v", err)
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to fetch order"})
		return
	}

	w.Header().Set("X-Data-Source", source)
	writeJSON(w, http.StatusOK, order)
}

func (a *App) handleHealth(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func ensureSchema(ctx context.Context, db *sql.DB) error {
	query := `
	CREATE TABLE IF NOT EXISTS orders (
		id SERIAL PRIMARY KEY,
		customer_name TEXT NOT NULL,
		amount NUMERIC(12,2) NOT NULL,
		status TEXT NOT NULL,
		created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
	);`

	_, err := db.ExecContext(ctx, query)
	return err
}

func createOrder(ctx context.Context, db *sql.DB, req CreateOrderRequest) (Order, error) {
	var order Order
	err := db.QueryRowContext(
		ctx,
		`INSERT INTO orders (customer_name, amount, status) VALUES ($1, $2, $3) RETURNING id, customer_name, amount, status, created_at`,
		req.CustomerName,
		req.Amount,
		req.Status,
	).Scan(&order.ID, &order.CustomerName, &order.Amount, &order.Status, &order.CreatedAt)
	return order, err
}

func getOrder(ctx context.Context, db *sql.DB, rdb *redis.Client, orderID int) (Order, string, error) {
	cacheKey := fmt.Sprintf("hour48:order:%d", orderID)

	cached, err := rdb.Get(ctx, cacheKey).Result()
	if err == nil {
		var order Order
		if err := json.Unmarshal([]byte(cached), &order); err != nil {
			return Order{}, "", err
		}
		return order, "redis-cache", nil
	}
	if err != nil && err != redis.Nil {
		return Order{}, "", err
	}

	var order Order
	err = db.QueryRowContext(
		ctx,
		`SELECT id, customer_name, amount, status, created_at FROM orders WHERE id = $1`,
		orderID,
	).Scan(&order.ID, &order.CustomerName, &order.Amount, &order.Status, &order.CreatedAt)
	if err != nil {
		return Order{}, "", err
	}

	if err := cacheOrder(ctx, rdb, order); err != nil {
		return order, "postgres", nil
	}

	return order, "postgres", nil
}

func cacheOrder(ctx context.Context, rdb *redis.Client, order Order) error {
	payload, err := json.Marshal(order)
	if err != nil {
		return err
	}

	cacheKey := fmt.Sprintf("hour48:order:%d", order.ID)
	return rdb.Set(ctx, cacheKey, payload, 2*time.Minute).Err()
}

func writeJSON(w http.ResponseWriter, statusCode int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(payload); err != nil {
		log.Printf("write json response: %v", err)
	}
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		log.Printf("%s %s completed in %s", r.Method, r.URL.Path, time.Since(start))
	})
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
