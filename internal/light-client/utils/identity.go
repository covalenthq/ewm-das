package utils

import (
	"crypto/ecdsa"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

type Identity struct {
	privKey *ecdsa.PrivateKey
}

func NewIdentity(privKeyHex string) (*Identity, error) {
	privKey, err := StringToPrivateKey(privKeyHex)
	if err != nil {
		return nil, err
	}
	return &Identity{privKey: privKey}, nil
}

func (i *Identity) GetPublicKey() *ecdsa.PublicKey {
	return &i.privKey.PublicKey
}

func (i *Identity) GetAddressHex() common.Address {
	return crypto.PubkeyToAddress(*i.GetPublicKey())
}

func (i *Identity) SignMessage(message []byte) (string, error) {
	signBytes, err := SignMessage(i.privKey, message)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", signBytes), nil
}
