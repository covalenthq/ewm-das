package main

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path/filepath"

	"github.com/covalenthq/das-ipfs-pinner/common"
	"github.com/spf13/cobra"
)

func main() {
	var addr string
	var data string

	rootCmd := &cobra.Command{
		Use:   common.BinaryName,
		Short: "CLI for interacting with the DAS IPFS Pinner daemon",
	}

	// Set addr flag as a persistent flag so it can be used across all commands
	rootCmd.PersistentFlags().StringVarP(&addr, "addr", "a", getEnv("DAEMON_ADDR", ":5080"), "Address of the daemon")

	uploadCmd := &cobra.Command{
		Use:   "upload",
		Short: "Upload data request to the daemon",
		Run: func(cmd *cobra.Command, args []string) {
			if data == "" {
				fmt.Println("File path is required for upload mode")
				os.Exit(1)
			}
			uploadData(addr, data)
		},
	}

	uploadCmd.Flags().StringVarP(&data, "data", "d", "", "Path to the binary file to send to the daemon")
	rootCmd.AddCommand(uploadCmd)

	downloadCmd := &cobra.Command{
		Use:   "download",
		Short: "Download data from the daemon using a CID",
		Run: func(cmd *cobra.Command, args []string) {
			if data == "" {
				fmt.Println("CID is required for download mode")
				os.Exit(1)
			}
			downloadData(addr, data)
		},
	}

	downloadCmd.Flags().StringVarP(&data, "data", "d", "", "CID for download")
	rootCmd.AddCommand(downloadCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func uploadData(addr, filePath string) {
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

	req, err := http.NewRequest("POST", addr+"/upload", &buf)
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

func downloadData(addr, cid string) {
	// Create a form with the CID
	formData := url.Values{}
	formData.Set("cid", cid)

	// Create a POST request with the form data
	resp, err := http.PostForm(fmt.Sprintf("%s/get", addr), formData)
	if err != nil {
		fmt.Printf("Error downloading data: %v\n", err)
		return
	}
	defer resp.Body.Close()

	// Read the response
	response, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Error reading response: %v\n", err)
		return
	}

	// Print the response
	fmt.Println(string(response))
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
