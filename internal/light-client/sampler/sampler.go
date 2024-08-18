package sampler

import (
	"encoding/base64"
	"fmt"
	"math/rand"

	logging "github.com/ipfs/go-log/v2"

	"github.com/ipfs/go-cid"
	ipfs "github.com/ipfs/go-ipfs-api"
)

var log = logging.Logger("light-client")

type Sampler struct {
	IPFSShell *ipfs.Shell
}

func NewSampler(ipfsAddr string) *Sampler {
	shell := ipfs.NewShell(ipfsAddr)
	return &Sampler{
		IPFSShell: shell,
	}
}

// ProcessEvent processes the event received from the EventListener
func (s *Sampler) ProcessEvent(cidStr string) {
	c, err := cid.Decode(cidStr)
	if err != nil {
		log.Errorf("Invalid CID: %v", err)
		return
	}

	rootData, err := s.retrieveDagData(c.String())
	if err != nil {
		log.Errorf("Failed to retrieve root DAG data: %v", err)
		return
	}

	nextData, err := s.retrieveDataFromRandomLink(rootData)
	if err != nil {
		log.Errorf("Failed to retrieve data from random link: %v", err)
		return
	}

	finalData, err := s.retrieveDataFromRandomIndex(nextData)
	if err != nil {
		log.Errorf("Failed to retrieve data from random index: %v", err)
		return
	}

	// Process commitments and then retrieve "proof" and "cell"
	s.processCommitments(finalData)
}

// retrieveDagData retrieves the DAG data for a given CID
func (s *Sampler) retrieveDagData(cidStr string) (interface{}, error) {
	var dagData interface{}
	if err := s.IPFSShell.DagGet(cidStr, &dagData); err != nil {
		return nil, err
	}

	return dagData, nil
}

// retrieveDataFromRandomLink selects a random link and retrieves the DAG data
func (s *Sampler) retrieveDataFromRandomLink(data interface{}) (interface{}, error) {
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
		return nil, fmt.Errorf("selectedLink is not a map[string]interface{}")
	}

	nextCIDStr, ok := linkMap["/"].(string)
	if !ok {
		return nil, fmt.Errorf(`linkMap["/"] is not a string`)
	}

	return s.retrieveDagData(nextCIDStr)
}

// retrieveDataFromRandomIndex selects a random index and retrieves the DAG data
func (s *Sampler) retrieveDataFromRandomIndex(data interface{}) (interface{}, error) {
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

		randomIndex := rand.Intn(int(length))

		if randomIndex >= len(links) {
			return nil, fmt.Errorf("randomIndex is out of bounds")
		}

		nextLink := links[randomIndex]

		linkMap, ok := nextLink.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("nextLink is not a map[string]interface{}")
		}

		finalCIDStr, ok := linkMap["/"].(string)
		if !ok {
			return nil, fmt.Errorf(`linkMap["/"] is not a string`)
		}

		return s.retrieveDagData(finalCIDStr)

	case []interface{}:
		randomIndex := rand.Intn(len(v))

		return v[randomIndex], nil

	default:
		return nil, fmt.Errorf("data is not a map[string]interface{} or []interface{}")
	}
}

// processCommitments retrieves data from a given CID, processes the commitments, and retrieves "proof" and "cell"
func (s *Sampler) processCommitments(data interface{}) {
	// Assume data is a CID string; retrieve the corresponding DAG data
	dataCID, ok := data.(map[string]interface{})["/"].(string)
	if !ok {
		log.Error("data is not a valid CID string")
		return
	}

	cellData, err := s.retrieveDagData(dataCID)
	if err != nil {
		log.Errorf("Failed to retrieve commitment data: %v", err)
		return
	}

	s.processProofAndCell(cellData)
}

// processProofAndCell processes the "proof" and "cell" fields from the DAG data
func (s *Sampler) processProofAndCell(data interface{}) {
	dataMap, ok := data.(map[string]interface{})
	if !ok {
		log.Error("proof and cell data is not a map[string]interface{}")
		return
	}

	// Process "proof"
	proofData, ok := dataMap["proof"].(map[string]interface{})
	if !ok {
		log.Error(`"proof" field is not a map[string]interface{}`)
		return
	}

	proofBytes, err := extractBytesFromNestedMap(proofData)
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

	cellBytes, err := extractBytesFromNestedMap(cellData)
	if err != nil {
		log.Errorf("Failed to extract cell bytes: %v", err)
		return
	}
	log.Debugf("Cell (bytes): %x\n", cellBytes)
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

	decodedBytes, err := base64.StdEncoding.DecodeString(bytesEncoded)
	if err != nil {
		return nil, fmt.Errorf("failed to decode base64: %w", err)
	}

	return decodedBytes, nil
}
