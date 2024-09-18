package internal

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"
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

type ScheduleRequest struct {
	Cid      string
	ClientId string
	ChainId  uint64
	BlockNum uint64
}
