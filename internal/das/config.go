package das

import (
	ckzgencoder "github.com/covalenthq/das-ipfs-pinner/internal/c-kzg-encoder"
)

func LoadConfig() interface{} {
	config := ckzgencoder.LoadConfig()
	return config
}
