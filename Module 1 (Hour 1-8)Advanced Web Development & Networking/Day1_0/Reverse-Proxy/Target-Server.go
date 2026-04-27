package main

import (
	"fmt"
	"net/http"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Hello from the target server!")
	})

	fmt.Println("Target server started on :9090")
	http.ListenAndServe(":9090", nil)
}

/*
Yes, there are **two separate code snippets** in the previous response:

1. **Reverse Proxy Server**:
   - This is the main reverse proxy implementation that listens on port `8080` and forwards incoming requests to a target server (e.g., `http://localhost:9090`).

2. **Target Server**:
   - This is a simple HTTP server that runs on port `9090` and acts as the backend server for the reverse proxy.

These two servers work together to demonstrate how a reverse proxy functions. Let me clarify their roles and how they interact.

---

### **1. Reverse Proxy Server**
This is the primary server that listens for incoming client requests and forwards them to the target server.

#### **Code for Reverse Proxy Server**
```go
package main

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// createReverseProxy creates a reverse proxy to forward requests to the target URL.
func createReverseProxy(targetURL string) (*httputil.ReverseProxy, error) {
	target, err := url.Parse(targetURL)
	if err != nil {
		return nil, err
	}

	// Create a reverse proxy
	proxy := httputil.NewSingleHostReverseProxy(target)

	// Modify the request before forwarding (optional)
	originalDirector := proxy.Director
	proxy.Director = func(req *http.Request) {
		originalDirector(req) // Call the original director
		req.Host = target.Host // Ensure the Host header matches the target
	}

	return proxy, nil
}

func main() {
	// Define the target server URL (e.g., "http://localhost:9090")
	targetURL := "http://localhost:9090"

	// Create the reverse proxy
	reverseProxy, err := createReverseProxy(targetURL)
	if err != nil {
		log.Fatalf("Failed to create reverse proxy: %s", err)
	}

	// Create a custom ServeMux
	mux := http.NewServeMux()

	// Route all requests to the reverse proxy
	mux.Handle("/", reverseProxy)

	// Create an HTTP server
	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	// Start the server in a goroutine
	go func() {
		log.Println("Reverse proxy server started on :8080")
		log.Printf("Forwarding requests to target server: %s\n", targetURL)

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

#### **What It Does**:
- Listens on port `8080`.
- Forwards all incoming requests to the target server running on `http://localhost:9090`.

---

### **2. Target Server**
This is the backend server that handles requests forwarded by the reverse proxy.

#### **Code for Target Server**
```go
package main

import (
	"fmt"
	"net/http"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Hello from the target server!")
	})

	fmt.Println("Target server started on :9090")
	http.ListenAndServe(":9090", nil)
}
```

#### **What It Does**:
- Listens on port `9090`.
- Responds with `"Hello from the target server!"` for all incoming requests.

---

### **How They Work Together**

1. **Start the Target Server**:
   - Run the target server on port `9090`:
     ```bash
     go run target-server.go
     ```
   - Logs:
     ```
     Target server started on :9090
     ```

2. **Start the Reverse Proxy**:
   - Run the reverse proxy on port `8080`:
     ```bash
     go run reverse-proxy.go
     ```
   - Logs:
     ```
     Reverse proxy server started on :8080
     Forwarding requests to target server: http://localhost:9090
     Press Ctrl+C to stop the server...
     ```

3. **Access the Reverse Proxy**:
   - Visit `http://localhost:8080` in your browser or use `curl`:
     ```bash
     curl http://localhost:8080
     ```
   - The reverse proxy forwards the request to the target server (`http://localhost:9090`), which responds with:
     ```
     Hello from the target server!
     ```

4. **Graceful Shutdown**:
   - Press `Ctrl+C` to stop the reverse proxy:
     ```
     Shutting down server...
     Server stopped gracefully.
     ```

---

### **Why Two Separate Servers?**

- **Separation of Concerns**:
  - The reverse proxy handles routing and forwarding, while the target server handles actual request processing.
  - This separation allows you to scale and manage each component independently.

- **Real-World Use Case**:
  - In production, the reverse proxy might forward requests to multiple backend servers (e.g., for load balancing).
  - The target server could be replaced with a more complex application (e.g., a REST API, a database-driven service, etc.).

---

### **Extending the Implementation**

1. **Load Balancing**:
   - Extend the reverse proxy to forward requests to multiple backend servers using a round-robin or weighted algorithm.

2. **SSL Termination**:
   - Configure the reverse proxy to handle HTTPS traffic using `ListenAndServeTLS`.

3. **Path-Based Routing**:
   - Route requests to different backend servers based on the URL path:
     ```go
     mux.Handle("/api/", reverseProxyForAPI)
     mux.Handle("/static/", reverseProxyForStaticFiles)
     ```

4. **Request Logging**:
   - Add logging middleware to log details about incoming requests and their responses.

---

*/