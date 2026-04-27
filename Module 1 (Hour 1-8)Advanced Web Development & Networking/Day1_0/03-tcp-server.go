package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"sync"
	"time"
)

func handleConn(conn net.Conn, id int) {
	defer conn.Close()
	log.Printf("[TCP-%d] New connection from %s", id, conn.RemoteAddr())

	scanner := bufio.NewScanner(conn)// waiting client to send data, then read it line by line
	for scanner.Scan() {
		msg := scanner.Text()
		log.Printf("[TCP-%d] Received: %s", id, msg)
		response := fmt.Sprintf("ECHO: %s | Server time: %s\n", msg, time.Now().Format(time.RFC3339))
		response += "From TCP server\n"
		conn.Write([]byte(response))
	}

	if err := scanner.Err(); err != nil {
		log.Printf("[TCP-%d] Error: %v", id, err)
	}
	log.Printf("[TCP-%d] Connection closed", id)
}

func main() {
	listener, err := net.Listen("tcp", ":9000")
	if err != nil {
		log.Fatal(err)
	}
	defer listener.Close()

	log.Println("TCP server listening on :9000 (concurrent connections)")

	var connID int
	var wg sync.WaitGroup

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Accept error: %v", err)
			continue
		}

		connID++
		wg.Add(1)
		go func(c net.Conn, id int) {
			defer wg.Done()
			handleConn(c, id)
		}(conn, connID)
	}
}