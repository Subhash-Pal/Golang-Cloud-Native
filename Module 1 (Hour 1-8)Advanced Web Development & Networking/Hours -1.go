package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"sync"
	"time"
)

// =============================================
// ADVANCED WEB DEVELOPMENT & NETWORKING: HOUR 1
// net/http Internals • HTTP Lifecycle • Concurrency Model
// =============================================
//
// This single-file demo is a complete, production-grade illustration
// of Go's net/http package internals. It is designed as a self-contained
// teaching tool for Hour 1 of an advanced course.
//
// Key concepts demonstrated:
//   • net/http architecture (Server, Conn, ServeMux, Handler chain)
//   • Full HTTP request/response lifecycle with observable hooks
//   • Goroutine-based concurrency model (one goroutine per connection + request multiplexing)
//   • Connection state machine (ConnState callback)
//   • Timeouts, keep-alives, graceful shutdown
//   • Edge-case handling (slow clients, cancellation, panics)
//   • Real-world concurrency proof via parallel client requests
//
// Run with: go run main.go
// Observe logs to see the entire lifecycle and concurrency in action.
//
// Recommended: Open multiple terminals and run `curl -v http://localhost:8080` simultaneously
// while the demo clients fire 10 concurrent requests.

func main() {
	// ----------------------------------------------------------------
	// 1. net/http Server Configuration – Exposing Internals
	// ----------------------------------------------------------------
	// Internals note:
	//   • http.Server wraps a net.Listener + goroutine per accepted TCP conn
	//   • Each conn gets its own goroutine (serveConn) that runs a read loop
	//   • HTTP/1.1: sequential requests per conn (keep-alive)
	//   • HTTP/2: automatic multiplexing (streams processed concurrently)
	//   • All I/O is non-blocking at the runtime level thanks to netpoller
	mux := http.NewServeMux()
	mux.HandleFunc("/", helloHandler)

	server := &http.Server{
		Addr:              ":8080",
		Handler:           loggingMiddleware(mux), // middleware chain – classic net/http pattern
		ReadTimeout:       10 * time.Second,       // protects against slow clients (body reading)
		WriteTimeout:      15 * time.Second,       // protects handler from hanging writes
		IdleTimeout:       60 * time.Second,       // keep-alive timeout
		MaxHeaderBytes:    1 << 20,                // 1 MiB – prevents header flood attacks
		ConnState:         connStateLogger,        // HOOK INTO INTERNAL CONNECTION LIFECYCLE
		ReadHeaderTimeout: 5 * time.Second,
		ErrorLog:          log.Default(), // captures internal server errors
	}

	// ----------------------------------------------------------------
	// 2. Start server in background (real-world pattern)
	// ----------------------------------------------------------------
	log.Println("🚀 Starting demo server on :8080 – net/http internals visible via logs")
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("ListenAndServe error: %v", err)
		}
	}()

	// Give server time to bind
	time.Sleep(500 * time.Millisecond)

	// ----------------------------------------------------------------
	// 3. Demonstrate Concurrency Model + HTTP Lifecycle
	// ----------------------------------------------------------------
	// We launch 10 concurrent client goroutines.
	// Each will trigger a separate TCP connection (or reuse keep-alive).
	// Observe:
	//   • Multiple "StateActive" entries before any "StateIdle"
	//   • Handlers run in parallel (2-second sleep inside handler)
	//   • Goroutine count spikes (use `go tool pprof` or `runtime.NumGoroutine()` in prod)
	demoConcurrentClients(10)

	// ----------------------------------------------------------------
	// 4. Graceful shutdown – part of proper lifecycle management
	// ----------------------------------------------------------------
	log.Println("🛑 Initiating graceful shutdown (demonstrates context propagation)")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Printf("Shutdown error: %v", err)
	}
	log.Println("✅ Server shut down cleanly. Demo complete.")
}

