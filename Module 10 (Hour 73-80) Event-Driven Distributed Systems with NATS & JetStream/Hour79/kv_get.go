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

	entry, err := kv.Get("service.timeout")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Value:", string(entry.Value()))
	fmt.Println("Revision:", entry.Revision())
}