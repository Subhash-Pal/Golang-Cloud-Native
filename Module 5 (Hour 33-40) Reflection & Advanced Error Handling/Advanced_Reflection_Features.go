package main

import (
	"fmt"
	"reflect"
)

// ExampleStruct represents a sample struct to demonstrate advanced reflection.
type ExampleStruct struct {
	Name     string
	Age      int
	IsActive bool
}

// UpdateField updates a field of a struct dynamically using reflection.
func UpdateField(obj interface{}, fieldName string, newValue interface{}) error {
	// Get the value of the object (must be a pointer to modify it)
	value := reflect.ValueOf(obj)
	if value.Kind() != reflect.Ptr || value.Elem().Kind() != reflect.Struct {
		return fmt.Errorf("obj must be a pointer to a struct")
	}

	// Dereference the pointer to get the actual struct
	structValue := value.Elem()

	// Find the field by name
	field := structValue.FieldByName(fieldName)
	if !field.IsValid() {
		return fmt.Errorf("field '%s' does not exist", fieldName)
	}

	// Ensure the field is settable
	if !field.CanSet() {
		return fmt.Errorf("field '%s' is not settable", fieldName)
	}

	// Convert the new value to the correct type
	newValueReflect := reflect.ValueOf(newValue)
	if field.Type() != newValueReflect.Type() {
		return fmt.Errorf("type mismatch: cannot assign %v to field of type %v", newValueReflect.Type(), field.Type())
	}

	// Set the new value
	field.Set(newValueReflect)
	return nil
}

// CallFunction dynamically calls a function using reflection.
func CallFunction(fn interface{}, args ...interface{}) ([]reflect.Value, error) {
	// Get the function value
	fnValue := reflect.ValueOf(fn)
	if fnValue.Kind() != reflect.Func {
		return nil, fmt.Errorf("provided argument is not a function")
	}

	// Prepare arguments for the function call
	in := make([]reflect.Value, len(args))
	for i, arg := range args {
		in[i] = reflect.ValueOf(arg)
	}

	// Call the function dynamically
	results := fnValue.Call(in)
	return results, nil
}

func main() {
	// Example 1: Dynamically update a struct field
	example := &ExampleStruct{
		Name:     "Shubh",
		Age:      25,
		IsActive: true,
	}
	fmt.Println("Before Update:", example)

	err := UpdateField(example, "Name", "John Doe")
	if err != nil {
		fmt.Println("Error updating field:", err)
	} else {
		fmt.Println("After Update:", example)
	}

	// Example 2: Dynamically call a function
	add := func(a, b int) int {
		return a + b
	}

	results, err := CallFunction(add, 10, 20)
	if err != nil {
		fmt.Println("Error calling function:", err)
	} else {
		fmt.Println("Function Result:", results[0].Interface())
	}
}
/*

Explanation of the Code
1. Dynamic Field Updates (UpdateField)
The UpdateField function uses reflection to dynamically update a field of a struct.
It ensures that:
The input is a pointer to a struct.
The specified field exists and is settable.
The new value matches the field's type.
If all conditions are met, the field's value is updated using field.Set().
2. Dynamic Function Invocation (CallFunction)
The CallFunction function uses reflection to dynamically invoke a function.
It:
Validates that the provided argument is a function.
Converts the input arguments into reflect.Value objects.
Calls the function dynamically using fnValue.Call().
The results of the function call are returned as a slice of reflect.Value.
Output Example
When you run the above code, the output will look something like this:


Before Update: &{Shubh 25 true}
After Update: &{John Doe 25 true}
Function Result: 30
*/