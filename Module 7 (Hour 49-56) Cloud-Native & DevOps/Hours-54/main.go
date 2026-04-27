package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"sync/atomic"
	"time"
)

var (
	startedAt    = time.Now()
	requestCount uint64
)

func startupDelay() time.Duration {
	raw := os.Getenv("STARTUP_DELAY_SECONDS")
	if raw == "" {
		return 5 * time.Second
	}
	seconds, err := strconv.Atoi(raw)
	if err != nil {
		return 5 * time.Second
	}
	return time.Duration(seconds) * time.Second
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func main() {
	delay := startupDelay()
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		total := atomic.AddUint64(&requestCount, 1)
		writeJSON(w, http.StatusOK, map[string]any{
			"hour":          54,
			"topic":         "Health probes and scaling",
			"startup_delay": delay.String(),
			"request_count": total,
		})
	})
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusOK, map[string]string{"status": "alive"})
	})
	mux.HandleFunc("/readyz", func(w http.ResponseWriter, r *http.Request) {
		if time.Since(startedAt) < delay {
			writeJSON(w, http.StatusServiceUnavailable, map[string]string{"status": "warming_up"})
			return
		}
		writeJSON(w, http.StatusOK, map[string]string{"status": "ready"})
	})
	mux.HandleFunc("/metrics", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain; version=0.0.4")
		fmt.Fprintf(w, "app_requests_total %d\n", atomic.LoadUint64(&requestCount))
		fmt.Fprintf(w, "app_uptime_seconds %.0f\n", time.Since(startedAt).Seconds())
	})

	addr := ":" + port
	log.Printf("hour 54 app listening on %s", addr)
	log.Fatal(http.ListenAndServe(addr, mux))
}
