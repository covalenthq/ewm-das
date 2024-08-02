package ckzgencoder

import (
	"fmt"

	ckzg4844 "github.com/ethereum/c-kzg-4844/bindings/go"
)

func (d *DataBlockImpl) Encode(data []byte) error {
	err := d.encodeBlobs(data)
	if err != nil {
		return err
	}

	err = d.commitToBlobs()
	if err != nil {
		return err
	}

	return nil
}

func (d *DataBlockImpl) Decode() ([]byte, error) {
	return nil, nil
}

func (d *DataBlockImpl) encodeBlobs(data []byte) error {
	var blobs []ckzg4844.Blob

	// Split data into blobs
	for i := 0; i < len(data); i += ckzg4844.BytesPerBlob {
		// if the remaining data is less than BytesPerBlob, copy the remaining data
		if len(data)-i < ckzg4844.BytesPerBlob {
			blob := ckzg4844.Blob{}
			copy(blob[:], data[i:])
			blobs = append(blobs, blob)
			break
		}

		blob := ckzg4844.Blob{}
		copy(blob[:], data[i:i+ckzg4844.BytesPerBlob])
		blobs = append(blobs, blob)
	}
	d.Blobs = blobs
	d.Size = uint64(len(data))
	return nil
}

func (d *DataBlockImpl) commitToBlobs() error {
	if d.Blobs == nil {
		return fmt.Errorf("no blobs to commit")
	}

	commitments := make([]ckzg4844.KZGCommitment, len(d.Blobs))
	for i, blob := range d.Blobs {
		commitment, err := ckzg4844.BlobToKZGCommitment(&blob)
		if err != nil {
			return err
		}
		commitments[i] = commitment
	}

	d.Commitments = commitments
	return nil
}
