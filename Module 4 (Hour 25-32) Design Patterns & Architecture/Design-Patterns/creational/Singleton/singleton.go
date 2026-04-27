// The Singleton Pattern ensures that a class (or struct in Golang) has only one instance

// Since Golang doesn’t have classes, we implement Singletons using structs and control instance creation using sync.Once or mutexes.

// Thread-Safe Database Singleton in Golang
// Using sync.Once and encapsulating the instance within a struct.

package main

import (
	"database/sql"
	"fmt"
	"log"
	"sync"

	_ "github.com/lib/pq" // PostgreSQL driver
)

// Database struct encapsulates the singleton instance
type Database struct {
	connection *sql.DB
}

var once sync.Once

// NewDatabase initializes the database connection
func NewDatabase() *Database {
	var dbInstance *Database

	once.Do(func() {
		fmt.Println("Initializing Database Connection...")
		conn, err := sql.Open("postgres", "user=postgres password=secret dbname=mydb sslmode=disable")
		if err != nil {
			log.Fatalf("Failed to connect to database: %v", err)
		}
		dbInstance = &Database{connection: conn}
	})
	return dbInstance
}

// GetDBConnection returns the database connection
func (d *Database) GetDBConnection() *sql.DB {
	return d.connection
}

func main() {
	db1 := NewDatabase()
	db2 := NewDatabase()

	fmt.Println(db1.GetDBConnection() == db2.GetDBConnection()) // Output: true (Same instance)
}
