package main

import (
	"context"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"golang.org/x/time/rate" //go get -u golang.org/x/time/rate
)

type client struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

var (
	clients = make(map[string]*client)
	mu      sync.Mutex
)

func getClient(ip string) *rate.Limiter {
	mu.Lock()
	defer mu.Unlock()

	v, exists := clients[ip]
	if !exists {
		limiter := rate.NewLimiter(2,5)
		clients[ip] = &client{limiter, time.Now()}
		return limiter
	}

	v.lastSeen = time.Now()
	return v.limiter
}

func cleanupClients(ctx context.Context) {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			mu.Lock()
			for ip, client := range clients {
				if time.Since(client.lastSeen) > 1*time.Minute {
					delete(clients, ip)
				}
			}
			mu.Unlock()
		case <-ctx.Done():
			log.Println("Cleanup worker stopping...")
			return
		}
	}
}

func limitMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			http.Error(w, "Internal Server Error", 500)
			return
		}

		limiter := getClient(ip)
		if !limiter.Allow() {
			w.Header().Set("Retry-After", "1")
			http.Error(w, "Too Many Requests", http.StatusTooManyRequests)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func main() {
	// 1. Create a context that listens for the interrupt signal (Ctrl+C)
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	// Start cleanup worker with the context
	go cleanupClients(ctx)

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{"status": "ok"}`))
	})

	server := &http.Server{
		Addr:    ":8080",
		Handler: limitMiddleware(mux),
	}

	// 2. Run the server in a goroutine so it doesn't block
	go func() {
		log.Println("Server starting on :8080. Press Ctrl+C to stop.")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Listen error: %s\n", err)
		}
	}()

	// 3. Wait for the Ctrl+C signal
	<-ctx.Done()
	log.Println("Shutting down gracefully...")

	// 4. Give the server 5 seconds to finish existing requests
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("Server exited.")
}
//Run the server and test with multiple requests to see the rate limiting in action.
//You can use the following command in PowerShell to send multiple requests:
//go run rate-limiter.go

//run the client test command in another terminal while the server is running:
// 1..10 | ForEach-Object { curl.exe -I http://localhost:8080 }


/*  ip is string 
   m[ip] = &client{limiter, time.Now()} // Create a new client with a rate limiter and store it in the map
/*