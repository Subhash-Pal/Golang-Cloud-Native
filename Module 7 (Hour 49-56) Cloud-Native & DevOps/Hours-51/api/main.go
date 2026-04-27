package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"
)

type quotePayload struct {
	Quote string `json:"quote"`
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func main() {
	quoteURL := os.Getenv("QUOTE_SERVICE_URL")
	if quoteURL == "" {
		quoteURL = "http://localhost:8081/quote"
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		resp, err := http.Get(quoteURL)
		if err != nil {
			writeJSON(w, http.StatusBadGateway, map[string]string{"error": "unable to reach quote-service"})
			return
		}
		defer resp.Body.Close()

		var quote quotePayload
		if err := json.NewDecoder(resp.Body).Decode(&quote); err != nil {
			writeJSON(w, http.StatusBadGateway, map[string]string{"error": "invalid response from quote-service"})
			return
		}

		writeJSON(w, http.StatusOK, map[string]any{
			"hour":          51,
			"topic":         "Docker Compose setup",
			"quote":         quote.Quote,
			"quote_service": quoteURL,
			"timestamp":     time.Now().Format(time.RFC3339),
		})
	})
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
	})

	log.Println("hour 51 api listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
