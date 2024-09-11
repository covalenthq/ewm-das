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
	"sync"
	"time"

	ipldencoder "github.com/covalenthq/das-ipfs-pinner/internal/pinner/ipld-encoder"
	logging "github.com/ipfs/go-log/v2"
	"github.com/ipld/go-ipld-prime/codec/dagjson"
	"golang.org/x/exp/rand"
)

var log = logging.Logger("das-gateway")

var DefaultGateways = []string{
	"https://w3s.link/",
	"https://trustless-gateway.link/",
	"https://dweb.link/",
	"https://ipfs.io/",
}

// Handler handles fetching data from IPFS gateways concurrently with a worker pool.
type Handler struct {
	gateways []string
	workers  int // Number of workers in the pool
}

// NewHandler creates a new GatewayHandler instance with a specified number of workers.
func NewHandler(gateways []string, workers int) *Handler {
	return &Handler{gateways: gateways, workers: workers}
}

// FetchFromGateways concurrently fetches data from multiple gateways using a worker pool.
// Returns the first successful result or an error if all attempts fail.
func (g *Handler) FetchFromGateways(ctx context.Context, cidStr string, data interface{}) error {
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel() // Ensures context cleanup

	results := make(chan error, 1) // Channel to capture the first result (success or failure)
	gatewayChan := make(chan string, len(g.gateways))

	// Start worker pool
	var wg sync.WaitGroup
	for i := 0; i < g.workers; i++ {
		wg.Add(1)
		go g.worker(ctx, cidStr, data, gatewayChan, results, &wg)
	}

	// Send gateways to the worker pool
	for _, gateway := range g.shufledGateways() {
		gatewayChan <- gateway
	}

	// Close gatewayChan as we have sent all gateways
	close(gatewayChan)

	// Wait for workers in a separate goroutine to close results when all are done
	go func() {
		wg.Wait()
		close(results)
	}()

	// Wait for the first successful result or context timeout
	for {
		select {
		case <-ctx.Done():
			return fmt.Errorf("operation timed out")
		case err := <-results:
			if err == nil {
				cancel() // Stop remaining workers after success
				return nil
			}
			// Log errors but continue to wait for a successful result
			log.Debugf("Error: %v", err)
		}
	}
}

// worker fetches data from a single gateway and returns a result.
// It stops if the context is canceled after one success.
func (g *Handler) worker(ctx context.Context, cidStr string, data interface{}, gatewayChan <-chan string, results chan<- error, wg *sync.WaitGroup) {
	defer wg.Done()
	for gateway := range gatewayChan {
		select {
		case <-ctx.Done():
			return // Stop if context is done or canceled
		default:
			// Try to fetch data from the current gateway
			if err := g.fetchFromGateway(ctx, gateway, cidStr, data, results); err == nil {
				// Successful fetch, signal success and exit
				results <- nil
				return
			}
		}
	}
}

func (g *Handler) shufledGateways() []string {
	gateways := make([]string, len(g.gateways))
	copy(gateways, g.gateways)
	for i := range gateways {
		j := i + rand.Intn(len(gateways)-i)
		gateways[i], gateways[j] = gateways[j], gateways[i]
	}
	return gateways
}

var mu sync.Mutex // Mutex to protect concurrent access to the gateways slice

// fetchFromGateway retrieves and processes data from a single gateway.
func (g *Handler) fetchFromGateway(ctx context.Context, gateway, cidStr string, data interface{}, results chan<- error) error {
	gatewayData, err := g.getDataFromGateway(ctx, gateway, cidStr)
	if err != nil {
		select {
		case results <- fmt.Errorf("gateway %s: %v", gateway, err):
		case <-ctx.Done():
		}
		return err
	}

	// Protect access to shared `data`
	mu.Lock()
	defer mu.Unlock()

	if err := g.decodeAndUnmarshal(gateway, gatewayData, data); err != nil {
		select {
		case results <- fmt.Errorf("gateway %s: %v", gateway, err):
		case <-ctx.Done():
		}
		return err
	}

	// Send success (nil error) if the data is fetched and processed correctly
	select {
	case results <- nil:
	case <-ctx.Done():
	}
	return nil
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
