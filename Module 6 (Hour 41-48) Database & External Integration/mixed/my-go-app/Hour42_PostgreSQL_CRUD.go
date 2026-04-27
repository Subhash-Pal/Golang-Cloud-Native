package main

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"log"
	// PostgreSQL driver
)

// User struct represents a record in our database
type User struct {
	ID    int
	Name  string
	Email string
}

func main() {
	// 1. Initial Connection to System DB to ensure 'mydb' exists
	adminConn := "host=localhost port=5432 user=postgres password=root dbname=postgres sslmode=disable"
	dbAdmin, err := sql.Open("postgres", adminConn)
	if err != nil {
		log.Fatal(err)
	}
	var exists bool
	dbAdmin.QueryRow("SELECT EXISTS(SELECT 1 FROM pg_database WHERE datname = 'mydb')").Scan(&exists)
	if !exists {
		dbAdmin.Exec("CREATE DATABASE mydb")
		fmt.Println("Created database: mydb")
	}
	dbAdmin.Close()

	// 2. Connect to the actual Application Database
	connStr := "host=localhost port=5432 user=postgres password=root dbname=mydb sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// 3. Setup Table
	createTable := `CREATE TABLE IF NOT EXISTS users (
		id SERIAL PRIMARY KEY,
		name TEXT NOT NULL,
		email TEXT UNIQUE NOT NULL
	);`
	db.Exec(createTable)

	// --- CRUD OPERATIONS ---

	// CREATE (Insert)
	var newID int
	err = db.QueryRow("INSERT INTO users (name, email) VALUES ($1, $2) RETURNING id",
		"Alice", "alice@example.com").Scan(&newID)
	if err != nil {
		fmt.Println("Create failed (likely duplicate email):", err)
	} else {
		fmt.Printf("CREATE: Inserted user with ID %d\n", newID)
	}

	// READ (Select)
	fmt.Println("\nREAD: Current Users:")
	rows, _ := db.Query("SELECT id, name, email FROM users")
	defer rows.Close()
	for rows.Next() {
		var u User
		rows.Scan(&u.ID, &u.Name, &u.Email)
		fmt.Printf("- [%d] %s (%s)\n", u.ID, u.Name, u.Email)
	}

	// UPDATE (Change)
	updateSQL := "UPDATE users SET name = $1 WHERE email = $2"
	res, _ := db.Exec(updateSQL, "Alice Wonderland", "alice@example.com")
	count, _ := res.RowsAffected()
	fmt.Printf("\nUPDATE: Modified %d row(s)\n", count)

	// DELETE (Remove)
	// Change the email here to test deleting a specific record
	deleteSQL := "DELETE FROM users WHERE email = $1"
	res, _ = db.Exec(deleteSQL, "delete-me@example.com")
	count, _ = res.RowsAffected()
	fmt.Printf("DELETE: Removed %d row(s)\n", count)
}

/*
Key CRUD Concepts:

Create: Used QueryRow with RETURNING id to get the unique ID assigned by the database.
Read: Used Query and a rows.Next() loop to fetch multiple records.
Update: Used Exec and checked RowsAffected() to confirm the change happened.
Delete: Used Exec with a WHERE clause to target specific records.

*/
