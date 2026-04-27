/*
Detecting goroutine leaks is crucial for maintaining the health and performance of Go applications. A goroutine leak occurs when a goroutine is unintentionally left running, consuming resources indefinitely. This can happen due to improper synchronization, unclosed channels, or forgotten defer statements.
Below are techniques and tools to detect and prevent goroutine leaks in Go.

1. Use runtime.NumGoroutine
The runtime.NumGoroutine function returns the number of currently active goroutines. By monitoring this value, you can identify unexpected increases in goroutine count.


Explanation:
The leakingGoroutine function creates a goroutine that waits indefinitely on a channel.
Monitoring runtime.NumGoroutine helps identify the increase in goroutine count.
*/

package main

import (
	"fmt"
	"runtime"
	"time"
)

func leakingGoroutine() {
	ch := make(chan int)
	go func() {
		<-ch // Goroutine waits indefinitely
	}()
}

func main() {
	fmt.Printf("Initial goroutines: %d\n", runtime.NumGoroutine())

	// Simulate a goroutine leak
	leakingGoroutine()

	// Check goroutine count after potential leak
	time.Sleep(1 * time.Second) // Allow time for goroutines to start
	fmt.Printf("Goroutines after leak: %d\n", runtime.NumGoroutine())
}