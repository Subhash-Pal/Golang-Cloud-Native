package main

import (
	"fmt"
	//"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Account struct {
	ID      uint `gorm:"primaryKey"`
	Name    string
	Balance float64
}

func main() {
	godotenv.Load()
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		os.Getenv("DB_HOST"), os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_NAME"), os.Getenv("DB_PORT"))

	db, _ := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	db.AutoMigrate(&Account{})

	// Setup: Create two accounts
	db.Save(&Account{ID: 1, Name: "Alice", Balance: 1000})
	db.Save(&Account{ID: 2, Name: "Bob", Balance: 500})

	// --- TRANSACTION START ---
	err := db.Transaction(func(tx *gorm.DB) error {
		// 1. Deduct from Alice
		if err := tx.Model(&Account{}).Where("id = ?", 1).Update("balance", gorm.Expr("balance - ?", 100)).Error; err != nil {
			return err // Returning error triggers ROLLBACK
		}

		// 2. Simulate an error (e.g., Bob's account is frozen or doesn't exist)
		// Change this to 'false' to see a successful COMMIT
		simulateError := false
		if simulateError {
			return fmt.Errorf("simulated network failure") // Rollback happens here!
		}

		// 3. Add to Bob
		if err := tx.Model(&Account{}).Where("id = ?", 2).Update("balance", gorm.Expr("balance + ?", 100)).Error; err != nil {
			return err
		}

		return nil // Returning nil triggers COMMIT
	})
	// --- TRANSACTION END ---

	if err != nil {
		fmt.Println("❌ Transaction Failed & Rolled Back:", err)
	} else {
		fmt.Println("✅ Transaction Committed Successfully!")
	}

	// Check final balances
	var accounts []Account
	db.Find(&accounts)
	fmt.Printf("Final States: %+v\n", accounts)
}

/*
How it works:
db.Transaction: This is the safest way to run transactions. It automatically handles BEGIN, COMMIT, and ROLLBACK.
The Trigger: If the function returns nil, GORM calls Commit. If it returns an error, GORM calls Rollback.
Atomic Power: In the example above, if simulateError is true, Alice’s balance won't actually decrease in the database, even though the code "ran" that line.
Manual Transactions (Advanced)
If you need more control, you can use manual triggers:
tx := db.Begin(): Start transaction.
tx.Rollback(): Undo changes.
tx.Commit(): Save changes forever.
*/
