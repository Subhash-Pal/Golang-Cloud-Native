package main

import (
	"errors"
	"fmt"
)

// Example of using errors.New to create a simple error.
func checkName(name string) error {
	if name == "" {
		return errors.New("name cannot be empty")
	}
	return nil
}

func main() {
	if err := checkName(""); err != nil {
		fmt.Println("Error:", err)
	}
}
