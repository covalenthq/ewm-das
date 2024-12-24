package main

import (
	"fmt"

	"github.com/covalenthq/das-ipfs-pinner/common"
	"github.com/covalenthq/das-ipfs-pinner/internal/light-client/apihandler"
	"github.com/covalenthq/das-ipfs-pinner/internal/light-client/poller"
	"github.com/covalenthq/das-ipfs-pinner/internal/light-client/sampler"
	"github.com/covalenthq/das-ipfs-pinner/internal/light-client/utils"
	"github.com/covalenthq/das-ipfs-pinner/internal/pinner/das"
	logging "github.com/ipfs/go-log/v2"
	"github.com/spf13/cobra"
)

var (
	loglevel      string
	rpcURL        string
	ipfsAddr      string
	privateKey    string
	samplingDelay uint
)

var greeting = `
███████╗██╗    ██╗███╗   ███╗     ██████╗██╗     ██╗███████╗███╗   ██╗████████╗
██╔════╝██║    ██║████╗ ████║    ██╔════╝██║     ██║██╔════╝████╗  ██║╚══██╔══╝
█████╗  ██║ █╗ ██║██╔████╔██║    ██║     ██║     ██║█████╗  ██╔██╗ ██║   ██║   
██╔══╝  ██║███╗██║██║╚██╔╝██║    ██║     ██║     ██║██╔══╝  ██║╚██╗██║   ██║   
███████╗╚███╔███╔╝██║ ╚═╝ ██║    ╚██████╗███████╗██║███████╗██║ ╚████║   ██║   
╚══════╝ ╚══╝╚══╝ ╚═╝     ╚═╝     ╚═════╝╚══════╝╚═╝╚══════╝╚═╝  ╚═══╝   ╚═╝   
                                                                                                                                                       
                                                                                                                                                                                              
`

var log = logging.Logger("light-client")

var rootCmd = &cobra.Command{
	Use:     "light-client",
	Short:   "A client to interact with blockchain events and IPFS",
	Long:    `This client listens for events from a smart contract on a specified chain, retrieves data from IPFS, and sends it to another service.`,
	Version: fmt.Sprintf("%s, commit %s", common.Version, common.GitCommit),
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		logging.SetLogLevel("*", loglevel)
		logging.SetLogLevel("rpc", "error")

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
	rootCmd.PersistentFlags().StringVar(&rpcURL, "rpc-url", "", "RPC URL for workload coordination")
	rootCmd.PersistentFlags().StringVar(&ipfsAddr, "ipfs-addr", ":5001", "IPFS node address")
	rootCmd.PersistentFlags().StringVar(&privateKey, "private-key", "", "Private key of the client")
	rootCmd.PersistentFlags().UintVar(&samplingDelay, "sampling-delay", 10, "Delay between sampling process and the receiving of the event")

	rootCmd.MarkPersistentFlagRequired("rpc-url")
	rootCmd.MarkPersistentFlagRequired("private-key")
}

func initConfig() {
	// Additional configuration initialization if needed
}

func startClient() {
	fmt.Println(greeting)
	fmt.Printf("Version: %s, commit: %s\n", common.Version, common.GitCommit)
	log.Info("Starting client...")

	identify, err := utils.NewIdentity(privateKey)
	if err != nil {
		log.Fatalf("Failed to create identity: %v", err)
	}
	log.Infof("Client idenity: %s", identify.GetAddress().Hex())

	api, err := apihandler.NewApiHandler(rpcURL, identify)
	if err != nil {
		log.Fatalf("Failed to create API handler: %v", err)
	}

	sampler, err := sampler.NewSampler(ipfsAddr, samplingDelay, api)
	if err != nil {
		log.Fatalf("Failed to initialize IPFS sampler: %v", err)
	}

	poller.NewWorkloadPoller(identify, sampler, api).Start()
}
