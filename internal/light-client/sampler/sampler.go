package sampler

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"math/rand"
	"strings"

	"github.com/ipfs/go-cid"
	ipfs "github.com/ipfs/go-ipfs-api"
	logging "github.com/ipfs/go-log/v2"
)

var log = logging.Logger("light-client")

type Sampler struct {
	IPFSShell *ipfs.Shell
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
func NewSampler(ipfsAddr string) (*Sampler, error) {
	shell := ipfs.NewShell(ipfsAddr)

	if _, _, err := shell.Version(); err != nil {
		return nil, fmt.Errorf("failed to connect to IPFS daemon: %w", err)
	}

	return &Sampler{
		IPFSShell: shell,
	}, nil
}

// ProcessEvent handles events asynchronously by processing the provided CID.
func (s *Sampler) ProcessEvent(cidStr string) {
	go func(cidStr string) {
		cidDecoded, err := cid.Decode(cidStr)
		if err != nil {
			log.Errorf("Invalid CID: %v", err)
			return
		}

		rootNode, err := s.getRootNode(cidDecoded.String())
		if err != nil {
			log.Errorf("Failed to fetch root DAG data: %v", err)
			return
		}

		randomLinkIndex := getRandomIndex(len(rootNode.Links))
		links, err := s.getLinks(rootNode.Links[randomLinkIndex].CID)
		if err != nil {
			log.Errorf("Failed to fetch link data: %v", err)
			return
		}

		randomLinkIndex = getRandomIndex(len(links))
		_, err = s.getDataNode(links[randomLinkIndex].CID)
		if err != nil {
			log.Errorf("Failed to fetch data node: %v", err)
			return
		}

		// Implement logic to send proof and cell to the service
		// sendProofAndCell(proof, cell)
	}(cidStr)
}

func (s *Sampler) getRootNode(cidStr string) (*RootNode, error) {
	var rootNode RootNode
	if err := s.IPFSShell.DagGet(cidStr, &rootNode); err != nil {
		return nil, err
	}

	return &rootNode, nil
}

// getLinks retrieves the DAG data for a given CID and unmarshals it into a slice of Link structs.
func (s *Sampler) getLinks(cidStr string) ([]Link, error) {
	var links []Link
	if err := s.IPFSShell.DagGet(cidStr, &links); err != nil {
		return nil, err
	}

	return links, nil
}

// getDataNode retrieves the DAG data for a given CID and unmarshals it into a DataMap struct.
func (s *Sampler) getDataNode(cidStr string) (*DataMap, error) {
	var data DataMap
	if err := s.IPFSShell.DagGet(cidStr, &data); err != nil {
		return nil, err
	}

	return &data, nil
}

// getRandomIndex returns a random index for a given length.
func getRandomIndex(length int) int {
	return rand.Intn(length)
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
