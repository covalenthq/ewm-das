package utils

import (
	"crypto/ecdsa"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

func StringToPrivateKey(key string) (*ecdsa.PrivateKey, error) {
	return crypto.HexToECDSA(key)
}

func ReadPrivateKey(file string) (*ecdsa.PrivateKey, error) {
	return crypto.LoadECDSA(file)
}

func GenerateKey() (*ecdsa.PrivateKey, error) {
	return crypto.GenerateKey()
}

func FromECDSAPub(pub *ecdsa.PublicKey) []byte {
	return crypto.FromECDSAPub(pub)
}

func FromECDSAPubToAddress(pub *ecdsa.PublicKey) []byte {
	return crypto.PubkeyToAddress(*pub).Bytes()
}

// SignMessage signs the message "hello" using Ethereum's ECC format
func SignMessage(privateKey *ecdsa.PrivateKey, message []byte) ([]byte, error) {
	// Hash the message with Ethereum-specific prefixing
	hash := crypto.Keccak256Hash([]byte(fmt.Sprintf("\x19Ethereum Signed Message:\n%d%s", len(message), message)))

	// Sign the hashed message using crypto.Sign
	return crypto.Sign(hash.Bytes(), privateKey)
}

// Verify the signature using the public key
func VerifySignature(message []byte, signature []byte) (bool, common.Address) {
	// Hash the message with Ethereum-specific prefixing
	hash := crypto.Keccak256Hash([]byte(fmt.Sprintf("\x19Ethereum Signed Message:\n%d%s", len(message), message)))

	// Remove recovery ID (last byte) from signature
	signatureNoRecovery := signature[:len(signature)-1]

	// Verify the signature
	publicKeyRecovered, err := crypto.Ecrecover(hash.Bytes(), signature)
	if err != nil {
		return false, common.Address{}
	}

	addressRecovered := common.BytesToAddress(crypto.Keccak256(publicKeyRecovered[1:])[12:])

	// Compare the recovered public key with the given one
	return crypto.VerifySignature(publicKeyRecovered, hash.Bytes(), signatureNoRecovery), addressRecovered
}
