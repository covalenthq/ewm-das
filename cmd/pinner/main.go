package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	"github.com/covalenthq/das-ipfs-pinner/common"
)

func main() {
	// Command-line flags for configuration
	debug := flag.Bool("debug", false, "Run in debug mode")
	addr := flag.String("addr", getEnv("DAEMON_ADDR", "localhost:8080"), "Address to run the daemon")
	flag.Parse()

	// Set up signal handling for graceful shutdown
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigs
		log.Println("Shutting down daemon...")
		// TODO: Add cleanup code here (e.g., close connections, etc.)
		os.Exit(0)
	}()

	if !*debug {
		// Only daemonize if not in debug mode
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

	// Set up HTTP handlers
	http.HandleFunc("/store", func(w http.ResponseWriter, r *http.Request) {
		// Handle storing data
		fmt.Fprintln(w, "Data stored")
	})

	http.HandleFunc("/extract", func(w http.ResponseWriter, r *http.Request) {
		// Handle extracting data
		fmt.Fprintln(w, "Data extracted")
	})

	// Start the HTTP server
	log.Printf("Starting daemon on %s...\n", *addr)
	if err := http.ListenAndServe(*addr, nil); err != nil {
		log.Fatalf("Could not start server: %v\n", err)
	}
}

func daemonize() {
	// Get the absolute path of the executable
	executablePath, err := filepath.Abs(os.Args[0])
	if err != nil {
		log.Fatalf("Error getting absolute path: %v\n", err)
	}
	args := os.Args[1:]

	// Set environment variable to prevent infinite loop
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

// getEnv gets the environment variable or returns a default value if not set
func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
