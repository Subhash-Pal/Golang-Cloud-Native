package main

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

var ctx = context.Background()

func main() {
	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	sub := rdb.Subscribe(ctx, "user_events")

	fmt.Println("📡 Listening for events...")

	for {
		msg, err := sub.ReceiveMessage(ctx)
		if err != nil {
			panic(err)
		}

		fmt.Println("🔥 Event received:", msg.Payload)

		// Simulated business logic
		if msg.Payload == "user_logged_in:shubh" {
			fmt.Println("📧 Send login notification email")
		}
	}
}