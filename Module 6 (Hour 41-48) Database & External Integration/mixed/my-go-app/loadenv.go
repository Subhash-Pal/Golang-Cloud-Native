package main

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	// Load the .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Retrieve credentials using native "os" handles
	host := os.Getenv("DB_HOST")
	user := os.Getenv("DB_USER")
	pass := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")
	port := os.Getenv("DB_PORT")

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		host, user, pass, dbname, port)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	// 1. Get the generic database interface from GORM
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatal("Failed to get database handle:", err)
	}

	// 2. Check the health/stats
	stats := sqlDB.Stats()

	// 3. Print a meaningful message
	fmt.Printf("✅ Connection Established!\n")
	fmt.Printf("Connected to: %s/%s\n", host, dbname)
	fmt.Printf("Total Connections in Pool: %d\n", stats.OpenConnections) // Total connections created
	fmt.Printf("Currently Busy (In Use): %d\n", stats.InUse)             // Connections doin
}
