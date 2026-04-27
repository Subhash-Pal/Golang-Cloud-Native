package main

import (
	"fmt"
	"log"
	"time"

	"github.com/nats-io/nats.go"
)

func main() {
	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		log.Fatal(err)
	}
	defer nc.Close()

	for i := 1; i <= 5; i++ {
		message := fmt.Sprintf("Order #%d created", i)

		err := nc.Publish("orders.created", []byte(message))
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println("[Publisher] Sent:", message)
		time.Sleep(1 * time.Second)
	}
}