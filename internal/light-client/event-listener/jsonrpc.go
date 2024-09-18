package eventlistener

import (
	"context"
	"crypto/ecdsa"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/covalenthq/das-ipfs-pinner/internal"
	"github.com/covalenthq/das-ipfs-pinner/internal/light-client/sampler"
	"github.com/covalenthq/das-ipfs-pinner/internal/light-client/utils"
	"github.com/ethereum/go-ethereum/common"
	"github.com/filecoin-project/go-jsonrpc"
)

// Listener represents the event listener with a private key and a handler
type Listener struct {
	privKey *ecdsa.PrivateKey
	sampler *sampler.Sampler
}

// NewListener creates a new Listener with the provided private key in hex format
func NewListener(hexPrivKey string, sampler *sampler.Sampler) (*Listener, error) {
	privKey, err := utils.StringToPrivateKey(hexPrivKey)
	if err != nil {
		return nil, err
	}

	return &Listener{privKey: privKey, sampler: sampler}, nil
}

// Id returns the address of the handler
func (h *Listener) Id() (string, error) {
	addr := utils.FromECDSAPubToAddress(&h.privKey.PublicKey)
	return common.Bytes2Hex(addr), nil
}

// Sample is a placeholder for implementing sampling logic
func (h *Listener) Sample(clientId, cid string, chainId, blockNum uint64, signature string) error {

	request := &internal.ScheduleRequest{
		ClientId: clientId,
		Cid:      cid,
		ChainId:  chainId,
		BlockNum: blockNum,
	}

	// Marshal the request into a JSON string
	requestBytes, err := json.Marshal(request)
	if err != nil {
		log.Errorf("Failed to marshal request: %v", err)
		return err
	}

	signatureBytes := common.Hex2Bytes(signature)

	ok, recoveredAddress := utils.VerifySignature(requestBytes, signatureBytes)
	if !ok {
		log.Errorf("Failed to verify signature")
		return fmt.Errorf("failed to verify signature")
	}

	log.Infof("Verified signature: %v", recoveredAddress.Hex())

	h.sampler.ProcessEvent(cid, blockNum)

	return nil
}

// Start initializes the listener, performs the subscription, and blocks until a shutdown signal is received
func (l *Listener) Start(addr string) error {
	var client struct {
		Subscribe func() error
	}

	signature, err := l.signMessage([]byte("hello"))
	if err != nil {
		return err
	}

	requestHeader := l.buildHeaders(signature)

	closer, err := jsonrpc.NewMergeClient(context.Background(), addr, "ServerHandler",
		[]interface{}{&client}, requestHeader, jsonrpc.WithClientHandler("Client", l))
	if err != nil {
		return err
	}
	defer closer()

	if err := client.Subscribe(); err != nil {
		return err
	}

	// Block until a termination signal is received
	l.waitForShutdown()

	fmt.Println("Shutting down listener...")
	return nil
}

// signMessage signs a message with the private key
func (l *Listener) signMessage(message []byte) ([]byte, error) {
	signature, err := utils.SignMessage(l.privKey, message)
	if err != nil {
		return nil, fmt.Errorf("failed to sign message: %w", err)
	}
	return signature, nil
}

// buildHeaders constructs the HTTP headers for the JSON-RPC request
func (l *Listener) buildHeaders(signature []byte) http.Header {
	id, _ := l.Id()

	requestHeader := http.Header{}
	requestHeader.Add("X-LC-Signature", fmt.Sprintf("%x", signature))
	requestHeader.Add("X-LC-Address", fmt.Sprintf("%x", id))

	return requestHeader
}

// waitForShutdown blocks the process until an interrupt or termination signal is received
func (l *Listener) waitForShutdown() {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
}
