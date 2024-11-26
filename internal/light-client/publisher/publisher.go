package publisher

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/covalenthq/das-ipfs-pinner/internal"
	"github.com/covalenthq/das-ipfs-pinner/internal/light-client/utils"
)

type Publisher struct {
	collectEndpoint string
	identity        *utils.Identity
}

// NewPublisher creates a new Publisher instance
func NewPublisher(apiUrl string, identity *utils.Identity) (*Publisher, error) {
	_, err := url.Parse(apiUrl)
	if err != nil {
		return nil, err
	}

	endpoint, err := url.JoinPath(apiUrl, "/samples")
	if err != nil {
		return nil, err
	}

	return &Publisher{
		collectEndpoint: endpoint,
		identity:        identity,
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

	req, err := http.NewRequestWithContext(ctx, "POST", p.collectEndpoint, bytes.NewBuffer(requestData))
	if err != nil {
		return err
	}

	// Set the headers
	req.Header.Set("X-LC-Signature", fmt.Sprintf("%x", signature))
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

func (p *Publisher) SendStoreRequest2(request *internal.StoreRequest2) error {
	ctx := context.Background()

	request.Timestamp = time.Now()

	// Marshal the request into JSON.
	requestData, err := json.Marshal(request)
	if err != nil {
		return err
	}

	signature, err := p.identity.SignMessage(requestData)
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", p.collectEndpoint, bytes.NewBuffer(requestData))
	if err != nil {
		return err
	}

	// Set the headers
	req.Header.Set("X-LC-Signature", fmt.Sprintf("%x", signature))
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
