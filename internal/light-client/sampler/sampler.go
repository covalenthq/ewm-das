package sampler

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"math/rand"
	"strings"

	logging "github.com/ipfs/go-log/v2"

	"github.com/ipfs/go-cid"
	ipfs "github.com/ipfs/go-ipfs-api"
)

var log = logging.Logger("light-client")

type Sampler struct {
	IPFSShell *ipfs.Shell
}

// Struct for the Link
type Link struct {
	CID string `json:"/"`
}

// Struct for the DAG node containing links
type RootNode struct {
	Version string `json:"version"`
	Size    int    `json:"size"`
	Length  int    `json:"length"`
	Links   []Link `json:"links"`
}

type NestedBytes struct {
	Bytes []byte `json:"bytes"`
}

type InnerMap struct {
	Nested NestedBytes `json:"/"`
}

type DataMap struct {
	Cell  InnerMap `json:"cell"`
	Proof InnerMap `json:"proof"`
}

// UnmarshalJSON custom unmarshaler to handle base64 decoding directly into the Bytes field.
func (n *NestedBytes) UnmarshalJSON(data []byte) error {
	// Create an anonymous struct to unmarshal the raw string
	var aux struct {
		Bytes string `json:"bytes"`
	}

	// Unmarshal the JSON into the auxiliary struct
	if err := json.Unmarshal(data, &aux); err != nil {
		return fmt.Errorf("error unmarshaling bytes: %w", err)
	}

	bytesEncoded := addBase64Padding(aux.Bytes)

	// Decode the base64 string
	decoded, err := base64.StdEncoding.DecodeString(bytesEncoded)
	if err != nil {
		return fmt.Errorf("error decoding base64 string: %w", err)
	}

	// Assign the decoded bytes to the Bytes field
	n.Bytes = decoded
	return nil
}

// NewSampler creates a new Sampler instance and checks if the IPFS daemon is running.
func NewSampler(ipfsAddr string) (*Sampler, error) {
	shell := ipfs.NewShell(ipfsAddr)

	// Check if the IPFS daemon is running by getting the node's version
	_, _, err := shell.Version()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to IPFS daemon: %w", err)
	}

	return &Sampler{
		IPFSShell: shell,
	}, nil
}

// ProcessEvent handles the event received from the EventListener
func (s *Sampler) ProcessEvent(cidStr string) {
	// Spawn a goroutine to handle the event processing asynchronously
	go func(cidStr string) {
		c, err := cid.Decode(cidStr)
		if err != nil {
			log.Errorf("Invalid CID: %v", err)
			return
		}

		rootNode, err := s.fetchRootNode(c.String())
		if err != nil {
			log.Errorf("Failed to fetch root DAG data: %v", err)
			return
		}

		cols, err := s.fetchRandomRow(rootNode)
		if err != nil {
			log.Errorf("Failed to fetch data from random link: %v", err)
			return
		}

		var randomIndex int
		_, err = s.fetchData(cols, &randomIndex)
		if err != nil {
			log.Errorf("Failed to fetch data from random index: %v", err)
			return
		}

		// Implement the logic to send the proof and cell to the service
		// sendProofAndCell(proof, cell)
	}(cidStr) // Pass cidStr to the anonymous function
}

func (s *Sampler) fetchRootNode(cidStr string) (*RootNode, error) {
	var root RootNode
	if err := s.IPFSShell.DagGet(cidStr, &root); err != nil {
		return nil, err
	}

	return &root, nil
}

// fetchDagData retrieves the DAG data for a given CID and unmarshals it into a slice of Link structs.
func (s *Sampler) fetchCols(cidStr string) ([]Link, error) {
	var dagData []Link
	if err := s.IPFSShell.DagGet(cidStr, &dagData); err != nil {
		return nil, err
	}

	return dagData, nil
}

// fetchDataNode retrieves the DAG data for a given CID
func (s *Sampler) fetchDataNode(cidStr string) (*DataMap, error) {
	var data DataMap
	if err := s.IPFSShell.DagGet(cidStr, &data); err != nil {
		return nil, err
	}

	return &data, nil
}

// fetchDataFromRandomLink selects a random link and retrieves the DAG data
func (s *Sampler) fetchRandomRow(root *RootNode) ([]Link, error) {
	selectedLink := root.Links[rand.Intn(len(root.Links))]

	return s.fetchCols(selectedLink.CID)
}

// fetchData selects a random index and retrieves the DAG data
func (s *Sampler) fetchData(links []Link, randomIndex *int) (*DataMap, error) {
	// Fetch the data from the random index
	*randomIndex = rand.Intn(int(len(links)))
	data, err := s.fetchDataNode(links[*randomIndex].CID)
	if err != nil {
		return nil, err
	}

	return data, nil
}

// addBase64Padding ensures the base64 string has correct padding
func addBase64Padding(encoded string) string {
	// Calculate the padding required
	missingPadding := len(encoded) % 4
	if missingPadding > 0 {
		padding := 4 - missingPadding
		encoded += strings.Repeat("=", padding)
	}
	return encoded
}
