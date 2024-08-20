package publisher

import (
	"encoding/json"
	"io/ioutil"
	"context"
	"time"
	logging "github.com/ipfs/go-log/v2"
	"cloud.google.com/go/pubsub"
	"google.golang.org/api/option"

)

var log = logging.Logger("light-client")

type Credentials struct {
    ClientEmail string `json:"client_email"`
}


// Publish to Pubsub
func Publishtocs(projectId string, topicId string, gcpcreds string, cid string, rowindex int, colindex int, booldec bool) {
	ctx := context.Background()

	// Read and parse the JSON file to get client_email
	data, err := ioutil.ReadFile(gcpcreds)
	if err != nil {
		log.Fatalf("Failed to read the credentials file: %v", err)
	}

	var creds Credentials
	if err := json.Unmarshal(data, &creds); err != nil {
		log.Fatalf("Failed to unmarshal the credentials JSON: %v", err)
	}

	// Print the client email
	log.Infof("Client Email: %s\n", creds.ClientEmail)

	// Create a Pub/Sub client.
	client, err := pubsub.NewClient(ctx, projectId, option.WithCredentialsFile(gcpcreds))
	if err != nil {
		log.Fatalf("Failed to create Pub/Sub client: %v", err)
	}
	defer client.Close()

	// Get a reference to the topic.
	topic := client.Topic(topicId)

	// Define the message payload with exported field names.
	message := struct {
		Email     string    `json:"email"`
		SignedAt  time.Time `json:"signed_at"`
		CID       string    `json:"cid"`
		RowIndex  int       `json:"rowindex"`
		ColumnIndex int     `json:"columnindex"`
		Status    bool      `json:"status"`
	}{
		Email:     creds.ClientEmail,
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