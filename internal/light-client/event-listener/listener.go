package eventlistener

import (
	"time"
	"sync"
	"context"

	logging "github.com/ipfs/go-log/v2"

	"github.com/covalenthq/das-ipfs-pinner/internal/light-client/contract"
	"github.com/covalenthq/das-ipfs-pinner/internal/light-client/sampler"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

var log = logging.Logger("light-client")

type EventListener struct {
	Client           *ethclient.Client
	ContractAddress  common.Address
	ContractInstance *contract.Contract
	Logs             chan types.Log
	Subscription     ethereum.Subscription
	Sampler          *sampler.Sampler
	RecentURLs       map[string]time.Time
	mu               sync.Mutex
	TimeWindow       time.Duration
}

// NewEventListener creates a new EventListener instance
func NewEventListener(clientURL, contractAddressHex string, sampler *sampler.Sampler, timeWindow time.Duration) *EventListener {
	client := connectToEthereumClient(clientURL)
	contractAddress := common.HexToAddress(contractAddressHex)
	contractInstance := loadContract(client, contractAddress)

	return &EventListener{
		Client:           client,
		ContractAddress:  contractAddress,
		ContractInstance: contractInstance,
		Logs:             make(chan types.Log),
		Sampler:          sampler,
		RecentURLs:       make(map[string]time.Time),
		TimeWindow:       timeWindow,
	}
}

// Connect to the Ethereum client
func connectToEthereumClient(url string) *ethclient.Client {
	client, err := ethclient.Dial(url)
	if err != nil {
		log.Fatalf("Failed to connect to the Ethereum client: %v", err)
	}
	return client
}

// Load the contract instance
func loadContract(client *ethclient.Client, address common.Address) *contract.Contract {
	contractInstance, err := contract.NewContract(address, client)
	if err != nil {
		log.Fatalf("Failed to load the contract: %v", err)
	}
	return contractInstance
}

// Subscribe to logs for the specified contract address
func (el *EventListener) SubscribeToLogs(ctx context.Context) {
	query := ethereum.FilterQuery{
		Addresses: []common.Address{el.ContractAddress},
	}

	sub, err := el.Client.SubscribeFilterLogs(ctx, query, el.Logs)
	if err != nil {
		log.Fatalf("Failed to subscribe to logs: %v", err)
	}

	el.Subscription = sub

	go func() {
		for err := range sub.Err() {
			log.Fatalf("Subscription error: %v", err)
		}
	}()

	log.Infof("Subscribed to logs for contract: %v", el.ContractAddress.Hex())
}

// Process incoming logs and handle events
func (el *EventListener) ProcessLogs(projectId string, topicId string) {
	for vLog := range el.Logs {
		log.Debugf("Log: %v\n", vLog)

		event, err := el.ContractInstance.ParseBlockResultProductionProofSubmitted(vLog)
		if err != nil {
			log.Warnf("Failed to parse log: %v", err)
			continue
		}

		log.Debugf("Event ChainID: %v\n", event.ChainId)
		log.Debugf("Event StorageURL: %v\n", event.StorageURL)

		url := event.StorageURL
		el.mu.Lock()
		isUnique := el.isUniqueURL(url)
		el.mu.Unlock()

		if isUnique {
			el.Sampler.ProcessEvent("bafyreiahay5quioczvzk5tdr7muuiyozmtsq6yizncwi6r6bst42v5jnqi", projectId, topicId)
		} else {
			log.Debugf("Skipping duplicate URL: %v\n", url)
		}
	}
		// el.Sampler.ProcessEvent(event.StorageURL)
		// el.Sampler.ProcessEvent("bafyreiahay5quioczvzk5tdr7muuiyozmtsq6yizncwi6r6bst42v5jnqi")
	
}

// Check if URL is unique within the time window
func (el *EventListener) isUniqueURL(url string) bool {
	now := time.Now()
	// Remove expired URLs
	for u, t := range el.RecentURLs {
		if now.Sub(t) > el.TimeWindow {
			delete(el.RecentURLs, u)
		}
	}
	// Check if the URL is in the map
	if _, exists := el.RecentURLs[url]; exists {
		return false
	}
	// Add the new URL with the current timestamp
	el.RecentURLs[url] = now
	return true
}