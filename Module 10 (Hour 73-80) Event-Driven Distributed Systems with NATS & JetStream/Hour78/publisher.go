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

	for i := 1; i <= 10; i++ {
		msg := fmt.Sprintf("order-%d", i)

		_, err := js.Publish("orders.created", []byte(msg))
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println("Published:", msg)
		time.Sleep(500 * time.Millisecond)
	}
}