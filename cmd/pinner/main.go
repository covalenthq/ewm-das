package main

import (
	"fmt"
	"os"
	"path/filepath"
	"syscall"

	"github.com/covalenthq/das-ipfs-pinner/api"
	"github.com/covalenthq/das-ipfs-pinner/common"
	"github.com/covalenthq/das-ipfs-pinner/internal/pinner/das"
	logging "github.com/ipfs/go-log/v2"
	"github.com/spf13/cobra"
)

var log = logging.Logger("das-pinner") // Initialize the logger

var (
	logLevel              string
	detached              bool
	addr                  string
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

		// Handle the debug flag or daemonize if not in debug mode
		if detached {
			if os.Getenv("GO_DETACHED") != "1" {
				daemonize()
			}
		}

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
	log.Info("Initializing root command...", os.Args)

	if err := rootCmd.Execute(); err != nil {
		log.Fatalf("Error executing command: %v", err)
		os.Exit(1)
	}
}

func init() {
	log.Infof("Version: %s, commit: %s", common.Version, common.GitCommit)

	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVar(&logLevel, "log-level", "info", "Log level")
	rootCmd.PersistentFlags().BoolVar(&detached, "detached", false, "Run in detached mode")
	rootCmd.PersistentFlags().StringVar(&addr, "addr", getEnv("PINNER_ADDR", "localhost:5080"), "Address to run the pinner service on")

	// W3 agent flags
	rootCmd.PersistentFlags().StringVar(&w3AgentKey, "w3-agent-key", "", "Key for the W3 agent")
	rootCmd.PersistentFlags().StringVar(&w3DelegationProofPath, "w3-delegation-proof-path", "", "Path to the W3 delegation proof")

	// Mark the flags as required
	rootCmd.MarkPersistentFlagRequired("w3-agent-key")
	rootCmd.MarkPersistentFlagRequired("w3-delegation-proof-path")

	// Check if we're in the child process (daemon)
	if os.Getenv("GO_DETACHED") == "1" {
		log.Info("Running in daemon mode.")
		// Do not reinitialize flags or commands here
		// Proceed directly to running the server or minimal initialization required
		return
	}

	log.Info("Running in non-daemon mode.")
}

func initConfig() {
	// Additional initialization if needed
}

func daemonize() {
	// Get the executable path
	executablePath, err := filepath.Abs(os.Args[0])
	if err != nil {
		log.Fatalf("Error getting absolute path: %v\n", err)
	}

	// Ensure arguments are correctly passed to the child process
	var args []string

	if len(os.Args) > 1 {
		args = os.Args[1:]
	} else {
		// Log a message if there are no arguments
		log.Warn("No arguments provided to the child process.")
	}

	// Set up the environment with a specific variable to identify the forked process
	env := append(os.Environ(), "GO_DETACHED=1")

	// Get the PINNER_DIR environment variable
	trustedDir := os.Getenv("PINNER_DIR")
	if trustedDir == "" {
		log.Fatalf("PINNER_DIR environment variable not set")
	}

	// Construct the log file path
	logFilePath := filepath.Join(trustedDir, "pinner.log")

	// Open the log file
	logFile, err := os.OpenFile(logFilePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("Error opening log file: %v\n", err)
	}
	defer logFile.Close()

	// Set up process attributes
	procAttr := &os.ProcAttr{
		Dir:   filepath.Dir(executablePath),
		Files: []*os.File{os.Stdin, logFile, logFile}, // Redirect stdout and stderr to log file
		Env:   env,
		Sys: &syscall.SysProcAttr{
			Setsid: true,
		},
	}

	// Start the new process
	process, err := os.StartProcess(executablePath, args, procAttr)
	if err != nil {
		log.Fatalf("Error starting process: %v\n", err)
	}
	// Release the parent process
	process.Release()
	os.Exit(0)
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
