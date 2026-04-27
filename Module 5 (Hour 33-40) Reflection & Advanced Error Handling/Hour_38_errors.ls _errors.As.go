/*
In Go, the errors.Is and errors.As functions (introduced in Go 1.13) are powerful tools for working with wrapped errors. They allow you to inspect and handle errors more effectively, especially when dealing with deeply nested error chains.
Below is an example that demonstrates how to use errors.Is and errors.As in practice.
Code: Using errors.Is and errors.As
*/

package main

import (
	"errors"
	"fmt"
)

// Sentinel Error
var ErrDatabase = errors.New("database error")

// CustomError represents a custom error type.
type CustomError struct {
	Code    int
	Message string
	Cause   error
}

// Implement the `error` interface by defining the `Error` method.
func (e *CustomError) Error() string {
	return fmt.Sprintf("Code: %d, Message: %s", e.Code, e.Message)
}

// Implement the `Unwrap` method for compatibility with `errors.Is` and `errors.As`.
func (e *CustomError) Unwrap() error {
	return e.Cause
}

// Simulate a function that returns a wrapped error.
func QueryDatabase(query string) error {
	if query == "" {
		// Wrap a sentinel error with additional context
		return &CustomError{
			Code:    500,
			Message: "Invalid query",
			Cause:   ErrDatabase,
		}
	}
	fmt.Println("Executing query:", query)
	return nil
}

func main() {
	// Example 1: Using errors.Is to check for a specific error
	err := QueryDatabase("")
	if err != nil {
		fmt.Println("Error occurred:", err)

		// Check if the error chain contains ErrDatabase
		if errors.Is(err, ErrDatabase) {
			fmt.Println("The error is caused by a database issue")
		} else {
			fmt.Println("The error is unrelated to the database")
		}
	}

	// Example 2: Using errors.As to extract a custom error
	var customErr *CustomError
	if errors.As(err, &customErr) {
		fmt.Println("Extracted Custom Error:")
		fmt.Println("Code:", customErr.Code)
		fmt.Println("Message:", customErr.Message)
		if customErr.Cause != nil {
			fmt.Println("Cause:", customErr.Cause)
		}
	} else {
		fmt.Println("The error is not a CustomError")
	}
}

/*
Explanation of the Code
1. Sentinel Error
ErrDatabase is a sentinel error, which is a predefined error used to identify specific error conditions.
Sentinel errors are useful for comparing errors using ==.
2. Custom Error Type
The CustomError struct includes:
Code: An integer representing an error code.
Message: A detailed error message.
Cause: An optional underlying error that caused this error.
It implements the error interface via the Error method.
It also implements the Unwrap method, which allows it to work with errors.Is and errors.As.
3. Error Wrapping
The QueryDatabase function simulates a database query.
If the query is invalid, it wraps the ErrDatabase sentinel error with additional context using CustomError.
4. Using errors.Is
errors.Is checks if a specific error (e.g., ErrDatabase) exists anywhere in the error chain.
This is useful for identifying root causes without manually unwrapping errors.
5. Using errors.As
errors.As extracts the first error in the chain that matches a specific type (e.g., *CustomError).
This allows you to access the fields of a custom error type directly.
Output Example
When you run the above code, the output will look something like this:

bash```
Error occurred: Code: 500, Message: Invalid query
The error is caused by a database issue
Extracted Custom Error:
Code: 500
Message: Invalid query
Cause: database error
```
Key Concepts Demonstrated
Sentinel Errors:
Sentinel errors are predefined errors used to identify specific error conditions.
They can be compared directly using == or checked using errors.Is.
Error Wrapping:
Wrapping errors with additional context helps preserve the root cause while adding meaningful information.
errors.Is:
errors.Is checks if a specific error exists anywhere in the error chain.
It simplifies error handling by avoiding manual unwrapping.
errors.As:
errors.As extracts the first error in the chain that matches a specific type.
It allows you to access the fields of custom error types directly.
Advanced Usage: Deeply Nested Errors
You can simulate deeply nested errors and still use errors.Is and errors.As effectively:

go

func deepError() error {
	err := QueryDatabase("")
	return fmt.Errorf("deep error: %w", err)
}

func main() {
	err := deepError()

	// Use errors.Is to check for a specific error
	if errors.Is(err, ErrDatabase) {
		fmt.Println("The error is caused by a database issue")
	}

	// Use errors.As to extract a custom error
	var customErr *CustomError
	if errors.As(err, &customErr) {
		fmt.Println("Extracted Custom Error:")
		fmt.Println("Code:", customErr.Code)
		fmt.Println("Message:", customErr.Message)
	}
}
	This will output:
	The error is caused by a database issue
Extracted Custom Error:
Code: 500
Message: Invalid query


Notes
Best Practices:
Use errors.Is for checking specific sentinel errors.
Use errors.As for extracting custom error types.
Always implement the Unwrap method for custom error types to make them compatible with errors.Is and errors.As.
Compatibility:
errors.Is and errors.As work seamlessly with errors wrapped using fmt.Errorf with %w.
Extensibility:
You can extend custom error types with additional fields or methods as needed.


*/