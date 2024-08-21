package publisher

import (
	"cloud.google.com/go/pubsub"
	"context"
	"encoding/base64"
	"encoding/json"
	logging "github.com/ipfs/go-log/v2"
	"google.golang.org/api/option"
	"time"
)

var log = logging.Logger("light-client")

type Publisher struct {
	ProjectID   string
	TopicID     string
	Credentials string
	Email       string
}

type Message struct {
	Email       string    `json:"email"`
	SignedAt    time.Time `json:"signed_at"`
	CID         string    `json:"cid"`
	RowIndex    int       `json:"rowindex"`
	ColumnIndex int       `json:"columnindex"`
	Status      bool      `json:"status"`
	Commitment  string    `json:"commitment"`
	Proof       string    `json:"proof"`
	Cell        string    `json:"cell"`
}

// NewPublisher creates a new Publisher instance
func NewPublisher(projectID, topicID, creds, email string) (*Publisher, error) {
	return &Publisher{
		ProjectID:   projectID,
		TopicID:     topicID,
		Credentials: creds,
		Email:       email,
	}, nil
}

// Publish to Pubsub
func (p *Publisher) PublishToCS(cid string, rowindex int, colindex int, booldec bool, commitment []byte, proof []byte, cell []byte) {
	ctx := context.Background()

	// Create a Pub/Sub client using the credentials
	client, err := pubsub.NewClient(ctx, p.ProjectID, option.WithCredentialsFile(p.Credentials))
	if err != nil {
		log.Errorf("Failed to create Pub/Sub client: %v", err)
	}
	defer client.Close()

	// Get a reference to the topic.
	topic := client.Topic(p.TopicID)

	message := Message{
		Email:       p.Email,
		SignedAt:    time.Now(),
		CID:         cid,
		RowIndex:    rowindex,
		ColumnIndex: colindex,
		Status:      booldec,
		Commitment:  base64.StdEncoding.EncodeToString(commitment),
		Proof:       base64.StdEncoding.EncodeToString(proof),
		Cell:        base64.StdEncoding.EncodeToString(cell),
	}

	// Marshal the message into JSON.
	messageData, err := json.Marshal(message)
	if err != nil {
		log.Errorf("Failed to marshal message: %v", err)
	}

	// Publish a message.
	result := topic.Publish(ctx, &pubsub.Message{
		Data: messageData,
	})

	// Block until the result is returned and a server-generated ID is returned for the published message.
	id, err := result.Get(ctx)
	if err != nil {
		log.Errorf("Failed to publish message: %v", err)
	} else {
		log.Infof("Published a message with a message ID: %s\n", id)
	}

}
