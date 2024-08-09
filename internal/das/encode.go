package das

import (
	"fmt"

	ckzgencoder "github.com/covalenthq/das-ipfs-pinner/internal/c-kzg-encoder"
	ipldencoder "github.com/covalenthq/das-ipfs-pinner/internal/ipld-encoder"
)

// Encode encodes the data.
func Encode(data []byte) (interface{}, error) {
	block, err := ckzgencoder.EncodeDatablock(data)
	if err != nil {
		return nil, err
	}

	ipldblock, err := ipldencoder.EncodeDatablock(block)
	if err != nil {
		return nil, err
	}

	fmt.Printf("Encoded data: %v\n", ipldblock)

	return nil, nil
}
