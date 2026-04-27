package main

import (
	"context"
	"log/slog"
	"sync"
	"time"
)

type job struct {
	ID      int
	Payload int
}

type result struct {
	JobID int
	Value int
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	jobs := make(chan job)
	results := make(chan result, 8)

	var workers sync.WaitGroup
	for workerID := range 4 {
		workers.Add(1)
		go func(id int) {
			defer workers.Done()
			worker(ctx, id+1, jobs, results)
		}(workerID)
	}

	go func() {
		defer close(jobs)
		for i := 1; i <= 12; i++ {
			jobs <- job{ID: i, Payload: i * 10}
		}
	}()

	go func() {
		workers.Wait()
		close(results)
	}()

	for item := range results {
		slog.Info("result received", "job_id", item.JobID, "value", item.Value)
	}
}

func worker(ctx context.Context, workerID int, jobs <-chan job, results chan<- result) {
	for {
		select {
		case <-ctx.Done():
			return
		case jobItem, ok := <-jobs:
			if !ok {
				return
			}

			time.Sleep(80 * time.Millisecond)
			results <- result{
				JobID: jobItem.ID,
				Value: jobItem.Payload * 2,
			}

			slog.Info("worker completed job", "worker_id", workerID, "job_id", jobItem.ID)
		}
	}
}
