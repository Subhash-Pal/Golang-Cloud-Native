package main

import (
	"fmt"
	"sync"
)

// Example 1: Map with channels as values
func mapWithChannels() {
	fmt.Println("\n=== Map with Channels ===")

	// Create a map where values are channels
	userChannels := make(map[string]chan string)

	// Initialize channels for each user
	userChannels["alice"] = make(chan string, 5)
	userChannels["bob"] = make(chan string, 5)

	// Send messages to channels
	userChannels["alice"] <- "Hello Alice!"
	userChannels["bob"] <- "Hello Bob!"

	// Receive messages
	fmt.Println("Alice's message:", <-userChannels["alice"])
	fmt.Println("Bob's message:", <-userChannels["bob"])
}

// Example 2: Channel with map as value
func channelWithMap() {
	fmt.Println("\n=== Channel with Map ===")

	// Create a channel that carries maps
	dataChan := make(chan map[string]int, 3)

	// Send maps through channel
	go func() {
		dataChan <- map[string]int{"a": 1, "b": 2}
		dataChan <- map[string]int{"c": 3, "d": 4}
		close(dataChan)
	}()

	// Receive maps from channel
	for data := range dataChan {
		fmt.Println("Received map:", data)
	}
}

// Example 3: Sync map with goroutines
func syncMapExample() {
	fmt.Println("\n=== Sync.Map with Goroutines ===")

	var syncMap sync.Map
	var wg sync.WaitGroup

	// Multiple goroutines writing to sync.Map
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			syncMap.Store(fmt.Sprintf("key%d", id), id*10)
		}(i)
	}

	wg.Wait()

	// Read all values
	fmt.Println("Sync.Map contents:")
	syncMap.Range(func(key, value interface{}) bool {
		fmt.Printf("  %s: %d\n", key, value)
		return true
	})
}

// Example 4: Channel for map updates (producer-consumer)
func mapUpdateChannel() {
	fmt.Println("\n=== Channel for Map Updates ===")

	// Channel to receive map updates (using interface{} for mixed types)
	updateChan := make(chan map[string]interface{}, 10)
	resultMap := make(map[string]interface{})

	// Consumer goroutine
	go func() {
		for update := range updateChan {
			for k, v := range update {
				resultMap[k] = v
			}
		}
		fmt.Println("Final map:", resultMap)
	}()

	// Producer sends updates
	updateChan <- map[string]interface{}{"name": "Alice", "age": 25}
	updateChan <- map[string]interface{}{"city": "NYC", "score": 95}
	updateChan <- map[string]interface{}{"active": 1}

	close(updateChan)
}

// Example 5: Select with map and channel
func selectWithMapChannel() {
	fmt.Println("\n=== Select with Map and Channel ===")

	// Channel for signals
	signalChan := make(chan string, 1)
	
	// Map to store status
	status := make(map[string]bool)
	status["ready"] = false

	// Goroutine to update status
	go func() {
		signalChan <- "update"
		status["ready"] = true
	}()

	// Select statement
	select {
	case msg := <-signalChan:
		fmt.Println("Received signal:", msg)
		fmt.Println("Status ready:", status["ready"])
	case <-make(chan string): // Non-blocking default case
		fmt.Println("No signal received")
	}
}

func main() {
	mapWithChannels()
	channelWithMap()
	syncMapExample()
	mapUpdateChannel()
	selectWithMapChannel()
}