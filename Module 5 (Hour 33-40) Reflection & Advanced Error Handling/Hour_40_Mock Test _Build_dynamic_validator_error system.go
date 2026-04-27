/*
To build a dynamic validator + error system as part of the mock test for Hour 40, we will combine several concepts from the previous hours, including:
Reflection: To dynamically validate struct fields.
Custom Error Types: To provide structured and detailed error messages.
Error Wrapping: To include context about validation failures.
Advanced Error Handling: Using errors.Is and errors.As for better error inspection.
Below is the implementation of a dynamic validator with an advanced error-handling system.
Code: Dynamic Validator + Error System
go

*/
package main

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
)

// CustomError represents a custom error type for validation failures.
type CustomError struct {
	Field   string // The field that failed validation
	Message string // Detailed error message
	Cause   error  // Underlying cause of the error (optional)
}

// Implement the `error` interface by defining the `Error` method.
func (e *CustomError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("Field: %s, Message: %s, Cause: %v", e.Field, e.Message, e.Cause)
	}
	return fmt.Sprintf("Field: %s, Message: %s", e.Field, e.Message)
}

// Implement the `Unwrap` method for compatibility with `errors.Is` and `errors.As`.
func (e *CustomError) Unwrap() error {
	return e.Cause
}

// NewCustomError creates a new instance of CustomError.
func NewCustomError(field, message string, cause error) error {
	return &CustomError{
		Field:   field,
		Message: message,
		Cause:   cause,
	}
}

// ExampleStruct demonstrates struct tags with validation rules.
type ExampleStruct struct {
	Name     string `validate:"required,min=3,max=50"`
	Age      int    `validate:"required,min=18,max=100"`
	Email    string `validate:"email"`
	IsActive bool   `validate:"-"`
}

// Validate performs generic validation on a struct using reflection.
func Validate(obj interface{}) error {
	// Ensure the input is a pointer to a struct
	value := reflect.ValueOf(obj)
	if value.Kind() != reflect.Ptr || value.Elem().Kind() != reflect.Struct {
		return errors.New("input must be a pointer to a struct")
	}

	// Dereference the pointer to get the actual struct
	structValue := value.Elem()
	typeOf := structValue.Type()

	var validationErrors []error

	// Iterate over the struct fields
	for i := 0; i < structValue.NumField(); i++ {
		field := typeOf.Field(i)
		tag := field.Tag.Get("validate")

		// Skip fields with no validation tag or "-" tag
		if tag == "" || tag == "-" {
			continue
		}

		// Get the field value
		fieldValue := structValue.Field(i)

		// Parse and apply validation rules
		if err := validateField(field.Name, fieldValue, tag); err != nil {
			validationErrors = append(validationErrors, err)
		}
	}

	// If there are validation errors, return a wrapped error
	if len(validationErrors) > 0 {
		return NewCustomError("Validation", "One or more fields failed validation", joinErrors(validationErrors))
	}

	return nil
}

// validateField validates a single field based on its tag.
func validateField(fieldName string, fieldValue reflect.Value, tag string) error {
	rules := strings.Split(tag, ",")
	for _, rule := range rules {
		switch {
		case rule == "required":
			if isEmpty(fieldValue) {
				return NewCustomError(fieldName, "is required", nil)
			}
		case strings.HasPrefix(rule, "min="):
			minValue := parseRuleValue(rule, "min=")
			if fieldValue.Kind() == reflect.String && len(fieldValue.String()) < minValue {
				return NewCustomError(fieldName, fmt.Sprintf("must have a minimum length of %d", minValue), nil)
			} else if fieldValue.Kind() == reflect.Int && int(fieldValue.Int()) < minValue {
				return NewCustomError(fieldName, fmt.Sprintf("must be at least %d", minValue), nil)
			}
		case strings.HasPrefix(rule, "max="):
			maxValue := parseRuleValue(rule, "max=")
			if fieldValue.Kind() == reflect.String && len(fieldValue.String()) > maxValue {
				return NewCustomError(fieldName, fmt.Sprintf("must have a maximum length of %d", maxValue), nil)
			} else if fieldValue.Kind() == reflect.Int && int(fieldValue.Int()) > maxValue {
				return NewCustomError(fieldName, fmt.Sprintf("must not exceed %d", maxValue), nil)
			}
		case rule == "email":
			if fieldValue.Kind() == reflect.String && !isValidEmail(fieldValue.String()) {
				return NewCustomError(fieldName, "must be a valid email address", nil)
			}
		}
	}
	return nil
}

// Helper function to check if a value is empty.
func isEmpty(value reflect.Value) bool {
	switch value.Kind() {
	case reflect.String:
		return value.String() == ""
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return value.Int() == 0
	case reflect.Bool:
		return !value.Bool()
	default:
		return false
	}
}

// Helper function to parse rule values (e.g., min=3).
func parseRuleValue(rule, prefix string) int {
	valueStr := strings.TrimPrefix(rule, prefix)
	var value int
	fmt.Sscanf(valueStr, "%d", &value)
	return value
}

// Helper function to validate email format (basic regex-like check).
func isValidEmail(email string) bool {
	return strings.Contains(email, "@") && strings.Contains(email, ".")
}

// Helper function to join multiple errors into a single error.
func joinErrors(errs []error) error {
	if len(errs) == 0 {
		return nil
	}
	messages := make([]string, len(errs))
	for i, err := range errs {
		messages[i] = err.Error()
	}
	return errors.New(strings.Join(messages, "; "))
}

func main() {
	// Create an instance of ExampleStruct
	example := &ExampleStruct{
		Name:     "", // Violates 'required' rule
		Age:      17, // Violates 'min=18' rule
		Email:    "invalid-email", // Violates 'email' rule
		IsActive: true,
	}

	// Validate the struct
	err := Validate(example)
	if err != nil {
		fmt.Println("Validation Error:", err)

		// Use errors.As to extract custom errors
		var customErr *CustomError
		if errors.As(err, &customErr) {
			fmt.Println("\nExtracted Custom Error:")
			fmt.Println("Field:", customErr.Field)
			fmt.Println("Message:", customErr.Message)
			if customErr.Cause != nil {
				fmt.Println("Cause:", customErr.Cause)
			}
		}
	} else {
		fmt.Println("Validation Passed!")
	}
}

/*
Explanation of the Code
1. Dynamic Validation
The Validate function uses reflection to iterate over struct fields and validate them based on their validate tags.
It collects all validation errors and returns them as a single wrapped error.
2. Custom Error Type
The CustomError struct includes:
Field: The name of the field that failed validation.
Message: A detailed error message.
Cause: An optional underlying error.
3. Error Wrapping
Multiple validation errors are joined into a single error using joinErrors.
The top-level error wraps all individual errors, providing context about the overall validation failure.
4. Error Inspection
The main function demonstrates how to use errors.As to extract and inspect custom errors.
Output Example
When you run the above code, the output will look something like this:

Validation Error: Field: Validation, Message: One or more fields failed validation, Cause: Field: Name, Message: is required; Field: Age, Message: must be at least 18; Field: Email, Message: must be a valid email address

Extracted Custom Error:
Field: Validation
Message: One or more fields failed validation
Cause: Field: Name, Message: is required; Field: Age, Message: must be at least 18; Field: Email, Message: must be a valid email address

Key Concepts Demonstrated
Dynamic Validation:
Reflection allows you to validate struct fields dynamically without hardcoding logic for each field.
Custom Error Types:
Structured errors provide detailed information about validation failures.
Error Wrapping:
Wrapping multiple errors into a single error simplifies error handling.
Error Inspection:
errors.As allows you to extract and handle custom error types specifically.

*/