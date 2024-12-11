package events

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/covalenthq/das-ipfs-pinner/common"
	"github.com/covalenthq/das-ipfs-pinner/internal/light-client/sampler"
	"github.com/covalenthq/das-ipfs-pinner/internal/light-client/utils"
	"github.com/filecoin-project/go-jsonrpc"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	logging "github.com/ipfs/go-log/v2"
)

var log = logging.Logger("event-listener")

// EventListener represents the event listener with a private key and a handler
type EventListener struct {
	identity *utils.Identity
	sampler  *sampler.Sampler
	id       uuid.UUID

	mu         sync.Mutex // Protects subscription status
	subscribed bool       // Tracks whether the client is subscribed

	networkStatus chan bool
}

// NewEventListener creates a new Listener with the provided private key in hex format
func NewEventListener(identity *utils.Identity, sampler *sampler.Sampler) *EventListener {
	return &EventListener{
		identity:      identity,
		sampler:       sampler,
		subscribed:    false,
		id:            uuid.New(),
		networkStatus: make(chan bool),
	}
}

// SessionId returns the unique identifier of the handler
func (h *EventListener) SessionId() (string, error) {
	return h.id.String(), nil
}

// Identity returns the address of the identity
func (h *EventListener) Identity() ([]byte, error) {
	return h.identity.GetAddress().Bytes(), nil
}

// Version returns the version of the handler
func (h *EventListener) Version() (string, error) {
	return fmt.Sprintf("%s-%s", common.Version, common.GitCommit), nil
}

// Start initializes the listener, performs the subscription, and blocks until a shutdown signal is received
func (l *EventListener) Start(addr string) error {
	go l.monitorNetworkRoutine()

	var rpcServer struct {
		Subscribe func() error
	}

	signature, err := l.identity.SignMessage([]byte(l.id.String()))
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
		[]interface{}{&rpcServer},
		requestHeader,
		clientHandlerOpt,
		proxyConnFactoryOpt,
	)
	if err != nil {
		return err
	}
	defer closer()

	// Handle subscription in a goroutine after client is initialized
	go l.handleSubscription(&rpcServer, reconnectNotify)

	log.Info("Client authenticated!")

	// Block until a termination signal is received
	l.waitForShutdown()

	fmt.Println("Shutting down listener...")
	return nil
}

// monitorNetworkRoutine checks the network availability every 10 seconds and sends the result through the networkStatus channel
func (l *EventListener) monitorNetworkRoutine() {
	var oldStatus bool
	for {
		status := isNetworkAvailable()

		if status != oldStatus {
			if status {
				log.Info("Network is available!")
			} else {
				log.Warn("Network is down!")
			}
		}

		oldStatus = status
		time.Sleep(10 * time.Second)
	}
}

// handleSubscription manages subscription on reconnection
func (l *EventListener) handleSubscription(rpcServer *struct{ Subscribe func() error }, reconnectNotify <-chan struct{}) {
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
		log.Info("Attempting to subscribe after [re]connection")
		if err := rpcServer.Subscribe(); err != nil {
			log.Errorf("Failed to subscribe after reconnection: %v", err)
		} else {
			l.mu.Lock()
			l.subscribed = true
			l.mu.Unlock()
			log.Info("Subscription successful after [re]connection")
		}
	}
}

// buildHeaders constructs the HTTP headers for the JSON-RPC request
func (l *EventListener) buildHeaders(signature []byte) http.Header {
	requestHeader := http.Header{}
	requestHeader.Add("X-LC-Signature", fmt.Sprintf("%x", signature))
	requestHeader.Add("X-LC-Address", l.identity.GetAddress().Hex())
	requestHeader.Add("X-LC-ID", l.id.String())
	return requestHeader
}

// waitForShutdown blocks the process until an interrupt or termination signal is received
func (l *EventListener) waitForShutdown() {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
}

// proxyConnFactory wraps the connection factory and notifies on reconnection
func proxyConnFactory(l *EventListener, reconnectNotify chan struct{}) func(func() (*websocket.Conn, error)) func() (*websocket.Conn, error) {
	return func(originalFactory func() (*websocket.Conn, error)) func() (*websocket.Conn, error) {
		return func() (*websocket.Conn, error) {
			// If we are here, it means no connection has been established yet and no subscription has been made
			l.mu.Lock()
			l.subscribed = false
			l.mu.Unlock()

			// Call the original connection factory
			conn, err := originalFactory()
			if err != nil {
				if errors.Is(err, websocket.ErrBadHandshake) {
					return nil, errors.New("authorization or authentication failed")
				}

				log.Debugw("%v", err)
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

// Function to check if network is available by pinging a known server (Google DNS)
func isNetworkAvailable() bool {
	_, err := net.DialTimeout("tcp", "8.8.8.8:53", time.Second*5)
	return err == nil
}
