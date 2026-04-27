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

	_, err = js.AddStream(&nats.StreamConfig{
		Name:      "ORDERS",
		Subjects:  []string{"orders.*"},
		Retention: nats.WorkQueuePolicy, // 🔥 Change to InterestPolicy / LimitsPolicy to test
		MaxAge:    1 * time.Hour,
		MaxMsgs:   100,
	})
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Stream ORDERS created")
}