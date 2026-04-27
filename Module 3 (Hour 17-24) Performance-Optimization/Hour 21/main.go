package main

import (
	"context"
	"log/slog"
	"time"
)

func main() {
	const totalMessages = 50_000

	unbuffered := timePipeline(totalMessages, false)
	buffered := timePipeline(totalMessages, true)

	slog.Info("channel optimization comparison",
		"messages", totalMessages,
		"unbuffered", unbuffered,
		"buffered_batching", buffered,
	)
}

func timePipeline(total int, optimized bool) time.Duration {
	start := time.Now()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	bufferSize := 0
	if optimized {
		bufferSize = 256
	}

	input := make(chan int, bufferSize)
	output := make(chan int, bufferSize)

	go producer(ctx, input, total)
	go processor(ctx, input, output)

	sum := 0
	for i := 0; i < total; i++ {
		sum += <-output
	}

	if sum == 0 {
		slog.Warn("unexpected zero sum")
	}

	return time.Since(start)
}

func producer(ctx context.Context, out chan<- int, total int) {
	defer close(out)
	for i := 0; i < total; i++ {
		select {
		case <-ctx.Done():
			return
		case out <- i: //write 
		}
	}
}

func processor(ctx context.Context, in <-chan int, out chan<- int) {
	defer close(out)
	for {
		select {
		case <-ctx.Done():
			return
		case value, ok := <-in:
			if !ok {
				return
			}
			out <- value * 2
		}
	}
}
