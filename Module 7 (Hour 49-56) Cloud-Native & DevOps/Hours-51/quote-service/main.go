package main

import (
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
)

var quotes = []string{
	"Containers standardize the runtime environment.",
	"Compose lets you boot an application stack with one command.",
	"DevOps improves feedback loops between code and operations.",
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/quote", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]string{
			"quote": quotes[rand.Intn(len(quotes))],
		})
	})
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
	})

	log.Println("hour 51 quote-service listening on :8081")
	log.Fatal(http.ListenAndServe(":8081", mux))
}
