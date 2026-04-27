package main

import (
	"errors"
	"fmt"
)

// MockDatabase simulates a resource that requires explicit cleanup
type MockDatabase struct {
	isOpen bool
}

func NewMockDatabase() (*MockDatabase, error) {
	fmt.Println("🌐 Establishing connection...")
	return &MockDatabase{isOpen: true}, nil
}

// Close cleans up the resource
func (db *MockDatabase) Close() error {
	if !db.isOpen {
		return fmt.Errorf("database already closed")
	}
	db.isOpen = false
	fmt.Println("🔌 Connection safely closed.")
	return nil
}

// ProcessRecord simulates work on a single record
func (db *MockDatabase) ProcessRecord(id int) error {
	if !db.isOpen {
		return fmt.Errorf("cannot process: connection is closed")
	}
	if id == 3 {
		return fmt.Errorf("record %d: data corruption detected", id)
	}
	fmt.Printf("✅ Record %d processed.\n", id)
	return nil
}

// runBatch demonstrates defer cleanup + error aggregation
func runBatch() error {
	db, err := NewMockDatabase()
	if err != nil {
		return fmt.Errorf("failed to initialize: %w", err)
	}

	// defer guarantees cleanup runs, even if we return early
	var closeErr error
	defer func() {
		if err := db.Close(); err != nil {
			closeErr = fmt.Errorf("cleanup failed: %w", err)
		}
	}()

	records := []int{1, 2, 3, 4, 5}
	var errs []error

	for _, id := range records {
		if err := db.ProcessRecord(id); err != nil {
			// Collect errors instead of failing immediately (resilience pattern)
			errs = append(errs, err)
			continue
		}
	}

	// If any records failed, aggregate them
	if len(errs) > 0 {
		// Go 1.20+ errors.Join creates a single error containing all failures
		aggErr := errors.Join(errs...)
		// Combine processing errors with potential cleanup errors
		return errors.Join(aggErr, closeErr)
	}

	return closeErr // Return cleanup error if no processing errors occurred
}

func main() {
	fmt.Println("🚀 Starting batch job...\n")
	err := runBatch()
	if err != nil {
		fmt.Printf("\n❌ Batch job completed with errors:\n%v\n", err)
	} else {
		fmt.Println("\n✅ Batch job completed successfully.")
	}
}