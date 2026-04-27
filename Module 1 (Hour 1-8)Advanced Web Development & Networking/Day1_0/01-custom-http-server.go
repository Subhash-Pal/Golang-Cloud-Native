package main

import (
	"fmt"
	"log"
	"net/http"
	//"strings"
	"sync/atomic"
	"time"
)

// Custom multiplexer (manual routing)
type CustomMux struct {
	routes map[string]http.HandlerFunc
}

func NewCustomMux() *CustomMux {
	return &CustomMux{routes: make(map[string]http.HandlerFunc)}
}

func (mux *CustomMux) HandleFunc(pattern string, handler http.HandlerFunc) {
	mux.routes[pattern] = handler
}

func (mux *CustomMux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Manual route matching (exact match for simplicity - you can extend to prefix/path params)
	if handler, ok := mux.routes[r.URL.Path]; ok {
		handler(w, r)
		return
	}
	http.NotFound(w, r)
}

var requestCounter int64

func helloHandler(w http.ResponseWriter, r *http.Request) {
	atomic.AddInt64(&requestCounter, 1)
	fmt.Fprintf(w, "Hello from custom server! Request #%d | Method: %s | Path: %s\n",
		atomic.LoadInt64(&requestCounter), r.Method, r.URL.Path)
}

func main() {
	mux := NewCustomMux()
	mux.HandleFunc("/", helloHandler)
	mux.HandleFunc("/api/v1/status", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{"status":"ok","timestamp":"`+time.Now().Format(time.RFC3339)+`"}`)
	})

	server := &http.Server{
		Addr:         ":8080",
		Handler:      mux,                    // our manual router
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	log.Println("Custom HTTP server starting on :8080 (manual routing + net/http internals)")
	log.Fatal(server.ListenAndServe())
}

/*

Hour 1–2: net/http Internals + Custom HTTP Server (No Framework) + Manual Routing
Key concepts covered

net/http server lifecycle (ListenAndServe → goroutine per request → ServeHTTP)
HTTP request parsing, header reading, body handling
Concurrency model (one goroutine per connection by default)
Manual routing (no http.Handle / http.ServeMux shortcuts)


Run:
Bash
````
go run 01-custom-http-server.go
````
Test:
Bash
curl http://localhost:8080/
curl http://localhost:8080/api/v1/status
curl http://localhost:8080/notfound   # 404

Nuances & Edge Cases

ServeHTTP is called in its own goroutine by the server → true concurrency.
If you block in one handler, other requests are unaffected (unlike some languages).
Production implication: Always set ReadTimeout/WriteTimeout to prevent Slowloris attacks.

*/