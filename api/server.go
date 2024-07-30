package api

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

// StartServer initializes and starts the HTTP server.
func StartServer(addr string) {
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
	http.HandleFunc("/store", storeHandler)
	http.HandleFunc("/extract", extractHandler)

	// Start the HTTP server
	log.Printf("Starting daemon on %s...\n", addr)
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatalf("Could not start server: %v\n", err)
	}
}

func storeHandler(w http.ResponseWriter, r *http.Request) {
	// Handle storing data
	fmt.Fprintln(w, "Data stored")
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
