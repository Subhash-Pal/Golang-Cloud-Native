package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"module6/hour45/broker"

	_ "github.com/lib/pq"
)

func main() {
	payload := fmt.Sprintf("order-created-%d", time.Now().Unix())
	if len(os.Args) > 1 {
		payload = os.Args[1]
	}

	db, err := broker.OpenDB()
	if err != nil {
		log.Fatalf("open database: %v", err)
	}
	defer db.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		log.Fatalf("ping database: %v", err)
	}

	if err := broker.PublishMessage(ctx, db, broker.Channel, payload); err != nil {
		log.Fatalf("publish message: %v", err)
	}

	log.Printf("Published message to channel %q: %s", broker.Channel, payload)
}
