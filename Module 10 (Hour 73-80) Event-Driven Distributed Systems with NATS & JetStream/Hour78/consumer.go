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

	_, err = js.Subscribe("orders.created", func(msg *nats.Msg) {
		fmt.Println("Processing:", string(msg.Data))

		// Simulate processing delay
		time.Sleep(1 * time.Second)

		// ACK is critical
		msg.Ack()
	},
		nats.Durable("ORDER_PROCESSOR"),
		nats.ManualAck(),
	)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Durable consumer running...")
	select {}
}