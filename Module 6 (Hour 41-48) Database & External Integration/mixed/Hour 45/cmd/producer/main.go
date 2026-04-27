package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/lib/pq"

	"hour45-message-broker/broker"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	db, err := broker.OpenDB()
	if err != nil {
		log.Fatalf("open database: %v", err)
	}
	defer db.Close()

	if err := db.PingContext(ctx); err != nil {
		log.Fatalf("ping database: %v", err)
	}

	if err := broker.EnsureSchema(ctx, db); err != nil {
		log.Fatalf("prepare schema: %v", err)
	}

	payload := fmt.Sprintf("order-created-%d", time.Now().Unix())
	if len(os.Args) > 1 {
		payload = os.Args[1]
	}

	if err := broker.PublishMessage(ctx, db, broker.Channel, payload); err != nil {
		log.Fatalf("publish message: %v", err)
	}

	log.Printf("Published message to channel %q: %s", broker.Channel, payload)
}
