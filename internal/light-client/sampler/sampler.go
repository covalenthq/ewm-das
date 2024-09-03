package sampler

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"net/url"
	"path"
	"strings"
	"time"

	verifier "github.com/covalenthq/das-ipfs-pinner/internal/light-client/c-kzg-verifier"
	publisher "github.com/covalenthq/das-ipfs-pinner/internal/light-client/publisher"
	ipldencoder "github.com/covalenthq/das-ipfs-pinner/internal/pinner/ipld-encoder"
	"github.com/ipfs/go-cid"
	ipfs "github.com/ipfs/go-ipfs-api"
	logging "github.com/ipfs/go-log/v2"
	"github.com/ipld/go-ipld-prime/codec/dagjson"
)

var log = logging.Logger("light-client")

var DefaultGateways = []string{
	"https://w3s.link/",
	"https://trustless-gateway.link/",
	"https://dweb.link/",
}

// Sampler is a struct that samples data from IPFS and verifies it.
type Sampler struct {
	ipfsShell     *ipfs.Shell
	gateways      []string
	pub           *publisher.Publisher
	samplingDelay uint
	samplingFn    func(int, int, float64) int
}

// Link represents a link to another CID in IPFS.
type Link struct {
	CID string `json:"/"`
}

// RootNode represents a DAG node containing metadata and links.
type RootNode struct {
	Version     string     `json:"version"`
	Size        int        `json:"size"`
	Length      int        `json:"length"`
	Links       []Link     `json:"links"`
	Commitments []InnerMap `json:"commitments"`
}

// NestedBytes holds the base64 decoded bytes.
type NestedBytes struct {
	Bytes []byte `json:"bytes"`
}

// InnerMap represents a nested structure containing bytes data.
type InnerMap struct {
	Nested NestedBytes `json:"/"`
}

// DataMap represents a cell and its proof, used in the DAG structure.
type DataMap struct {
	Cell  InnerMap `json:"cell"`
	Proof InnerMap `json:"proof"`
}

type errorContext struct {
	Err     error
	Context string
}

type resultContext struct {
	Result  interface{}
	Context string
}

// UnmarshalJSON handles base64 decoding directly into the Bytes field.
func (n *NestedBytes) UnmarshalJSON(data []byte) error {
	var aux struct {
		Bytes string `json:"bytes"`
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return fmt.Errorf("error unmarshaling bytes: %w", err)
	}

	bytesWithPadding := ensureBase64Padding(aux.Bytes)

	decodedBytes, err := base64.StdEncoding.DecodeString(bytesWithPadding)
	if err != nil {
		return fmt.Errorf("error decoding base64 string: %w", err)
	}

	n.Bytes = decodedBytes
	return nil
}

// NewSampler creates a new Sampler instance and checks the connection to the IPFS daemon.
func NewSampler(ipfsAddr string, samplingDelay uint, pub *publisher.Publisher) (*Sampler, error) {
	shell := ipfs.NewShell(ipfsAddr)

	if _, _, err := shell.Version(); err != nil {
		return nil, fmt.Errorf("failed to connect to IPFS daemon: %w", err)
	}

	return &Sampler{
		ipfsShell:     shell,
		gateways:      DefaultGateways,
		pub:           pub,
		samplingDelay: samplingDelay,
		samplingFn:    CalculateSamplesNeeded,
	}, nil
}

// ProcessEvent handles events asynchronously by processing the provided CID.
func (s *Sampler) ProcessEvent(cidStr string, blockHeight uint64) {
	go func(cidStr string) {
		rawCid, err := cid.Decode(cidStr)
		if err != nil {
			log.Errorf("Invalid CID: %v", err)
			return
		}

		if rawCid.Prefix().Codec != cid.DagCBOR {
			log.Debugf("Unsupported CID codec: %v. Skipping", rawCid.Prefix().Codec)
			return
		}

		log.Debugf("Processing event for CID [%s] is deferred for %d min", cidStr, s.samplingDelay/60)
		time.Sleep(time.Duration(s.samplingDelay) * time.Second)
		log.Debugf("Processing event for CID [%s] ...", cidStr)

		var rootNode RootNode
		if err := s.GetData(cidStr, &rootNode); err != nil {
			log.Errorf("Failed to fetch root DAG data: %v", err)
			return
		}

		sampleIterations := s.samplingFn(rootNode.Size, rootNode.Size/2, 0.95)

		for rowIndex, rowLink := range rootNode.Links {
			var links []Link
			if err := s.GetData(rowLink.CID, &links); err != nil {
				log.Errorf("Failed to fetch link data: %v", err)
				return
			}

			// Track sampled column indices to avoid duplicates
			sampledCols := make(map[int]bool)

			for i := 0; i < sampleIterations; i++ {
				// Find a unique column index that hasn't been sampled yet
				var colIndex int
				for {
					colIndex = rand.Intn(len(links))
					if !sampledCols[colIndex] {
						sampledCols[colIndex] = true
						break
					}
				}

				var data DataMap
				if err := s.GetData(links[colIndex].CID, &data); err != nil {
					log.Errorf("Failed to fetch data node: %v", err)
					return
				}

				commitment := rootNode.Commitments[rowIndex].Nested.Bytes
				proof := data.Proof.Nested.Bytes
				cell := data.Cell.Nested.Bytes
				res, err := verifier.NewKZGVerifier(commitment, proof, cell, uint64(colIndex)).Verify()
				if err != nil {
					log.Errorf("Failed to verify proof and cell: %v", err)
					return
				}

				log.Infof("cell=[%2d,%3d], verified=%-5v, cid=%-40v", rowIndex, colIndex, res, cidStr)

				if err := s.pub.PublishToCS(cidStr, rowIndex, colIndex, res, commitment, proof, cell, blockHeight); err != nil {
					log.Errorf("Failed to publish to Cloud Storage: %v", err)
					return
				}
			}
		}
	}(cidStr)
}

