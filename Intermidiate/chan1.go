package main

import (
	"fmt"
	"time"
)

func sendData(ch chan<- int) {
	for i := 1; i <= 5; i++ {
		ch <- i // Send data to the channel
		time.Sleep(5 * time.Millisecond)
	}
	close(ch) // Close the channel when done
}

func main() {
	ch := make(chan int) // Create a channel

	go sendData(ch) // Launch a goroutine to send data

	// Receive data from the channel
	for val := range ch {
		fmt.Println("Received:", val)
	}
}