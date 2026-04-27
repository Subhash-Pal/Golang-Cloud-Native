/*
4. Minimize Channel Size
Using excessively large buffers can waste memory. Choose a buffer size that balances throughput and memory usage.
Example:


Optimization:
A smaller buffer size reduces memory overhead while still providing enough capacity to avoid excessive blocking.



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

	// Use a small buffer size to minimize memory usage
	tasks := make(chan int, 2) // Small buffer size
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
	close(tasks)

	// Wait for all workers to finish
	wg.Wait()
	fmt.Println("All tasks completed")
}