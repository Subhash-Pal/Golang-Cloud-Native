package main

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/lib/pq" // driver
)

// User struct to hold database records
type User struct {
	ID        int
	Name      string
	Email     string
	CreatedAt time.Time
}

func main() {
	connStr := "host=localhost port=5432 user=postgres password=root dbname=mydb sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// 1. Execute the SELECT query
	rows, err := db.Query("SELECT id, name, email, created_at FROM users")
	if err != nil {
		log.Fatal("Query error:", err)
	}
	defer rows.Close() // Crucial: Always close rows to free connections

	fmt.Println("--- User List ---")

	// 2. Iterate through the result set
	for rows.Next() {
		var u User
		// Scan columns into the struct fields (must match query order)
		err := rows.Scan(&u.ID, &u.Name, &u.Email, &u.CreatedAt)
		if err != nil {
			log.Fatal("Scan error:", err)
		}

		// 3. Display the data
		fmt.Printf("ID: %d | Name: %-10s | Email: %-20s | Joined: %s\n",
			u.ID, u.Name, u.Email, u.CreatedAt.Format("2006-01-02"))
	}

	// 4. Check for errors that occurred during iteration
	if err = rows.Err(); err != nil {
		log.Fatal(err)
	}
}

/*
Important Patterns Used:
rows.Close(): Result sets hold a database connection open until closed. Using defer rows.Close() ensures the connection returns to the pool even if a scan fails.
rows.Next(): This loop continues as long as there is another row to process.
rows.Scan(): This copies data from the current row into the provided pointers. You must provide one pointer for every column in your SELECT statement.
rows.Err(): Always check this after the loop. If the connection drops midway through a large result set, rows.Next() will simply return false, and only rows.Err() will tell you something went wrong.
/*
