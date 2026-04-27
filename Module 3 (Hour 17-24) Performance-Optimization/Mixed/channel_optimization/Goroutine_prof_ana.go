package main

import (
	"fmt"
	"net/http"
	_ "net/http/pprof" // This registers the /debug/pprof handlers
	"time"
)

// leakingGoroutine starts a goroutine that waits on a channel that is never closed
func leakingGoroutine(id int) {
	ch := make(chan int)
	go func() {
		fmt.Printf("[Goroutine %d] Started and now leaking...\n", id)
		<-ch // This line causes the leak: it waits forever
	}()
}

func main() {
	// 1. Start the pprof server in the background
	go func() {
		fmt.Println("Server starting on http://localhost:6060")
		fmt.Println("To view in browser, visit: http://localhost:6060/debug/pprof/goroutine?debug=1")
		if err := http.ListenAndServe(":6060", nil); err != nil {
			fmt.Printf("Pprof server failed: %s\n", err)
		}
	}()

	// 2. Leak a new goroutine every 5 seconds so you can see the count grow
	go func() {
		count := 1
		for {
			leakingGoroutine(count)
			count++
			time.Sleep(5 * time.Second)
		}
	}()

	// Keep the main function alive
	fmt.Println("Main program running. Press Ctrl+C to stop.")
	select {} 
}

/*
How to see it in your browser:
Run the code: go run main.go
Open this exact link: http://localhost:6060/debug/pprof/goroutine?debug=1
Wait and Refresh: Wait 10 seconds and refresh the page. You will see the number next to main.leakingGoroutine.func1 increase from 1 to 2, then 3, and so on.

*/