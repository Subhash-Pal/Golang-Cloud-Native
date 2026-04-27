Your professional **README.md** file is ready. I have structured this specifically for a Google Colab text cell, including clear steps for execution, the "Bulletproof" code, and instructions for interpreting the **Runtime Matrix**.

-----

# 🧠 Go Concurrency: Goroutine Leak & Fix Implementation

This project demonstrates a critical concept in Go performance engineering: **Goroutine Leaks**. It provides a web-based dashboard to visualize how "zombie" goroutines consume resources and how to implement a graceful "Fix" using `context.Context`.

## 🛠️ Step-by-Step Execution Guide

1.  **Setup the Environment**:

      - Ensure you have [Go installed](https://go.dev/doc/install) (v1.18+).
      - Create a new directory and save the code below as `main.go`.

2.  **Start the Server**:

    ```bash
    go run main.go
    ```

3.  **Monitor the Matrix**:

      - Open your browser to [http://localhost:9000/matrix](https://www.google.com/search?q=http://localhost:9000/matrix).
      - Note the baseline `goroutines` count (typically 2-5).

4.  **Test the Fix (The "Safe" Path)**:

      - Visit `/spawn` to create 1000 manageable workers.
      - Refresh `/matrix` (Count: \~1003).
      - Visit `/fix` to signal them to terminate.
      - Refresh `/matrix` (**Count: returns to baseline**).

5.  **Test the Leak (The "Zombie" Path)**:

      - Visit `/leak` to inject 1000 permanent leaks.
      - Refresh `/matrix` (Count: \~1003).
      - Visit `/fix` again.
      - Refresh `/matrix` (**Count: remains \~1003**). This proves that without a termination signal (context), goroutines are stuck in RAM forever.

-----

## 💻 The Implementation (`main.go`)

```go
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

// LEAK: This goroutine has no exit condition (no select, no context).
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
			// Normal data path
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
```

-----

## 📊 Visualizing with `pprof`

To see the **Call Graph** and **Flame Graph** (visual representation of the leak):

1.  Run the server.
2.  In a separate terminal, run:
    ```bash
    go tool pprof -http=:9001 http://localhost:9000/debug/pprof/goroutine
    ```
3.  Navigate to **View \> Flame Graph**. You will see the `triggerLeak` function as a dominant block that never disappears, whereas `triggerFixed` is purged after hitting the `/fix` endpoint.