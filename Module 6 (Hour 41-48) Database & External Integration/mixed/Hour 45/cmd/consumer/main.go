package main

import (
	"context"
	"log"
	"time"

	"github.com/lib/pq"
	_ "github.com/lib/pq"

	"hour45-message-broker/broker"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
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

	listener := pq.NewListener(broker.ConnStr, 10*time.Second, time.Minute, nil)
	defer listener.Close()

	if err := listener.Listen(broker.Channel); err != nil {
		log.Fatalf("listen on channel %q: %v", broker.Channel, err)
	}

	log.Printf("Consumer is listening on channel %q", broker.Channel)

	notification, err := broker.WaitForNotification(listener, 20*time.Second)
	if err != nil {
		log.Fatalf("receive message: %v", err)
	}

	if err := broker.StoreMessage(ctx, db, notification.Channel, notification.Extra); err != nil {
		log.Fatalf("store received message: %v", err)
	}

	log.Printf("Consumed and stored message from %q: %s", notification.Channel, notification.Extra)
}
