package ckzgencoder

import (
	"github.com/covalenthq/das-ipfs-pinner/internal"
	ckzg4844 "github.com/ethereum/c-kzg-4844/v2/bindings/go"
)

// EncodeDatablock encodes the data block.
func EncodeDatablock(data []byte) (internal.DataBlock, error) {
	datablock := NewDataBlock()

	err := datablock.Encode(data)
	if err != nil {
		return nil, err
	}

	return datablock, nil
}

// Encode encodes the data block and computes the cells and KZG proofs.
func (d *DataBlockImpl) Encode(data []byte) error {
	if err := d.encodeBlobs(data); err != nil {
		return err
	}
	return d.computeCellsAndKZGProofs()
}

// Decode decodes the data block.
func (d *DataBlockImpl) Decode() ([]byte, error) {
	if d.cells != nil {
		return d.decodeCells()
	}
	return d.decodeBlobs()
}

// encodeBlobs encodes the data into blobs and commitments.
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

	d.size = uint64(len(data))
	d.blobs = blobs
	d.commitments = commitments
	return nil
}

// addBlobAndCommitment adds a blob and its corresponding commitment.
func (d *DataBlockImpl) addBlobAndCommitment(blobs *[]*ckzg4844.Blob, commitments *[]ckzg4844.KZGCommitment, blob *ckzg4844.Blob) error {
	commitment, err := ckzg4844.BlobToKZGCommitment(blob)
	if err != nil {
		return err
	}
	*blobs = append(*blobs, blob)
	*commitments = append(*commitments, commitment)
	return nil
}

// decodeBlobs decodes the blobs back into the original data.
func (d *DataBlockImpl) decodeBlobs() ([]byte, error) {
	data := make([]byte, d.size)
	j := 0

	for _, blob := range d.blobs {
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

// decodeCells decodes the cells back into the original data.
func (d *DataBlockImpl) decodeCells() ([]byte, error) {
	data := make([]byte, d.size)
	j := 0
	count := uint64(0)

	for _, cells := range d.cells {
		for _, cell := range cells[:64] {
			if count >= d.size {
				return data, nil
			}

			for k := 0; k < len(cell); k += 32 {
				remaining := int(d.size) - j
				if remaining < 31 {
					copy(data[j:], cell[k+1:k+1+remaining])
					j += remaining
					break
				}

				copy(data[j:j+31], cell[k+1:k+32])
				j += 31
			}
		}
	}

	return data, nil
}

// computeCellsAndKZGProofs computes the cells and KZG proofs for the encoded blobs.
func (d *DataBlockImpl) computeCellsAndKZGProofs() error {
	if d.blobs == nil {
		return nil
	}

	d.cells = make([][ckzg4844.CellsPerExtBlob]ckzg4844.Cell, len(d.blobs))
	d.proofs = make([][ckzg4844.CellsPerExtBlob]ckzg4844.KZGProof, len(d.blobs))

	for i, blob := range d.blobs {
		cells, proofs, err := ckzg4844.ComputeCellsAndKZGProofs(blob)
		if err != nil {
			return err
		}
		d.cells[i] = cells
		d.proofs[i] = proofs
	}

	return nil
}
