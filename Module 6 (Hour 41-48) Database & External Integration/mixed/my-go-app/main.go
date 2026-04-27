package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq" // driver
)

func main() {
	// 1. Connect to the default 'postgres' database first
	adminConn := "host=localhost port=5432 user=postgres password=root dbname=postgres sslmode=disable"
	dbAdmin, err := sql.Open("postgres", adminConn)
	if err != nil {
		log.Fatal(err)
	}
	defer dbAdmin.Close()

	// 2. Check if 'mydb' exists
	var exists bool
	query := "SELECT EXISTS(SELECT 1 FROM pg_database WHERE datname = 'mydb')"
	err = dbAdmin.QueryRow(query).Scan(&exists)
	if err != nil {
		log.Fatal(err)
	}

	// 3. Create it if it doesn't exist
	if !exists {
		fmt.Println("Database 'mydb' not found. Creating it...")
		// Note: CREATE DATABASE cannot run inside a transaction
		_, err = dbAdmin.Exec("CREATE DATABASE mydb")
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("Database 'mydb' created successfully.")
	} else {
		fmt.Println("Database 'mydb' already exists.")
	}

	// 4. Now connect to the actual 'mydb' database
	connStr := "host=localhost port=5432 user=postgres password=root dbname=mydb sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Verify the final connection
	if err = db.Ping(); err != nil {
		log.Fatal("Could not connect to mydb:", err)
	}
	fmt.Println("Successfully connected to mydb!")
}

/*
Key steps in this code:
Initial Admin Connection: It connects to dbname=postgres (the default system DB) because you cannot check for a database's existence while trying to connect to it.
System Catalog Query: It queries pg_database, which is a system table containing metadata for all databases on the server.
Dynamic Creation: If the name is missing, it executes CREATE DATABASE mydb.
Clean Exit: Uses log.Fatal instead of panic for cleaner error messages during development. 
/*

 