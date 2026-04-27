package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

func worker(ctx context.Context, id int, high, low <-chan int, wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		select {
		case <-ctx.Done(): // Check if main said "stop"
			fmt.Printf("Worker %d shutting down safely...\n", id)
			return
		case task, ok := <-high:
			if !ok { high = nil; break }
			fmt.Printf("Worker %d [HIGH] task %d\n", id, task)
			time.Sleep(100 * time.Millisecond) // Simulate heavy work
		case task, ok := <-low:
			if !ok { low = nil; break }
			fmt.Printf("Worker %d [LOW ] task %d\n", id, task)
			time.Sleep(100 * time.Millisecond) // Simulate heavy work
		}

		if high == nil && low == nil {
			return
		}
	}
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	highPriority := make(chan int, 20)
	lowPriority := make(chan int, 20)
	var wg sync.WaitGroup

	for i := 1; i <= 3; i++ {
		wg.Add(1)
		go worker(ctx, i, highPriority, lowPriority, &wg)
	}

	// Feed tasks
	for i := 1; i <= 10; i++ {
		highPriority <- i + 100
		lowPriority <- i
	}

	// Wait a bit, then trigger an "Emergency Shutdown" 
	// before all tasks are even finished
	time.Sleep(250 * time.Millisecond)
	fmt.Println("--- SHUTDOWN SIGNAL SENT ---")
	cancel() 

	wg.Wait()
	fmt.Println("All workers exited cleanly.")
}
