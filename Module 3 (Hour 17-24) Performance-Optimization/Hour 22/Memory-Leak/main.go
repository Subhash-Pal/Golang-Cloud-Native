package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"runtime"
	"sync"
	"time"

	_ "net/http/pprof"
)

/*
GLOBAL TRACKING
We track all "leaked" goroutines so we can cancel them later.
*/
var (
	mu      sync.Mutex
	cancels []context.CancelFunc
)

/*
========================
LEAK HANDLER
========================
Creates goroutines that WAIT FOREVER unless cancelled.
*/
func leakHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithCancel(context.Background())

	mu.Lock()
	cancels = append(cancels, cancel)
	mu.Unlock()

	go func(ctx context.Context) {
		select {
		case <-ctx.Done():
			return // cleanup when cancelled
		}
	}(ctx)

	fmt.Fprintln(w, "leak created")
}

/*
========================
FIX HANDLER (REAL FIX)
========================
Cancels ALL leaked goroutines.
*/
func fixHandler(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	for _, cancel := range cancels {
		cancel() // ✅ stop goroutines
	}
	cancels = nil
	mu.Unlock()

	// allow scheduler + GC to clean up
	time.Sleep(200 * time.Millisecond)

	fmt.Fprintln(w, "all leaks cleaned")
}

/*
========================
STATS
========================
*/
func statsHandler(w http.ResponseWriter, r *http.Request) {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	fmt.Fprintf(w,
		"Goroutines: %d\nHeapAlloc: %d bytes\nHeapObjects: %d\nNumGC: %d\n",
		runtime.NumGoroutine(),
		m.HeapAlloc,
		m.HeapObjects,
		m.NumGC,
	)
}

func main() {
	http.HandleFunc("/leak", leakHandler)
	http.HandleFunc("/fix", fixHandler)
	http.HandleFunc("/stats", statsHandler)

	log.Println("Server running at :8080")
	log.Println("Endpoints: /leak | /fix | /stats")

	log.Fatal(http.ListenAndServe(":8080", nil))
}