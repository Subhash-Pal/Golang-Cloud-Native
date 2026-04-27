package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

func handler(w http.ResponseWriter, r *http.Request) {
	url := os.Getenv("QUOTE_SERVICE_URL")

	resp, err := http.Get(url)
	if err != nil {
		http.Error(w, "error calling quote service", 500)
		return
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)

	var data map[string]interface{}
	json.Unmarshal(body, &data)

	response := map[string]interface{}{
		"source":    "api",
		"quote":     data,
		"timestamp": time.Now().UTC(),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func main() {
	http.HandleFunc("/", handler)


	err := http.ListenAndServe(":8080", nil)
if err != nil {
    panic(err)  // forces visible failure
}

}