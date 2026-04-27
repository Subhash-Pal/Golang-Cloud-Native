/*

3. Use Multiple Channels for Different Priorities
If tasks have different priorities, use separate channels for high-priority and low-priority tasks.
Workers can prioritize consuming from the high-priority channel.



Why this results in "Mixed" output:
No default block: The worker doesn't "check" High and then "skip" to Low. It waits for either to be ready.
Concurrency: Because you have 3 workers, one might grab a High task while the other grabs a Low task at the exact same millisecond.
Go's Select Logic: When both channels are ready, Go's select chooses a case at random. This naturally "mixes" the output while ensuring both streams keep moving.

Try this code: You should see the [HIGH] and [LOW] labels appearing out of order (e.g., LOW 1, HIGH 101, LOW 2, HIGH 102).
Does this interleaved behavior match your expectations for a priority system that still allows some low-priority tasks to run without being completely blocked by high-priority ones?

*/


package main


import (
	"fmt"
	"sync"
	"time"
)

func worker(id int, highPriority, lowPriority <-chan int, wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		// A single select allows Go to pick between them.
		// If both have data, it's a 50/50 split (Pseudo-random).
		// This ensures LOW priority tasks aren't totally blocked by HIGH ones.
		select {
		case task, ok := <-highPriority:
			if !ok {
				highPriority = nil
			} else {
				fmt.Printf("Worker %d [HIGH] task %d\n", id, task)
				time.Sleep(20 * time.Millisecond) // Simulate work
			}
		case task, ok := <-lowPriority:
			if !ok {
				lowPriority = nil
			} else {
				fmt.Printf("Worker %d [LOW ] task %d\n", id, task)
				time.Sleep(20 * time.Millisecond) // Simulate work
			}
		}

		if highPriority == nil && lowPriority == nil {
			return
		}
	}
}

func main() {
	numWorkers := 3
	highPriority := make(chan int, 20)
	lowPriority := make(chan int, 20)
	var wg sync.WaitGroup

	for i := 1; i <= numWorkers; i++ {
		wg.Add(1)
		go worker(i, highPriority, lowPriority, &wg)
	}

	// Flood both channels at the same time
	go func() {
		for i := 1; i <= 10; i++ {
			highPriority <- i + 100 // High tasks: 101-110
			lowPriority <- i        // Low tasks: 1-10
		}
		close(highPriority)
		close(lowPriority)
	}()

	wg.Wait()
	fmt.Println("Done.")
}
