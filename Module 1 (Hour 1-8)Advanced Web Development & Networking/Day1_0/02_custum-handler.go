
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

// Middleware for logging requests
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Request received: %s %s", r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
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

	// Wrap the mux with logging middleware
	handlerWithMiddleware := loggingMiddleware(mux)

	// Create an HTTP server
	server := &http.Server{
		Addr:    ":8080",
		Handler: handlerWithMiddleware,
	}

	// Start the server in a goroutine
	go func() {
		log.Println("Server started on :8080")
		log.Println("Access the following URLs:")
		log.Println("- http://localhost:8080/hello")
		log.Println("- http://localhost:8080/custom")

		// ListenAndServe blocks until the server is shut down
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Error starting server: %s\n", err)
		}
	}()

	// Wait for interrupt signal to gracefully shut down
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM) // Listen for SIGINT (Ctrl+C) and SIGTERM

	// Block until a signal is received
	log.Println("Press Ctrl+C to stop the server...")
	<-sig

	log.Println("Shutting down server...")

	// Create a context with a timeout for graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second) // Increased timeout
	defer cancel()

	// Attempt to shut down the server gracefully
	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Error during shutdown: %s\n", err)
	}

	log.Println("Server stopped gracefully.")
}
/*

Certainly! Below is a detailed **description of the code**, explaining **what it does** and **why each part is necessary**. This will help you understand the purpose and functionality of the program.

---

### **Code Description**

#### **1. Purpose of the Program**
The program creates an HTTP server in Go that:
- Listens on port `8080`.
- Handles incoming HTTP requests for two routes: `/hello` and `/custom`.
- Logs details about incoming requests using middleware.
- Gracefully shuts down when interrupted (e.g., via `Ctrl+C`).

This implementation is designed to be robust, ensuring that active connections are allowed to complete before the server stops.

---

#### **2. Key Components of the Code**

##### **a. Middleware (`loggingMiddleware`)**
```go
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Request received: %s %s", r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
	})
}
```
- **What It Does**:
  - Logs details about each incoming request, including the HTTP method (`GET`, `POST`, etc.) and the URL path.
  - Passes the request to the next handler in the chain using `next.ServeHTTP`.

- **Why It’s Necessary**:
  - Provides visibility into incoming requests, which is essential for debugging and monitoring.

---

##### **b. Handlers (`helloHandler` and `customHandler`)**
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
  - `helloHandler`: Responds with a greeting message if the request method is `GET`. Returns a `405 Method Not Allowed` error for other methods.
  - `customHandler`: Responds with a simple message indicating that this is a custom route.

- **Why They’re Necessary**:
  - Demonstrate how to handle specific routes and respond to client requests.

---

##### **c. Custom ServeMux**
```go
mux := http.NewServeMux()
mux.HandleFunc("/hello", helloHandler)
mux.HandleFunc("/custom", customHandler)
```
- **What It Does**:
  - Creates a custom `ServeMux` (multiplexer) to map URL paths (`/hello`, `/custom`) to their respective handlers.

- **Why It’s Necessary**:
  - Provides more control over routing compared to the default `http.DefaultServeMux`.

---

##### **d. Middleware Wrapping**
```go
handlerWithMiddleware := loggingMiddleware(mux)
```
- **What It Does**:
  - Wraps the `ServeMux` with the `loggingMiddleware` to ensure all incoming requests pass through the middleware.

- **Why It’s Necessary**:
  - Ensures consistent logging for all requests, regardless of the route.

---

##### **e. HTTP Server Configuration**
```go
server := &http.Server{
	Addr:    ":8080",
	Handler: handlerWithMiddleware,
}
```
- **What It Does**:
  - Configures the HTTP server to listen on `:8080` and use the middleware-wrapped `ServeMux` as the request handler.

- **Why It’s Necessary**:
  - Centralizes server configuration, making it easier to manage and extend.

---

##### **f. Starting the Server**
```go
go func() {
	log.Println("Server started on :8080")
	log.Println("Access the following URLs:")
	log.Println("- http://localhost:8080/hello")
	log.Println("- http://localhost:8080/custom")

	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Error starting server: %s\n", err)
	}
}()
```
- **What It Does**:
  - Starts the server in a separate goroutine using `ListenAndServe`.
  - Logs messages indicating that the server has started and provides URLs for accessing the handlers.

- **Why It’s Necessary**:
  - Runs the server asynchronously, allowing the main goroutine to handle shutdown signals.
  - Provides clear instructions for testing the server.

---

##### **g. Signal Handling**
```go
sig := make(chan os.Signal, 1)
signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

log.Println("Press Ctrl+C to stop the server...")
<-sig
```
- **What It Does**:
  - Sets up a channel (`sig`) to capture interrupt signals (`SIGINT` for `Ctrl+C` and `SIGTERM` for termination).
  - Blocks the main goroutine until a signal is received.

- **Why It’s Necessary**:
  - Ensures the program can gracefully shut down when interrupted.

---

##### **h. Graceful Shutdown**
```go
log.Println("Shutting down server...")

ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
defer cancel()

if err := server.Shutdown(ctx); err != nil {
	log.Fatalf("Error during shutdown: %s\n", err)
}

log.Println("Server stopped gracefully.")
```
- **What It Does**:
  - Initiates a graceful shutdown of the server using `server.Shutdown`.
  - Allows up to 10 seconds for active connections to complete before forcefully terminating them.

- **Why It’s Necessary**:
  - Prevents abrupt termination of active requests, ensuring a smooth user experience.
  - Logs confirmation that the server has stopped gracefully.

---

#### **3. Why This Implementation Works**

1. **Concurrency**:
   - The server runs in its own goroutine, allowing the main goroutine to handle shutdown signals independently.

2. **Graceful Shutdown**:
   - The `server.Shutdown` method ensures that:
     - No new requests are accepted after shutdown begins.
     - Active requests are allowed to complete within the specified timeout.

3. **Signal Handling**:
   - Capturing `SIGINT` and `SIGTERM` ensures the program responds to both manual interruptions (`Ctrl+C`) and system-level termination signals.

4. **Middleware**:
   - Logging middleware provides visibility into incoming requests without modifying individual handlers.

5. **Timeout Management**:
   - A 10-second timeout ensures sufficient time for active connections to complete while preventing indefinite blocking.

---

#### **4. Expected Behavior**

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
     ```
     And responds with:
     ```
     Hello, Shubh! Welcome to your Go HTTP server.
     ```

   - Visiting `http://localhost:8080/custom` logs:
     ```
     Request received: GET /custom
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

#### **5. Common Issues Addressed**

1. **Blocking Handlers**:
   - Ensured no handler blocks indefinitely, allowing the server to shut down properly.

2. **Signal Handling**:
   - Verified that `signal.Notify` captures `SIGINT` and `SIGTERM` correctly.

3. **Timeout Too Short**:
   - Increased the timeout to 10 seconds to allow sufficient time for active connections to complete.

4. **Improper Goroutine Management**:
   - Ensured the server runs in a separate goroutine, allowing the main goroutine to handle shutdown signals.

---

This description explains the purpose and functionality of each part of the code, ensuring you understand why the program behaves as expected. 


*/