package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	_ "net/http/pprof" // Import pprof for profiling
)

// MockResponse represents the structure of the API response.
type MockResponse struct {
	Message string `json:"message"`
}

// slowAPIHandler simulates a slow API endpoint.
func slowAPIHandler(w http.ResponseWriter, r *http.Request) {
	time.Sleep(3 * time.Second) // Simulate a slow API call
	response := MockResponse{Message: "This is a slow API response."}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// optimizedAPIHandler simulates an optimized API endpoint.
func optimizedAPIHandler(w http.ResponseWriter, r *http.Request) {
	// Simulate optimization by reducing sleep time.
	time.Sleep(200 * time.Millisecond) // Faster response
	response := MockResponse{Message: "This is an optimized API response."}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func main() {
	// Start a separate HTTP server for profiling.
	go func() {
		log.Println("Starting profiling server on :6060")
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()

	// Define routes for the mock API.
	http.HandleFunc("/slow-api", slowAPIHandler)
	http.HandleFunc("/optimized-api", optimizedAPIHandler)

	// Start the main HTTP server.
	port := ":8080"
	fmt.Printf("Mock API server is running on http://localhost%s\n", port)
	fmt.Println("Access /slow-api for the slow API and /optimized-api for the optimized API.")
	log.Fatal(http.ListenAndServe(port, nil))
}


/*
Below is a Go code snippet that demonstrates how to mock an API, optimize it by simulating a faster response, and perform profiling analysis using Go's built-in `net/http/pprof` package for profiling.

```go
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	_ "net/http/pprof" // Import pprof for profiling
)

// MockResponse represents the structure of the API response.
type MockResponse struct {
	Message string `json:"message"`
}

// slowAPIHandler simulates a slow API endpoint.
func slowAPIHandler(w http.ResponseWriter, r *http.Request) {
	time.Sleep(3 * time.Second) // Simulate a slow API call
	response := MockResponse{Message: "This is a slow API response."}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// optimizedAPIHandler simulates an optimized API endpoint.
func optimizedAPIHandler(w http.ResponseWriter, r *http.Request) {
	// Simulate optimization by reducing sleep time.
	time.Sleep(200 * time.Millisecond) // Faster response
	response := MockResponse{Message: "This is an optimized API response."}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func main() {
	// Start a separate HTTP server for profiling.
	go func() {
		log.Println("Starting profiling server on :6060")
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()

	// Define routes for the mock API.
	http.HandleFunc("/slow-api", slowAPIHandler)
	http.HandleFunc("/optimized-api", optimizedAPIHandler)

	// Start the main HTTP server.
	port := ":8080"
	fmt.Printf("Mock API server is running on http://localhost%s\n", port)
	fmt.Println("Access /slow-api for the slow API and /optimized-api for the optimized API.")
	log.Fatal(http.ListenAndServe(port, nil))
}
```

### Explanation of the Code:

1. **Mocking the API**:
   - The `/slow-api` endpoint simulates a slow API by introducing a 3-second delay using `time.Sleep`.
   - The `/optimized-api` endpoint simulates an optimized version of the API with a much shorter delay (200ms).

2. **Profiling Analysis**:
   - The `_ "net/http/pprof"` import enables profiling capabilities in Go.
   - A separate HTTP server is started on `localhost:6060` to expose profiling data. You can access this data using tools like `go tool pprof` or visualize it using a browser at `http://localhost:6060/debug/pprof/`.

3. **Running the Server**:
   - The main HTTP server runs on `localhost:8080`, exposing the `/slow-api` and `/optimized-api` endpoints.

### Steps to Profile the Application:

1. Run the program:
   ```bash
   go run main.go
   ```

2. Access the endpoints:
   - Slow API: `http://localhost:8080/slow-api`
   - Optimized API: `http://localhost:8080/optimized-api`

3. Use the profiling server:
   - Open `http://localhost:6060/debug/pprof/` in your browser to view profiling data.
   - Use `go tool pprof` to analyze CPU or memory usage:
     ```bash
     go tool pprof http://localhost:6060/debug/pprof/profile?seconds=30
     ```

4. Optimize further based on profiling results.

This setup allows you to test, optimize, and profile your API effectively.

*/