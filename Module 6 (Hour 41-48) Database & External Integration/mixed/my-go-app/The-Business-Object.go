package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"

	_ "github.com/lib/pq" // drive
)

// User represents our internal Business Object
// It maps to both a Database Schema and an External API Schema
type User struct {
	ID    int    `db:"id"    json:"user_id"`   // Database 'id' -> JSON 'user_id'
	Name  string `db:"name"  json:"full_name"` // Database 'name' -> JSON 'full_name'
	Email string `db:"email" json:"contact"`   // Database 'email' -> JSON 'contact'
}

func main() {
	// Current local configuration
	connStr := "host=localhost port=5432 user=postgres password=root dbname=mydb sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// --- 1. MAPPING FROM LOCAL DATABASE ---
	dbUser := User{}
	err = db.QueryRow("SELECT id, name, email FROM users LIMIT 1").
		Scan(&dbUser.ID, &dbUser.Name, &dbUser.Email)

	if err == nil {
		fmt.Printf("Mapped from DB: %+v\n", dbUser)
	}

	// --- 2. MAPPING FROM EXTERNAL API (Example JSON) ---
	externalData := `{"user_id": 99, "full_name": "API User", "contact": "api@external.com"}`

	var apiUser User
	err = json.Unmarshal([]byte(externalData), &apiUser)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Mapped from API: %+v\n", apiUser)
}

/*
To align with your current PostgreSQL setup (mydb, users table), we use Struct Tags. These map database columns to your Go struct fields so the application can interact with the data as native objects.
Here is a simple example showing how to map a database table and a JSON API response to the same Business Object.
1. The Business Object (Entity)
We define a User struct. The db:"..." tags map to your Postgres columns, and json:"..." tags map to external API fields.

Why this is helpful:
Native Interaction: Your business logic only talks to the User struct. It doesn't care if the data came from a SELECT statement or a GET request.
Schema Decoupling: If the external API changes a field name from contact to email_address, you only change the struct tag in one place—your internal code logic remains untouched.
Consistency: You ensure that "Email" is always treated as a string across your entire system, regardless of where it lives externally.
2. Mapping from Database
When you query the database, you scan the results directly into your User struct. This way, you can work with dbUser as a native Go object in your application.
3. Mapping from External API
When you receive JSON data from an API, you unmarshal it into the same User struct. This allows you to treat API data as native objects without writing separate parsing logic for each source.
This pattern of using a single Business Object with struct tags for mapping is a common and powerful way to manage data across different layers of an application while keeping your code clean and maintainable.
*/
