package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

var ctx = context.Background()

func main() {
	// Create Redis client
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379", // Docker mapped port
		Password: "",               // no password
		DB:       0,                // default DB
	})

	// Test connection
	pong, err := rdb.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("Redis connection failed: %v", err)
	}

	fmt.Println("Connected to Redis:", pong)

	// Simple SET
	err = rdb.Set(ctx, "name", "Shubh", 10*time.Minute).Err()
	if err != nil {
		log.Fatalf("SET failed: %v", err)
	}

	// Simple GET
	val, err := rdb.Get(ctx, "name").Result()
	if err != nil {
		log.Fatalf("GET failed: %v", err)
	}

	fmt.Println("Value from Redis:", val)
}

/*
# Redis + Golang Quick Start

1. Run Redis: `docker run -d -p 6379:6379 --name redis-local redis:7-alpine`
2. Verify container: `docker ps`
3. Initialize Go module: `go mod init redis-demo`
4. Install dependency: `go get github.com/redis/go-redis/v9`
5. Run app: `go run Connect-redis.go`

*/
