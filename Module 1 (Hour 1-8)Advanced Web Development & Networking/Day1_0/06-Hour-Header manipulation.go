package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"
)

// Middleware to validate and manipulate headers
func headerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Log incoming headers
		log.Printf("Incoming headers: %v", r.Header)

		// Validate the "Authorization" header
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			http.Error(w, "Unauthorized: Missing or invalid Authorization header", http.StatusUnauthorized)
			return
		}

		// Add or modify a custom header
		r.Header.Set("X-Request-ID", "generated-id-12345")

		// Pass the updated request to the next handler
		next.ServeHTTP(w, r)
	})
}

// Handler to demonstrate header manipulation in responses
func headerHandler(w http.ResponseWriter, r *http.Request) {
	// Retrieve a header from the request
	requestID := r.Header.Get("X-Request-ID")
	if requestID == "" {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Set headers in the response
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("X-Response-ID", requestID)

	// Remove a header (optional)
	w.Header().Del("X-Unnecessary-Header")

	// Write the response body
	response := fmt.Sprintf(`{"message": "Hello!", "requestID": "%s"}`, requestID)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(response))
}

func main() {
	// Create a custom ServeMux
	mux := http.NewServeMux()

	// Register the "/header" route
	mux.HandleFunc("/header", headerHandler)

	// Wrap the mux with the header middleware
	handlerWithMiddleware := headerMiddleware(mux)

	// Create an HTTP server
	server := &http.Server{
		Addr:    ":8080",
		Handler: handlerWithMiddleware,
	}

	// Start the server
	log.Println("Server started on :8080")
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Error starting server: %s\n", err)
	}
}

/*
Header manipulation in Go is a common task when working with HTTP servers and clients. 
Headers are key-value pairs that provide metadata about the HTTP request or response. In Go, you can manipulate headers using the `http.Request` and `http.ResponseWriter` objects.

Below is an example of **header manipulation** in Go, demonstrating how to:

1. Add, modify, or remove headers in an HTTP request.
2. Set headers in an HTTP response.
3. Validate headers in middleware.

---

### **Code Example: Header Manipulation**

```go
package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"
)

// Middleware to validate and manipulate headers
func headerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Log incoming headers
		log.Printf("Incoming headers: %v", r.Header)

		// Validate the "Authorization" header
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			http.Error(w, "Unauthorized: Missing or invalid Authorization header", http.StatusUnauthorized)
			return
		}

		// Add or modify a custom header
		r.Header.Set("X-Request-ID", "generated-id-12345")

		// Pass the updated request to the next handler
		next.ServeHTTP(w, r)
	})
}

// Handler to demonstrate header manipulation in responses
func headerHandler(w http.ResponseWriter, r *http.Request) {
	// Retrieve a header from the request
	requestID := r.Header.Get("X-Request-ID")
	if requestID == "" {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Set headers in the response
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("X-Response-ID", requestID)

	// Remove a header (optional)
	w.Header().Del("X-Unnecessary-Header")

	// Write the response body
	response := fmt.Sprintf(`{"message": "Hello!", "requestID": "%s"}`, requestID)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(response))
}

func main() {
	// Create a custom ServeMux
	mux := http.NewServeMux()

	// Register the "/header" route
	mux.HandleFunc("/header", headerHandler)

	// Wrap the mux with the header middleware
	handlerWithMiddleware := headerMiddleware(mux)

	// Create an HTTP server
	server := &http.Server{
		Addr:    ":8080",
		Handler: handlerWithMiddleware,
	}

	// Start the server
	log.Println("Server started on :8080")
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Error starting server: %s\n", err)
	}
}
```

---

### **Explanation of the Code**

#### **1. Validating and Manipulating Request Headers**
```go
authHeader := r.Header.Get("Authorization")
if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
	http.Error(w, "Unauthorized: Missing or invalid Authorization header", http.StatusUnauthorized)
	return
}

r.Header.Set("X-Request-ID", "generated-id-12345")
```
- **What It Does**:
  - Retrieves the `Authorization` header using `r.Header.Get`.
  - Validates the header to ensure it starts with `"Bearer "`.
  - Adds or modifies a custom header (`X-Request-ID`) using `r.Header.Set`.

- **Why It’s Necessary**:
  - Ensures the request contains required headers and adds additional metadata for downstream processing.

---

#### **2. Setting Response Headers**
```go
w.Header().Set("Content-Type", "application/json")
w.Header().Set("X-Response-ID", requestID)

w.Header().Del("X-Unnecessary-Header")
```
- **What It Does**:
  - Sets headers in the response using `w.Header().Set`.
  - Removes an unnecessary header using `w.Header().Del`.

- **Why It’s Necessary**:
  - Allows you to control the metadata sent back to the client.

---

#### **3. Writing the Response Body**
```go
response := fmt.Sprintf(`{"message": "Hello!", "requestID": "%s"}`, requestID)
w.WriteHeader(http.StatusOK)
w.Write([]byte(response))
```
- **What It Does**:
  - Constructs a JSON response containing the `requestID`.
  - Writes the response body using `w.Write`.

- **Why It’s Necessary**:
  - Demonstrates how to send structured data in the response.

---

### **Expected Behavior**

1. **Start the Server**:
   - Run the server:
     ```bash
     go run header-manipulation.go
     ```
   - Logs:
     ```
     Server started on :8080
     ```

2. **Send Requests**:
   - Use `curl` to send a request with valid headers:
     ```bash
     curl -H "Authorization: Bearer secret-token" http://localhost:8080/header
     ```
   - Response:
     ```json
     {"message": "Hello!", "requestID": "generated-id-12345"}
     ```

   - Send a request without the `Authorization` header:
     ```bash
     curl http://localhost:8080/header
     ```
   - Response:
     ```
     Unauthorized: Missing or invalid Authorization header
     ```

3. **Inspect Headers**:
   - Use tools like `curl -v` to inspect the response headers:
     ```bash
     curl -v -H "Authorization: Bearer secret-token" http://localhost:8080/header
     ```
   - Observe the `Content-Type`, `X-Response-ID`, and other headers in the response.

---

### **Advantages of This Approach**

1. **Header Validation**:
   - Ensures that requests contain required headers before processing.

2. **Custom Metadata**:
   - Allows you to add or modify headers to propagate metadata (e.g., request IDs).

3. **Security**:
   - Protects against unauthorized access by validating sensitive headers (e.g., `Authorization`).

4. **Flexibility**:
   - Provides fine-grained control over both request and response headers.

---

### **Extending the Implementation**

Here are some ways to enhance header manipulation:

1. **CORS Support**:
   - Add Cross-Origin Resource Sharing (CORS) headers to allow requests from specific origins:
     ```go
     w.Header().Set("Access-Control-Allow-Origin", "*")
     ```

2. **Rate Limiting**:
   - Use headers to implement rate limiting (e.g., `X-RateLimit-Limit`, `X-RateLimit-Remaining`).

3. **Logging Enhancements**:
   - Log all incoming and outgoing headers for debugging purposes.

4. **Custom Middleware**:
   - Create reusable middleware for common header manipulations (e.g., adding trace IDs).

5. **Security Headers**:
   - Add security-related headers such as `Strict-Transport-Security`, `Content-Security-Policy`, or `X-Frame-Options`.

---

This implementation demonstrates how to manipulate headers in Go for both requests and responses. 


*/