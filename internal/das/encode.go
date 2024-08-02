package das

import ckzgencoder "github.com/covalenthq/das-ipfs-pinner/internal/c-kzg-encoder"

func Encode(data []byte) (interface{}, error) {
	datablock := ckzgencoder.DataBlockImpl{}

	err := datablock.Encode(data)
	if err != nil {
		return nil, err
	}

	return nil, nil
}
