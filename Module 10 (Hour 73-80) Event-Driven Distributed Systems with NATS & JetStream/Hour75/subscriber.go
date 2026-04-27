package main

import (
	"fmt"
	"log"

	"github.com/nats-io/nats.go"
)

func main() {
	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		log.Fatal(err)
	}
	defer nc.Close()

	fmt.Println("Subscriber listening on 'orders.created'...")

	_, err = nc.Subscribe("orders.created", func(msg *nats.Msg) {
		fmt.Printf("[Subscriber] Received: %s\n", string(msg.Data))
	})
	if err != nil {
		log.Fatal(err)
	}

	select {} // keep running
}