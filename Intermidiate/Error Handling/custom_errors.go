package main

import (
	"fmt"
)

type RequestError struct {
	StatusCode int
	Message    string
}

// Implement the Error() method
func (e *RequestError) Error() string {
	return fmt.Sprintf("status %d: %s", e.StatusCode, e.Message)
}

func doRequest() error {
	return &RequestError{StatusCode: 404, Message: "API missing"}
}

func main() {
	if err := doRequest(); err != nil {
		fmt.Println("Error:", err)
	}
}
