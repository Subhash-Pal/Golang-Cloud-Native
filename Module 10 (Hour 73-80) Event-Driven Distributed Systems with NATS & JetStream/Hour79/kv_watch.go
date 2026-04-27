package main

import (
	"fmt"
	"log"

	"github.com/nats-io/nats.go"
)

func main() {
	nc, _ := nats.Connect(nats.DefaultURL)
	defer nc.Close()

	js, _ := nc.JetStream()

	kv, _ := js.KeyValue("CONFIG")

	watcher, err := kv.Watch("service.timeout")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Watching for changes...")

	for update := range watcher.Updates() {
		if update != nil {
			fmt.Printf("Updated: %s = %s (rev=%d)\n",
				update.Key(),
				string(update.Value()),
				update.Revision())
		}
	}
}