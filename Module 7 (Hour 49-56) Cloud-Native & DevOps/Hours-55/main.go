package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
)

func newMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]any{
			"hour":  55,
			"topic": "CI/CD pipeline setup",
			"checks": []string{
				"go test",
				"go build",
				"docker build",
			},
		})
	})
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})
	return mux
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	addr := ":" + port
	log.Printf("hour 55 app listening on %s", addr)
	log.Fatal(http.ListenAndServe(addr, newMux()))
}
