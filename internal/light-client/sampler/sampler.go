package sampler

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"math/rand"
	"strings"
	"context"
	"time"
	
	verifier "github.com/covalenthq/das-ipfs-pinner/internal/light-client/c-kzg-verifier"
	"github.com/ipfs/go-cid"
	ipfs "github.com/ipfs/go-ipfs-api"
	logging "github.com/ipfs/go-log/v2"
	"cloud.google.com/go/pubsub"

)

var log = logging.Logger("light-client")

// Sampler is a struct that samples data from IPFS and verifies it.
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
func (s *Sampler) ProcessEvent(cidStr string, projectId string, topicId string) {
	go func(cidStr string) {
		_, err := cid.Decode(cidStr)
		if err != nil {
			log.Errorf("Invalid CID: %v", err)
			return
		}

		var rootNode RootNode
		if err := s.IPFSShell.DagGet(cidStr, &rootNode); err != nil {
			log.Errorf("Failed to fetch root DAG data: %v", err)
			return
		}

		rowindex := rand.Intn(len(rootNode.Links))
		var links []Link
		if err := s.IPFSShell.DagGet(rootNode.Links[rowindex].CID, &links); err != nil {
			log.Errorf("Failed to fetch link data: %v", err)
			return
		}

		var data DataMap
		colindex := rand.Intn(len(links))
		if err := s.IPFSShell.DagGet(links[colindex].CID, &data); err != nil {
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

		log.Infof("Verification result: %v", res)

		//publish message here
		publishtocs(projectId, topicId, cidStr, rowindex, colindex, res)

	}(cidStr)
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


// Publish to Pubsub
func publishtocs(projectId string, topicId string, cid string, rowindex int, colindex int, booldec bool) {
	ctx := context.Background()

	// Create a Pub/Sub client.
	client, err := pubsub.NewClient(ctx, projectId)
	if err != nil {
		log.Fatalf("Failed to create Pub/Sub client: %v", err)
	}
	defer client.Close()

	// Get a reference to the topic.
	topic := client.Topic(topicId)

	// Define the message payload with exported field names.
	message := struct {
		SignedAt  time.Time `json:"signed_at"`
		CID       string    `json:"cid"`
		RowIndex  int       `json:"rowindex"`
		ColumnIndex int     `json:"columnindex"`
		Status    bool      `json:"status"`
	}{
		SignedAt:  time.Now(), // Current timestamp
		CID:       cid,
		RowIndex:  rowindex,
		ColumnIndex: colindex,
		Status:    booldec,
	}

	// Marshal the message into JSON.
	messageData, err := json.Marshal(message)
	if err != nil {
		log.Fatalf("Failed to marshal message: %v", err)
	}

	// Publish a message.
	result := topic.Publish(ctx, &pubsub.Message{
		Data: messageData,
	})

	// Block until the result is returned and a server-generated ID is returned for the published message.
	id, err := result.Get(ctx)
	if err != nil {
		log.Fatalf("Failed to publish message: %v", err)
	} else {
		log.Infof("Published a message with a message ID: %s\n", id)
	}
}