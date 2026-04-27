package main

import (
	"errors"
	"fmt"
)

// CustomError represents a custom error type with additional fields.
type CustomError struct {
	Code    int    // HTTP status code or error code
	Message string // Detailed error message
	Cause   error  // Underlying cause of the error (optional)
}

// Implement the `error` interface by defining the `Error` method.
func (e *CustomError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("Code: %d, Message: %s, Cause: %v", e.Code, e.Message, e.Cause)
	}
	return fmt.Sprintf("Code: %d, Message: %s", e.Code, e.Message)
}

// NewCustomError creates a new instance of CustomError.
func NewCustomError(code int, message string, cause error) error {
	return &CustomError{
		Code:    code,
		Message: message,
		Cause:   cause,
	}
}

// Example function that may return a custom error.
func Divide(a, b float64) (float64, error) {
	if b == 0 {
		// Return a custom error for division by zero
		return 0, NewCustomError(400, "Division by zero is not allowed", nil)
	}
	return a / b, nil
}

// Example function that wraps an error with additional context.
func ProcessData(data string) error {
	if data == "" {
		// Create a base error
		baseErr := errors.New("data cannot be empty")
		// Wrap it in a custom error
		return NewCustomError(500, "Failed to process data", baseErr)
	}
	fmt.Println("Processing data:", data)
	return nil
}

func main() {
	// Example 1: Division by zero
	result, err := Divide(10, 0)
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println("Result:", result)
	}

	// Example 2: Processing empty data
	err = ProcessData("")
	if err != nil {
		fmt.Println("Error:", err)
	}

	// Example 3: Type assertion to access custom error fields
	err = ProcessData("")
	if customErr, ok := err.(*CustomError); ok {
		fmt.Println("Custom Error Details:")
		fmt.Println("Code:", customErr.Code)
		fmt.Println("Message:", customErr.Message)
		if customErr.Cause != nil {
			fmt.Println("Cause:", customErr.Cause)
		}
	}
}
/*

Explanation of the Code
1. Custom Error Type
The CustomError struct contains:
Code: An integer representing an error code (e.g., HTTP status code).
Message: A detailed error message.
Cause: An optional underlying error that caused this error.

2. Implementing the error Interface
To make CustomError compatible with Go's error handling, it implements the Error method from the error interface.
The Error method formats the error details into a string.
3. Creating Custom Errors
The NewCustomError function is a constructor for creating instances of CustomError.
It allows you to specify the error code, message, and optional cause.
4. Using Custom Errors
The Divide function demonstrates returning a custom error when dividing by zero.
The ProcessData function demonstrates wrapping an existing error with additional context using CustomError.
5. Type Assertion
In the main function, type assertion (err.(*CustomError)) is used to extract the fields of a custom error for further inspection.
Output Example
When you run the above code, the output will look something like this:

Error: Code: 400, Message: Division by zero is not allowed
Error: Code: 500, Message: Failed to process data, Cause: data cannot be empty
Custom Error Details:
Code: 500
Message: Failed to process data
Cause: data cannot be empty



Key Concepts Demonstrated
Custom Error Types:
Custom error types allow you to include additional metadata (e.g., error codes, causes) beyond a simple string message.
Error Wrapping:
Wrapping errors with additional context helps preserve the root cause while adding meaningful information.
Type Assertion:
Type assertion allows you to inspect and handle custom error types specifically.
Use Cases:
APIs: Return structured errors with HTTP status codes.
Logging: Include detailed error information for debugging.
Error Handling: Handle different error types differently based on their structure.



Advanced Usage: Error Unwrapping
Go 1.13 introduced the errors.Is and errors.As functions for working with wrapped errors. You can extend the example to support these features:
go

import "errors"

// Unwrap implements the `Unwrap` method for compatibility with `errors.Is` and `errors.As`.
func (e *CustomError) Unwrap() error {
	return e.Cause
}

func main() {
	// Example: Using errors.Is to check for a specific error
	baseErr := errors.New("data cannot be empty")
	wrappedErr := NewCustomError(500, "Failed to process data", baseErr)

	if errors.Is(wrappedErr, baseErr) {
		fmt.Println("The error is caused by 'data cannot be empty'")
	}

	// Example: Using errors.As to extract a custom error
	var customErr *CustomError
	if errors.As(wrappedErr, &customErr) {
		fmt.Println("Extracted Custom Error:", customErr)
	}
}

This will output:
The error is caused by 'data cannot be empty'
Extracted Custom Error: Code: 500, Message: Failed to process data, Cause: data cannot be empty
Notes
Error Wrapping Best Practices:
Use fmt.Errorf with %w for simple wrapping.
Use custom error types for structured errors with additional fields.
Compatibility:
Ensure your custom error types implement the Unwrap method if you want them to work with errors.Is and errors.As.
Extensibility:
You can extend CustomError with additional fields or methods as needed.

*/