package main
/*

public class Myclass{

}
Myclass ob=new Myclass()
===
type mystruct struct{

}
ob:=mystruct{}

ob:=&mystruc{}

p:=new(mystruct)

*/
import (
	"fmt"
	"net"
	"net/http"
	"sync"
	"time"

	"golang.org/x/time/rate"
)

// visitor stores the rate limiter and the last time it was used for a specific IP.
type visitor struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

var (
	visitors = make(map[string]*visitor)
	mu       sync.Mutex
)

func init() {
	// Background routine to remove inactive visitors every minute
	go cleanupVisitors()
}

// getVisitor retrieves or creates a rate limiter for a given IP address.
func getVisitor(ip string) *rate.Limiter {
	mu.Lock()
	defer mu.Unlock()

	v, exists := visitors[ip]
	if !exists {
		// Allow 2 requests per second with a burst of 5
		limiter := rate.NewLimiter(1, 2)
		visitors[ip] = &visitor{limiter, time.Now()}
		return limiter
	}

	v.lastSeen = time.Now()
	return v.limiter
}

func cleanupVisitors() {
	for {
		time.Sleep(time.Minute)
		mu.Lock()
		for ip, v := range visitors {
			// Remove visitors inactive for more than 3 minutes
			if time.Since(v.lastSeen) > 3*time.Minute {
				delete(visitors, ip)
			}
		}
		mu.Unlock()
	}
}

// limitMiddleware wraps an http.Handler to apply the rate limiting.
func limitMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Extract IP address from the request
		ip, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		limiter := getVisitor(ip)
		if !limiter.Allow() {
			// Return 429 Too Many Requests if the limit is exceeded
			w.Header().Set("X-RateLimit-Retry-After", "wait a bit")
			http.Error(w, http.StatusText(http.StatusTooManyRequests), http.StatusTooManyRequests)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, you are within the rate limit!"))
	})

	fmt.Println("Server starting on :8080...")
	// Wrap the mux with our middleware
	http.ListenAndServe(":8080", limitMiddleware(mux))
}
/*
 Invoke-WebRequest -Uri "http://localhost:8080" -UseBasicParsing | Select-Object -Property StatusCode
>>     Start-Sleep -Milliseconds 600
>> }
>> */