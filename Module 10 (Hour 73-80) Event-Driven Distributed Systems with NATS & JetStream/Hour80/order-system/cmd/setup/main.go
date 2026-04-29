package main

import (
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

	js, err := nc.JetStream()
	if err != nil {
		log.Fatal(err)
	}

	// Reset system
	_ = js.DeleteStream("ORDERS")

	// Create Stream
	_, err = js.AddStream(&nats.StreamConfig{
		Name:      "ORDERS",
		Subjects:  []string{"orders.created"},
		Retention: nats.WorkQueuePolicy,
		MaxAge:    1 * time.Hour,
	})
	if err != nil {
		log.Fatal(err)
	}

	// Create KV bucket
	_, err = js.CreateKeyValue(&nats.KeyValueConfig{
		Bucket: "ORDER_STATUS",
	})
	if err != nil && err.Error() != "bucket name already in use" {
		log.Fatal(err)
	}

	log.Println("Setup completed")
}