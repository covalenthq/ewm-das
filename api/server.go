package api

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/covalenthq/das-ipfs-pinner/internal/das"
	ipfsnode "github.com/covalenthq/das-ipfs-pinner/internal/ipfs-node"
)

// StartServer initializes and starts the HTTP server.
func StartServer(addr string) {
	ipfsnode, err := ipfsnode.NewIPFSNode()
	if err != nil {
		log.Fatalf("Failed to initialize IPFS node: %v", err)
	}

	// Set up signal handling for graceful shutdown
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigs
		log.Println("Shutting down daemon...")
		// Perform cleanup here if needed
		os.Exit(0)
	}()

	// Set up HTTP handlers
	http.HandleFunc("/store", createStoreHandler(ipfsnode))
	http.HandleFunc("/extract", extractHandler)

	// Start the HTTP server
	log.Printf("Starting daemon on %s...\n", addr)
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatalf("Could not start server: %v\n", err)
	}
}

func storeHandler(w http.ResponseWriter, r *http.Request) {
	// Ensure that the request is a POST method
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST method is allowed", http.StatusMethodNotAllowed)
		return
	}

	// Read the binary data from the request body
	data, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	// Handle the binary data (e.g., save it to a file or a database)
	// For demonstration purposes, we'll just print the data length
	log.Printf("Received %d bytes of data\n", len(data))

	_, err = das.Encode(data)
	if err != nil {
		log.Printf("Failed to encode data: %v\n", err)
		http.Error(w, "Failed to store data", http.StatusInternalServerError)
		return
	}

	// Respond to the client
	fmt.Fprintln(w, "Data stored successfully")
}

func createStoreHandler(ipfsNode *ipfsnode.IPFSNode) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Ensure that the request is a POST method
		if r.Method != http.MethodPost {
			http.Error(w, "Only POST method is allowed", http.StatusMethodNotAllowed)
			return
		}

		// Read the binary data from the request body
		data, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Failed to read request body", http.StatusInternalServerError)
			return
		}
		defer r.Body.Close()

		// Handle the binary data (e.g., encode and store it using IPFSNode)
		log.Printf("Received %d bytes of data\n", len(data))

		block, err := das.Encode(data)
		if err != nil {
			log.Printf("Failed to encode data: %v\n", err)
			http.Error(w, "Failed to store data", http.StatusInternalServerError)
			return
		}

		// Store the encoded block to IPFS
		err = ipfsNode.PublishBlock(block)
		if err != nil {
			log.Printf("Failed to store data to IPFS: %v\n", err)
			http.Error(w, "Failed to store data", http.StatusInternalServerError)
			return
		}

		// Respond to the client
		fmt.Fprintln(w, "Data stored successfully")
	}
}

func extractHandler(w http.ResponseWriter, r *http.Request) {
	// Handle extracting data
	cid := r.URL.Query().Get("cid")
	if cid == "" {
		http.Error(w, "CID is required", http.StatusBadRequest)
		return
	}
	fmt.Fprintf(w, "Extracting data for CID: %s\n", cid)
}
