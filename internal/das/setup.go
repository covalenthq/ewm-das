package das

import ckzgencoder "github.com/covalenthq/das-ipfs-pinner/internal/c-kzg-encoder"

func InitializeTrustedSetup(config interface{}) error {
	setup := ckzgencoder.NewTrustedSetup()
	ckzgConfig := config.(ckzgencoder.Config)

	return setup.LoadTrustedSetup(ckzgConfig)
}
