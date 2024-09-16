package publisher

import (
	"io"
	"net/http"
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/covalenthq/das-ipfs-pinner/common"
	"time"
	"github.com/ethereum/go-ethereum/crypto"
)

type Publisher struct {
	apiUrl   	  string
	privateKeyStr string
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


// NewPublisher creates a new Publisher instance for the API
func NewPublisher(apiUrl, privateKeyStr string) (*Publisher, error) {
	return &Publisher{
		apiUrl:   apiUrl,
		privateKeyStr: privateKeyStr,
	}, nil
}


func getPublicAddressFromPrivateKey(privateKeyHex string) (string, error) {
	privateKey, err := crypto.HexToECDSA(privateKeyHex)
	if err != nil {
		return "", err
	}
	publicKey := privateKey.PublicKey
	publicAddress := crypto.PubkeyToAddress(publicKey).Hex()
	return publicAddress, nil
}


func signMessage(messageData []byte, privateKeyHex string) (string, error) {
	privateKey, err := crypto.HexToECDSA(privateKeyHex)
	if err != nil {
		return "", err
	}
	signature, err := crypto.Sign(crypto.Keccak256(messageData), privateKey)
	if err != nil {
		return "", err
	}
	return "0x" + fmt.Sprintf("%x", signature), nil
}


// Publish to Pubsub
func (p *Publisher) PublishToCS(cid string, rowIndex int, colIndex int, status bool, commitment []byte, proof []byte, cell []byte, blockHeight uint64) error {
	ctx := context.Background()

	publicAddress, err := getPublicAddressFromPrivateKey(p.privateKeyStr)
	if err != nil {
		return err
	}

	message := message{
		ClientId:    publicAddress,
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
	
	// Sign the message
	signature, err := signMessage(messageData, p.privateKeyStr)
	if err != nil {
		return err
	}

	// Create the HTTP request to the API
	req, err := http.NewRequestWithContext(ctx, "POST", p.apiUrl, bytes.NewBuffer(messageData))
	if err != nil {
		return err
	}

	// Set the headers
	req.Header.Set("signature", signature)
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