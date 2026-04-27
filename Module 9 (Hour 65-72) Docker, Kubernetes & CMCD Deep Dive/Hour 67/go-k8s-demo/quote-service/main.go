package main

import (
	"encoding/json"
	"net/http"
	"time"
)

func quoteHandler(w http.ResponseWriter, r *http.Request) {
	response := map[string]interface{}{
		"quote":     "Kubernetes brings orchestration to your containers.",
		"timestamp": time.Now().UTC(),
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func main() {
	http.HandleFunc("/quote", quoteHandler)
	http.ListenAndServe(":8081", nil)
}