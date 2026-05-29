package main

import (
	"fmt"
	"os"

	"github.com/covalenthq/das-ipfs-pinner/api"
	"github.com/covalenthq/das-ipfs-pinner/common"
	"github.com/covalenthq/das-ipfs-pinner/internal/pinner/das"
	ipfsnode "github.com/covalenthq/das-ipfs-pinner/internal/pinner/ipfs-node"
	logging "github.com/ipfs/go-log/v2"
	"github.com/spf13/cobra"
)

var greeting = `
██████   █████  ███████     ██████  ██ ███    ██ ███    ██ ███████ ██████
██   ██ ██   ██ ██          ██   ██ ██ ████   ██ ████   ██ ██      ██   ██
██   ██ ███████ ███████     ██████  ██ ██ ██  ██ ██ ██  ██ █████   ██████
██   ██ ██   ██      ██     ██      ██ ██  ██ ██ ██  ██ ██ ██      ██   ██
██████  ██   ██ ███████     ██      ██ ██   ████ ██   ████ ███████ ██   ██


`

var log = logging.Logger("das-pinner") // Initialize the logger

var (
	addr     string
	logLevel string
)

// rootCmd represents the base command
var rootCmd = &cobra.Command{
	Use:     common.BinaryName,
	Short:   "A service for storing and retrieving DAS-data, backed by IPFS",
	Long:    `Pinner is a service that handles storing binary data and extracting it via HTTP. It is backed by IPFS and uses KZG commitments for data integrity.`,
	Version: fmt.Sprintf("%s, commit %s", common.Version, common.GitCommit),
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		logging.SetLogLevel("das-pinner", logLevel)

		// Load the configuration
		config := das.LoadConfig()
		// Initialize the KZG trusted setup
		log.Info("Initializing trusted setup...")
		if err := das.InitializeTrustedSetup(config); err != nil {
			log.Fatalf("Failed to initialize trusted setup: %v", err)
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		ipfsURL := os.Getenv("IPFS_RPC_URL")
		ipfsToken := os.Getenv("IPFS_RPC_TOKEN")
		fbToken := os.Getenv("FILEBASE_RPC_TOKEN")

		if ipfsURL == "" && fbToken == "" {
			log.Fatalf("must set IPFS_RPC_URL (self-hosted IPFS) or FILEBASE_RPC_TOKEN (Filebase)")
		}

		config := api.ServerConfig{
			Addr: addr,
			Pinner: ipfsnode.PinnerConfig{
				IPFSRPC:  ipfsnode.IPFSRPCConfig{RPCURL: ipfsURL, RPCToken: ipfsToken},
				Filebase: ipfsnode.FilebaseConfig{RPCToken: fbToken},
				Gateway:  os.Getenv("DEDICATED_GATEWAY"),
			},
		}
		api.StartServer(config)
	},
}

func main() {
	fmt.Print(greeting)
	fmt.Printf("Version: %s, commit: %s\n", common.Version, common.GitCommit)

	if err := rootCmd.Execute(); err != nil {
		log.Fatalf("Error executing command: %v", err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVar(&logLevel, "log-level", "info", "Log level")
	rootCmd.PersistentFlags().StringVar(&addr, "addr", getEnv("PINNER_ADDR", ":5080"),
		"Address to run the pinner service on")
}

func initConfig() {
	// Additional initialization if needed
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() error {
	return rootCmd.Execute()
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
