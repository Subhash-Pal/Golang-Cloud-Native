package main

import (
	"fmt"
	"net/http"
)
//
func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, Optimized Production Docker Build!")
}

func main() {
	http.HandleFunc("/", handler)
	http.ListenAndServe(":8080", nil)
}

//docker build -t hour66-go-secure .
//docker run --rm -d --name hour66-go-secure-app -p 8081:8080 --read-only --cap-drop=ALL --security-opt no-new-privileges:true hour66-go-secure
//docker stop hour66-go-secure-app 