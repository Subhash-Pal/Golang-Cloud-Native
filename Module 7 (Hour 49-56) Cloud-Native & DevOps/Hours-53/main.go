package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
)

func getenv(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func main() {
	port := getenv("PORT", "8080")
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusOK, map[string]any{
			"hour":  53,
			"topic": "Writing Deployment and Service YAML",
			"config": map[string]string{
				"app_mode":  getenv("APP_MODE", "development"),
				"log_level": getenv("LOG_LEVEL", "debug"),
			},
		})
	})
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
	})
	mux.HandleFunc("/config", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusOK, map[string]string{
			"app_mode":  getenv("APP_MODE", "development"),
			"log_level": getenv("LOG_LEVEL", "debug"),
		})
	})

	addr := ":" + port
	log.Printf("hour 53 app listening on %s", addr)
	log.Fatal(http.ListenAndServe(addr, mux))
}
