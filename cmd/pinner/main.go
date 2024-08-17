package main

import (
	"os"

	logging "github.com/ipfs/go-log/v2"
)

var log = logging.Logger("das-pinner") // Initialize the logger

func main() {
	if err := execute(); err != nil {
		log.Fatalf("Error executing command: %v", err)
		os.Exit(1)
	}
}
