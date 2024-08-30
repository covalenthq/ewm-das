package eventlistener

import (
	"context"
	"net/url"
	"time"

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

	// Initial subscription
	sub, err := el.client.SubscribeFilterLogs(ctx, query, el.logs)
	if err != nil {
		log.Errorf("Failed to subscribe to logs: %v", err)
		el.retrySubscription(ctx, query) // Try to recover by retrying the subscription
		return
	}

	el.subscription = sub

	go func() {
		for {
			select {
			case err := <-sub.Err():
				if err != nil {
					log.Errorf("Subscription error: %v", err)
					el.retrySubscription(ctx, query) // Try to recover by retrying the subscription
					return
				}
			case <-ctx.Done():
				log.Infof("Context canceled, stopping log subscription.")
				return
			}
		}
	}()

	log.Infof("Subscribed to logs for contract: %v", el.contractAddress.Hex())
}

func (el *EventListener) retrySubscription(ctx context.Context, query ethereum.FilterQuery) {
	backoff := 2 * time.Second    // Initial backoff duration
	maxBackoff := 1 * time.Minute // Maximum backoff duration

	for {
		select {
		case <-ctx.Done():
			log.Infof("Context canceled, aborting subscription retry.")
			return
		default:
			sub, err := el.client.SubscribeFilterLogs(ctx, query, el.logs)
			if err != nil {
				log.Errorf("Retrying subscription to logs failed: %v", err)

				// Increase backoff duration, but don't exceed maxBackoff
				time.Sleep(backoff)
				backoff *= 2
				if backoff > maxBackoff {
					backoff = maxBackoff
				}
				continue
			}

			el.subscription = sub

			go func() {
				for err := range sub.Err() {
					if err != nil {
						log.Errorf("Subscription error after retry: %v", err)
						el.retrySubscription(ctx, query) // Retry again if needed
						return
					}
				}
			}()

			log.Infof("Successfully resubscribed to logs for contract: %v", el.contractAddress.Hex())
			return
		}
	}
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

		log.Debugf("chain-id=%v, storage-url=%v, block-height=%v", event.ChainId, event.StorageURL, event.BlockHeight)

		// strip the ipfs://
		parsedURL, err := url.Parse(event.StorageURL)
		if err != nil {
			log.Errorf("Failed to parse URL: %v", err)
			continue
		}

		el.sampler.ProcessEvent(parsedURL.Host, event.BlockHeight)
	}
}
