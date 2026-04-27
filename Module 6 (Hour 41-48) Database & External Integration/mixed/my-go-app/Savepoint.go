/*
In complex workflows, you may want to commit parts of a transaction while discarding others if specific steps fail. GORM provides SavePoint(name string) and RollbackTo(name string) to create these intermediate checkpoints.
GORM
GORM
 +2
Multi-Stage Transaction Example
This script simulates a "Course Enrollment" system:
Mandatory: Create the User.
Checkpoint: Set a SavePoint.
Optional: Attempt to grant a "Promotion Code." If this fails, we only roll back the promo, keeping the user created.
If the promo succeeds, we commit everything.
*/
package main

import (
	"errors"
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Student struct {
	ID   uint `gorm:"primaryKey"`
	Name string
}

type Promo struct {
	ID     uint `gorm:"primaryKey"`
	Code   string
	UsedBy uint
}

func main() {
	dsn := "host=localhost user=postgres password=root dbname=mydb port=5432 sslmode=disable"
	db, _ := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	db.AutoMigrate(&Student{}, &Promo{})

	// Start manual transaction
	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// STEP 1: Mandatory Operation
	student := Student{Name: "Rahul"}
	if err := tx.Create(&student).Error; err != nil {
		tx.Rollback()
		log.Fatal("Critical Error: Could not create student")
	}
	fmt.Println("Step 1: Student created.")

	// STEP 2: Set SavePoint
	tx.SavePoint("before_promo")
	fmt.Println("Checkpoint: 'before_promo' set.")

	// STEP 3: Optional Operation (The "Risky" part)
	promo := Promo{Code: "FREE100", UsedBy: student.ID}
	// Let's simulate a failure here (e.g., a database constraint violation)
	err := errors.New("promo code expired")

	if err != nil {
		fmt.Printf("Step 3 Failed: %v. Rolling back to checkpoint...\n", err)
		// Undo ONLY Step 3
		tx.RollbackTo("before_promo")
	} else {
		tx.Create(&promo)
		fmt.Println("Step 3: Promotion applied.")
	}

	// Finalize: Everything before the rollback (or everything if no failure) is saved
	tx.Commit()
	fmt.Println("Transaction Finished. Check your DB!")
}

/*
Critical Rules for SavePoints:
Manual Control: Unlike the standard db.Transaction block which is "all or nothing," SavePoint is typically used with db.Begin() for fine-grained manual control.
Release Memory: While not always required by drivers, some databases prefer you eventually commit or roll back the full transaction to clear resources.
Nested Transactions: If you use db.Transaction inside another db.Transaction block, GORM actually creates a SavePoint automatically behind the scenes to handle the nesting
*/
