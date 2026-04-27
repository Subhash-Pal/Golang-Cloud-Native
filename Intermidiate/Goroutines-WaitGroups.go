package main

import (
	"fmt"
	"sync"
	"time"
)

func worker(id int, wg *sync.WaitGroup) {
	defer wg.Done() // Notify the WaitGroup when done

	fmt.Printf("Worker %d starting\n", id)
	time.Sleep(time.Duration(id) * time.Second) // Simulate work
	fmt.Printf("Worker %d done\n", id)
}

func main() {
	var wg sync.WaitGroup

	for i := 1; i <= 3; i++ {
		wg.Add(1)                // Increment the WaitGroup counter
		go worker(i, &wg)        // Launch a worker goroutine
	}

	wg.Wait() // Wait for all workers to finish
	fmt.Println("All workers completed")
}