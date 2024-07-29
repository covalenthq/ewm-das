// cli/cli.go
package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/covalenthq/das-ipfs-pinner/common"
)

func main() {
	mode := flag.String("mode", "", "Mode of operation: store or extract")
	data := flag.String("data", "", "Data to send to the daemon")
	addr := flag.String("addr", getEnv("DAEMON_ADDR", "http://localhost:8080"), "Address of the daemon")
	flag.Parse()

	switch *mode {
	case "store":
		if *data == "" {
			fmt.Println("Data is required for store mode")
			os.Exit(1)
		}
		storeData(*addr, *data)
	case "extract":
		extractData(*addr)
	default:
		fmt.Printf("Invalid mode. Use %s -mode=store or -mode=extract\n", common.BinaryName)
	}
}

func storeData(addr, data string) {
	resp, err := http.Post(addr+"/store", "application/json", bytes.NewReader([]byte(data)))
	if err != nil {
		fmt.Printf("Error storing data: %v\n", err)
		return
	}
	defer resp.Body.Close()
	response, _ := io.ReadAll(resp.Body)
	fmt.Println(string(response))
}

func extractData(addr string) {
	resp, err := http.Get(addr + "/extract")
	if err != nil {
		fmt.Printf("Error extracting data: %v\n", err)
		return
	}
	defer resp.Body.Close()
	response, _ := io.ReadAll(resp.Body)
	fmt.Println(string(response))
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
