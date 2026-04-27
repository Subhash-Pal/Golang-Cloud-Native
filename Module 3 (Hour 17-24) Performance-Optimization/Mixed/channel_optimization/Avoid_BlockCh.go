/*
2. Avoid Blocking Channels
Blocking channels can cause goroutines to wait unnecessarily. Use select with a default case to avoid blocking when sending or receiving data.

Optimization:
Using select with a default case prevents blocking when the channel is full or empty.
Useful in scenarios where you want to prioritize other work over waiting on a channel.

*/

package main

import (
	"fmt"
	"time"
)

func main() {
	ch := make(chan int, 1)

	// Non-blocking send
	select {
	case ch <- 42:
		fmt.Println("Sent value 42")
	default:
		fmt.Println("Channel is full, skipping send")
	}

	// Non-blocking receive
	select {
	case value := <-ch:
		fmt.Printf("Received value %d\n", value)
	default:
		fmt.Println("Channel is empty, skipping receive")
	}

	time.Sleep(1 * time.Second)
}