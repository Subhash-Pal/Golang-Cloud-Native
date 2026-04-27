
package main

import (
	"errors"
	"fmt"
)

// safeExecute wraps a function that might panic
func safeExecute(fn func()) (recovered error) {
	defer func() {
		if r := recover(); r != nil {
			if err, ok := r.(error); ok {
				recovered = fmt.Errorf("recovered panic: %w", err)
			} else {
				recovered = fmt.Errorf("recovered panic: %v", r)
			}
		}
	}()

	fn()//risky operation that may panic
	return nil
}

func riskyOperation() {
	// Simulate unexpected runtime failure
	panic(errors.New("unexpected nil pointer dereference"))
}

func main() {
	fmt.Println("🛡️ Testing panic recovery...")
	
	err := safeExecute(riskyOperation)
	if err != nil {
		fmt.Printf("⚠️ Caught & converted to error: %v\n", err)
	} else {
		fmt.Println("✅ No panic occurred.")
	}
}
	