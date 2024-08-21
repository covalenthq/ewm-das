package sampler

import (
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

	publisher "github.com/covalenthq/das-ipfs-pinner/internal/light-client/publisher"
	verifier "github.com/covalenthq/das-ipfs-pinner/internal/light-client/c-kzg-verifier"
	"github.com/ipfs/go-cid"
	ipfs "github.com/ipfs/go-ipfs-api"
	logging "github.com/ipfs/go-log/v2"
)

var log = logging.Logger("light-client")

var DefaultGateways = []string{
	"https://w3s.link/",
	"https://trustless-gateway.link/",
	"https://dweb.link/",
}

// Sampler is a struct that samples data from IPFS and verifies it.
type Sampler struct {
	IPFSShell *ipfs.Shell
	Gateways  []string
	pub *publisher.Publisher

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
func NewSampler(ipfsAddr string, pub *publisher.Publisher) (*Sampler, error) {
	shell := ipfs.NewShell(ipfsAddr)

	if _, _, err := shell.Version(); err != nil {
		return nil, fmt.Errorf("failed to connect to IPFS daemon: %w", err)
	}

	return &Sampler{
		IPFSShell: shell,
		Gateways:  DefaultGateways,
	}, nil
}

// ProcessEvent handles events asynchronously by processing the provided CID.
func (s *Sampler) ProcessEvent(cidStr string) {
	go func(cidStr string) {
		_, err := cid.Decode(cidStr)
		if err != nil {
			log.Errorf("Invalid CID: %v", err)
			return
		}

		var rootNode RootNode
		if err := s.GetData(cidStr, &rootNode); err != nil {
			log.Errorf("Failed to fetch root DAG data: %v", err)
			return
		}

		rowindex := rand.Intn(len(rootNode.Links))
		var links []Link
		if err := s.GetData(rootNode.Links[rowindex].CID, &links); err != nil {
			log.Errorf("Failed to fetch link data: %v", err)
			return
		}

		var data DataMap
		colindex := rand.Intn(len(links))
		if err := s.GetData(links[colindex].CID, &data); err != nil {
			log.Errorf("Failed to fetch data node: %v", err)
			return
		}

		commitment := rootNode.Commitments[rowindex].Nested.Bytes
		proof := data.Proof.Nested.Bytes
		cell := data.Cell.Nested.Bytes
		res, err := verifier.NewKZGVerifier(commitment, proof, cell, uint64(colindex)).Verify()
		if err != nil {
			log.Errorf("Failed to verify proof and cell: %v", err)
			return
		}

		log.Infof("Verification result for [%d, %d]: %v", rowindex, colindex, res)
		s.pub.Publishtocs(cidStr, rowindex, colindex, res, commitment, proof, cell)

	}(cidStr)
}

func (s *Sampler) GetData(cidStr string, data interface{}) error {
	cid, err := cid.Decode(cidStr)
	if err != nil {
		return err
	}

	resultChan := make(chan interface{})
	errorChan := make(chan error)

	// Define a context with timeout to prevent hanging
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Start a goroutine to get data from the IPFS node
	go func() {
		if err := s.IPFSShell.DagGet(cid.String(), &data); err != nil {
			errorChan <- err
			return
		}
		select {
		case resultChan <- data:
		case <-ctx.Done():
		}
	}()

	// Start goroutines to get data from each public gateway
	for _, gateway := range s.Gateways {
		go func(gateway string) {
			gatewayData, err := s.getDataFromGateway(gateway, cid.String())
			if err != nil {
				errorChan <- err
				return
			}

			// Unmarshal JSON data into the provided data interface
			if err := json.Unmarshal(gatewayData, &data); err != nil {
				errorChan <- err
				return
			}

			select {
			case resultChan <- data:
			case <-ctx.Done():
			}

			// populate ipfs node with data
			storedCid, err := s.IPFSShell.DagPut(data, "dag-cbor", "dag-cbor")
			if err != nil {
				errorChan <- err
			}

			if storedCid != cid.String() {
				errorChan <- fmt.Errorf("IPFS node returned different CID: %s", storedCid)
			}
		}(gateway)
	}

	// Wait for the first successful response or all errors
	var finalError error
	successCount := 0
	totalCount := len(s.Gateways) + 1 // +1 for the IPFS node

	for i := 0; i < totalCount; i++ {
		select {
		case <-ctx.Done():
			if successCount > 0 {
				return nil
			}
			return fmt.Errorf("timeout exceeded")

		case result := <-resultChan:
			data = result
			successCount++
			if successCount == 1 {
				// If we get at least one successful response, return nil
				return nil
			}

		case err := <-errorChan:
			finalError = err
		}
	}

	// If we reach here, it means all attempts failed
	if successCount == 0 {
		return finalError
	}

	return nil
}

func (s *Sampler) getDataFromGateway(gateway, cid string) ([]byte, error) {
	// Parse the base gateway URL
	baseURL, err := url.Parse(gateway)
	if err != nil {
		return nil, fmt.Errorf("invalid gateway URL: %v", err)
	}

	// Append the IPFS path and CID
	baseURL.Path = path.Join(baseURL.Path, "ipfs", cid)

	// Add the raw format query parameter
	query := baseURL.Query()
	query.Set("format", "raw")
	baseURL.RawQuery = query.Encode()

	// Perform the HTTP GET request
	resp, err := http.Get(baseURL.String())
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
