package main

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"time"

	"github.com/jackc/pgconn" // For PgError type
	//go get github.com/jackc/pgx/v5/pgconn
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	dsn := "host=localhost user=postgres password=root dbname=mydb port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Printf("Failed to connect: %v\n", err)
		return
	}

	// 1. Create a Context with a 5-second deadline
	// This will force the entire process to stop if it exceeds 5 seconds.
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	maxRetries := 3
	baseDelay := 100 * time.Millisecond

	// 2. Pass the context into our retry function
	err = retryWithContext(ctx, db, maxRetries, baseDelay)

	if err != nil {
		fmt.Printf("\n❌ Final Status: %v\n", err)
	} else {
		fmt.Println("\n✅ Final Status: Transaction Succeeded!")
	}
}

func retryWithContext(ctx context.Context, db *gorm.DB, maxAttempts int, delay time.Duration) error {
	for i := 1; i <= maxAttempts; i++ {
		// 3. CHECK: Has the context already been cancelled or timed out?
		if err := ctx.Err(); err != nil {
			return fmt.Errorf("context stopped: %w", err)
		}

		// 4. BIND CONTEXT: Use WithContext to link the transaction to the context
		err := db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
			fmt.Printf("🚀 Attempt %d: Executing database operation...\n", i)

			// SIMULATION: Fail with a "Serialization Failure" (40001) 60% of the time
			if rand.Intn(10) < 6 {
				return &pgconn.PgError{
					Code:    "40001",
					Message: "could not serialize access due to concurrent update",
				}
			}
			return nil
		})

		// Success!
		if err == nil {
			return nil
		}

		// 5. ANALYSIS: Determine if the error is retryable
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && (pgErr.Code == "40001" || pgErr.Code == "40P01") {

			// 6. WAIT OR CANCEL: Use a select block to sleep or exit early if context expires
			jitter := time.Duration(rand.Intn(50)) * time.Millisecond
			currentDelay := delay + jitter

			fmt.Printf("⚠️  Retryable Error [%s]. Waiting %v...\n", pgErr.Code, currentDelay)

			select {
			case <-time.After(currentDelay):
				// Delay finished, loop continues to next attempt
				delay *= 2
			case <-ctx.Done():
				// Context expired while waiting
				return ctx.Err()
			}
		} else {
			// Non-retryable error (or context timeout error from tx)
			return fmt.Errorf("transaction aborted: %v", err)
		}
	}
	return fmt.Errorf("max retries (%d) reached", maxAttempts)
}

/*
Output Example (will vary due to randomness):
go run .\Identify-Retriable-Errors.go
🚀 Attempt 1: Executing database operation...
⚠️  Retryable Error [40001]. Waiting 102ms...
🚀 Attempt 2: Executing database operation...
⚠️  Retryable Error [40001]. Waiting 223ms...
🚀 Attempt 3: Executing database operation...

✅ Final Status: Transaction Succeeded!
*/

/*
Why Context is Important here:
Early Exit: If you have a 2-second retry delay, but the context expires in 0.5 seconds, the select block will trigger <-ctx.Done() and exit immediately instead of waiting for a useless retry.
Native GORM Support: By using db.WithContext(ctx), GORM automatically attaches the context to its internal database calls. If the context times out mid-query, the query is cancelled by the database driver.
Traceability: Context is the standard way to pass "request-scoped" data (like trace IDs) through your application layers alongside your database transactions

*/
