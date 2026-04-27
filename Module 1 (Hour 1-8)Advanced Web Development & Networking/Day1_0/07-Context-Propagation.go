package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"
	"os"
	"os/signal"
	"syscall"
)

// Middleware to inject request-scoped data into the context
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Generate a unique request ID
		requestID := r.Header.Get("X-Request-ID")
		if requestID == "" {
			requestID = "generated-" + fmt.Sprintf("%d", time.Now().UnixNano())
		}

		// Create a new context with the request ID
		ctx := context.WithValue(r.Context(), "requestID", requestID)

		// Log the start of the request
		log.Printf("Request received: %s %s (Request ID: %s)", r.Method, r.URL.Path, requestID)

		// Pass the updated context to the next handler
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// Handler that uses the propagated context
func helloHandler(w http.ResponseWriter, r *http.Request) {
	// Retrieve the request ID from the context
	ctx := r.Context()
	requestID, ok := ctx.Value("requestID").(string)
	if !ok {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Simulate some work with a delay
	select {
	case <-time.After(2 * time.Second): // Simulate processing
		fmt.Fprintf(w, "Hello! Your Request ID is: %s\n", requestID)
	case <-ctx.Done(): // Handle cancellation
		log.Printf("Request cancelled (Request ID: %s)", requestID)
		http.Error(w, "Request cancelled", http.StatusRequestTimeout)
	}
}

func main() {
	// Create a custom ServeMux
	mux := http.NewServeMux()

	// Register the "/hello" route
	mux.HandleFunc("/hello", helloHandler)

	// Wrap the mux with the logging middleware
	handlerWithMiddleware := loggingMiddleware(mux)

	// Create an HTTP server
	server := &http.Server{
		Addr:    ":8080",
		Handler: handlerWithMiddleware,
	}

	// Start the server in a goroutine
	go func() {
		log.Println("Server started on :8080")
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
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Attempt to shut down the server gracefully
	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Error during shutdown: %s\n", err)
	}

	log.Println("Server stopped gracefully.")
}

/*
Context propagation in Go is a powerful mechanism for managing request-scoped data, cancellation signals, and deadlines across goroutines. 
It is particularly useful in intermediate-level applications where you need to pass contextual information (e.g., request IDs, user authentication tokens, or timeouts)
between different layers of your application.

Below is an **intermediate-level example** of context propagation in Go, demonstrating how to use `context.Context` to propagate request-scoped data and handle cancellations.

---

### **Code Example: Context Propagation**

```go
package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"
)

// Middleware to inject request-scoped data into the context
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Generate a unique request ID
		requestID := r.Header.Get("X-Request-ID")
		if requestID == "" {
			requestID = "generated-" + fmt.Sprintf("%d", time.Now().UnixNano())
		}

		// Create a new context with the request ID
		ctx := context.WithValue(r.Context(), "requestID", requestID)

		// Log the start of the request
		log.Printf("Request received: %s %s (Request ID: %s)", r.Method, r.URL.Path, requestID)

		// Pass the updated context to the next handler
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// Handler that uses the propagated context
func helloHandler(w http.ResponseWriter, r *http.Request) {
	// Retrieve the request ID from the context
	ctx := r.Context()
	requestID, ok := ctx.Value("requestID").(string)
	if !ok {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Simulate some work with a delay
	select {
	case <-time.After(2 * time.Second): // Simulate processing
		fmt.Fprintf(w, "Hello! Your Request ID is: %s\n", requestID)
	case <-ctx.Done(): // Handle cancellation
		log.Printf("Request cancelled (Request ID: %s)", requestID)
		http.Error(w, "Request cancelled", http.StatusRequestTimeout)
	}
}

func main() {
	// Create a custom ServeMux
	mux := http.NewServeMux()

	// Register the "/hello" route
	mux.HandleFunc("/hello", helloHandler)

	// Wrap the mux with the logging middleware
	handlerWithMiddleware := loggingMiddleware(mux)

	// Create an HTTP server
	server := &http.Server{
		Addr:    ":8080",
		Handler: handlerWithMiddleware,
	}

	// Start the server in a goroutine
	go func() {
		log.Println("Server started on :8080")
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
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
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

#### **1. Context Creation**
```go
ctx := context.WithValue(r.Context(), "requestID", requestID)
```
- **What It Does**:
  - Creates a new context with a key-value pair (`"requestID": requestID`) using `context.WithValue`.
  - The original context (`r.Context()`) is derived from the HTTP request.

- **Why It’s Necessary**:
  - Allows request-scoped data (e.g., request IDs) to be propagated across middleware and handlers.

---

#### **2. Context Propagation**
```go
next.ServeHTTP(w, r.WithContext(ctx))
```
- **What It Does**:
  - Updates the HTTP request with the new context containing the request ID.
  - Passes the updated request to the next handler in the chain.

- **Why It’s Necessary**:
  - Ensures that downstream handlers have access to the propagated context.

---

#### **3. Context Usage**
```go
requestID, ok := ctx.Value("requestID").(string)
if !ok {
	http.Error(w, "Internal server error", http.StatusInternalServerError)
	return
}
```
- **What It Does**:
  - Retrieves the request ID from the context using `ctx.Value`.
  - Handles cases where the value is missing or of the wrong type.

- **Why It’s Necessary**:
  - Demonstrates how to safely access context values in handlers.

---

#### **4. Cancellation Handling**
```go
select {
case <-time.After(2 * time.Second): // Simulate processing
	fmt.Fprintf(w, "Hello! Your Request ID is: %s\n", requestID)
case <-ctx.Done(): // Handle cancellation
	log.Printf("Request cancelled (Request ID: %s)", requestID)
	http.Error(w, "Request cancelled", http.StatusRequestTimeout)
}
```
- **What It Does**:
  - Simulates a long-running operation with a 2-second delay.
  - Checks if the context has been canceled using `ctx.Done()` and handles it appropriately.

- **Why It’s Necessary**:
  - Demonstrates how to handle cancellation signals (e.g., when a client disconnects or a timeout occurs).

---

#### **5. Graceful Shutdown**
```go
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()

if err := server.Shutdown(ctx); err != nil {
	log.Fatalf("Error during shutdown: %s\n", err)
}
```
- **What It Does**:
  - Shuts down the server gracefully with a 5-second timeout.
  - Cancels any active requests after the timeout.

- **Why It’s Necessary**:
  - Ensures the server stops cleanly without abruptly terminating active connections.

---

### **Expected Behavior**

1. **Start the Server**:
   - Run the server:
     ```bash
     go run context-propagation.go
     ```
   - Logs:
     ```
     Server started on :8080
     Press Ctrl+C to stop the server...
     ```

2. **Send Requests**:
   - Use `curl` to send a request:
     ```bash
     curl -H "X-Request-ID: custom-id" http://localhost:8080/hello
     ```
   - Response:
     ```
     Hello! Your Request ID is: custom-id
     ```

   - If no `X-Request-ID` header is provided:
     ```bash
     curl http://localhost:8080/hello
     ```
   - Response:
     ```
     Hello! Your Request ID is: generated-<timestamp>
     ```

3. **Simulate Cancellation**:
   - Send a request and cancel it before the 2-second delay completes:
     ```bash
     curl http://localhost:8080/hello
     ```
   - Stop the request manually (e.g., by pressing `Ctrl+C` in the terminal).
   - Logs:
     ```
     Request cancelled (Request ID: <id>)
     ```

4. **Graceful Shutdown**:
   - Press `Ctrl+C` to stop the server:
     ```
     Shutting down server...
     Server stopped gracefully.
     ```

---

### **Advantages of This Approach**

1. **Request-Scoped Data**:
   - Context allows you to propagate request-specific data (e.g., request IDs, user tokens) across middleware and handlers.

2. **Cancellation Support**:
   - Context provides built-in support for cancellation, making it easy to handle scenarios like client disconnections or timeouts.

3. **Timeout Management**:
   - You can set deadlines or timeouts for operations using `context.WithTimeout`.

4. **Scalability**:
   - Context is lightweight and designed for concurrent use, making it suitable for high-performance applications.

---

### **Extending the Implementation**

Here are some ways to enhance the context propagation example:

1. **Authentication**:
   - Use context to propagate user authentication tokens and validate them in middleware.

2. **Logging Enhancements**:
   - Include the request ID in all logs to improve traceability.

3. **Custom Keys**:
   - Define custom types for context keys to avoid collisions:
     ```go
     type contextKey string
     const requestIDKey contextKey = "requestID"
     ```

4. **Distributed Tracing**:
   - Integrate with distributed tracing systems (e.g., OpenTelemetry) to track requests across services.

5. **Background Tasks**:
   - Use context to manage background tasks and ensure they respect cancellation signals.

---

This implementation demonstrates how to use context propagation effectively in Go.

*/