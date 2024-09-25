package utils

import (
	"crypto/ecdsa"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

type Identity struct {
	privKey *ecdsa.PrivateKey
}

func NewIdentity(privKeyHex string) (*Identity, error) {
	privKey, err := HexToPrivateKey(privKeyHex)
	if err != nil {
		return nil, err
	}
	return &Identity{privKey: privKey}, nil
}

func (i *Identity) GetPublicKey() *ecdsa.PublicKey {
	return &i.privKey.PublicKey
}

func (i *Identity) GetAddress() common.Address {
	return crypto.PubkeyToAddress(*i.GetPublicKey())
}

func (i *Identity) SignMessage(message []byte) ([]byte, error) {
	return SignMessage(i.privKey, message)
}
