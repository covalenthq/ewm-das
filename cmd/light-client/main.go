package main

import (
	"context"
	"fmt"

	"github.com/covalenthq/das-ipfs-pinner/common"
	eventlistener "github.com/covalenthq/das-ipfs-pinner/internal/light-client/event-listener"
	publisher "github.com/covalenthq/das-ipfs-pinner/internal/light-client/publisher"
	"github.com/covalenthq/das-ipfs-pinner/internal/light-client/sampler"
	"github.com/covalenthq/das-ipfs-pinner/internal/pinner/das"
	logging "github.com/ipfs/go-log/v2"
	"github.com/spf13/cobra"
)

var (
	loglevel      string
	rpcURL        string
	contract      string
	ipfsAddr      string
	gcpTopicId    string
	gcpCredsFile  string
	clientId      string
	samplingDelay uint
)

var greeting = `
░▒▓█▓▒░      ░▒▓█▓▒░░▒▓██████▓▒░░▒▓█▓▒░░▒▓█▓▒░▒▓████████▓▒░       ░▒▓██████▓▒░░▒▓█▓▒░      ░▒▓█▓▒░▒▓████████▓▒░▒▓███████▓▒░▒▓████████▓▒░ 
░▒▓█▓▒░      ░▒▓█▓▒░▒▓█▓▒░░▒▓█▓▒░▒▓█▓▒░░▒▓█▓▒░  ░▒▓█▓▒░          ░▒▓█▓▒░░▒▓█▓▒░▒▓█▓▒░      ░▒▓█▓▒░▒▓█▓▒░      ░▒▓█▓▒░░▒▓█▓▒░ ░▒▓█▓▒░     
░▒▓█▓▒░      ░▒▓█▓▒░▒▓█▓▒░      ░▒▓█▓▒░░▒▓█▓▒░  ░▒▓█▓▒░          ░▒▓█▓▒░      ░▒▓█▓▒░      ░▒▓█▓▒░▒▓█▓▒░      ░▒▓█▓▒░░▒▓█▓▒░ ░▒▓█▓▒░     
░▒▓█▓▒░      ░▒▓█▓▒░▒▓█▓▒▒▓███▓▒░▒▓████████▓▒░  ░▒▓█▓▒░          ░▒▓█▓▒░      ░▒▓█▓▒░      ░▒▓█▓▒░▒▓██████▓▒░ ░▒▓█▓▒░░▒▓█▓▒░ ░▒▓█▓▒░     
░▒▓█▓▒░      ░▒▓█▓▒░▒▓█▓▒░░▒▓█▓▒░▒▓█▓▒░░▒▓█▓▒░  ░▒▓█▓▒░          ░▒▓█▓▒░      ░▒▓█▓▒░      ░▒▓█▓▒░▒▓█▓▒░      ░▒▓█▓▒░░▒▓█▓▒░ ░▒▓█▓▒░     
░▒▓█▓▒░      ░▒▓█▓▒░▒▓█▓▒░░▒▓█▓▒░▒▓█▓▒░░▒▓█▓▒░  ░▒▓█▓▒░          ░▒▓█▓▒░░▒▓█▓▒░▒▓█▓▒░      ░▒▓█▓▒░▒▓█▓▒░      ░▒▓█▓▒░░▒▓█▓▒░ ░▒▓█▓▒░     
░▒▓████████▓▒░▒▓█▓▒░░▒▓██████▓▒░░▒▓█▓▒░░▒▓█▓▒░  ░▒▓█▓▒░           ░▒▓██████▓▒░░▒▓████████▓▒░▒▓█▓▒░▒▓████████▓▒░▒▓█▓▒░░▒▓█▓▒░ ░▒▓█▓▒░     
                                                                                                                                         
                                                                                                                                         
`

var log = logging.Logger("light-client")

var rootCmd = &cobra.Command{
	Use:     "light-client",
	Short:   "A client to interact with blockchain events and IPFS",
	Long:    `This client listens for events from a smart contract on a specified chain, retrieves data from IPFS, and sends it to another service.`,
	Version: fmt.Sprintf("%s, commit %s", common.Version, common.GitCommit),
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		logging.SetLogLevel("light-client", loglevel)

		// Load the configuration
		config := das.LoadConfig()
		// Initialize the KZG trusted setup
		if err := das.InitializeTrustedSetup(config); err != nil {
			log.Fatalf("Failed to initialize trusted setup: %v", err)
		}
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
	rootCmd.PersistentFlags().StringVar(&gcpTopicId, "topic-id", "", "Topic name of Pub Sub")
	rootCmd.PersistentFlags().StringVar(&gcpCredsFile, "gcp-creds-file", "", "Path of GCP credential json file")
	rootCmd.PersistentFlags().StringVar(&clientId, "client-id", "", "arbitrary client ID, used to identify the client")
	rootCmd.PersistentFlags().UintVar(&samplingDelay, "sampling-delay", 120, "Delay between sampling process and the receiving of the event")

	rootCmd.MarkPersistentFlagRequired("rpc-url")
	rootCmd.MarkPersistentFlagRequired("contract")
	rootCmd.MarkPersistentFlagRequired("project-id")
	rootCmd.MarkPersistentFlagRequired("topic-id")
	rootCmd.MarkPersistentFlagRequired("gcp-creds-file")
	rootCmd.MarkPersistentFlagRequired("client-id")
}

func initConfig() {
	// Additional configuration initialization if needed
}

func startClient() {
	fmt.Println(greeting)
	fmt.Printf("Version: %s, commit: %s\n", common.Version, common.GitCommit)
	log.Info("Starting client...")

	pub, err := publisher.NewPublisher(gcpTopicId, gcpCredsFile, clientId)
	if err != nil {
		log.Fatalf("Failed to create publisher: %v", err)
	}

	sampler, err := sampler.NewSampler(ipfsAddr, samplingDelay, pub)
	if err != nil {
		log.Fatalf("Failed to initialize IPFS sampler: %v", err)
	}

	eventlistener := eventlistener.NewEventListener(rpcURL, contract, sampler)
	eventlistener.SubscribeToLogs(context.Background())
	eventlistener.ProcessLogs()
}
