//Hour 3: Manual Middleware Chaining + Logging Wrapper
package main

import (
	"log"
	"net/http"
	"time"
)

// Middleware signature
type Middleware func(http.Handler) http.Handler

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		log.Printf("START %s %s from %s", r.Method, r.URL.Path, r.RemoteAddr)
		next.ServeHTTP(w, r)
		log.Printf("END %s %s - took %v", r.Method, r.URL.Path, time.Since(start))
	})
}

func RecoveryMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("PANIC recovered: %v", err)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			}
		}()
		next.ServeHTTP(w, r)
	})
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})

	// Chain middleware (order matters!)
	handler := RecoveryMiddleware(LoggingMiddleware(mux))

	server := &http.Server{Addr: ":8081", Handler: handler}
	log.Println("Middleware server on :8081")
	log.Fatal(server.ListenAndServe())
}

/*
Run & test the same way as above. Watch the console logs.
Production note: This is exactly how popular frameworks (Gin, Echo, etc.) 
implement middleware under the hood.
*/