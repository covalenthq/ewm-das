package publisher

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/covalenthq/das-ipfs-pinner/internal"
	"github.com/covalenthq/das-ipfs-pinner/internal/light-client/utils"
)

type Publisher struct {
	collectApi string
	identity   *utils.Identity
}

// NewPublisher creates a new Publisher instance
func NewPublisher(collectionApi string, identity *utils.Identity) (*Publisher, error) {
	return &Publisher{
		collectApi: collectionApi,
		identity:   identity,
	}, nil
}

// Publish to Pubsub
func (p *Publisher) SendStoreRequest(request *internal.StoreRequest) error {
	ctx := context.Background()

	request.SignedAt = time.Now()

	// Marshal the request into JSON.
	requestData, err := json.Marshal(request)
	if err != nil {
		return err
	}

	signature, err := p.identity.SignMessage(requestData)
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", p.collectApi, bytes.NewBuffer(requestData))
	if err != nil {
		return err
	}

	// Set the headers
	req.Header.Set("X-LC-Signature", signature)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Check the response status code
	if resp.StatusCode != http.StatusOK {
		responseBody, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("API request failed with status: %s, response: %s", resp.Status, responseBody)
	}

	return nil
}
