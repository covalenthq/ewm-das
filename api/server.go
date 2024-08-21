package api

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	ipfsnode "github.com/covalenthq/das-ipfs-pinner/internal/pinner/ipfs-node"
	logging "github.com/ipfs/go-log/v2"
)

var log = logging.Logger("das-pinner") // Initialize the logger

// MaxMultipartMemory is the maximum memory that the server will use to parse multipart form data.
const MaxMultipartMemory = 10 << 20 // 10 MB

// ServerConfig contains the configuration for the HTTP server.
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
	mux.HandleFunc("/api/v1/upload", createUploadHandler(ipfsNode))
	mux.HandleFunc("/api/v1/download", createDownloadHandler(ipfsNode))

	// Deprecated endpoints - same behavior, with deprecation notice in headers
	mux.HandleFunc("/upload", deprecatedHandler(createUploadHandler(ipfsNode), "/api/v1/upload"))
	mux.HandleFunc("/get", deprecatedHandler(createDownloadHandler(ipfsNode), "/api/v1/download"))

	server := &http.Server{
		Addr:    config.Addr,
		Handler: mux,
	}

	// Set up signal handling for graceful shutdown
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigs
		log.Info("Shutting down server...")

		// Create a context with timeout for the shutdown process
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := server.Shutdown(ctx); err != nil {
			log.Fatalf("Server Shutdown Failed: %+v", err)
		}
	}()

	log.Infof("Starting server on %s...", config.Addr)
	if err := server.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatalf("Could not start server: %v", err)
	}
}

func handleError(w http.ResponseWriter, errMsg string, statusCode int) {
	log.Infof("%s: %v", errMsg, statusCode)
	http.Error(w, errMsg, statusCode)
}
