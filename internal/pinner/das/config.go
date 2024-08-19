package das

import (
	ckzgencoder "github.com/covalenthq/das-ipfs-pinner/internal/pinner/c-kzg-encoder"
)

// LoadConfig loads the configuration.
func LoadConfig() interface{} {
	config := ckzgencoder.LoadConfig()
	return config
}
