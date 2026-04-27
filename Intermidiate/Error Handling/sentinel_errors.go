package main

import (
	"errors"
	"fmt"
)

// Defining a sentinel error
var ErrNotFound = errors.New("resource not found")

func fetchData(id int) error {
	if id == 0 {
		return ErrNotFound
	}
	return nil
}

func main() {
	err := fetchData(0)
	// Using errors.Is to check for a specific error type
	if errors.Is(err, ErrNotFound) {
		fmt.Println("Handle 404: Not Found")
	}
}
