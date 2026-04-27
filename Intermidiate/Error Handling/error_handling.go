package main

import (
	"errors"
	"fmt"
	//"os"
	"strconv"
)

// ---------------------------------------------------------
// 1. CUSTOM ERROR TYPE
// Implementing the built-in `error` interface (Error() string)
// ---------------------------------------------------------
type DatabaseError struct {
	Operation string
	Code      int
	Message   string
}

// Error() makes DatabaseError satisfy the error interface
func (e *DatabaseError) Error() string {
	return fmt.Sprintf("db error [%s] (code %d): %s", e.Operation, e.Code, e.Message)
}

// ---------------------------------------------------------
// 2. ERROR-RETURNING FUNCTIONS
// ---------------------------------------------------------

// validateInput demonstrates returning a custom error
func validateInput(ageStr string) (int, error) {
	age, err := strconv.Atoi(ageStr)
	if err != nil {
		// Wrap standard library error with context
		return 0, fmt.Errorf("failed to parse age '%s': %w", ageStr, err)
	}
	if age < 0 {
		return 0, &DatabaseError{
			Operation: "validate",
			Code:      400,
			Message:   "age cannot be negative",
		}
	}
	return age, nil
}

// simulateDBOp demonstrates error wrapping & context chaining
func simulateDBOp(age int) error {
	if age > 120 {
		return fmt.Errorf("simulate db insert failed: %w", &DatabaseError{
			Operation: "insert",
			Code:      500,
			Message:   "age exceeds maximum allowed",
		})
	}
	fmt.Printf("✅ User with age %d saved successfully.\n", age)
	return nil
}

// ---------------------------------------------------------
// 3. MAIN: DEMONSTRATING ERROR INSPECTION PATTERNS
// ---------------------------------------------------------
func main() {
	testCases := []string{"25", "not_a_number", "-5", "130"}

	for _, tc := range testCases {
		fmt.Printf("\n🔍 Testing input: %q\n", tc)

		age, err := validateInput(tc)
		if err != nil {
			handleError(err)
			continue
		}

		if err := simulateDBOp(age); err != nil {
			handleError(err)
		}
	}
}

// ---------------------------------------------------------
// 4. CENTRALIZED ERROR HANDLER
// Shows how to inspect wrapped & custom errors properly
// ---------------------------------------------------------
func handleError(err error) {
	fmt.Println("❌ Error occurred:", err)

	// errors.Is: checks if a specific error is anywhere in the chain
	if errors.Is(err, strconv.ErrSyntax) {
		fmt.Println("   ↳ Root cause: Invalid number format")
	}

	// errors.As: extracts a specific custom error type from the chain
	var dbErr *DatabaseError
	if errors.As(err, &dbErr) {
		fmt.Printf("   ↳ Custom DB Error -> Operation: %s | Code: %d | Msg: %s\n",
			dbErr.Operation, dbErr.Code, dbErr.Message)
	}
}