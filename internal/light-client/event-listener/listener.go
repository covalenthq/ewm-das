package eventlistener

import (
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

// EventListener listens for events emitted by the contract
type EventListener struct {
	Client           *ethclient.Client
	ContractAddress  common.Address
	ContractInstance *contract.Contract
	Logs             chan types.Log
	Subscription     ethereum.Subscription
	Sampler          *sampler.Sampler
}

// NewEventListener creates a new EventListener instance
func NewEventListener(clientURL, contractAddressHex string, sampler *sampler.Sampler) *EventListener {
	client := connectToEthereumClient(clientURL)
	contractAddress := common.HexToAddress(contractAddressHex)
	contractInstance := loadContract(client, contractAddress)

	return &EventListener{
		Client:           client,
		ContractAddress:  contractAddress,
		ContractInstance: contractInstance,
		Logs:             make(chan types.Log),
		Sampler:          sampler,
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

// SubscribeToLogs subscribes to logs emitted by the contract
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

// ProcessLogs processes the logs emitted by the contract
func (el *EventListener) ProcessLogs() {
	for vLog := range el.Logs {
		log.Debugf("Log: %v\n", vLog)

		event, err := el.ContractInstance.ParseBlockResultProductionProofSubmitted(vLog)
		if err != nil {
			log.Warnf("Failed to parse log: %v", err)
			continue
		}

		log.Debugf("Event ChainID: %v\n", event.ChainId)
		log.Debugf("Event StorageURL: %v\n", event.StorageURL)

		// el.Sampler.ProcessEvent(event.StorageURL)
		el.Sampler.ProcessEvent("bafyreiahay5quioczvzk5tdr7muuiyozmtsq6yizncwi6r6bst42v5jnqi")
	}
}
