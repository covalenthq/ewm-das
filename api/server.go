package api

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	ipfsnode "github.com/covalenthq/das-ipfs-pinner/internal/ipfs-node"
)

const MaxMultipartMemory = 10 << 20 // 10 MB

type ServerConfig struct {
	Addr                  string
	W3AgentKey            string
	W3DelegationProofPath string
}

// StartServer initializes and starts the HTTP server.
func StartServer(config ServerConfig) {
	ipfsNode, err := ipfsnode.NewIPFSNode(config.W3AgentKey, config.W3DelegationProofPath)
	if err != nil {
		log.Fatalf("Failed to initialize IPFS node: %v", err)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/store", createStoreHandler(ipfsNode))
	mux.HandleFunc("/extract", extractHandler)

	server := &http.Server{
		Addr:    config.Addr,
		Handler: mux,
	}

	// Set up signal handling for graceful shutdown
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigs
		log.Println("Shutting down server...")

		// Create a context with timeout for the shutdown process
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := server.Shutdown(ctx); err != nil {
			log.Fatalf("Server Shutdown Failed: %+v", err)
		}
	}()

	log.Printf("Starting server on %s...\n", config.Addr)
	if err := server.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatalf("Could not start server: %v\n", err)
	}
}

func handleError(w http.ResponseWriter, errMsg string, statusCode int) {
	log.Printf("%s: %v", errMsg, statusCode)
	http.Error(w, errMsg, statusCode)
}
