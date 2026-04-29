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

	js, _ := nc.JetStream()

	for i := 1; i <= 5; i++ {
		orderID := fmt.Sprintf("order-%d", i)

		_, err := js.Publish("orders.created", []byte(orderID))
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println("Created:", orderID)
		time.Sleep(500 * time.Millisecond)
	}
}