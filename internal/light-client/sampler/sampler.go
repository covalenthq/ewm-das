package sampler

import (
	"encoding/base64"
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

		rootNodeData, err := s.fetchDagData(c.String())
		if err != nil {
			log.Errorf("Failed to fetch root DAG data: %v", err)
			return
		}

		randomLinkData, err := s.fetchDataFromRandomLink(rootNodeData)
		if err != nil {
			log.Errorf("Failed to fetch data from random link: %v", err)
			return
		}

		var randomIndex int
		randomIndexData, err := s.fetchDataFromRandomIndex(randomLinkData, &randomIndex)
		if err != nil {
			log.Errorf("Failed to fetch data from random index: %v", err)
			return
		}

		// Process commitments and then fetch "proof" and "cell"
		s.handleProofAndCell(randomIndexData)
	}(cidStr) // Pass cidStr to the anonymous function
}

// fetchDagData retrieves the DAG data for a given CID
func (s *Sampler) fetchDagData(cidStr string) (interface{}, error) {
	var dagData interface{}
	if err := s.IPFSShell.DagGet(cidStr, &dagData); err != nil {
		return nil, err
	}

	return dagData, nil
}

// fetchDataFromRandomLink selects a random link and retrieves the DAG data
func (s *Sampler) fetchDataFromRandomLink(data interface{}) (interface{}, error) {
	dataMap, ok := data.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("data is not a map[string]interface{}")
	}

	links, ok := dataMap["links"].([]interface{})
	if !ok || len(links) == 0 {
		return nil, fmt.Errorf("links is not a valid []interface{} or is empty")
	}

	selectedLink := links[rand.Intn(len(links))]

	linkMap, ok := selectedLink.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("selected link is not a map[string]interface{}")
	}

	nextCIDStr, ok := linkMap["/"].(string)
	if !ok {
		return nil, fmt.Errorf(`linkMap["/"] is not a string`)
	}

	return s.fetchDagData(nextCIDStr)
}

// fetchDataFromRandomIndex selects a random index and retrieves the DAG data
func (s *Sampler) fetchDataFromRandomIndex(data interface{}, randomIndex *int) (interface{}, error) {
	switch v := data.(type) {
	case map[string]interface{}:
		length, ok := v["length"].(float64)
		if !ok {
			return nil, fmt.Errorf("length is not a float64")
		}

		links, ok := v["links"].([]interface{})
		if !ok || len(links) == 0 {
			return nil, fmt.Errorf("links is not a valid []interface{} or is empty")
		}

		*randomIndex = rand.Intn(int(length))

		if *randomIndex >= len(links) {
			return nil, fmt.Errorf("random index is out of bounds")
		}

		selectedLink := links[*randomIndex]

		linkMap, ok := selectedLink.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("selected link is not a map[string]interface{}")
		}

		finalCIDStr, ok := linkMap["/"].(string)
		if !ok {
			return nil, fmt.Errorf(`linkMap["/"] is not a string`)
		}

		return s.fetchDagData(finalCIDStr)

	case []interface{}:
		randomIndex := rand.Intn(len(v))
		return v[randomIndex], nil

	default:
		return nil, fmt.Errorf("data is not a map[string]interface{} or []interface{}")
	}
}

func (s *Sampler) handleProofAndCell(data interface{}) (proofBytes, cellBytes []byte, err error) {
	dataCID, ok := data.(map[string]interface{})["/"].(string)
	if !ok {
		log.Error("Data is not a valid CID string")
		return
	}

	commitmentData, err := s.fetchDagData(dataCID)
	if err != nil {
		log.Errorf("Failed to retrieve commitment data: %v", err)
		return
	}

	return s.processProofAndCell(commitmentData)
}

// processProofAndCell handles the "proof" and "cell" fields from the DAG data
func (s *Sampler) processProofAndCell(data interface{}) (proofBytes, cellBytes []byte, err error) {
	dataMap, ok := data.(map[string]interface{})
	if !ok {
		log.Error("Proof and cell data is not a map[string]interface{}")
		return
	}

	// Process "proof"
	proofData, ok := dataMap["proof"].(map[string]interface{})
	if !ok {
		log.Error(`"proof" field is not a map[string]interface{}`)
		return
	}

	proofBytes, err = extractBytesFromNestedMap(proofData)
	if err != nil {
		log.Errorf("Failed to extract proof bytes: %v", err)
		return
	}
	log.Debugf("Proof (bytes): %x\n", proofBytes)

	// Process "cell"
	cellData, ok := dataMap["cell"].(map[string]interface{})
	if !ok {
		log.Error(`"cell" field is not a map[string]interface{}`)
		return
	}

	cellBytes, err = extractBytesFromNestedMap(cellData)
	if err != nil {
		log.Errorf("Failed to extract cell bytes: %v", err)
		return
	}
	log.Debugf("Cell (bytes): %x\n", cellBytes)

	return proofBytes, cellBytes, nil
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

// extractBytesFromNestedMap navigates through the nested map structure and extracts the bytes.
func extractBytesFromNestedMap(data map[string]interface{}) ([]byte, error) {
	innerMap, ok := data["/"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf(`expected map with key "/" but got %T`, data["/"])
	}

	bytesEncoded, ok := innerMap["bytes"].(string)
	if !ok {
		return nil, fmt.Errorf(`expected "bytes" field to be a string but got %T`, innerMap["bytes"])
	}

	// Ensure correct padding before decoding
	bytesEncoded = addBase64Padding(bytesEncoded)

	decodedBytes, err := base64.StdEncoding.DecodeString(bytesEncoded)
	if err != nil {
		return nil, fmt.Errorf("failed to decode base64: %w", err)
	}

	return decodedBytes, nil
}
