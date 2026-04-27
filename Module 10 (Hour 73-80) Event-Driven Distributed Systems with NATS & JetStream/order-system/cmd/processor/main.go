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

	kv, err := js.KeyValue("ORDER_STATUS")
	if err != nil {
		log.Fatal(err)
	}

	_, err = js.Subscribe("orders.created", func(msg *nats.Msg) {
		orderID := string(msg.Data)

		fmt.Println("Processing:", orderID)

		// Simulate processing
		time.Sleep(1 * time.Second)

		// Update KV state
		_, err := kv.Put(orderID, []byte("PROCESSED"))
		if err != nil {
			log.Println("KV error:", err)
			return
		}

		// ACK message
		msg.Ack()

		fmt.Println("Completed:", orderID)

	},
		nats.Durable("ORDER_WORKER"),
		nats.ManualAck(),
	)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Processor running...")
	select {}
}