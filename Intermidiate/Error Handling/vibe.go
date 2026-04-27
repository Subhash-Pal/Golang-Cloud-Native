//write golang code to demonstrate interface which can store all other golang data types and also demonstrate error handling in golang
package main

import (
	"errors"
	"fmt"
)	

// Define a generic interface that can hold any data type
type Any interface{}

// Function to demonstrate error handling
func demonstrateErrorHandling() error {
	// Simulate an error condition
	if false {
		return errors.New("an error occurred")
	}
	return nil
}	

func checkType(data Any) {
	switch v := data.(type) {
	case int:		
		fmt.Println("Data is an integer:", v)
	case string:
		fmt.Println("Data is a string:", v)
	case []int:
		fmt.Println("Data is a slice of integers:", v)	
	default:
		fmt.Println("Data is of an unknown type:", v)
	}
}

func main() {
	// Using the Any interface to store different data types
	var data Any		
	data = 42
	fmt.Println("Integer:", data)	
	data = "Hello, World!"
	fmt.Println("String:", data)	
	data = []int{1, 2, 3, 4, 5}
	fmt.Println("Slice:", data)
	// Demonstrate error handling
	err := demonstrateErrorHandling()
	if err != nil {
		fmt.Println("Error:", err)
	}

	// Simulate an error condition
	err = errors.New("something went wrong")
	if err != nil {	
		fmt.Println("Error:", err)
	}	
	data = err
	fmt.Println("Stored error in Any interface:", data)	
	//reflecting on the type of data stored in Any interface
	

	checkType(10)
	
}
