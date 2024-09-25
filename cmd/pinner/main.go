package main

import (
	"fmt"
	"os"

	"github.com/covalenthq/das-ipfs-pinner/api"
	"github.com/covalenthq/das-ipfs-pinner/common"
	"github.com/covalenthq/das-ipfs-pinner/internal/pinner/das"
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
	addr                  string
	logLevel              string
	w3AgentKey            string
	w3DelegationProofPath string
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
		// Populate the ServerConfig struct
		config := api.ServerConfig{
			Addr:                  addr,
			W3AgentKey:            w3AgentKey,
			W3DelegationProofPath: w3DelegationProofPath,
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
	rootCmd.PersistentFlags().StringVar(&addr, "addr", getEnv("PINNER_ADDR", ":5080"), "Address to run the pinner service on")

	// W3 agent flags
	rootCmd.PersistentFlags().StringVar(&w3AgentKey, "w3-agent-key", "", "Key for the W3 agent")
	rootCmd.PersistentFlags().StringVar(&w3DelegationProofPath, "w3-delegation-proof-path", "", "Path to the W3 delegation proof")

	// Mark the flags as required
	rootCmd.MarkPersistentFlagRequired("w3-agent-key")
	rootCmd.MarkPersistentFlagRequired("w3-delegation-proof-path")
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
