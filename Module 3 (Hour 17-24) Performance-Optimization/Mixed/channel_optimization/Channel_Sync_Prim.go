/*
6. Combine Channels with Sync Primitives
In some cases, combining channels with synchronization primitives like sync.Mutex or sync.Cond can improve performance by reducing contention.

Optimization:
A sync.Mutex ensures safe access to shared resources without relying solely on channels.
Summary of Best Practices
Use buffered channels to reduce contention.
Avoid blocking by using select with a default case.
Use multiple channels for prioritized tasks.
Minimize channel buffer sizes to save memory.
Use context.Context for graceful cancellation.
Combine channels with synchronization primitives when necessary.
By applying these optimizations, you can improve the performance and scalability of your Go programs that rely on channels for communication.
*/

package main

import (
	"fmt"
	"sync"
)

func worker(id int, tasks <-chan int, wg *sync.WaitGroup, mu *sync.Mutex) {
	defer wg.Done()
	for task := range tasks {
		mu.Lock()
		fmt.Printf("Worker %d processing task %d\n", id, task)
		mu.Unlock()
	}
}

func main() {
	numWorkers := 3
	numTasks := 10

	tasks := make(chan int, numTasks)
	var wg sync.WaitGroup
	var mu sync.Mutex

	// Start workers
	for i := 1; i <= numWorkers; i++ {
		wg.Add(1)
		go worker(i, tasks, &wg, &mu)
	}

	// Send tasks to the channel
	for i := 1; i <= numTasks; i++ {
		tasks <- i
	}
	close(tasks)

	// Wait for all workers to finish
	wg.Wait()
	fmt.Println("All tasks completed")
}