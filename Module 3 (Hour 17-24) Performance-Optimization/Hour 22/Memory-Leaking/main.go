package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	_ "net/http/pprof" // Standard profiling at /debug/pprof/
	"runtime"
	"sync"
)

var (
	// Root context for broadcasting the shutdown signal
	rootCtx, cancelAll = context.WithCancel(context.Background())
	wg                 sync.WaitGroup
)

// LEAK: This goroutine has no exit condition.
func triggerLeak() {
	go func() {
		ch := make(chan int)
		<-ch // Stuck here forever (ZOMBIE)
	}()
}

// FIXED: This goroutine monitors the context for a shutdown signal.
func triggerFixed() {
	wg.Add(1)
	go func() {
		defer wg.Done()
		ch := make(chan int)
		select {
		case <-ch:
			// Success path
		case <-rootCtx.Done():
			// The FIX: Function returns, killing the goroutine
			return
		}
	}()
}

func main() {
	mux := http.NewServeMux()

	// 1. THE MATRIX: View current runtime metrics
	mux.HandleFunc("/matrix", func(w http.ResponseWriter, r *http.Request) {
		runtime.GC() // Clean up dead goroutines before reporting
		var m runtime.MemStats
		runtime.ReadMemStats(&m)
		
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"goroutines":    runtime.NumGoroutine(),
			"heap_alloc_mb": m.Alloc / 1024 / 1024,
			"sys_memory_mb": m.Sys / 1024 / 1024,
		})
	})

	// 2. TRIGGER LEAK: Inject 1000 permanent zombies
	mux.HandleFunc("/leak", func(w http.ResponseWriter, r *http.Request) {
		for i := 0; i < 1000; i++ { triggerLeak() }
		fmt.Fprint(w, "Injected 1000 leaky goroutines.")
		log.Println("Action: 1000 LEAKS injected.")
	})

	// 3. SPAWN FIXED: Inject 1000 manageable workers
	mux.HandleFunc("/spawn", func(w http.ResponseWriter, r *http.Request) {
		for i := 0; i < 1000; i++ { triggerFixed() }
		fmt.Fprint(w, "Spawned 1000 fixed goroutines.")
		log.Println("Action: 1000 FIXED workers spawned.")
	})

	// 4. APPLY THE FIX: Close the context to kill all fixed workers
	mux.HandleFunc("/fix", func(w http.ResponseWriter, r *http.Request) {
		cancelAll() // Send signal
		wg.Wait()   // Block until they actually exit
		
		// Re-initialize context for the next test cycle
		rootCtx, cancelAll = context.WithCancel(context.Background())
		
		fmt.Fprint(w, "FIXED: All manageable goroutines terminated.")
		log.Println("Action: FIX applied. Context cancelled.")
	})

	// Standard pprof handler for visual profiling
	mux.Handle("/debug/pprof/", http.DefaultServeMux)

	fmt.Println("--------------------------------------------------")
	fmt.Println("🚀 GO RUNTIME MONITORING SERVER ACTIVE")
	fmt.Println("--------------------------------------------------")
	fmt.Println("📊 MATRIX: http://localhost:9000/matrix")
	fmt.Println("🔴 LEAK:   http://localhost:9000/leak")
	fmt.Println("🟢 SPAWN:  http://localhost:9000/spawn")
	fmt.Println("🛠️  FIX:    http://localhost:9000/fix")
	fmt.Println("📈 PPROF:  http://localhost:9000/debug/pprof/")
	fmt.Println("--------------------------------------------------")

	log.Fatal(http.ListenAndServe(":9000", mux))
}