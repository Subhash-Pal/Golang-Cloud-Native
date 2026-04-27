package main

import (
	"fmt"
	"log"
	"time"

	"github.com/nats-io/nats.go" // Correct: no :// allowed here
)
/*
func main() {
	// Connect to your local NATS server (ensure nats-server -js is running)
	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		log.Fatal(err)
	}
	defer nc.Close()

	fmt.Println("Connected to NATS successfully!")

	// Simple Pub/Sub test
	nc.Subscribe("test", func(m *nats.Msg) {
		fmt.Printf("Received message: %s\n", string(m.Data))
	})

	nc.Publish("test", []byte("Hello from Go!"))

	// Keep alive briefly to see the output
	time.Sleep(500 * time.Millisecond)
}
	*/

	func main() {
	// Connect to NATS
	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		log.Fatal(err)
	}
	defer nc.Close()

	fmt.Println("Connected to NATS!")

	// Simple Pub/Sub Test
	nc.Subscribe("test", func(m *nats.Msg) {
		fmt.Printf("Received: %s\n", string(m.Data))
	})
	nc.Publish("test", []byte("Hello NATS"))

	time.Sleep(500 * time.Millisecond)
}
//How to run above code ?

