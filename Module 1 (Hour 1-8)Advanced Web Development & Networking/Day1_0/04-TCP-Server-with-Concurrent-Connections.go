package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

// Client represents a connected client.
type Client struct {
	conn net.Conn
}

// handleConnection handles communication with a single client.
func handleConnection(client *Client, wg *sync.WaitGroup) {
	defer wg.Done() // Notify the WaitGroup when the goroutine finishes
	defer client.conn.Close()

	log.Printf("Client connected: %s", client.conn.RemoteAddr())

	// Create a buffered reader for reading data from the client
	reader := bufio.NewReader(client.conn)

	for {
		// Read data from the client
		message, err := reader.ReadString('\n')
		if err != nil {
			log.Printf("Client disconnected: %s", client.conn.RemoteAddr())
			return
		}

		// Log and process the received message
		log.Printf("Received from %s: %s", client.conn.RemoteAddr(), message)

		// Echo the message back to the client
		response := fmt.Sprintf("Echo: %s", message)
		_, err = client.conn.Write([]byte(response))
		if err != nil {
			log.Printf("Error writing to client %s: %v", client.conn.RemoteAddr(), err)
			return
		}
	}
}

func main() {
	// Define the address and port to listen on
	address := ":8081"

	// Start listening for incoming TCP connections
	listener, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatalf("Failed to start TCP server: %v", err)
	}
	defer listener.Close()

	log.Printf("TCP server started on %s", address)

	// Wait for interrupt signal to gracefully shut down
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

	// Use a WaitGroup to track active connections
	var wg sync.WaitGroup

	// Start accepting connections in a loop
	go func() {
		for {
			// Accept a new connection
			conn, err := listener.Accept()
			if err != nil {
				select {
				case <-sig: // Check if shutdown signal was received
					log.Println("Stopping listener...")
					return
				default:
					log.Printf("Error accepting connection: %v", err)
					continue
				}
			}

			// Increment the WaitGroup counter for the new connection
			wg.Add(1)

			// Handle the connection in a separate goroutine
			client := &Client{conn: conn}
			go handleConnection(client, &wg)
		}
	}()

	log.Println("Press Ctrl+C to stop the server...")
	<-sig

	log.Println("Shutting down server...")

	// Stop accepting new connections
	err = listener.Close()
	if err != nil {
		log.Printf("Error closing listener: %v", err)
	}

	// Wait for all active connections to finish
	wg.Wait()

	log.Println("Server stopped gracefully.")
}

