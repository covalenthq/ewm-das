package gateway

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

var log = logging.Logger("das-gateway")

var DefaultGateways = []string{
	"https://w3s.link/",
	"https://trustless-gateway.link/",
	"https://dweb.link/",
}

// Handler handles fetching data from IPFS gateways concurrently.
type Handler struct {
	gateways []string
}

// NewHandler creates a new GatewayHandler instance.
func NewHandler(gateways []string) *Handler {
	return &Handler{gateways: gateways}
}

// FetchFromGateways concurrently fetches data from multiple gateways.
// Returns the first successful result or an error if all attempts fail.
func (g *Handler) FetchFromGateways(ctx context.Context, cidStr string, data interface{}) error {
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	results := make(chan error, 1) // Only need 1 result, so buffer size is 1
	for _, gateway := range g.gateways {
		go g.fetchFromGateway(ctx, gateway, cidStr, data, results)
	}

	// Return immediately after the first successful result
	select {
	case <-ctx.Done():
		return fmt.Errorf("operation timed out")
	case err := <-results:
		if err != nil {
			log.Debugf("Error: %v", err)
			return err
		}
		return nil
	}
}

// fetchFromGateway retrieves and processes data from a single gateway.
func (g *Handler) fetchFromGateway(ctx context.Context, gateway, cidStr string, data interface{}, results chan<- error) {
	gatewayData, err := g.getDataFromGateway(ctx, gateway, cidStr)
	if err != nil {
		select {
		case results <- fmt.Errorf("gateway %s: %v", gateway, err):
		case <-ctx.Done():
		}
		return
	}

	if err := g.decodeAndUnmarshal(gateway, gatewayData, data); err != nil {
		select {
		case results <- fmt.Errorf("gateway %s: %v", gateway, err):
		case <-ctx.Done():
		}
		return
	}

	// Send success (nil error) if the data is fetched and processed correctly
	select {
	case results <- nil:
	case <-ctx.Done():
	}
}

// getDataFromGateway fetches raw data from a specified gateway.
func (g *Handler) getDataFromGateway(ctx context.Context, gateway, cid string) ([]byte, error) {
	baseURL, err := url.Parse(gateway)
	if err != nil {
		return nil, err
	}
	baseURL.Path = path.Join(baseURL.Path, "ipfs", cid)
	baseURL.RawQuery = url.Values{"format": []string{"raw"}}.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, baseURL.String(), nil)
	if err != nil {
		return nil, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		if urlErr, ok := err.(*url.Error); ok && urlErr.Err == context.Canceled {
			return nil, context.Canceled
		}
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("received status %d", resp.StatusCode)
	}

	return io.ReadAll(resp.Body)
}

// decodeAndUnmarshal decodes IPLD data and unmarshals it into the provided structure.
func (g *Handler) decodeAndUnmarshal(gateway string, gatewayData []byte, data interface{}) error {
	node, err := ipldencoder.DecodeNode(gatewayData)
	if err != nil {
		return fmt.Errorf("failed to decode IPLD data from %s: %w", gateway, err)
	}

	var jsonData bytes.Buffer
	if err := dagjson.Encode(node, &jsonData); err != nil {
		return fmt.Errorf("failed to encode IPLD node into JSON from %s: %w", gateway, err)
	}

	if err := json.Unmarshal(jsonData.Bytes(), data); err != nil {
		return fmt.Errorf("failed to unmarshal JSON data from %s: %w", gateway, err)
	}

	return nil
}
