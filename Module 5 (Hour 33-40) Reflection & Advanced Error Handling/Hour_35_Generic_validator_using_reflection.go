package main

import (
	"fmt"
	"reflect"
	"strings"
)

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
		return fmt.Errorf("input must be a pointer to a struct")
	}

	// Dereference the pointer to get the actual struct
	structValue := value.Elem()
	typeOf := structValue.Type()

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
			return err
		}
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
				return fmt.Errorf("field '%s' is required", fieldName)
			}
		case strings.HasPrefix(rule, "min="):
			minValue := parseRuleValue(rule, "min=")
			if fieldValue.Kind() == reflect.String && len(fieldValue.String()) < minValue {
				return fmt.Errorf("field '%s' must have a minimum length of %d", fieldName, minValue)
			} else if fieldValue.Kind() == reflect.Int && int(fieldValue.Int()) < minValue {
				return fmt.Errorf("field '%s' must be at least %d", fieldName, minValue)
			}
		case strings.HasPrefix(rule, "max="):
			maxValue := parseRuleValue(rule, "max=")
			if fieldValue.Kind() == reflect.String && len(fieldValue.String()) > maxValue {
				return fmt.Errorf("field '%s' must have a maximum length of %d", fieldName, maxValue)
			} else if fieldValue.Kind() == reflect.Int && int(fieldValue.Int()) > maxValue {
				return fmt.Errorf("field '%s' must not exceed %d", fieldName, maxValue)
			}
		case rule == "email":
			if fieldValue.Kind() == reflect.String && !isValidEmail(fieldValue.String()) {
				return fmt.Errorf("field '%s' must be a valid email address", fieldName)
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

func main() {
	// Create an instance of ExampleStruct
	example := &ExampleStruct{
		Name:     "Shubh",
		Age:      25,
		Email:    "shubh@example.com",
		IsActive: true,
	}

	// Validate the struct
	if err := Validate(example); err != nil {
		fmt.Println("Validation Error:", err)
	} else {
		fmt.Println("Validation Passed!")
	}
}

/*
Explanation of the Code
1. Struct Definition
The ExampleStruct contains fields annotated with validate tags.
Each tag specifies validation rules such as required, min, max, and email.
2. Validation Logic
The Validate function uses reflection to iterate over the struct fields.
It extracts the validate tag for each field and applies the corresponding validation rules.
3. Field Validation
The validateField function parses the validation rules and applies them dynamically.
Supported rules include:
required: Ensures the field is not empty.
min: Specifies the minimum length (for strings) or value (for integers).
max: Specifies the maximum length (for strings) or value (for integers).
email: Validates the field as an email address (basic check).
4. Helper Functions
isEmpty: Checks if a field is empty based on its type.
parseRuleValue: Extracts numeric values from rules like min=3.
isValidEmail: Performs a basic check for email validity.
Output Example
Validation Error: field 'Name' is required

Key Concepts Demonstrated
Dynamic Validation:
Reflection allows you to apply validation rules dynamically without hardcoding logic for each field.
Struct Tags:
Struct tags are used to define validation rules in a declarative manner.
Extensibility:
You can easily extend the validator to support additional rules (e.g., regex, custom functions).
Error Handling:
The validator returns detailed error messages indicating which field failed validation and why.

*/