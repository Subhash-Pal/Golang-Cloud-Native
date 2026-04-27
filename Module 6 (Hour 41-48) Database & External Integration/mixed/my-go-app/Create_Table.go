package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq" // driver
)

func main() {
	// 1. Initial admin connection to ensure 'mydb' exists
	adminConn := "host=localhost port=5432 user=postgres password=root dbname=postgres sslmode=disable"
	dbAdmin, err := sql.Open("postgres", adminConn)
	if err != nil {
		log.Fatal(err)
	}

	var exists bool
	dbAdmin.QueryRow("SELECT EXISTS(SELECT 1 FROM pg_database WHERE datname = 'mydb')").Scan(&exists)
	if !exists {
		_, err = dbAdmin.Exec("CREATE DATABASE mydb")
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("Database 'mydb' created.")
	}
	dbAdmin.Close()

	// 2. Connect to the application database 'mydb'
	connStr := "host=localhost port=5432 user=postgres password=root dbname=mydb sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// 3. Create a table using Exec()
	// 'IF NOT EXISTS' prevents errors on subsequent runs
	createTableSQL := `
	CREATE TABLE IF NOT EXISTS users (
		id SERIAL PRIMARY KEY,
		name TEXT NOT NULL,
		email TEXT UNIQUE NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);`

	_, err = db.Exec(createTableSQL)
	if err != nil {
		log.Fatal("Error creating table:", err)
	}
	fmt.Println("Table 'users' is ready!")

	// 4. Insert data with placeholders ($1, $2) to prevent SQL injection
	insertSQL := `INSERT INTO users (name, email) VALUES ($1, $2) ON CONFLICT (email) DO NOTHING;`
	_, err = db.Exec(insertSQL, "John Doe", "john@example.com")
	if err != nil {
		log.Fatal("Error inserting data:", err)
	}
	fmt.Println("Sample data inserted successfully!")
}
