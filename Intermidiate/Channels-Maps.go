package main

import (
	"fmt"
	"sync"
)
//<-ch chan 
func processOrders(orderMap map[string]int, ch chan<- string, wg *sync.WaitGroup) {
	defer wg.Done()

	for product, quantity := range orderMap {
		ch <- fmt.Sprintf("Processed %d units of %s", quantity, product)
	}
}

func main() {
	orderMap := map[string]int{
		"Apples":  10,
		"Bananas": 5,
		"Grapes":  7,
	}

	ch := make(chan string,1) // Channel to receive processed messages
	var wg sync.WaitGroup   // WaitGroup to wait for goroutines


	// Read messages from the channel
	go func() {
		for msg := range ch {
			fmt.Println(msg)
		}
	}()

	wg.Add(1)
	go processOrders(orderMap, ch, &wg)

	

	wg.Wait()       // Wait for the goroutine to finish
	close(ch)       // Close the channel after processing
}