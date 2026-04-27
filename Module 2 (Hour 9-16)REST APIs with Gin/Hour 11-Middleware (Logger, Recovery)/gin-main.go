package main

import (
	"github.com/gin-gonic/gin" //go get -u github.com/gin-gonic/gin
	"net/http"
)

func main() {
	// Create a new Gin router with default middleware (Logger and Recovery)
	router := gin.Default()

	// Define a route that simulates a panic
	router.GET("/panic", func(c *gin.Context) {
		// Simulate a runtime panic
		panic("Something went wrong!")
	})

	// Define a normal route
	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Hello, World!",
		})
	})

	// Start the server on port 8080
	port := ":8080"
	router.Run(port)
}

/*
Certainly! Let's look at a real-world example of middleware usage in Go by examining how the popular web framework **Gin** implements logging and recovery middleware. Gin is widely used in production-grade applications for building RESTful APIs, microservices, and more.

### Real-World Example: Middleware in Gin Framework

The Gin framework provides built-in middleware for logging (`Logger`) and recovery (`Recovery`). These are commonly used in production applications to handle request logging and gracefully recover from panics.

Here’s an example of how you might use Gin's `Logger` and `Recovery` middleware in a real-world application:

---

#### Code Example: Using Gin Middleware for Logging and Recovery

```go
package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func main() {
	// Create a new Gin router with default middleware (Logger and Recovery)
	router := gin.Default()

	// Define a route that simulates a panic
	router.GET("/panic", func(c *gin.Context) {
		// Simulate a runtime panic
		panic("Something went wrong!")
	})

	// Define a normal route
	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Hello, World!",
		})
	})

	// Start the server on port 8080
	port := ":8080"
	router.Run(port)
}
```

---

### Explanation of the Code:
1. **Gin Default Middleware**:
   - When you call `gin.Default()`, it automatically includes two middleware:
     - **Logger**: Logs HTTP requests (method, path, status code, response time, etc.).
     - **Recovery**: Recovers from panics and prevents the server from crashing.

2. **Panic Simulation**:
   - The `/panic` route intentionally triggers a panic using `panic("Something went wrong!")`.
   - The `Recovery` middleware catches this panic, logs the error, and returns a `500 Internal Server Error` response to the client.

3. **Normal Route**:
   - The `/` route demonstrates a typical API endpoint that responds with JSON data.

4. **Running the Server**:
   - The server listens on port `8080`. You can test it by visiting `http://localhost:8080` or `http://localhost:8080/panic`.

---

### Output Examples:

#### 1. Accessing `/`:
When you visit `http://localhost:8080/`, the response will be:
```json
{
  "message": "Hello, World!"
}
```

The console will log something like:
```
[GIN] 2026/04/18 - 14:32:10 | 200 |      123.456µs |       127.0.0.1 | GET      "/"
```

#### 2. Accessing `/panic`:
When you visit `http://localhost:8080/panic`, the server will catch the panic and respond with:
```json
{
  "error": "Internal Server Error"
}
```

The console will log something like:
```
[GIN] 2026/04/18 - 14:32:20 | 500 |       1.234ms |       127.0.0.1 | GET      "/panic"
[Recovery] 2026/04/18 - 14:32:20 panic recovered:
Something went wrong!
...
```

---

### Why This is a Real-World Example:
1. **Production-Grade Middleware**:
   - The `Logger` middleware helps developers monitor incoming requests and debug issues.
   - The `Recovery` middleware ensures that the server remains stable even if unexpected errors occur.

2. **Scalability**:
   - Gin is designed for high-performance applications and is often used in microservices architectures.
   - Middleware like `Logger` and `Recovery` can be extended or replaced with custom implementations as needed.

3. **Error Handling**:
   - By handling panics gracefully, the application avoids downtime and provides meaningful error responses to clients.

4. **Ease of Use**:
   - Gin's middleware system is simple yet powerful, making it easy to integrate into real-world projects.

---

### Customizing Middleware in Gin:
If you want to customize the logging format or recovery behavior, you can replace the default middleware with your own implementation. For example:

```go
router := gin.New() // Create a router without default middleware
router.Use(gin.LoggerWithFormatter(func(params gin.LogFormatterParams) string {
	return fmt.Sprintf("%s - [%s] \"%s %s %s %d %s\"\n",
		params.ClientIP,
		params.TimeStamp.Format(time.RFC1123),
		params.Method,
		params.Path,
		params.Request.Proto,
		params.StatusCode,
		params.Latency,
	)
}))
router.Use(gin.Recovery())
```

This allows you to tailor the middleware to your specific requirements while maintaining the robustness of the Gin framework.

---

### Conclusion:
This example demonstrates how middleware like logging and recovery is implemented in a real-world Go application using the Gin framework. These middleware components are essential for monitoring, debugging, and ensuring the stability of production-grade web services.
*/