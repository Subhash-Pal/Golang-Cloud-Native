package main

import (
	"fmt"
	"sync"
	"time"
)

// Task represents a job that needs to be executed by a worker.
type Task struct {
	ID       int
	Duration time.Duration // Simulated task duration
}

// WorkerPool manages a pool of workers and distributes tasks among them.
type WorkerPool struct {
	Tasks    chan Task      // Channel for tasks
	Workers  int            // Number of workers
	wg       sync.WaitGroup // WaitGroup to wait for all workers to finish
	stopChan chan struct{}  // Channel to signal workers to stop
}

// NewWorkerPool creates a new WorkerPool with the specified number of workers.
func NewWorkerPool(workers int) *WorkerPool {
	return &WorkerPool{
		Tasks:    make(chan Task, 100), // Buffered channel for tasks
		Workers:  workers,
		stopChan: make(chan struct{}),
	}
}

// Start initializes and starts the worker pool.
func (wp *WorkerPool) Start() {
	for i := 0; i < wp.Workers; i++ {
		wp.wg.Add(1)
		go func(workerID int) {
			defer wp.wg.Done()
			fmt.Printf("Worker %d started\n", workerID)
			for {
				select {
				case task, ok := <-wp.Tasks:
					if !ok { // Channel closed, exit
						fmt.Printf("Worker %d stopping\n", workerID)
						return
					}
					fmt.Printf("Worker %d processing Task %d\n", workerID, task.ID)
					time.Sleep(task.Duration) // Simulate task processing
					fmt.Printf("Worker %d completed Task %d\n", workerID, task.ID)
				case <-wp.stopChan: // Stop signal received
					fmt.Printf("Worker %d stopping\n", workerID)
					return
				}
			}
		}(i + 1) // Pass worker ID
	}
}

// AddTask adds a new task to the worker pool.
func (wp *WorkerPool) AddTask(task Task) {
	wp.Tasks <- task
}

// Stop gracefully stops the worker pool.
func (wp *WorkerPool) Stop() {
	close(wp.Tasks) // Close the task channel to signal workers to exit
	wp.wg.Wait()    // Wait for all workers to finish
	close(wp.stopChan)
	fmt.Println("Worker pool stopped")
}

func main() {
	// Create a worker pool with 3 workers
	pool := NewWorkerPool(3)
	pool.Start()

	// Add some tasks to the pool
	for i := 1; i <= 10; i++ {
		task := Task{
			ID:       i,
			Duration: time.Duration(i*100) * time.Millisecond, // Simulated task duration
		}
		pool.AddTask(task)
	}

	// Let the workers process the tasks
	time.Sleep(2 * time.Second)

	// Stop the worker pool
	pool.Stop()
}

/*
Explanation of the Code
Task Struct:
Represents a unit of work with an ID and a simulated Duration to mimic processing time.
WorkerPool Struct:
Manages the worker pool with:
Tasks: A buffered channel to hold tasks.
Workers: The number of workers (goroutines) in the pool.
wg: A sync.WaitGroup to ensure all workers finish before stopping the pool.
stopChan: A channel to signal workers to stop gracefully.
Start Method:
Launches the specified number of workers as goroutines.
Each worker listens on the Tasks channel for incoming tasks and processes them.
Workers exit when the Tasks channel is closed or when a stop signal is received.
AddTask Method:
Adds a new task to the Tasks channel for processing.
Stop Method:
Closes the Tasks channel to signal workers to stop.
Waits for all workers to finish using the WaitGroup.
Main Function:
Creates a worker pool with 3 workers.
Adds 10 tasks with varying durations.
Stops the pool after allowing time for task processing.


Output Example
When you run the program, you might see output similar to this:
1234567891011121314151617
Key Features
Concurrency Control:
Limits the number of concurrent goroutines to the size of the worker pool.
Graceful Shutdown:
Ensures all workers finish their current tasks before shutting down.
Scalability:
You can easily adjust the number of workers based on the workload.
Task Buffering:
The buffered Tasks channel allows tasks to be queued without blocking the main thread.


*/