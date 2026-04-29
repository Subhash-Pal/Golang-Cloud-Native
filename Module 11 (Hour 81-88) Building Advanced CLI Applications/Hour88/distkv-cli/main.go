package main

import (
	"log"

	"github.com/Subhash-Pal/distkv-cli/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
