package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

//////////////////////
// Router (Manual)
//////////////////////

type HandlerFunc func(http.ResponseWriter, *http.Request)

type Router struct {
	routes map[string]map[string]HandlerFunc
}

func NewRouter() *Router {
	return &Router{
		routes: make(map[string]map[string]HandlerFunc),
	}
}

func (r *Router) Handle(method, path string, handler HandlerFunc) {
	if _, exists := r.routes[path]; !exists {
		r.routes[path] = make(map[string]HandlerFunc)
	}
	r.routes[path][method] = handler
}

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if methods, ok := r.routes[req.URL.Path]; ok {
		if handler, ok := methods[req.Method]; ok {
			handler(w, req)
			return
		}
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	http.NotFound(w, req)
}

//////////////////////
// Middleware System
//////////////////////

type Middleware func(HandlerFunc) HandlerFunc

func Chain(h HandlerFunc, m ...Middleware) HandlerFunc {
	for i := len(m) - 1; i >= 0; i-- {
		h = m[i](h)
	}
	return h
}

//////////////////////
// Logging Middleware
//////////////////////

func LoggingMiddleware(next HandlerFunc) HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		log.Printf("Started %s %s", r.Method, r.URL.Path)

		next(w, r)

		log.Printf("Completed in %v", time.Since(start))
	}
}

//////////////////////
// Auth Middleware
//////////////////////

func AuthMiddleware(next HandlerFunc) HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")

		if token != "Bearer secret-token" {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("Unauthorized\n"))
			return
		}

		next(w, r)
	}
}

//////////////////////
// Token Bucket Rate Limiter
//////////////////////

type TokenBucket struct {
	capacity int
	tokens   int
	rate     int
	last     time.Time
	mu       sync.Mutex
}

func NewTokenBucket(capacity, rate int) *TokenBucket {
	return &TokenBucket{
		capacity: capacity,
		tokens:   capacity,
		rate:     rate,
		last:     time.Now(),
	}
}

func (tb *TokenBucket) Allow() bool {
	tb.mu.Lock()
	defer tb.mu.Unlock()

	now := time.Now()
	elapsed := now.Sub(tb.last).Seconds()

	// refill tokens
	tb.tokens += int(elapsed * float64(tb.rate))
	if tb.tokens > tb.capacity {
		tb.tokens = tb.capacity
	}

	tb.last = now

	if tb.tokens > 0 {
		tb.tokens--
		return true
	}

	return false
}

func RateLimitMiddleware(tb *TokenBucket) Middleware {
	return func(next HandlerFunc) HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			if !tb.Allow() {
				w.WriteHeader(http.StatusTooManyRequests)
				w.Write([]byte("Rate limit exceeded\n"))
				return
			}
			next(w, r)
		}
	}
}

//////////////////////
// Handlers
//////////////////////

func HelloHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello, World!\n"))
}

func SecureHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Secure endpoint accessed\n"))
}

//////////////////////
// Main (Server + Graceful Shutdown)
//////////////////////

func main() {
	router := NewRouter()

	// Rate limiter: capacity=5, refill=1 token/sec
	tb := NewTokenBucket(5, 1)

	// Public route
	router.Handle("GET", "/", Chain(
		HelloHandler,
		LoggingMiddleware,
		RateLimitMiddleware(tb),
	))

	// Protected route
	router.Handle("GET", "/secure", Chain(
		SecureHandler,
		LoggingMiddleware,
		AuthMiddleware,
		RateLimitMiddleware(tb),
	))

	server := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	// Run server in goroutine
	go func() {
		fmt.Println("Server running on http://localhost:8080")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Listen error: %v", err)
		}
	}()

	// Graceful shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	<-stop
	fmt.Println("\nShutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Shutdown failed: %v", err)
	}

	fmt.Println("Server gracefully stopped")
}


/*
LAB 1: Custom HTTP Engine (Go)

DESCRIPTION:
- Manual HTTP server using net/http
- Custom router (no frameworks)
- Middleware chaining (logging, auth)
- Token bucket rate limiter
- Graceful shutdown (SIGINT/SIGTERM)

----------------------------------------
PREREQUISITES:
1. Install Go (>= 1.20)
2. Verify installation:
   go version

3. Open this folder in VS Code

----------------------------------------
RUN SERVER:
1. Open terminal in VS Code
2. Execute:
   go run main.go

3. Server will start:
   http://localhost:8080

----------------------------------------
TEST ENDPOINTS:

[1] Public Endpoint:
- URL: http://localhost:8080/
- Command (PowerShell):
  curl http://localhost:8080/

- Expected Output:
  Hello, World!

----------------------------------------

[2] Secure Endpoint (Auth Required):

- Without token:
  curl http://localhost:8080/secure
  → Response: Unauthorized

- With token (PowerShell):
  curl http://localhost:8080/secure -Headers @{"Authorization"="Bearer secret-token"}

- OR (recommended real curl):
  curl.exe -H "Authorization: Bearer secret-token" http://localhost:8080/secure

- Expected Output:
  Secure endpoint accessed

----------------------------------------

[3] Rate Limiting Test:
- Rapidly call endpoint >5 times:
  → Response: Rate limit exceeded (HTTP 429)

----------------------------------------

GRACEFUL SHUTDOWN:
- Press CTRL + C
- Server waits up to 5 seconds for active requests
- Then exits cleanly

----------------------------------------

ARCHITECTURE OVERVIEW:

Router:
- Map-based routing: path → method → handler
- O(1) lookup

Middleware Chain:
Request Flow:
RateLimit → Auth → Logging → Handler

Components:
- LoggingMiddleware: logs request lifecycle
- AuthMiddleware: checks Authorization header
- RateLimitMiddleware: token bucket control

----------------------------------------

TOKEN BUCKET CONFIG:
- Capacity: 5 tokens
- Refill Rate: 1 token/sec

----------------------------------------

DEBUG TIP:
To inspect headers, add:
log.Println(r.Header)

----------------------------------------

COMMON ISSUE (Windows PowerShell):
- curl is alias of Invoke-WebRequest
- Use:
  -Headers @{"Authorization"="Bearer secret-token"}

----------------------------------------
*/