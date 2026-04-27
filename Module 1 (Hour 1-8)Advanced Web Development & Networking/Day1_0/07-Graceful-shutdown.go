package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// Handler for "/hello" route
func helloHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Simulate some work with a delay
	time.Sleep(5 * time.Second) // Simulate a long-running request
	fmt.Fprintln(w, "Hello! Your request has been processed.")
}

func main() {
	// Create a custom ServeMux
	mux := http.NewServeMux()

	// Register the "/hello" route
	mux.HandleFunc("/hello", helloHandler)

	// Create an HTTP server
	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	// Start the server in a goroutine
	go func() {
		log.Println("Server started on :8080")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Error starting server: %s\n", err)
		}
	}()

	// Wait for interrupt signals to gracefully shut down
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM) // Listen for SIGINT (Ctrl+C) and SIGTERM

	log.Println("Press Ctrl+C to stop the server...")
	<-sig

	log.Println("Shutting down server...")

	// Create a context with a timeout for graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Attempt to shut down the server gracefully
	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Error during shutdown: %s\n", err)
	}

	log.Println("Server stopped gracefully.")
}