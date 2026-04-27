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

// Middleware 1: Logging Middleware
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Logging Middleware: Request received: %s %s", r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
	})
}

// Middleware 2: Authentication Middleware
func authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader != "Bearer secret-token" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		log.Println("Auth Middleware: User authenticated")
		next.ServeHTTP(w, r)
	})
}

// Middleware 3: Timing Middleware
func timingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		duration := time.Since(start)
		log.Printf("Timing Middleware: Request processed in %v", duration)
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

func main() {
	// Create a custom ServeMux
	mux := http.NewServeMux()

	// Register the "/hello" route with the handler
	mux.HandleFunc("/hello", helloHandler)

	// Chain middlewares manually
	finalHandler := loggingMiddleware(
		authMiddleware(
			timingMiddleware(mux),
		),
	)

	// Create an HTTP server
	server := &http.Server{
		Addr:    ":8080",
		Handler: finalHandler,
	}

	// Start the server in a goroutine
	go func() {
		log.Println("Server started on :8080")
		log.Println("Access the following URL:")
		log.Println("- http://localhost:8080/hello")

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
}// curl.exe -H "Authorization: Bearer secret-token" http://localhost:8080/hello

/*
Middleware chaining is a powerful concept in Go's `net/http` package that allows you to apply multiple layers of processing to incoming HTTP requests. Each middleware function can perform tasks like logging, authentication, request validation, etc., before passing the request to the next handler in the chain.

Below, I'll show you how to **manually chain middleware** without relying on external libraries or frameworks. This approach gives you full control over the middleware execution order and behavior.

---

### **Code Example: Middleware Chaining Manually**

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

// Middleware 1: Logging Middleware
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Logging Middleware: Request received: %s %s", r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
		ServeHTTP dispatches the request to the handler whose pattern most closely matches the request URL.
	})
}

// Middleware 2: Authentication Middleware
func authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader != "Bearer secret-token" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		log.Println("Auth Middleware: User authenticated")
		next.ServeHTTP(w, r)
	})
}

// Middleware 3: Timing Middleware
func timingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		duration := time.Since(start)
		log.Printf("Timing Middleware: Request processed in %v", duration)
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

func main() {
	// Create a custom ServeMux
	mux := http.NewServeMux()

	// Register the "/hello" route with the handler
	mux.HandleFunc("/hello", helloHandler)

	// Chain middlewares manually
	finalHandler := loggingMiddleware(
		authMiddleware(
			timingMiddleware(mux),
		),
	)

	// Create an HTTP server
	server := &http.Server{
		Addr:    ":8080",
		Handler: finalHandler,
	}

	// Start the server in a goroutine
	go func() {
		log.Println("Server started on :8080")
		log.Println("Access the following URL:")
		log.Println("- http://localhost:8080/hello")

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

### **Explanation of Middleware Chaining**

#### **1. Middleware Functions**
Each middleware function wraps the next handler in the chain. Here are the three middleware functions used in this example:

- **Logging Middleware**:
  ```go
  func loggingMiddleware(next http.Handler) http.Handler {
      return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
          log.Printf("Logging Middleware: Request received: %s %s", r.Method, r.URL.Path)
          next.ServeHTTP(w, r)
      })
  }
  ```
  - Logs details about the incoming request.
  - Calls `next.ServeHTTP` to pass the request to the next handler.

- **Authentication Middleware**:
  ```go
  func authMiddleware(next http.Handler) http.Handler {
      return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
          authHeader := r.Header.Get("Authorization")
          if authHeader != "Bearer secret-token" {
              http.Error(w, "Unauthorized", http.StatusUnauthorized)
              return
          }
          log.Println("Auth Middleware: User authenticated")
          next.ServeHTTP(w, r)
      })
  }
  ```
  - Checks for a valid `Authorization` header.
  - If authentication fails, it responds with a `401 Unauthorized` error.
  - Otherwise, it passes the request to the next handler.

- **Timing Middleware**:
  ```go
  func timingMiddleware(next http.Handler) http.Handler {
      return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
          start := time.Now()
          next.ServeHTTP(w, r)
          duration := time.Since(start)
          log.Printf("Timing Middleware: Request processed in %v", duration)
      })
  }
  ```
  - Measures the time taken to process the request.
  - Logs the duration after the request is handled.

---

#### **2. Middleware Chaining**
Middleware chaining is achieved by wrapping one middleware around another. The order of chaining determines the execution order of the middleware.

In this example:
```go
finalHandler := loggingMiddleware(
    authMiddleware(
        timingMiddleware(mux),
    ),
)
```
- The `timingMiddleware` is the innermost layer, so it executes first.
- The `authMiddleware` executes next, after the timing middleware.
- The `loggingMiddleware` executes last, logging the request details.

When a request is received:
1. The `loggingMiddleware` logs the request.
2. The `authMiddleware` checks for authentication.
3. The `timingMiddleware` measures the processing time.
4. Finally, the request is passed to the actual handler (`mux`).

---

#### **3. Expected Behavior**

1. **Valid Request**:
   - Send a request with the correct `Authorization` header:
     ```bash
     curl -H "Authorization: Bearer secret-token" http://localhost:8080/hello
     ```
   - Logs:
     ```
     Logging Middleware: Request received: GET /hello
     Auth Middleware: User authenticated
     Timing Middleware: Request processed in 1.234ms
     ```
   - Response:
     ```
     Hello, Shubh! Welcome to your Go HTTP server.
     ```

2. **Unauthorized Request**:
   - Send a request without the `Authorization` header:
     ```bash
     curl http://localhost:8080/hello
     ```
   - Logs:
     ```
     Logging Middleware: Request received: GET /hello
     ```
   - Response:
     ```
     Unauthorized
     ```

3. **Graceful Shutdown**:
   - Press `Ctrl+C` to stop the server:
     ```
     Shutting down server...
     Server stopped gracefully.
     ```

---

### **Why This Approach Works**

1. **Manual Control**:
   - By chaining middleware manually, you have full control over the order of execution.

2. **Modularity**:
   - Each middleware function is independent and reusable, making the code easier to maintain.

3. **Flexibility**:
   - You can add, remove, or reorder middleware as needed without modifying the core logic.

4. **Graceful Shutdown**:
   - The server shuts down gracefully, ensuring active connections are completed before termination.

---

This implementation demonstrates how to manually chain middleware in Go. 


*/