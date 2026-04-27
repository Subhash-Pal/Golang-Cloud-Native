package main
import (
	"log"

	"github.com/nats-io/nats.go"
)
//what is this above


func main() {
	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		log.Fatal(err)
	}
	defer nc.Close()

	js, _ := nc.JetStream()

	kv, err := js.CreateKeyValue(&nats.KeyValueConfig{
		Bucket: "CONFIG",
	})
	if err != nil {
		log.Fatal(err)
	}

	log.Println("KV Bucket created:", kv.Bucket())
}
//Explain me above code 



