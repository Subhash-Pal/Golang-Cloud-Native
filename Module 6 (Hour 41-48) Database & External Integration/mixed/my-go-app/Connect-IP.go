package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq" //go get 
)

func main() {
	// Using 127.0.0.1 (Standard Local IP) instead of "localhost"
	// Ensure your password matches what you set during Postgres installation
	connStr := "host=127.0.0.1 port=5432 user=postgres password=root dbname=postgres sslmode=disable"

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Connection setup failed:", err)
	}
	defer db.Close()

	// Test the connection
	err = db.Ping()
	if err != nil {
		log.Fatal("Failed to connect to local IP:", err)
	}

	fmt.Println("Successfully connected to PostgreSQL at 127.0.0.1!")
}

/*
Why use 127.0.0.1?
Speed: It bypasses the DNS lookup required for the word "localhost".
Consistency: It behaves exactly like a remote IP address, which helps you test if your code is ready for a real server.
Troubleshooting Local IP Connection
If you get a "connection refused" error while using the IP:
Check Port: Ensure PostgreSQL is actually running on 5432 (this is the default).
Verify Password: If you see password authentication failed, double-check that your password is root.
Firewall: On Windows, sometimes the built-in firewall blocks local IP traffic if the network profile is set to "Public".
*/
