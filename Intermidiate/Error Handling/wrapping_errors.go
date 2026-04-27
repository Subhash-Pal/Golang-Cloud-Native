package main

import (
	"errors"
	"fmt"
)

var ErrDatabase = errors.New("database failure")

func queryDB() error {
	// Wrap the original error with more context
	return fmt.Errorf("failed to get user: %w", ErrDatabase)
}

func main() {
	err := queryDB()
	fmt.Println("Full Error:", err)

	// Check if the wrapped error contains the sentinel error
	if errors.Is(err, ErrDatabase) {
		fmt.Println("Caught wrapped database error")
	}
}
