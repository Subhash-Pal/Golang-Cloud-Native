
package main

import (
	"context"
	"errors"
	"fmt"
	"time"
)

// TransientError marks errors that are safe to retry
type TransientError struct{ msg string }
func (e *TransientError) Error() string { return e.msg }

// RetryOperation executes an action with exponential backoff & context awareness
func RetryOperation(ctx context.Context, maxRetries int, baseDelay time.Duration, op func() error) error {
	var err error
	for i := 0; i <= maxRetries; i++ {
		err = op()
		if err == nil {
			return nil
		}

		// Only retry if error is marked transient
		var tErr *TransientError
		if !errors.As(err, &tErr) {
			return fmt.Errorf("non-retryable error: %w", err)
		}

		if i == maxRetries {
			return fmt.Errorf("max retries (%d) exceeded: %w", maxRetries, err)
		}

		delay := baseDelay * time.Duration(1<<uint(i)) // 1s, 2s, 4s...
		fmt.Printf("⏳ Attempt %d failed: %v. Retrying in %v...\n", i+1, err, delay)

		select {
		case <-time.After(delay):
			// Wait passed, continue loop
		case <-ctx.Done():
			return fmt.Errorf("context cancelled during retry: %w", ctx.Err())
		}
	}
	return err
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 8*time.Second)
	defer cancel()

	attempt := 0
	err := RetryOperation(ctx, 3, time.Second, func() error {
		attempt++
		if attempt < 3 {
			return &TransientError{msg: "upstream service timeout"}
		}
		fmt.Println("✅ Operation succeeded on attempt", attempt)
		return nil
	})

	if err != nil {
		fmt.Printf("❌ Final error: %v\n", err)
	}
}
