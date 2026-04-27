package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	// 1. Load Secrets
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		os.Getenv("DB_HOST"), os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_NAME"), os.Getenv("DB_PORT"))

	// 2. Open Connection with GORM
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}

	// 3. Performance Tuning (Connection Pool)
	sqlDB, _ := db.DB()
	sqlDB.SetMaxIdleConns(10)           // Keeps 10 connections "warmed up"
	sqlDB.SetMaxOpenConns(100)          // Max total connections allowed
	sqlDB.SetConnMaxLifetime(time.Hour) // Close connections after 1 hour to refresh memory

	// 4. Meaningful Health Check
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	if err := sqlDB.PingContext(ctx); err != nil {
		log.Fatalf("Database is unreachable: %v", err)
	}

	// 5. Display Detailed Status
	stats := sqlDB.Stats()
	fmt.Println("🚀 Database Engine Ready")
	fmt.Printf("--- Status ---\n")
	fmt.Printf("Host: %s | DB: %s\n", os.Getenv("DB_HOST"), os.Getenv("DB_NAME"))
	fmt.Printf("Connection Pool: %d Open | %d Idle\n", stats.OpenConnections, stats.Idle)

	// Bonus: Get Postgres version for verification
	var version string
	sqlDB.QueryRow("SELECT version()").Scan(&version)
	fmt.Printf("Postgres Info: %s\n", version[:30]+"...")
}
