package main

import (
	"log"
	"os"

	"github.com/covalenthq/das-ipfs-pinner/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		log.Fatalf("Error executing command: %v", err)
		os.Exit(1)
	}
}
