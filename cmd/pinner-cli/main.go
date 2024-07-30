// cli/cli.go
package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"

	"github.com/covalenthq/das-ipfs-pinner/common"
)

func main() {
	mode := flag.String("mode", "", "Mode of operation: store or extract")
	data := flag.String("data", "", "Path to the binary file to send to the daemon or CID for extraction")
	addr := flag.String("addr", getEnv("DAEMON_ADDR", "http://localhost:5080"), "Address of the daemon")
	flag.Parse()

	switch *mode {
	case "store":
		if *data == "" {
			fmt.Println("File path is required for store mode")
			os.Exit(1)
		}
		storeData(*addr, *data)
	case "extract":
		if *data == "" {
			fmt.Println("CID is required for extract mode")
			os.Exit(1)
		}
		extractData(*addr, *data)
	default:
		fmt.Printf("Invalid mode. Use %s -mode=store -data=path/to/file or %s -mode=extract -data=CID\n", common.BinaryName, common.BinaryName)
	}
}

func storeData(addr, filePath string) {
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Printf("Error opening file: %v\n", err)
		return
	}
	defer file.Close()

	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)

	part, err := writer.CreateFormFile("file", filepath.Base(filePath))
	if err != nil {
		fmt.Printf("Error creating form file: %v\n", err)
		return
	}

	_, err = io.Copy(part, file)
	if err != nil {
		fmt.Printf("Error copying file content: %v\n", err)
		return
	}

	err = writer.Close()
	if err != nil {
		fmt.Printf("Error closing writer: %v\n", err)
		return
	}

	req, err := http.NewRequest("POST", addr+"/store", &buf)
	if err != nil {
		fmt.Printf("Error creating request: %v\n", err)
		return
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Error sending request: %v\n", err)
		return
	}
	defer resp.Body.Close()

	response, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Error reading response: %v\n", err)
		return
	}

	fmt.Println("Response from server:", string(response))
}

func extractData(addr, cid string) {
	resp, err := http.Get(fmt.Sprintf("%s/extract?cid=%s", addr, cid))
	if err != nil {
		fmt.Printf("Error extracting data: %v\n", err)
		return
	}
	defer resp.Body.Close()

	response, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Error reading response: %v\n", err)
		return
	}
	fmt.Println(string(response))
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
