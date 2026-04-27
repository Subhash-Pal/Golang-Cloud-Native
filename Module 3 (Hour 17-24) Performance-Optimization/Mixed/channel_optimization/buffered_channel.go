/*

1. Use Buffered Channels
Buffered channels reduce contention by allowing senders to queue messages without waiting for receivers to be ready immediately. 
This avoids blocking and improves throughput.

Optimization:
The buffered channel (make(chan int, numTasks)) allows the sender to queue all tasks without waiting for workers to process them immediately.
Reduces contention and improves throughput.

*/
package main

import (
	"fmt"
	"sync"
)

func worker(id int, tasks <-chan int, wg *sync.WaitGroup) {
	defer wg.Done()
	for task := range tasks {
		fmt.Printf("Worker %d processing task %d\n", id, task)
	}
}

func main() {
	numWorkers := 3
	numTasks := 10

	// Create a buffered channel with capacity equal to the number of tasks
	tasks := make(chan int, numTasks)
	var wg sync.WaitGroup

	// Start workers
	for i := 1; i <= numWorkers; i++ {
		wg.Add(1)
		go worker(i, tasks, &wg)
	}

	// Send tasks to the channel
	for i := 1; i <= numTasks; i++ {
		tasks <- i
	}
	close(tasks) // Close the channel to signal no more tasks

	// Wait for all workers to finish
	wg.Wait()
	fmt.Println("All tasks completed")
}

