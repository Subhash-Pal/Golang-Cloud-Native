package main

import (
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	// Removed lib/pq; GORM uses pgx driver by default which is faster and cleaner
)

type User struct {
	ID   uint   `gorm:"primaryKey"`
	Name string `gorm:"column:name"`
	// Use 'unique' instead of 'uniqueIndex' for standard constraint mapping
	Email string `gorm:"unique"`
}

func main() {
	dsn := "host=localhost user=postgres password=root dbname=mydb port=5432 sslmode=disable"

	// Open connection
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		// This suppresses the "DROP CONSTRAINT" error logs during migration
		SkipDefaultTransaction: true,
	})
	if err != nil {
		log.Fatal(err)
	}

	// 1. Migrate the schema
	db.AutoMigrate(&User{})

	// 2. SMART CREATE (Upsert)
	// If the email exists, it updates the name. If not, it creates a new record.
	user := User{Name: "Gopher", Email: "gopher@golang.org"}
	db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "email"}},
		DoUpdates: clause.AssignmentColumns([]string{"name"}),
	}).Create(&user)

	fmt.Printf("User Status: ID %d, Name %s\n", user.ID, user.Name)

	// 3. SAFE READ
	var foundUser User
	if err := db.Where("email = ?", "gopher@golang.org").First(&foundUser).Error; err != nil {
		fmt.Println("User not found")
	} else {
		fmt.Printf("Read from DB: %s (%s)\n", foundUser.Name, foundUser.Email)
	}

	// 4. SAFE UPDATE
	if foundUser.ID != 0 {
		db.Model(&foundUser).Update("Name", "Master Gopher")
		fmt.Println("Update successful!")
	}
}
