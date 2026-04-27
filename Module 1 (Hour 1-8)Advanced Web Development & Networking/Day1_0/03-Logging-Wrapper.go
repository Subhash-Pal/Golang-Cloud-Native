package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"context"
)

// LoggingWrapper wraps an http.Handler and logs details about incoming requests.
func LoggingWrapper(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Log request details
		log.Printf("Request received: %s %s", r.Method, r.URL.Path)
		// Call the next handler in the chain
		next.ServeHTTP(w, r)
		// Log response details
		duration := time.Since(start)
		log.Printf("Request completed: %s %s (%v)", r.Method, r.URL.Path, duration)
	})
}

// Handler for "/hello" route
func helloHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	fmt.Fprintln(w, "Hello, Shubh! Welcome to your Go HTTP server.")
}

// Handler for "/custom" route
func customHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "This is a custom route!")
}

func main() {
	// Create a custom ServeMux
	mux := http.NewServeMux()

	
	// Register routes
	mux.HandleFunc("/hello", helloHandler)
	mux.HandleFunc("/custom", customHandler)

	// Wrap the entire ServeMux with the LoggingWrapper
	finalHandler := LoggingWrapper(mux)

	// Create an HTTP server
	server := &http.Server{
		Addr:    ":8080",
		Handler: finalHandler,
	}

	// Start the server in a goroutine
	go func() {
		log.Println("Server started on :8080")
		log.Println("Access the following URLs:")
		log.Println("- http://localhost:8080/hello")
		log.Println("- http://localhost:8080/custom")

		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Error starting server: %s\n", err)
		}
	}()

	// Wait for interrupt signal to gracefully shut down
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

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

/*
Certainly! A **logging wrapper** is a common pattern in Go for adding logging functionality to HTTP handlers. Instead of manually wrapping each handler with middleware, you can create a reusable `LoggingHandler` function that wraps any `http.Handler` and logs details about incoming requests.

Below is an example of a **logging wrapper** implementation in Go:

---

### **Code Example: Logging Wrapper**

```go
package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// LoggingWrapper wraps an http.Handler and logs details about incoming requests.
func LoggingWrapper(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Log request details
		log.Printf("Request received: %s %s", r.Method, r.URL.Path)

		// Call the next handler in the chain
		next.ServeHTTP(w, r)

		// Log response details
		duration := time.Since(start)
		log.Printf("Request completed: %s %s (%v)", r.Method, r.URL.Path, duration)
	})
}

// Handler for "/hello" route
func helloHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	fmt.Fprintln(w, "Hello, Shubh! Welcome to your Go HTTP server.")
}

// Handler for "/custom" route
func customHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "This is a custom route!")
}

func main() {
	// Create a custom ServeMux
	mux := http.NewServeMux()

	// Register routes
	mux.HandleFunc("/hello", helloHandler)
	mux.HandleFunc("/custom", customHandler)

	// Wrap the entire ServeMux with the LoggingWrapper
	finalHandler := LoggingWrapper(mux)

	// Create an HTTP server
	server := &http.Server{
		Addr:    ":8080",
		Handler: finalHandler,
	}

	// Start the server in a goroutine
	go func() {
		log.Println("Server started on :8080")
		log.Println("Access the following URLs:")
		log.Println("- http://localhost:8080/hello")
		log.Println("- http://localhost:8080/custom")

		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Error starting server: %s\n", err)
		}
	}()

	// Wait for interrupt signal to gracefully shut down
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

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
```

---

### **Explanation of the Code**

#### **1. LoggingWrapper Function**
```go
func LoggingWrapper(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Log request details
		log.Printf("Request received: %s %s", r.Method, r.URL.Path)

		// Call the next handler in the chain
		next.ServeHTTP(w, r)

		// Log response details
		duration := time.Since(start)
		log.Printf("Request completed: %s %s (%v)", r.Method, r.URL.Path, duration)
	})
}
```
- **What It Does**:
  - Logs the start of the request, including the HTTP method and URL path.
  - Calls the next handler (`next.ServeHTTP`) to process the request.
  - Logs the completion of the request, including the total processing time.

- **Why It’s Necessary**:
  - Provides visibility into incoming requests and their processing times.
  - Centralizes logging logic, making it reusable across all handlers.

---

#### **2. Wrapping the ServeMux**
```go
finalHandler := LoggingWrapper(mux)
```
- **What It Does**:
  - Wraps the entire `ServeMux` with the `LoggingWrapper`.
  - Ensures that all routes registered with the `ServeMux` are logged automatically.

- **Why It’s Necessary**:
  - Avoids the need to manually wrap individual handlers with logging logic.

---

#### **3. Handlers**
```go
func helloHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	fmt.Fprintln(w, "Hello, Shubh! Welcome to your Go HTTP server.")
}

func customHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "This is a custom route!")
}
```
- **What They Do**:
  - Handle requests for `/hello` and `/custom` routes.
  - Return appropriate responses based on the request.

- **Why They’re Necessary**:
  - Demonstrate how the logging wrapper works with multiple handlers.

---

#### **4. Graceful Shutdown**
The graceful shutdown logic ensures that active connections are allowed to complete before the server stops. This part of the code remains unchanged from previous examples.

---

### **Expected Behavior**

1. **Server Start**:
   - When the program starts, it logs:
     ```
     Server started on :8080
     Access the following URLs:
     - http://localhost:8080/hello
     - http://localhost:8080/custom
     Press Ctrl+C to stop the server...
     ```

2. **Handling Requests**:
   - Visiting `http://localhost:8080/hello` logs:
     ```
     Request received: GET /hello
     Request completed: GET /hello (1.234ms)
     ```
     And responds with:
     ```
     Hello, Shubh! Welcome to your Go HTTP server.
     ```

   - Visiting `http://localhost:8080/custom` logs:
     ```
     Request received: GET /custom
     Request completed: GET /custom (567µs)
     ```
     And responds with:
     ```
     This is a custom route!
     ```

3. **Graceful Shutdown**:
   - Pressing `Ctrl+C` logs:
     ```
     Shutting down server...
     Server stopped gracefully.
     ```

---

### **Advantages of This Approach**

1. **Reusability**:
   - The `LoggingWrapper` can be reused for any `http.Handler`, making it easy to apply consistent logging across your application.

2. **Centralized Logging**:
   - All logging logic is centralized in one place, reducing duplication and improving maintainability.

3. **Flexibility**:
   - You can extend the `LoggingWrapper` to include additional details, such as request headers, query parameters, or response status codes.

4. **Performance Monitoring**:
   - By logging the processing time, you can monitor the performance of your handlers and identify bottlenecks.

---

### **Extending the LoggingWrapper**

You can enhance the `LoggingWrapper` to log more details, such as:
- **Request Headers**:
  ```go
  log.Printf("Headers: %v", r.Header)
  ```

- **Query Parameters**:
  ```go
  log.Printf("Query Params: %v", r.URL.Query())
  ```

- **Response Status Code**:
  To log the response status code, you’ll need to use a custom `ResponseWriter` wrapper. Let me know if you'd like an example of this!

---

This implementation demonstrates how to create a reusable logging wrapper in Go. 

*/