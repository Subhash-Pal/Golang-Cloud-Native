package main

import (
	"api-gateway/config"
	"api-gateway/gateway"
)

func main() {

	cfg := config.LoadConfig()

	r := gateway.SetupRouter(cfg)

	r.Run(":8080")
}