func (s *Sampler) GetData(cidStr string, data interface{}) error {
	cid, err := cid.Decode(cidStr)
	if err != nil {
		return err
	}

	resultChan := make(chan resultContext)
	errorChan := make(chan errorContext)

	// Define a context with timeout to prevent hanging
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Start a goroutine to get data from the IPFS node
	go func() {
		if err := s.ipfsShell.DagGet(cid.String(), &data); err != nil {
			errorChan <- errorContext{Err: err, Context: "IPFS node"}
			return
		}
		select {
		case resultChan <- resultContext{Result: data, Context: "IPFS node"}:
			cancel()
		case <-ctx.Done():
		}
	}()

	// Start goroutines to get data from each public gateway
	for _, gateway := range s.gateways {
		go func(gateway string) {
			gatewayData, err := s.getDataFromGateway(ctx, gateway, cid.String())
			if err != nil {
				errorChan <- errorContext{Err: err, Context: gateway}
				return
			}

			// Decode the data into an IPLD node from CBOR
			node, err := ipldencoder.DecodeNode(gatewayData)
			if err != nil {
				errorChan <- errorContext{Err: err, Context: gateway}
				return
			}

			// Encode the IPLD node into JSON
			var jsonData bytes.Buffer
			if err := dagjson.Encode(node, &jsonData); err != nil {
				errorChan <- errorContext{Err: err, Context: gateway}
				return
			}

			// Unmarshal JSON data into the provided data interface
			if err := json.Unmarshal(jsonData.Bytes(), data); err != nil {
				errorChan <- errorContext{Err: err, Context: gateway}
				return
			}

			select {
			case resultChan <- resultContext{Result: data, Context: fmt.Sprintf("gateway %s for %s", gateway, cid.String())}:
				cancel()
			case <-ctx.Done():
			}

			// populate ipfs node with data
			// TODO: deduce the correct format from the CID
			storedCid, err := s.ipfsShell.DagPut(data, "dag-cbor", "dag-cbor")
			if err != nil {
				errorChan <- errorContext{Err: err, Context: "IPFS node"}
			}

			if storedCid != cid.String() {
				errorChan <- errorContext{Err: fmt.Errorf("IPFS node returned different CID: %s", storedCid), Context: "Result CID"}
			}
		}(gateway)
	}

	// Wait for the first successful response or all errors
	var finalError error
	successCount := 0
	totalCount := len(s.gateways) + 1 // +1 for the IPFS node

	for i := 0; i < totalCount; i++ {
		select {
		case <-ctx.Done():
			if successCount > 0 {
				return nil
			}
			return fmt.Errorf("timeout exceeded")

		case result := <-resultChan:
			log.Debugf("Data fetched from %s", result.Context)
			data = result.Result
			successCount++
			if successCount == 1 {
				// If we get at least one successful response, return nil
				return nil
			}

		case errCxt := <-errorChan:
			log.Debugf("Error getting data from %s: %v", errCxt.Context, errCxt.Err)
			finalError = errCxt.Err
		}
	}

	// If we reach here, it means all attempts failed
	if successCount == 0 {
		return finalError
	}

	return nil
}

func (s *Sampler) getDataFromGateway(ctx context.Context, gateway, cid string) ([]byte, error) {
	// Parse the base gateway URL
	baseURL, err := url.Parse(gateway)
	if err != nil {
		return nil, err
	}

	// Append the IPFS path and CID
	baseURL.Path = path.Join(baseURL.Path, "ipfs", cid)

	// Add the raw format query parameter
	query := baseURL.Query()
	query.Set("format", "raw")
	baseURL.RawQuery = query.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, baseURL.String(), nil)
	if err != nil {
		return nil, err
	}

	// Perform the HTTP GET request
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("gateway %s returned status %d", gateway, resp.StatusCode)
	}

	return io.ReadAll(resp.Body)
}

// ensureBase64Padding ensures the base64 string has correct padding.
func ensureBase64Padding(encoded string) string {
	paddingNeeded := len(encoded) % 4
	if paddingNeeded > 0 {
		padding := 4 - paddingNeeded
		encoded += strings.Repeat("=", padding)
	}
	return encoded
}
