package internal

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

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

// ensureBase64Padding ensures the base64 string has correct padding.
func ensureBase64Padding(encoded string) string {
	paddingNeeded := len(encoded) % 4
	if paddingNeeded > 0 {
		padding := 4 - paddingNeeded
		encoded += strings.Repeat("=", padding)
	}
	return encoded
}

type SamplingRequest struct {
	Cid            string `json:"cid"`
	ClientIdentity string `json:"client_identity"`
	ChainId        uint64 `json:"chain_id"`
	BlockNum       uint64 `json:"block_num"`
}

type StoreRequest struct {
	SamplingReqest    SamplingRequest
	SamplingSignature string    `json:"sampling_signature"`
	SignedAt          time.Time `json:"signed_at"`
	BlobIndex         int       `json:"blob_index"`
	CellIndex         int       `json:"cell_index"`
	Status            bool      `json:"status"`
	Commitment        string    `json:"commitment"`
	Proof             string    `json:"proof"`
	Cell              string    `json:"cell"`
	Version           string    `json:"version"`
}

type Workload struct {
	Hash         string `json:"hash"`
	ChainID      int    `json:"chain_id"`
	BlockHeight  int    `json:"block_height"`
	BlockHash    string `json:"block_hash"`
	SpecimenHash string `json:"specimen_hash"`
	StorageURL   string `json:"storage_url"`
	Challenge    string `json:"challenge"`
}

// Define the top-level struct
type WorkloadResponse struct {
	NextUpdate time.Time  `json:"next_update"`
	Workloads  []Workload `json:"workloads"`
}

type StoreRequest2 struct {
	WorkloadRequest Workload  `json:"workload"`
	Timestamp       time.Time `json:"timestamp"`
	BlobIndex       int       `json:"blob_index"`
	CellIndex       int       `json:"cell_index"`
	Commitment      string    `json:"commitment"`
	Proof           string    `json:"proof"`
	Cell            string    `json:"cell"`
	Version         string    `json:"version"`
}
