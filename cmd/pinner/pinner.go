package main

import (
	"fmt"
	"os"
	"path/filepath"
	"syscall"

	logging "github.com/ipfs/go-log/v2"
	"github.com/spf13/cobra"

	"github.com/covalenthq/das-ipfs-pinner/api"
	"github.com/covalenthq/das-ipfs-pinner/common"
	"github.com/covalenthq/das-ipfs-pinner/internal/das"
)

var (
	debug                 bool
	addr                  string
	w3AgentKey            string
	w3DelegationProofPath string
)

// rootCmd represents the base command
var rootCmd = &cobra.Command{
	Use:     common.BinaryName,
	Short:   "A daemon for storing and retrieving data",
	Long:    `Pinner is a daemon that handles storing binary data and extracting it via HTTP.`,
	Version: fmt.Sprintf("%s, commit %s", common.Version, common.GitCommit),
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		if debug {
			logging.SetLogLevel("das-pinner", "debug")
		} else {
			logging.SetLogLevel("das-pinner", "info")
		}

		// Load the configuration
		config := das.LoadConfig()
		// Initialize the KZG trusted setup
		if err := das.InitializeTrustedSetup(config); err != nil {
			log.Fatalf("Failed to initialize trusted setup: %v", err)
		}

		// Handle the debug flag or daemonize if not in debug mode
		if !debug {
			if os.Getenv("GO_DAEMON") != "1" {
				daemonize()
			}
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

func init() {
	log.Info("Initializing root command...", os.Args)

	// Check if we're in the child process (daemon)
	if os.Getenv("GO_DAEMON") == "1" {
		log.Info("Running in daemon mode.")
		// Do not reinitialize flags or commands here
		// Proceed directly to running the server or minimal initialization required
		return
	}

	log.Info("Running in non-daemon mode.")

	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().BoolVar(&debug, "debug", false, "Run in debug mode")
	rootCmd.PersistentFlags().StringVar(&addr, "addr", getEnv("DAEMON_ADDR", "localhost:5080"), "Address to run the daemon")

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
	env := append(os.Environ(), "GO_DAEMON=1")

	// Set up process attributes
	procAttr := &os.ProcAttr{
		Dir:   filepath.Dir(executablePath),
		Files: []*os.File{os.Stdin, os.Stdout, os.Stderr},
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
func execute() error {
	return rootCmd.Execute()
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