/*
Handling **concurrent connections** is one of the key strengths of Go's concurrency model. In a TCP server, each client connection can be handled in its own goroutine, allowing the server to manage multiple clients simultaneously without blocking.

Below, I'll explain how to handle concurrent connections in a TCP server and provide an enhanced version of the previous TCP server example to demonstrate this concept.

---

### **Key Concepts for Concurrent Connections**

1. **Goroutines**:
   - Each client connection is handled in a separate goroutine.
   - Goroutines are lightweight threads managed by Go's runtime, making it efficient to handle thousands of concurrent connections.

2. **Non-Blocking Accept**:
   - The `listener.Accept()` method is non-blocking when used with goroutines, allowing the server to accept new connections while processing existing ones.

3. **Connection Pooling**:
   - You can maintain a pool of active connections if you need to broadcast messages or manage clients globally.

4. **Graceful Shutdown**:
   - When shutting down the server, ensure all active connections are closed gracefully.

---

### **Enhanced Code Example: TCP Server with Concurrent Connections**

```go
package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

// Client represents a connected client.
type Client struct {
	conn net.Conn
}

// handleConnection handles communication with a single client.
func handleConnection(client *Client, wg *sync.WaitGroup) {
	defer wg.Done() // Notify the WaitGroup when the goroutine finishes
	defer client.conn.Close()

	log.Printf("Client connected: %s", client.conn.RemoteAddr())

	// Create a buffered reader for reading data from the client
	reader := bufio.NewReader(client.conn)

	for {
		// Read data from the client
		message, err := reader.ReadString('\n')
		if err != nil {
			log.Printf("Client disconnected: %s", client.conn.RemoteAddr())
			return
		}

		// Log and process the received message
		log.Printf("Received from %s: %s", client.conn.RemoteAddr(), message)

		// Echo the message back to the client
		response := fmt.Sprintf("Echo: %s", message)
		_, err = client.conn.Write([]byte(response))
		if err != nil {
			log.Printf("Error writing to client %s: %v", client.conn.RemoteAddr(), err)
			return
		}
	}
}

func main() {
	// Define the address and port to listen on
	address := ":8081"

	// Start listening for incoming TCP connections
	listener, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatalf("Failed to start TCP server: %v", err)
	}
	defer listener.Close()

	log.Printf("TCP server started on %s", address)

	// Wait for interrupt signal to gracefully shut down
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

	// Use a WaitGroup to track active connections
	var wg sync.WaitGroup

	// Start accepting connections in a loop
	go func() {
		for {
			// Accept a new connection
			conn, err := listener.Accept()
			if err != nil {
				select {
				case <-sig: // Check if shutdown signal was received
					log.Println("Stopping listener...")
					return
				default:
					log.Printf("Error accepting connection: %v", err)
					continue
				}
			}

			// Increment the WaitGroup counter for the new connection
			wg.Add(1)

			// Handle the connection in a separate goroutine
			client := &Client{conn: conn}
			go handleConnection(client, &wg)
		}
	}()

	log.Println("Press Ctrl+C to stop the server...")
	<-sig

	log.Println("Shutting down server...")

	// Stop accepting new connections
	err = listener.Close()
	if err != nil {
		log.Printf("Error closing listener: %v", err)
	}

	// Wait for all active connections to finish
	wg.Wait()

	log.Println("Server stopped gracefully.")
}
```

---

### **Explanation of the Code**

#### **1. Handling Concurrent Connections**
```go
go handleConnection(client, &wg)
```
- **What It Does**:
  - Spawns a new goroutine to handle each client connection.
  - Uses a `sync.WaitGroup` to track active connections and ensure they complete before shutting down.

- **Why It’s Necessary**:
  - Allows the server to handle multiple clients concurrently without blocking.

---

#### **2. Graceful Shutdown**
```go
<-sig

log.Println("Shutting down server...")

err = listener.Close()
if err != nil {
	log.Printf("Error closing listener: %v", err)
}

wg.Wait()
```
- **What It Does**:
  - Waits for a shutdown signal (`Ctrl+C`) to stop accepting new connections.
  - Closes the listener to prevent new connections.
  - Waits for all active connections to finish using the `WaitGroup`.

- **Why It’s Necessary**:
  - Ensures the server shuts down gracefully without abruptly terminating active connections.

---

#### **3. Connection Tracking**
```go
var wg sync.WaitGroup

wg.Add(1)
go handleConnection(client, &wg)

wg.Wait()
```
- **What It Does**:
  - Uses a `sync.WaitGroup` to track active connections.
  - Increments the counter for each new connection and decrements it when the connection is closed.

- **Why It’s Necessary**:
  - Ensures the server waits for all active connections to finish before exiting.

---

### **Expected Behavior**

1. **Start the TCP Server**:
   - Run the server:
     ```bash
     go run tcp-server.go
     ```
   - Logs:
     ```
     TCP server started on :8081
     Press Ctrl+C to stop the server...
     ```

2. **Connect Multiple Clients**:
   - Use multiple terminal windows to connect to the server (e.g., using `telnet` or `nc`):
     ```bash
     telnet localhost 8081
     ```
   - Or use `nc` (Netcat):
     ```bash
     nc localhost 8081
     ```

3. **Send Messages**:
   - Type a message and press Enter in each client:
     ```
     Hello from Client 1!
     ```
   - The server responds with:
     ```
     Echo: Hello from Client 1!
     ```

4. **Graceful Shutdown**:
   - Press `Ctrl+C` to stop the server:
     ```
     Shutting down server...
     Server stopped gracefully.
     ```

---

### **Advantages of This Approach**

1. **Concurrency**:
   - Each client connection is handled in its own goroutine, enabling high concurrency.

2. **Graceful Shutdown**:
   - Ensures all active connections are completed before the server exits.

3. **Scalability**:
   - The server can handle thousands of concurrent connections efficiently.

4. **Extensibility**:
   - You can extend the server to implement features like broadcasting, authentication, or timeouts.

---

### **Extending the Implementation**

Here are some ways to enhance the TCP server:

1. **Broadcasting**:
   - Maintain a global list of active clients and broadcast messages to all connected clients.

2. **Timeouts**:
   - Set read/write timeouts to prevent connections from hanging indefinitely:
     ```go
     conn.SetReadDeadline(time.Now().Add(10 * time.Second))
     ```

3. **Authentication**:
   - Add authentication logic to verify client credentials before processing requests.

4. **TLS Encryption**:
   - Use `tls.Listen` to enable secure communication over TLS.

5. **Message Framing**:
   - Implement a more robust framing mechanism (e.g., length-prefixed messages) instead of relying on `\n` as a delimiter.

---

This implementation demonstrates how to handle concurrent connections in a TCP server using Go. 
*/