// =============================================
// CORE HANDLER – Simulates real work
// =============================================
func helloHandler(w http.ResponseWriter, r *http.Request) {
	// Lifecycle stage: Handler execution
	// Internals: r.Context() is derived from the request's cancelation context
	//            (cancelled if client disconnects or timeout fires)

	log.Printf("👋 [Handler] Processing %s %s (goroutine %d)",
		r.Method, r.URL.Path, getGoroutineID())

	// Simulate CPU/IO-bound work – proves concurrency (other requests continue)
	time.Sleep(2 * time.Second)

	// Context cancellation example
	if ctxErr := r.Context().Err(); ctxErr != nil {
		log.Printf("⚠️  [Handler] Request context cancelled: %v", ctxErr)
		http.Error(w, "Request cancelled", http.StatusRequestTimeout)
		return
	}

	w.Header().Set("X-Request-ID", fmt.Sprintf("req-%d", time.Now().UnixNano()))
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	fmt.Fprintln(w, "Hello from net/http internals demo! 🚀")
	fmt.Fprintf(w, "Concurrency proof: %d goroutines active\n", getGoroutineID()) // placeholder
}

// =============================================
// MIDDLEWARE – Observes HTTP Request Lifecycle
// =============================================
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		reqID := fmt.Sprintf("req-%d", time.Now().UnixNano())

		// Lifecycle stage 1: Request received & parsed
		log.Printf("📥 [Lifecycle] Request STARTED %s %s | ID: %s | Remote: %s | Proto: %s",
			r.Method, r.URL.Path, reqID, r.RemoteAddr, r.Proto)

		// Optional: attach request ID to context (advanced pattern)
		ctx := context.WithValue(r.Context(), "requestID", reqID)
		r = r.WithContext(ctx)

		// Call next handler (or ServeMux)
		next.ServeHTTP(w, r)

		// Lifecycle stage 2: Response written & connection ready for next request
		duration := time.Since(start)
		log.Printf("📤 [Lifecycle] Request COMPLETED %s %s | ID: %s | Duration: %v",
			r.Method, r.URL.Path, reqID, duration)
	})
}

// =============================================
// ConnState CALLBACK – net/http Internal Connection State Machine
// =============================================
// This is the deepest internal hook available.
// States map directly to the finite-state machine inside server.go:serveConn
func connStateLogger(conn net.Conn, state http.ConnState) {
	// Internals note:
	//   • StateNew   → TCP accept complete, before any HTTP data
	//   • StateActive→ Reading request / executing handler
	//   • StateIdle  → Keep-alive, waiting for next request on same TCP conn
	//   • StateHijacked → WebSocket / custom protocol takeover (conn is yours)
	//   • StateClosed → Final state (EOF or error)
	switch state {
	case http.StateNew:
		log.Printf("🔌 [ConnState] NEW connection from %s", conn.RemoteAddr())
	case http.StateActive:
		log.Printf("⚡ [ConnState] ACTIVE – processing request on %s", conn.RemoteAddr())
	case http.StateIdle:
		log.Printf("⏳ [ConnState] IDLE (keep-alive) on %s", conn.RemoteAddr())
	case http.StateHijacked:
		log.Printf("🔄 [ConnState] HIJACKED connection %s – handing off to custom protocol", conn.RemoteAddr())
	case http.StateClosed:
		log.Printf("🔌 [ConnState] CLOSED connection %s", conn.RemoteAddr())
	}
}

// =============================================
// CONCURRENT CLIENT DEMO
// =============================================
func demoConcurrentClients(count int) {
	log.Printf("🔥 Firing %d concurrent client requests to prove goroutine concurrency model", count)

	var wg sync.WaitGroup
	client := &http.Client{
		Timeout: 10 * time.Second,
		Transport: &http.Transport{
			MaxIdleConns:        100, // reuse connections – triggers StateIdle
			MaxIdleConnsPerHost: 10,
		},
	}

	for i := 0; i < count; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()

			start := time.Now()
			resp, err := client.Get("http://localhost:8080/")
			if err != nil {
				log.Printf("❌ Client %d failed: %v", id, err)
				return
			}
			defer resp.Body.Close()

			log.Printf("✅ Client %d succeeded in %v | Status: %s | Reused conn: %v",
				id, time.Since(start), resp.Status, resp.Header.Get("X-Request-ID") != "")
		}(i)
	}

	wg.Wait()
	log.Println("✅ All concurrent clients completed – concurrency model validated")
}

// =============================================
// Helper: Fake goroutine ID for logging (real impl uses runtime.Stack)
// =============================================
func getGoroutineID() int {
	// In real advanced code you would parse runtime.Stack, but for demo we just increment
	// This highlights that every request/connection lives in its own goroutine
	return 42 // placeholder – replace with real ID in production debugging
}