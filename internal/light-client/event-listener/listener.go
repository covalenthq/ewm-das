package eventlistener

import (
	"context"
	"net/url"

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
	client           *ethclient.Client
	contractAddress  common.Address
	contractInstance *contract.Contract
	logs             chan types.Log
	subscription     ethereum.Subscription
	sampler          *sampler.Sampler
}

// NewEventListener creates a new EventListener instance
func NewEventListener(clientURL, contractAddressHex string, sampler *sampler.Sampler) *EventListener {
	client := connectToEthereumClient(clientURL)
	contractAddress := common.HexToAddress(contractAddressHex)
	contractInstance := loadContract(client, contractAddress)

	return &EventListener{
		client:           client,
		contractAddress:  contractAddress,
		contractInstance: contractInstance,
		logs:             make(chan types.Log),
		sampler:          sampler,
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
		Addresses: []common.Address{el.contractAddress},
	}

	sub, err := el.client.SubscribeFilterLogs(ctx, query, el.logs)
	if err != nil {
		log.Fatalf("Failed to subscribe to logs: %v", err)
	}

	el.subscription = sub

	go func() {
		for err := range sub.Err() {
			log.Fatalf("Subscription error: %v", err)
		}
	}()

	log.Infof("Subscribed to logs for contract: %v", el.contractAddress.Hex())
}

// ProcessLogs processes the logs emitted by the contract
func (el *EventListener) ProcessLogs() {
	for vLog := range el.logs {
		log.Debugf("Log Event: %v", vLog.Topics)

		event, err := el.contractInstance.ParseBlockSpecimenProductionProofSubmitted(vLog)
		if err != nil {
			if err.Error() == "event signature mismatch" {
				log.Debug("Event signature mismatch")
				continue
			}

			log.Errorf("Failed to parse log: %v", err)
		}

		log.Debugf("Event ChainID: %v", event.ChainId)
		log.Debugf("Event StorageURL: %v", event.StorageURL)
		log.Debugf("Event BlockHeight: %v", event.BlockHeight)

		// strip the ipfs://
		parsedURL, err := url.Parse(event.StorageURL)
		if err != nil {
			log.Errorf("Failed to parse URL: %v", err)
			continue
		}

		el.sampler.ProcessEvent(parsedURL.Host, uint64(event.BlockHeight))
	}
}
