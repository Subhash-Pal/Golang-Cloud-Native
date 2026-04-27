package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

type Product struct {
	ID    int     `json:"id"`
	Name  string  `json:"name"`
	Price float64 `json:"price"`
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	rdb := redis.NewClient(&redis.Options{
		Addr:     envOrDefault("REDIS_ADDR", "127.0.0.1:6379"),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       envIntOrDefault("REDIS_DB", 0),
	})
	defer rdb.Close()

	if err := rdb.Ping(ctx).Err(); err != nil {
		log.Fatalf("ping redis: %v", err)
	}

	repo := map[int]Product{
		102: {ID: 102, Name: "Mechanical Keyboard", Price: 3499.00},
	}

	productID := 102//key
	product, source, err := getProductWithCache(ctx, rdb, repo, productID)
	if err != nil {
		log.Fatalf("first fetch: %v", err)
	}
	fmt.Printf("First fetch came from %s: %+v\n", source, product)

	product, source, err = getProductWithCache(ctx, rdb, repo, productID)
	if err != nil {
		log.Fatalf("second fetch: %v", err)
	}
	fmt.Printf("Second fetch came from %s: %+v\n", source, product)
}

func getProductWithCache(ctx context.Context, rdb *redis.Client, repo map[int]Product, productID int) (Product, string, error) {
	cacheKey := fmt.Sprintf("hour44:product:%d", productID)

	cached, err := rdb.Get(ctx, cacheKey).Result()
	if err == nil {
		var product Product
		if err := json.Unmarshal([]byte(cached), &product); err != nil {
			return Product{}, "", err
		}
		return product, "redis-cache", nil
	}

	if err != redis.Nil {
		return Product{}, "", err
	}

	product, ok := repo[productID]
	if !ok {
		return Product{}, "", fmt.Errorf("product %d not found", productID)
	}

	payload, err := json.Marshal(product)
	if err != nil {
		return Product{}, "", err
	}

	if err := rdb.Set(ctx, cacheKey, payload, 30*time.Second).Err(); err != nil {
		return Product{}, "", err
	}

	return product, "repository", nil
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

