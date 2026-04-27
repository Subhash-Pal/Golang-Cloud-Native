package main

import (
	"context"
	"log"
	"time"

	"module6/hour45/broker"

	"github.com/lib/pq"
	_ "github.com/lib/pq"
)

func main() {
	db, err := broker.OpenDB()
	if err != nil {
		log.Fatalf("open database: %v", err)
	}
	defer db.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		log.Fatalf("ping database: %v", err)
	}

	if err := broker.EnsureSchema(ctx, db); err != nil {
		log.Fatalf("ensure schema: %v", err)
	}

	listener := pq.NewListener(broker.ConnString(), 10*time.Second, time.Minute, nil)
	defer listener.Close()

	if err := listener.Listen(broker.Channel); err != nil {
		log.Fatalf("listen on channel: %v", err)
	}

	log.Printf("Consumer is listening on channel %q", broker.Channel)

	notification, err := broker.WaitForNotification(listener, 20*time.Second)
	if err != nil {
		log.Fatalf("wait for notification: %v", err)
	}

	storeCtx, storeCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer storeCancel()

	if err := broker.StoreMessage(storeCtx, db, notification.Channel, notification.Extra); err != nil {
		log.Fatalf("store message: %v", err)
	}

	log.Printf("Consumed and stored message from %q: %s", notification.Channel, notification.Extra)
}
