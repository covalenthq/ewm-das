package main

import (
	"context"

	eventlistener "github.com/covalenthq/das-ipfs-pinner/internal/light-client/event-listener"
	"github.com/covalenthq/das-ipfs-pinner/internal/light-client/sampler"
	logging "github.com/ipfs/go-log/v2"
	"github.com/spf13/cobra"
)

var (
	loglevel   string
	rpcURL     string
	contract   string
	ipfsAddr   string
	serviceURL string
)

var log = logging.Logger("light-client")

var rootCmd = &cobra.Command{
	Use:   "my-client",
	Short: "A client to interact with blockchain events and IPFS",
	Long:  `This client listens for events from a smart contract on a specified chain, retrieves data from IPFS, and sends it to another service.`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		logging.SetLogLevel("light-client", loglevel)
	},
	Run: func(cmd *cobra.Command, args []string) {
		// Implement the main logic here
		startClient()
	},
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatalf("Error executing command: %v", err)
	}
}

func init() {
	log.Info("Initializing client...")

	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVar(&loglevel, "loglevel", "info", "Log level (debug, info, warn, error, fatal, panic)")
	rootCmd.PersistentFlags().StringVar(&rpcURL, "rpc-url", "", "RPC URL of the blockchain node")
	rootCmd.PersistentFlags().StringVar(&contract, "contract", "", "Contract address to listen for events")
	rootCmd.PersistentFlags().StringVar(&ipfsAddr, "ipfs-addr", "http://localhost:5001", "IPFS node address")
	rootCmd.PersistentFlags().StringVar(&serviceURL, "service-url", "", "URL of the service to send data to")

	rootCmd.MarkPersistentFlagRequired("rpc-url")
	rootCmd.MarkPersistentFlagRequired("contract")
	rootCmd.MarkPersistentFlagRequired("service-url")
}

func initConfig() {
	// Additional configuration initialization if needed
}

func startClient() {
	log.Info("Starting client...")

	sampler, err := sampler.NewSampler(ipfsAddr)
	if err != nil {
		log.Fatalf("Failed to initialize IPFS sampler: %v", err)
	}

	eventlistener := eventlistener.NewEventListener(rpcURL, contract, sampler)
	eventlistener.SubscribeToLogs(context.Background())
	eventlistener.ProcessLogs()
}
