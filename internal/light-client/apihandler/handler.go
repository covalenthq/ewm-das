package apihandler

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/cenkalti/backoff/v4"
	pb "github.com/covalenthq/das-ipfs-pinner/internal/light-client/schemapb"
	"github.com/covalenthq/das-ipfs-pinner/internal/light-client/utils"
	"golang.org/x/crypto/sha3"
	"google.golang.org/protobuf/proto"

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

	binWorkloadEndpoint, err := url.JoinPath(apiUrl, "/bin-workloads")
	if err != nil {
		return nil, err
	}

	binSamplesEndpoint, err := url.JoinPath(apiUrl, "/bin-samples")
	if err != nil {
		return nil, err
	}

	log.Infof("API URL: %s", apiUrl)

	return &ApiHandler{
		workloadEndpoint: binWorkloadEndpoint,
		samplesEndpoint:  binSamplesEndpoint,
		identity:         identity,
	}, nil
}

func (p *ApiHandler) GetWorkload() (*pb.WorkloadsResponse, error) {
	ctx := context.Background()
	endpoint := p.workloadEndpoint

	timestamp := time.Now().Unix()
	url, err := url.Parse(endpoint)
	if err != nil {
		return nil, err
	}
	message := constructMessage("GET", url.Path, "", timestamp, nil)
	signature, err := p.identity.SignMessage([]byte(message))
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, "GET", endpoint, nil)
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

	var response pb.WorkloadsResponse
	err = proto.Unmarshal([]byte(body), &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

func (p *ApiHandler) SendSampleVerifyRequest(request *pb.SampleVerifyRequest) error {
	return retryWithBackoff(func() error {
		ctx := context.Background()

		request.Timestamp = uint64(time.Now().Unix())
		endpoint := p.samplesEndpoint

		// Marshal the request into JSON.
		requestData, err := proto.Marshal(request)
		if err != nil {
			return err
		}

		url, err := url.Parse(endpoint)
		if err != nil {
			return err
		}

		message := constructMessage("POST", url.Path, "", int64(request.Timestamp), requestData)
		signature, err := p.identity.SignMessage([]byte(message))
		if err != nil {
			return err
		}

		req, err := http.NewRequestWithContext(ctx, "POST", endpoint, bytes.NewBuffer(requestData))
		if err != nil {
			return err
		}

		// Set the headers
		req.Header.Set("X-ETH-ADDRESS", p.identity.GetAddress().Hex())
		req.Header.Set("X-SIGNATURE", fmt.Sprintf("%x", signature))
		req.Header.Set("X-TIMESTAMP", fmt.Sprintf("%d", request.Timestamp))
		req.Header.Set("Content-Type", "application/octet-stream")

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
	}, 3)
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

// retryWithBackoff executes a function with retry logic and exponential backoff.
func retryWithBackoff(operation func() error, maxRetries int) error {
	// Configure exponential backoff
	bo := backoff.NewExponentialBackOff()
	bo.InitialInterval = 500 * time.Millisecond // Start with 500ms delay
	bo.RandomizationFactor = 0.5                // Add randomness to the delay
	bo.Multiplier = 2                           // Double the delay for each retry
	bo.MaxInterval = 5 * time.Second            // Maximum delay between retries
	bo.MaxElapsedTime = 20 * time.Second        // Stop retrying after 20s

	// Limit retries with the maximum retries setting
	retryWithLimit := backoff.WithMaxRetries(bo, uint64(maxRetries))

	// Retry the operation with the backoff strategy
	err := backoff.Retry(func() error {
		err := operation()
		if err != nil {
			log.Warnf("Retrying due to:", err)
			return err
		}
		log.Debug("Operation succeeded!")
		return nil
	}, retryWithLimit)

	if err != nil {
		return fmt.Errorf("operation failed after %d retries: %w", maxRetries, err)
	}
	return nil
}
