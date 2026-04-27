package main

import (
	"fmt"
	"net/http"
	_ "net/http/pprof" // Register pprof handlers
	"sync"
)

func worker(id int, wg *sync.WaitGroup) {
	defer wg.Done()
	// Increased iterations so there is actually something to profile
	for i := 0; i < 1e8; i++ {
		_ = i * i
	}
}

func main() {
	// Start pprof in a background goroutine
	go func() {
		fmt.Println("Pprof server running at http://localhost:6060/debug/pprof/")
		if err := http.ListenAndServe(":6060", nil); err != nil {
			fmt.Printf("Pprof server failed: %v\n", err)
		}
	}()

	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go worker(i, &wg)
	}

	wg.Wait()
	fmt.Println("All workers completed.")

	// CRITICAL: Keep the process alive so you can actually visit the URL 
	// and the pprof server doesn't die immediately.
	fmt.Println("Keeping server alive for profiling. Press Ctrl+C to exit.")
	select {} 
}
