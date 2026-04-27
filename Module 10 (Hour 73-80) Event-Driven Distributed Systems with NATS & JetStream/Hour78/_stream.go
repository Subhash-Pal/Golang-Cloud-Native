package main

import (
	"log"
	"time"

	"github.com/nats-io/nats.go"
)

func main() {
	// Step 1: Connect
	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		log.Fatal("Connection error:", err)
	}
	defer nc.Close()

	js, err := nc.JetStream()
	if err != nil {
		log.Fatal("JetStream error:", err)
	}

	streamName := "ORDERS"

	// Step 2: Delete existing stream (reset-safe)
	_ = js.DeleteStream(streamName) // ignore error if not exists

	log.Println("Old stream (if any) deleted")

	// Step 3: Create stream with configurable retention
	cfg := &nats.StreamConfig{
		Name:     streamName,
		Subjects: []string{"orders.*"},

		// 🔥 CHANGE THIS LINE TO TEST DIFFERENT POLICIES
		Retention: nats.WorkQueuePolicy,
		// Retention: nats.InterestPolicy,
		// Retention: nats.LimitsPolicy,

		MaxAge:  1 * time.Hour,
		MaxMsgs: 100,
	}

	_, err = js.AddStream(cfg)
	if err != nil {
		log.Fatal("Stream creation failed:", err)
	}

	log.Println("Stream ORDERS created with config:")
	log.Println("Retention:", cfg.Retention)
	log.Println("MaxMsgs:", cfg.MaxMsgs)
	log.Println("MaxAge:", cfg.MaxAge)
}