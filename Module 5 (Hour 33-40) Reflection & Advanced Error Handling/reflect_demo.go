package main

import (
	"fmt"
	"reflect"
)

// ExampleStruct represents a sample struct to demonstrate reflection.
type ExampleStruct struct {
	Name     string
	Age      int
	IsActive bool
}

func main() {
	// Create an instance of ExampleStruct
	example := ExampleStruct{
		Name:     "Shubh",
		Age:      25,
		IsActive: true,
	}

	// Reflect on the type and value of the struct
	fmt.Println("=== Reflecting on Struct ===")
	reflectOnStruct(example)

	// Reflect on a slice
	fmt.Println("\n=== Reflecting on Slice ===")
	numbers := []int{1, 2, 3, 4, 5}
	reflectOnSlice(numbers)

	// Reflect on a map
	fmt.Println("\n=== Reflecting on Map ===")
	userInfo := map[string]interface{}{
		"Name": "Shubh",
		"Age":  25,
	}
	reflectOnMap(userInfo)
}

// reflectOnStruct demonstrates reflection on a struct.
func reflectOnStruct(input interface{}) {
	value := reflect.ValueOf(input)
	typeOf := value.Type()

	fmt.Printf("Type: %s\n", typeOf)
	fmt.Printf("Number of fields: %d\n", value.NumField())

	for i := 0; i < value.NumField(); i++ {
		field := typeOf.Field(i)
		fieldValue := value.Field(i)

		fmt.Printf("Field Name: %s, Type: %s, Value: %v\n", field.Name, field.Type, fieldValue)
	}
}

// reflectOnSlice demonstrates reflection on a slice.
func reflectOnSlice(input interface{}) {
	value := reflect.ValueOf(input)
	typeOf := value.Type()

	fmt.Printf("Type: %s\n", typeOf)
	fmt.Printf("Length: %d, Capacity: %d\n", value.Len(), value.Cap())

	for i := 0; i < value.Len(); i++ {
		element := value.Index(i)
		fmt.Printf("Index %d: Value: %v, Type: %s\n", i, element.Interface(), element.Type())
	}
}

// reflectOnMap demonstrates reflection on a map.
func reflectOnMap(input interface{}) {
	value := reflect.ValueOf(input)
	typeOf := value.Type()

	fmt.Printf("Type: %s\n", typeOf)
	fmt.Printf("Number of keys: %d\n", value.Len())

	for _, key := range value.MapKeys() {
		mapValue := value.MapIndex(key)
		fmt.Printf("Key: %v (%s), Value: %v (%s)\n", key.Interface(), key.Type(), mapValue.Interface(), mapValue.Type())
	}
}
/*

Explanation of the Code
Struct Reflection:
The reflectOnStruct function uses reflect.ValueOf and reflect.TypeOf to inspect the fields of a struct.
It iterates over the fields using NumField() and retrieves their names, types, and values.
Slice Reflection:
The reflectOnSlice function reflects on a slice, printing its type, length, and capacity.
It iterates through the elements of the slice using Index() and prints their values and types.
Map Reflection:
The reflectOnMap function reflects on a map, printing its type and the number of keys.
It iterates over the keys using MapKeys() and retrieves the corresponding values using MapIndex().
Output Example
When you run the above code, the output will look something like this:


=== Reflecting on Struct ===
Type: main.ExampleStruct
Number of fields: 3
Field Name: Name, Type: string, Value: Shubh
Field Name: Age, Type: int, Value: 25
Field Name: IsActive, Type: bool, Value: true

=== Reflecting on Slice ===
Type: []int
Length: 5, Capacity: 5
Index 0: Value: 1, Type: int
Index 1: Value: 2, Type: int
Index 2: Value: 3, Type: int
Index 3: Value: 4, Type: int
Index 4: Value: 5, Type: int

=== Reflecting on Map ===
Type: map[string]interface {}
Number of keys: 2
Key: Name (string), Value: Shubh (string)
Key: Age (string), Value: 25 (int)
*/