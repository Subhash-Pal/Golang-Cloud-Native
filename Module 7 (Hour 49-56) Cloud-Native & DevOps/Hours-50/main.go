package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
)

var (
	version = "dev"
	commit  = "local"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		_ = json.NewEncoder(w).Encode(map[string]string{
			"hour":        "50",
			"topic":       "Multi-stage Docker builds",
			"description": "This image is built in one stage and run from a smaller runtime stage.",
		})
	})
	mux.HandleFunc("/build", func(w http.ResponseWriter, r *http.Request) {
		_ = json.NewEncoder(w).Encode(map[string]string{
			"version": version,
			"commit":  commit,
		})
	})
	mux.HandleFunc("/text", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "Hello from the multi-stage container build.")
	})

	addr := ":" + port
	log.Printf("hour 50 server listening on %s", addr)
	log.Fatal(http.ListenAndServe(addr, mux))
}
