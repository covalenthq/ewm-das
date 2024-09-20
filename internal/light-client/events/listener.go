package events

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/covalenthq/das-ipfs-pinner/internal"
	"github.com/covalenthq/das-ipfs-pinner/internal/light-client/sampler"
	"github.com/covalenthq/das-ipfs-pinner/internal/light-client/utils"
	"github.com/ethereum/go-ethereum/common"
	"github.com/filecoin-project/go-jsonrpc"
	"github.com/gorilla/websocket"
	logging "github.com/ipfs/go-log/v2"
)

var log = logging.Logger("event-listener")

// EventListener represents the event listener with a private key and a handler
type EventListener struct {
	identity *utils.Identity
	sampler  *sampler.Sampler

	mu         sync.Mutex // Protects subscription status
	subscribed bool       // Tracks whether the client is subscribed
}

// NewEventListener creates a new Listener with the provided private key in hex format
func NewEventListener(identity *utils.Identity, sampler *sampler.Sampler) *EventListener {
	return &EventListener{
		identity:   identity,
		sampler:    sampler,
		subscribed: false,
	}
}

// Id returns the address of the handler
func (h *EventListener) Id() (string, error) {
	return h.identity.GetAddress().Hex(), nil
}

// Sample is a placeholder for implementing sampling logic
func (h *EventListener) Sample(clientId, cid string, chainId, blockNum uint64, signature string) error {
	request := &internal.ScheduleRequest{
		ClientId: clientId,
		Cid:      cid,
		ChainId:  chainId,
		BlockNum: blockNum,
	}

	err := h.verifyRequest(request, signature)
	if err != nil {
		return err
	}

	h.sampler.ProcessEvent(cid, blockNum)

	return nil
}

// Start initializes the listener, performs the subscription, and blocks until a shutdown signal is received
func (l *EventListener) Start(addr string) error {
	var client struct {
		Subscribe func() error
	}

	signature, err := l.identity.SignMessage([]byte("hello"))
	if err != nil {
		return err
	}

	requestHeader := l.buildHeaders(signature)

	// Channel to notify reconnections
	reconnectNotify := make(chan struct{})

	clientHandlerOpt := jsonrpc.WithClientHandler("Client", l)
	proxyConnFactoryOpt := jsonrpc.WithProxyConnFactory(proxyConnFactory(l, reconnectNotify))

	closer, err := jsonrpc.NewMergeClient(
		context.Background(),
		addr,
		"ServerHandler",
		[]interface{}{&client},
		requestHeader,
		clientHandlerOpt,
		proxyConnFactoryOpt,
	)
	if err != nil {
		return err
	}
	defer closer()

	// Handle subscription in a goroutine after client is initialized
	go l.handleSubscription(&client, reconnectNotify)

	log.Info("Client authenticated")

	// Block until a termination signal is received
	l.waitForShutdown()

	fmt.Println("Shutting down listener...")
	return nil
}

// handleSubscription manages subscription on reconnection
func (l *EventListener) handleSubscription(client *struct{ Subscribe func() error }, reconnectNotify <-chan struct{}) {
	for {
		<-reconnectNotify // Wait for reconnection notification

		l.mu.Lock()
		if l.subscribed {
			log.Info("Already subscribed, skipping resubscription")
			l.mu.Unlock()
			continue
		}
		l.mu.Unlock()

		// Attempt to subscribe
		log.Info("Attempting to subscribe after reconnection")
		if err := client.Subscribe(); err != nil {
			log.Errorf("Failed to subscribe after reconnection: %v", err)
		} else {
			l.mu.Lock()
			l.subscribed = true
			l.mu.Unlock()
			log.Info("Subscription successful after reconnection")
		}
	}
}

// buildHeaders constructs the HTTP headers for the JSON-RPC request
func (l *EventListener) buildHeaders(signature string) http.Header {
	requestHeader := http.Header{}
	requestHeader.Add("X-LC-Signature", signature)
	requestHeader.Add("X-LC-Address", l.identity.GetAddress().Hex())
	return requestHeader
}

// waitForShutdown blocks the process until an interrupt or termination signal is received
func (l *EventListener) waitForShutdown() {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
}

func (l *EventListener) verifyRequest(request *internal.ScheduleRequest, signature string) error {
	requestBytes, err := json.Marshal(request)
	if err != nil {
		return fmt.Errorf("failed to marshal request: %w", err)
	}

	signatureBytes := common.Hex2Bytes(signature)

	ok, recoveredAddress := utils.VerifySignature(requestBytes, signatureBytes)
	if !ok {
		return fmt.Errorf("failed to verify signature")
	}

	log.Infof("Verified signature: %v", recoveredAddress.Hex())
	return nil
}

// proxyConnFactory wraps the connection factory and notifies on reconnection
func proxyConnFactory(l *EventListener, reconnectNotify chan struct{}) func(func() (*websocket.Conn, error)) func() (*websocket.Conn, error) {
	return func(originalFactory func() (*websocket.Conn, error)) func() (*websocket.Conn, error) {
		return func() (*websocket.Conn, error) {
			// Call the original connection factory
			conn, err := originalFactory()
			if err != nil {
				log.Debug(fmt.Sprintf("Connection failed: %v", err))
				l.mu.Lock()
				l.subscribed = false // Reset subscription status on connection failure
				l.mu.Unlock()
				return nil, err
			}

			log.Debug("Connection established!")

			// Notify that the connection is established
			go func() {
				reconnectNotify <- struct{}{}
			}()

			return conn, nil
		}
	}
}
