/*
5. Use Context for Cancellation
Use context.Context to cancel goroutines gracefully when they are no longer needed. This prevents goroutines from leaking and wasting resources.


Optimization:
context.Context allows you to cancel goroutines cleanly, preventing resource leaks.
*/

package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

func worker(ctx context.Context, id int, tasks <-chan int, wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		select {
		case <-ctx.Done():
			fmt.Printf("Worker %d stopping\n", id)
			return
		case task, ok := <-tasks:
			if !ok {
				fmt.Printf("Worker %d stopping (channel closed)\n", id)
				return
			}
			fmt.Printf("Worker %d processing task %d\n", id, task)
		}
	}
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	numWorkers := 3
	numTasks := 10

	tasks := make(chan int, numTasks)
	var wg sync.WaitGroup

	// Start workers
	for i := 1; i <= numWorkers; i++ {
		wg.Add(1)
		go worker(ctx, i, tasks, &wg)
	}

	// Send tasks to the channel
	for i := 1; i <= numTasks; i++ {
		tasks <- i
	}
	close(tasks)

	// Simulate cancellation after some time
	time.Sleep(2 * time.Second)
	cancel()

	// Wait for all workers to finish
	wg.Wait()
	fmt.Println("All workers stopped")
}