package das

import ckzgencoder "github.com/covalenthq/das-ipfs-pinner/internal/pinner/c-kzg-encoder"

// InitializeTrustedSetup initializes the trusted setup.
func InitializeTrustedSetup(config interface{}) error {
	setup := ckzgencoder.NewTrustedSetup()
	ckzgConfig := config.(ckzgencoder.Config)

	return setup.LoadTrustedSetup(ckzgConfig)
}
