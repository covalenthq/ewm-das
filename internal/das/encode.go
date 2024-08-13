package das

import (
	ckzgencoder "github.com/covalenthq/das-ipfs-pinner/internal/c-kzg-encoder"
	ipldencoder "github.com/covalenthq/das-ipfs-pinner/internal/ipld-encoder"
)

// Encode encodes the data.
func Encode(data []byte) (*ipldencoder.IPLDDataBlock, error) {
	block, err := ckzgencoder.EncodeDatablock(data)
	if err != nil {
		return nil, err
	}

	return ipldencoder.EncodeDatablock(block)
}
