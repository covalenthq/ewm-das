package gatewayhandler

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"
	"time"

	ipldencoder "github.com/covalenthq/das-ipfs-pinner/internal/pinner/ipld-encoder"
	logging "github.com/ipfs/go-log/v2"
	"github.com/ipld/go-ipld-prime/codec/dagjson"
)

var log = logging.Logger("light-client")

var DefaultGateways = []string{
	"https://w3s.link/",
	"https://trustless-gateway.link/",
	"https://dweb.link/",
}

type errorContext struct {
	Err     error
	Context string
}

type resultContext struct {
	Result  interface{}
	Context string
}

// GatewayHandler handles fetching data from IPFS gateways concurrently.
type GatewayHandler struct {
	gateways []string
}

// NewGatewayHandler creates a new GatewayHandler instance.
func NewGatewayHandler(gateways []string) *GatewayHandler {
	return &GatewayHandler{
		gateways: gateways,
	}
}

// FetchFromGateways tries to fetch data concurrently from all gateways.
// It returns the first successful result or an error if all attempts fail.
func (g *GatewayHandler) FetchFromGateways(ctx context.Context, cidStr string, data interface{}) error {
	// Channels to capture results and errors
	resultChan := make(chan resultContext)
	errorChan := make(chan errorContext)

	// Set up context with timeout to prevent indefinite waiting
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	// Start goroutines for each gateway
	for _, gateway := range g.gateways {
		go g.fetchFromGateway(ctx, gateway, cidStr, data, resultChan, errorChan)
	}

	// Wait for the first successful result or all errors
	return g.waitForFirstSuccess(ctx, resultChan, errorChan)
}

// fetchFromGateway fetches data from a gateway concurrently.
func (g *GatewayHandler) fetchFromGateway(ctx context.Context, gateway, cidStr string, data interface{}, resultChan chan<- resultContext, errorChan chan<- errorContext) {
	// Fetch data from gateway
	gatewayData, err := g.getDataFromGateway(ctx, gateway, cidStr)
	if err != nil {
		if err == context.Canceled {
			return // Don't report the cancellation error
		}
		errorChan <- errorContext{Err: err, Context: gateway}
		return
	}

	// Decode and process the data
	if err := g.decodeAndUnmarshal(gateway, gatewayData, data); err != nil {
		errorChan <- errorContext{Err: err, Context: gateway}
		return
	}

	select {
	case resultChan <- resultContext{Result: data, Context: fmt.Sprintf("gateway %s", gateway)}:
	case <-ctx.Done():
		// Context canceled, stop further actions
	}
}

// getDataFromGateway fetches raw data from a gateway.
func (g *GatewayHandler) getDataFromGateway(ctx context.Context, gateway, cid string) ([]byte, error) {
	baseURL, err := url.Parse(gateway)
	if err != nil {
		return nil, err
	}

	baseURL.Path = path.Join(baseURL.Path, "ipfs", cid)

	query := baseURL.Query()
	query.Set("format", "raw")
	baseURL.RawQuery = query.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, baseURL.String(), nil)
	if err != nil {
		return nil, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		if urlErr, ok := err.(*url.Error); ok {
			// Handle specific URL errors like context cancellation
			if urlErr.Err.Error() == "context canceled" {
				return nil, context.Canceled
			}
		}
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("gateway %s returned status %d", gateway, resp.StatusCode)
	}

	return io.ReadAll(resp.Body)
}

// decodeAndUnmarshal decodes the IPLD data and unmarshals it into the provided interface.
func (g *GatewayHandler) decodeAndUnmarshal(gateway string, gatewayData []byte, data interface{}) error {
	// Decode the data into an IPLD node from CBOR
	node, err := ipldencoder.DecodeNode(gatewayData)
	if err != nil {
		return fmt.Errorf("failed to decode IPLD data from %s: %w", gateway, err)
	}

	// Encode the IPLD node into JSON
	var jsonData bytes.Buffer
	if err := dagjson.Encode(node, &jsonData); err != nil {
		return fmt.Errorf("failed to encode IPLD node into JSON from %s: %w", gateway, err)
	}

	// Unmarshal JSON data into the provided data interface
	if err := json.Unmarshal(jsonData.Bytes(), data); err != nil {
		return fmt.Errorf("failed to unmarshal JSON data from %s: %w", gateway, err)
	}

	return nil
}

// waitForFirstSuccess waits for the first successful result or returns an error if all fail.
func (g *GatewayHandler) waitForFirstSuccess(ctx context.Context, resultChan <-chan resultContext, errorChan <-chan errorContext) error {
	successCount := 0
	totalCount := len(g.gateways)

	for i := 0; i < totalCount; i++ {
		select {
		case <-ctx.Done():
			// If the context is canceled and we've had no success, return a timeout error
			if successCount > 0 {
				return nil
			}
			return fmt.Errorf("timeout exceeded")

		case result := <-resultChan:
			log.Debugf("Data fetched from %s", result.Context)
			successCount++
			// Return immediately after the first successful result
			if successCount == 1 {
				return nil
			}

		case errCxt := <-errorChan:
			log.Debugf("Error getting data from %s: %v", errCxt.Context, errCxt.Err)
		}
	}

	// If no successful results, return an error
	if successCount == 0 {
		return fmt.Errorf("failed to fetch data from any gateway")
	}

	return nil
}
