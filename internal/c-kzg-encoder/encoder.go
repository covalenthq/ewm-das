package ckzgencoder

import (
	ckzg4844 "github.com/ethereum/c-kzg-4844/bindings/go"
)

// Encode encodes the data block.
func (d *DataBlockImpl) Encode(data []byte) error {
	return d.encodeBlobs(data)
}

// Decode decodes the data block.
func (d *DataBlockImpl) Decode() ([]byte, error) {
	return d.decodeBlobs()
}

func (d *DataBlockImpl) encodeBlobs(data []byte) error {
	var (
		blobs       []*ckzg4844.Blob
		commitments []ckzg4844.KZGCommitment
		j           int
		blob        = new(ckzg4844.Blob)
	)

	for i := 0; i < len(data); i += 31 {
		if j == ckzg4844.BytesPerBlob {
			if err := d.addBlobAndCommitment(&blobs, &commitments, blob); err != nil {
				return err
			}
			blob = new(ckzg4844.Blob)
			j = 0
		}

		copy(blob[j+1:j+min(32, len(data)+1)], data[i:min(i+31, len(data))])
		j += 32
	}

	if j > 0 {
		if err := d.addBlobAndCommitment(&blobs, &commitments, blob); err != nil {
			return err
		}
	}

	d.Size = uint64(len(data))
	d.Blobs = blobs
	d.Commitments = commitments
	return nil
}

func (d *DataBlockImpl) addBlobAndCommitment(blobs *[]*ckzg4844.Blob, commitments *[]ckzg4844.KZGCommitment, blob *ckzg4844.Blob) error {
	commitment, err := ckzg4844.BlobToKZGCommitment(blob)
	if err != nil {
		return err
	}
	*blobs = append(*blobs, blob)
	*commitments = append(*commitments, commitment)
	return nil
}

func (d *DataBlockImpl) decodeBlobs() ([]byte, error) {
	data := make([]byte, d.Size)
	j := 0

	for _, blob := range d.Blobs {
		for k := 0; k < len(blob); k += 32 {
			remaining := len(data) - j
			if remaining < 31 {
				copy(data[j:], blob[k+1:k+1+remaining])
				j += remaining
				break
			}

			copy(data[j:j+31], blob[k+1:k+32])
			j += 31
		}
	}

	return data, nil
}
