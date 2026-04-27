package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	//"time"
)

// handleConnection handles communication with a single client.
func handleConnection(conn net.Conn) {
	defer conn.Close() // Ensure the connection is closed when done

	log.Printf("Client connected: %s", conn.RemoteAddr())

	// Create a buffered reader for reading data from the client
	reader := bufio.NewReader(conn)

	for {
		// Read data from the client
		message, err := reader.ReadString('\n')
		if err != nil {
			log.Printf("Client disconnected: %s", conn.RemoteAddr())
			return
		}

		// Log and process the received message
		log.Printf("Received from %s: %s", conn.RemoteAddr(), message)

		// Echo the message back to the client
		response := fmt.Sprintf("Echo: %s", message)
		_, err = conn.Write([]byte(response))
		if err != nil {
			log.Printf("Error writing to client %s: %v", conn.RemoteAddr(), err)
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

	// Start accepting connections in a goroutine
	go func() {
		for {
			// Accept a new connection
			conn, err := listener.Accept()
			if err != nil {
				log.Printf("Error accepting connection: %v", err)
				continue
			}

			// Handle the connection in a separate goroutine
			go handleConnection(conn)
		}
	}()

	log.Println("Press Ctrl+C to stop the server...")
	<-sig

	log.Println("Shutting down server...")

	// Gracefully close the listener
	err = listener.Close()
	if err != nil {
		log.Printf("Error closing listener: %v", err)
	}

	log.Println("Server stopped gracefully.")
}
/*

A **TCP server** in Go allows you to handle raw TCP connections. Unlike HTTP servers, which are built on top of the `net/http` package, a TCP server operates at a lower level using the `net` package. This makes it suitable for custom protocols or scenarios where you need fine-grained control over communication.

Below is an example of a **TCP server implementation** in Go:

---

### **Code Example: TCP Server Implementation**

```go
package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// handleConnection handles communication with a single client.
func handleConnection(conn net.Conn) {
	defer conn.Close() // Ensure the connection is closed when done

	log.Printf("Client connected: %s", conn.RemoteAddr())

	// Create a buffered reader for reading data from the client
	reader := bufio.NewReader(conn)

	for {
		// Read data from the client
		message, err := reader.ReadString('\n')
		if err != nil {
			log.Printf("Client disconnected: %s", conn.RemoteAddr())
			return
		}

		// Log and process the received message
		log.Printf("Received from %s: %s", conn.RemoteAddr(), message)

		// Echo the message back to the client
		response := fmt.Sprintf("Echo: %s", message)
		_, err = conn.Write([]byte(response))
		if err != nil {
			log.Printf("Error writing to client %s: %v", conn.RemoteAddr(), err)
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

	// Start accepting connections in a goroutine
	go func() {
		for {
			// Accept a new connection
			conn, err := listener.Accept()
			if err != nil {
				log.Printf("Error accepting connection: %v", err)
				continue
			}

			// Handle the connection in a separate goroutine
			go handleConnection(conn)
		}
	}()

	log.Println("Press Ctrl+C to stop the server...")
	<-sig

	log.Println("Shutting down server...")

	// Gracefully close the listener
	err = listener.Close()
	if err != nil {
		log.Printf("Error closing listener: %v", err)
	}

	log.Println("Server stopped gracefully.")
}
```

---

### **Explanation of the Code**

#### **1. Starting the TCP Server**
```go
listener, err := net.Listen("tcp", address)
if err != nil {
	log.Fatalf("Failed to start TCP server: %v", err)
}
defer listener.Close()
```
- **What It Does**:
  - Starts a TCP listener on the specified address (e.g., `:8081`).
  - Uses `net.Listen` to create a listener that accepts incoming TCP connections.

- **Why It’s Necessary**:
  - The listener is the entry point for all incoming connections.

---

#### **2. Handling Connections**
```go
conn, err := listener.Accept()
if err != nil {
	log.Printf("Error accepting connection: %v", err)
	continue
}
go handleConnection(conn)
```
- **What It Does**:
  - Accepts a new TCP connection using `listener.Accept`.
  - Handles each connection in a separate goroutine using `handleConnection`.

- **Why It’s Necessary**:
  - Each client connection is handled independently, allowing the server to handle multiple clients concurrently.

---

#### **3. Processing Client Data**
```go
reader := bufio.NewReader(conn)
message, err := reader.ReadString('\n')
if err != nil {
	log.Printf("Client disconnected: %s", conn.RemoteAddr())
	return
}
```
- **What It Does**:
  - Reads data from the client using a buffered reader (`bufio.NewReader`).
  - Waits for a newline character (`\n`) to mark the end of a message.

- **Why It’s Necessary**:
  - Provides a simple way to read text-based messages from the client.

---

#### **4. Responding to the Client**
```go
response := fmt.Sprintf("Echo: %s", message)
_, err = conn.Write([]byte(response))
if err != nil {
	log.Printf("Error writing to client %s: %v", conn.RemoteAddr(), err)
	return
}
```
- **What It Does**:
  - Sends a response back to the client by writing to the connection.

- **Why It’s Necessary**:
  - Demonstrates how to send data back to the client after processing.

---

#### **5. Graceful Shutdown**
```go
sig := make(chan os.Signal, 1)
signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

<-sig

log.Println("Shutting down server...")
err = listener.Close()
if err != nil {
	log.Printf("Error closing listener: %v", err)
}
```
- **What It Does**:
  - Listens for interrupt signals (`Ctrl+C`) to gracefully shut down the server.
  - Closes the listener to stop accepting new connections.

- **Why It’s Necessary**:
  - Ensures the server can be stopped cleanly without abruptly terminating active connections.

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

2. **Connect to the Server**:
   - Use a TCP client (e.g., `telnet` or `nc`) to connect to the server:
     ```bash
     telnet localhost 8081
     ```
   - Or use `nc` (Netcat):
     ```bash
     nc localhost 8081
     ```

3. **Send Messages**:
   - Type a message and press Enter:
     ```
     Hello, TCP server!
     ```
   - The server responds with:
     ```
     Echo: Hello, TCP server!
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

2. **Custom Protocols**:
   - You can implement custom communication protocols instead of relying on HTTP.

3. **Low-Level Control**:
   - Direct access to TCP sockets allows fine-grained control over data transmission.

4. **Scalability**:
   - The server can handle multiple clients simultaneously without blocking.

---

### **Extending the Implementation**

Here are some ways to enhance the TCP server:

1. **Message Framing**:
   - Instead of using `\n` as a delimiter, implement a more robust framing mechanism (e.g., length-prefixed messages).

2. **Authentication**:
   - Add authentication logic to verify client credentials before processing requests.

3. **Timeouts**:
   - Set read/write timeouts to prevent connections from hanging indefinitely:
     ```go
     conn.SetReadDeadline(time.Now().Add(10 * time.Second))
     ```

4. **Broadcasting**:
   - Implement a broadcast mechanism to send messages to all connected clients.

5. **TLS Encryption**:
   - Use `tls.Listen` to enable secure communication over TLS.

---

This implementation demonstrates how to create a basic TCP server in Go.

*/
