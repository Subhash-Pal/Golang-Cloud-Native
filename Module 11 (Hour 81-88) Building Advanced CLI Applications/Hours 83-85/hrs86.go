package main

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

// Event structure (production-style)
type Event struct {
	Type      string            `json:"type"`
	Timestamp string            `json:"timestamp"`
	Payload   map[string]string `json:"payload"`
}

// Parse key=value,key=value → map
func parseData(input string) map[string]string {
	data := make(map[string]string)

	pairs := strings.Split(input, ",")
	for _, p := range pairs {
		kv := strings.SplitN(p, "=", 2)
		if len(kv) == 2 {
			data[kv[0]] = kv[1]
		}
	}
	return data
}

// Mock publisher (replace later with Kafka/RabbitMQ)
func publishEvent(event Event) {
	bytes, err := json.MarshalIndent(event, "", "  ")
	if err != nil {
		log.Fatal("Failed to serialize event:", err)
	}

	fmt.Println("📤 Event Published:")
	fmt.Println(string(bytes))
}

var eventType string
var eventData string

func main() {

	var rootCmd = &cobra.Command{
		Use: "app",
	}

	var eventCmd = &cobra.Command{
		Use:   "event",
		Short: "Event operations",
	}

	var publishCmd = &cobra.Command{
		Use:   "publish",
		Short: "Publish an event",
		Run: func(cmd *cobra.Command, args []string) {

			if eventType == "" {
				log.Fatal("Missing --type")
			}

			payload := parseData(eventData)

			event := Event{
				Type:      eventType,
				Timestamp: time.Now().Format(time.RFC3339),
				Payload:   payload,
			}

			publishEvent(event)
		},
	}

	// Flags
	publishCmd.Flags().StringVar(&eventType, "type", "", "Event type")
	publishCmd.Flags().StringVar(&eventData, "data", "", "Event data (key=value,key=value)")

	// Hierarchy
	rootCmd.AddCommand(eventCmd)
	eventCmd.AddCommand(publishCmd)

	rootCmd.Execute()
}
//go run hrs86.go event publish --type user_created --data "id=1,name=Shubh"