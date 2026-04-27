/*
Hour 49:
Docker Fundamentals (Hands-On)
Learning Objectives
Understand the basics of Docker and containerization.
Learn how to create, run, and manage Docker containers.
Build a simple Golang application and containerize it.
Hands-On Steps
Install Docker:
Download and install Docker Desktop for Windows from the official website.
Verify installation by running docker --version in your terminal.
Create a Simple Golang Application:
Write a basic "Hello, World!" HTTP server in Go:

*/

package main

import (
    "fmt"
    "net/http"
)

func handler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "Hello, Docker!")
}

func main() {
    http.HandleFunc("/", handler)
    http.ListenAndServe(":8080", nil)
}