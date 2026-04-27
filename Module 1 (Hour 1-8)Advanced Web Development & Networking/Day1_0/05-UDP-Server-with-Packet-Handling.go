package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	// Define the address and port to listen on
	address := ":8081"

	// Resolve the UDP address
	udpAddr, err := net.ResolveUDPAddr("udp", address)
	if err != nil {
		log.Fatalf("Failed to resolve UDP address: %v", err)
	}

	// Start listening for incoming UDP packets
	conn, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		log.Fatalf("Failed to start UDP server: %v", err)
	}
	defer conn.Close()

	log.Printf("UDP server started on %s", address)

	// Wait for interrupt signal to gracefully shut down
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

	// Buffer to store incoming packets
	buffer := make([]byte, 1024)

	// Start processing packets in a loop
	go func() {
		for {
			// Read a packet from the connection
			n, addr, err := conn.ReadFromUDP(buffer)
			if err != nil {
				select {
				case <-sig: // Check if shutdown signal was received
					log.Println("Stopping listener...")
					return
				default:
					log.Printf("Error reading packet: %v", err)
					continue
				}
			}

			// Extract the message from the buffer
			message := string(buffer[:n])
			log.Printf("Received from %s: %s", addr, message)

			// Process the packet (e.g., echo the message back)
			response := fmt.Sprintf("Echo: %s", message)
			_, err = conn.WriteToUDP([]byte(response), addr)
			if err != nil {
				log.Printf("Error sending response to %s: %v", addr, err)
			}
		}
	}()

	log.Println("Press Ctrl+C to stop the server...")
	<-sig

	log.Println("Shutting down server...")
	log.Println("Server stopped gracefully.")

}

/*
$udpClient = New-Object System.Net.Sockets.UdpClient
$udpClient.Connect("localhost", 8081)
$encodedMessage = [System.Text.Encoding]::ASCII.GetBytes("Hello, UDP server!")
$udpClient.Send($encodedMessage, $encodedMessage.Length)
$udpClient.Close()


If you prefer a direct command-line tool similar to nc, you can install Ncat (part of the Nmap project). 
Download: Available via the Nmap Download Page.
Usage: On Windows, the command is typically ncat instead of nc:
cmd
echo Hello, UDP server! | ncat -u localhost 8081

*/