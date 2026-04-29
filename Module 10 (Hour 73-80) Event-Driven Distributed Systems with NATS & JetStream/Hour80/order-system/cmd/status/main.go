package main

import (
	"fmt"
	"log"

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

	keys, err := kv.Keys()
	if err != nil {
		log.Fatal(err)
	}

	for _, k := range keys {
		entry, err := kv.Get(k)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("%s → %s\n", k, string(entry.Value()))
	}
}