package main

import (
	"fmt"
	"net/http"
	"time"
)
// int a , FileReqqader fr
// w http.ResponseWriter
func helloHandler(w http.ResponseWriter, r *http.Request) {
	// Ensure the request method is GET
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	// Write a response to the client
	fmt.Fprintf(w, "Hello, Shubh! Welcome to your Go HTTP server.")
}

func healthStatus(w http.ResponseWriter, r *http.Request) {
	// Ensure the request method is GET
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	now := time.Now()
	// Write a response to the client
	fmt.Fprintf(w, "Server is Running fine and healthy."+now.String())
}

func main() {
	// Define the route and its handler
	http.HandleFunc("/hello", helloHandler)
	http.HandleFunc("/health", healthStatus)

	// Start the HTTP server on port 8080
	fmt.Println("Starting server on :8080...")
	fmt.Println("http://localhost:8080/hello")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Printf("Error starting server: %s\n", err)
	}
}