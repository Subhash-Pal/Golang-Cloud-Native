package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

// LoggerMiddleware logs the incoming request details.
func LoggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		log.Printf("Started %s %s", r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
		log.Printf("Completed %s in %v", r.URL.Path, time.Since(start))
	})
}

// RecoveryMiddleware recovers from any panics in the application and returns a 500 error.
func RecoveryMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("Recovered from panic: %v", err)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			}
		}()
		next.ServeHTTP(w, r)
	})
}

// ExampleHandler is a simple handler to demonstrate middleware usage.
func ExampleHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/panic" {
		panic("Simulated panic!")
	}
	w.Write([]byte("Hello, World!"))
}

func main() {
	// Define the main handler
	mainHandler := http.HandlerFunc(ExampleHandler)

	// Wrap the handler with middleware
	handlerWithMiddleware := LoggerMiddleware(RecoveryMiddleware(mainHandler))

	// Start the HTTP server
	port := ":8080"
	fmt.Printf("Server is running on http://localhost%s\n", port)
	log.Fatal(http.ListenAndServe(port, handlerWithMiddleware))
}
/*Here's an example of Go (Golang) code that demonstrates how to implement middleware for logging and recovery in an HTTP server. Middleware is a common pattern in Go web applications to handle cross-cutting concerns like logging requests and recovering from panics.

```go
package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

// LoggerMiddleware logs the incoming request details.
func LoggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		log.Printf("Started %s %s", r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
		log.Printf("Completed %s in %v", r.URL.Path, time.Since(start))
	})
}

// RecoveryMiddleware recovers from any panics in the application and returns a 500 error.
func RecoveryMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("Recovered from panic: %v", err)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			}
		}()
		next.ServeHTTP(w, r)
	})
}

// ExampleHandler is a simple handler to demonstrate middleware usage.
func ExampleHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/panic" {
		panic("Simulated panic!")
	}
	w.Write([]byte("Hello, World!"))
}

func main() {
	// Define the main handler
	mainHandler := http.HandlerFunc(ExampleHandler)

	// Wrap the handler with middleware
	handlerWithMiddleware := LoggerMiddleware(RecoveryMiddleware(mainHandler))

	// Start the HTTP server
	port := ":8080"
	fmt.Printf("Server is running on http://localhost%s\n", port)
	log.Fatal(http.ListenAndServe(port, handlerWithMiddleware))
}
```

### Explanation:
1. **LoggerMiddleware**:
   - Logs the start and completion of each request.
   - Includes the HTTP method, URL path, and the time taken to process the request.

2. **RecoveryMiddleware**:
   - Recovers from any panics that occur during request handling.
   - Logs the panic message and responds with a `500 Internal Server Error`.

3. **ExampleHandler**:
   - A simple handler that writes "Hello, World!" to the response.
   - Simulates a panic when the `/panic` route is accessed.

4. **Middleware Chaining**:
   - Middleware functions are applied in the order they are wrapped around the handler.
   - In this case, `LoggerMiddleware` is applied first, followed by `RecoveryMiddleware`.

### How to Test:
- Run the program and visit `http://localhost:8080` in your browser or use `curl`.
- Access `http://localhost:8080/panic` to trigger a simulated panic and observe how the recovery middleware handles it.

This is a clean and modular way to handle logging and error recovery in Go web applications.
*/