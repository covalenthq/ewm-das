package cmd

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	"github.com/spf13/cobra"

	"github.com/covalenthq/das-ipfs-pinner/common"
)

var (
	debug bool
	addr  string
)

// rootCmd represents the base command
var rootCmd = &cobra.Command{
	Use:     common.BinaryName,
	Short:   "A daemon for storing and retrieving data",
	Long:    `Pinner is a daemon that handles storing binary data and extracting it via HTTP.`,
	Version: fmt.Sprintf("%s, commit %s", common.Version, common.GitCommit),
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		// Handle the debug flag or daemonize if not in debug mode
		if !debug {
			if os.Getenv("GO_DAEMON") != "1" {
				daemonize()
			}
			// Set up logging to a file in daemon mode
			logFilePath := common.LogFileName()
			logFile, err := os.OpenFile(logFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
			if err != nil {
				log.Fatalf("Error opening log file: %v", err)
			}
			defer logFile.Close()
			log.SetOutput(logFile)
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		// Default action, start the server
		startServer()
	},
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().BoolVar(&debug, "debug", false, "Run in debug mode")
	rootCmd.PersistentFlags().StringVar(&addr, "addr", getEnv("DAEMON_ADDR", "localhost:5080"), "Address to run the daemon")
}

func initConfig() {
	// Additional initialization if needed
}

func daemonize() {
	executablePath, err := filepath.Abs(os.Args[0])
	if err != nil {
		log.Fatalf("Error getting absolute path: %v\n", err)
	}
	args := os.Args[1:]

	env := append(os.Environ(), "GO_DAEMON=1")

	procAttr := &os.ProcAttr{
		Dir: filepath.Dir(executablePath),
		Files: []*os.File{
			os.Stdin,
			os.Stdout,
			os.Stderr,
		},
		Env: env,
		Sys: &syscall.SysProcAttr{
			Setsid: true,
		},
	}

	process, err := os.StartProcess(executablePath, args, procAttr)
	if err != nil {
		log.Fatalf("Error starting process: %v\n", err)
	}
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

func startServer() {
	// Set up signal handling for graceful shutdown
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigs
		log.Println("Shutting down daemon...")
		// TODO: Add cleanup code here (e.g., close connections, etc.)
		os.Exit(0)
	}()

	// Set up HTTP handlers
	http.HandleFunc("/store", func(w http.ResponseWriter, r *http.Request) {
		// Handle storing data
		fmt.Fprintln(w, "Data stored")
	})

	http.HandleFunc("/extract", func(w http.ResponseWriter, r *http.Request) {
		// Handle extracting data
		cid := r.URL.Query().Get("cid")
		if cid == "" {
			http.Error(w, "CID is required", http.StatusBadRequest)
			return
		}
		fmt.Fprintf(w, "Extracting data for CID: %s\n", cid)
	})

	// Start the HTTP server
	log.Printf("Starting daemon on %s...\n", addr)
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatalf("Could not start server: %v\n", err)
	}
}
