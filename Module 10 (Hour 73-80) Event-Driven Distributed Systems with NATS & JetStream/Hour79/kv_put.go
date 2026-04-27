package main

import (
	"log"

	"github.com/nats-io/nats.go"
)

func main() {
	nc, _ := nats.Connect(nats.DefaultURL)
	defer nc.Close()

	js, _ := nc.JetStream()

	kv, _ := js.KeyValue("CONFIG")

	_, err := kv.Put("service.timeout", []byte("30"))
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Key updated: service.timeout = 30")
}