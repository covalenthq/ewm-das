package publisher

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"time"

	"cloud.google.com/go/pubsub"
	"github.com/covalenthq/das-ipfs-pinner/common"
	"github.com/covalenthq/das-ipfs-pinner/internal/light-client/utils"
	"google.golang.org/api/option"
)

type Publisher struct {
	projectID       string
	topicID         string
	credentialsFile string
	identity        *utils.Identity
}

type message struct {
	ClientId    string    `json:"client_id"`
	SignedAt    time.Time `json:"signed_at"`
	CID         string    `json:"cid"`
	RowIndex    int       `json:"rowindex"`
	ColumnIndex int       `json:"columnindex"`
	Status      bool      `json:"status"`
	Commitment  string    `json:"commitment"`
	Proof       string    `json:"proof"`
	Cell        string    `json:"cell"`
	BlockHeight uint64    `json:"block_height"`
	Version     string    `json:"version"`
}

// Define a struct with only the `project_id` field
type serviceAccount struct {
	ProjectID string `json:"project_id"`
}

// NewPublisher creates a new Publisher instance
func NewPublisher(topicID, credsFile string, identity *utils.Identity) (*Publisher, error) {
	file, err := os.Open(credsFile)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// Read the file contents into a byte slice
	data, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	// Unmarshal only the project_id field
	var account serviceAccount
	err = json.Unmarshal(data, &account)
	if err != nil {
		return nil, err
	}

	return &Publisher{
		projectID:       account.ProjectID,
		topicID:         topicID,
		credentialsFile: credsFile,
		identity:        identity,
	}, nil
}

// Publish to Pubsub
func (p *Publisher) PublishToCS(cid string, rowIndex int, colIndex int, status bool, commitment []byte, proof []byte, cell []byte, blockHeight uint64) error {
	ctx := context.Background()

	// Create a Pub/Sub client using the credentials
	client, err := pubsub.NewClient(ctx, p.projectID, option.WithCredentialsFile(p.credentialsFile))
	if err != nil {
		return err
	}
	defer client.Close()

	// Get a reference to the topic.
	topic := client.Topic(p.topicID)

	message := message{
		ClientId:    p.identity.GetAddress().Hex(),
		SignedAt:    time.Now(),
		CID:         cid,
		RowIndex:    rowIndex,
		ColumnIndex: colIndex,
		Status:      status,
		Commitment:  base64.StdEncoding.EncodeToString(commitment),
		Proof:       base64.StdEncoding.EncodeToString(proof),
		Cell:        base64.StdEncoding.EncodeToString(cell),
		BlockHeight: blockHeight,
		Version:     fmt.Sprintf("%s-%s", common.Version, common.GitCommit),
	}

	// Marshal the message into JSON.
	messageData, err := json.Marshal(message)
	if err != nil {
		return err
	}

	// Publish a message.
	result := topic.Publish(ctx, &pubsub.Message{
		Data: messageData,
	})

	// Block until the result is returned and a server-generated ID is returned for the published message.
	if _, err = result.Get(ctx); err != nil {
		return err
	}

	return nil
}
