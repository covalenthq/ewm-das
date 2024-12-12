package apihandler

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
	"golang.org/x/crypto/sha3"

	logging "github.com/ipfs/go-log/v2"
)

var log = logging.Logger("api-handler")

type ApiHandler struct {
	workloadEndpoint string
	samplesEndpoint  string
	identity         *utils.Identity
}

// NewApiHandler creates a new API handler instance
func NewApiHandler(apiUrl string, identity *utils.Identity) (*ApiHandler, error) {
	_, err := url.Parse(apiUrl)
	if err != nil {
		return nil, err
	}

	workloadEndpoint, err := url.JoinPath(apiUrl, "/workloads")
	if err != nil {
		return nil, err
	}

	samplesEndpoint, err := url.JoinPath(apiUrl, "/samples")
	if err != nil {
		return nil, err
	}

	log.Infof("API URL: %s", apiUrl)

	return &ApiHandler{
		workloadEndpoint: workloadEndpoint,
		samplesEndpoint:  samplesEndpoint,
		identity:         identity,
	}, nil
}

func (p *ApiHandler) GetWorkload() (*internal.WorkloadResponse, error) {
	ctx := context.Background()

	timestamp := time.Now().Unix()
	url, err := url.Parse(p.workloadEndpoint)
	if err != nil {
		return nil, err
	}
	message := constructMessage("GET", url.Path, "", timestamp, nil)
	signature, err := p.identity.SignMessage([]byte(message))
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, "GET", p.workloadEndpoint, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("X-ETH-ADDRESS", p.identity.GetAddress().Hex())
	req.Header.Set("X-SIGNATURE", fmt.Sprintf("%x", signature))
	req.Header.Set("X-TIMESTAMP", fmt.Sprintf("%d", timestamp))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Check the response status code
	if resp.StatusCode != http.StatusOK {
		responseBody, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API request failed with status: %s, response: %s", resp.Status, responseBody)
	}

	// Read the response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var response internal.WorkloadResponse
	err = json.Unmarshal([]byte(body), &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

func (p *ApiHandler) SendStoreRequest(request *internal.StoreRequest) error {
	ctx := context.Background()

	request.Timestamp = time.Now()

	// Marshal the request into JSON.
	requestData, err := json.Marshal(request)
	if err != nil {
		return err
	}

	timestamp := request.Timestamp.Unix()
	url, err := url.Parse(p.samplesEndpoint)
	if err != nil {
		return err
	}

	message := constructMessage("POST", url.Path, "", timestamp, requestData)
	signature, err := p.identity.SignMessage([]byte(message))
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", p.samplesEndpoint, bytes.NewBuffer(requestData))
	if err != nil {
		return err
	}

	// Set the headers
	req.Header.Set("X-ETH-ADDRESS", p.identity.GetAddress().Hex())
	req.Header.Set("X-SIGNATURE", fmt.Sprintf("%x", signature))
	req.Header.Set("X-TIMESTAMP", fmt.Sprintf("%d", timestamp))
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

// constructMessage builds the canonical message
func constructMessage(method, path, queryString string, timestamp int64, body []byte) string {
	if method == "POST" {
		// Hash the body for POST requests
		bodyHash := hashKeccak256(body)
		return fmt.Sprintf("POST\n%s\n%s\ntimestamp: %d", path, bodyHash, timestamp)
	}
	// Use query string for GET requests
	return fmt.Sprintf("GET\n%s\n%s\ntimestamp: %d", path, queryString, timestamp)
}

// hashKeccak256 hashes the input string using Keccak256
func hashKeccak256(data []byte) string {
	hasher := sha3.NewLegacyKeccak256()
	hasher.Write(data)
	return fmt.Sprintf("0x%x", hasher.Sum(nil))
}
