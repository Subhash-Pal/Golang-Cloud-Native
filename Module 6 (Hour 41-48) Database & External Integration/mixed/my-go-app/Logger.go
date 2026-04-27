package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger" // Import the logger module
)

func main() {
	godotenv.Load()

	// 1. Configure the Logger
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Second, // Log SQL if it takes longer than 1s
			LogLevel:                  logger.Info, // Options: Silent, Error, Warn, Info
			IgnoreRecordNotFoundError: true,        // Don't log "record not found" as an error
			Colorful:                  true,        // Enable colors in terminal
		},
	)

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		os.Getenv("DB_HOST"), os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_NAME"), os.Getenv("DB_PORT"))

	// 2. Pass the logger into GORM config
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: newLogger,
	})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("🔎 Logger is active. Watch the SQL output below:")

	// 3. Run a query to see the logger in action
	var result string
	db.Raw("SELECT 'Hello GORM Logger!'").Scan(&result)

	fmt.Printf("Result: %s\n", result)
}

/*
What you will see in your Terminal:
When you run this, GORM will print a line like this for every database action:
[0.452ms] [rows:1] SELECT 'Hello GORM Logger!'
Why use different LogLevels?
logger.Info: Shows everything (best for development).
logger.Warn: Only shows slow queries and warnings.
logger.Error: Only shows queries that actually failed (best for production to keep logs small).
SlowThreshold: This is a "performance handle"—it helps you find which queries are making your app slow.
*/
