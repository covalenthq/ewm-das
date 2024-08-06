package ckzgencoder

import (
	"errors"

	ckzg4844 "github.com/ethereum/c-kzg-4844/bindings/go"
)

var (
	ErrCellsOrProofsMissing = errors.New("cells or proofs missing")
	ErrVerificationFailed   = errors.New("verification failed")
)

// Encode encodes the data block.
func (d *DataBlockImpl) Encode(data []byte) error {
	// Step 1: Encode the data into blobs and commitments
	if err := d.encodeBlobs(data); err != nil {
		return err
	}
	// Step 2: Compute cells and KZG proofs
	if err := d.computeCellsAndKZGProofs(); err != nil {
		return err
	}
	return nil
}

// Decode decodes the data block.
func (d *DataBlockImpl) Decode() ([]byte, error) {
	return d.decodeBlobs()
}

func (d *DataBlockImpl) Verify() error {
	if d.Blobs == nil {
		return nil
	}

	if d.Cells == nil || d.Proofs == nil {
		return ErrCellsOrProofsMissing
	}

	for i, commitment := range d.Commitments {
		// Prepare the commitment
		var commitments [1]ckzg4844.Bytes48
		copy(commitments[0][:], commitment[:])

		for j, cell := range d.Cells[i] {
			// Prepare the cell, proof, and index for verification
			var subCells [1]ckzg4844.Cell
			copy(subCells[0][:], cell[:])

			var proofs [1]ckzg4844.Bytes48
			copy(proofs[0][:], d.Proofs[i][j][:])

			indexes := [1]uint64{uint64(j)}

			// Verify the cell KZG proof
			ok, err := ckzg4844.VerifyCellKZGProofBatch(commitments[:], indexes[:], subCells[:], proofs[:])
			if err != nil {
				return err
			}
			if !ok {
				return ErrVerificationFailed
			}
		}
	}

	return nil
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

func (d *DataBlockImpl) computeCellsAndKZGProofs() error {
	if d.Blobs == nil {
		return nil
	}

	d.Cells = make([][ckzg4844.CellsPerExtBlob]ckzg4844.Cell, len(d.Blobs))
	d.Proofs = make([][ckzg4844.CellsPerExtBlob]ckzg4844.KZGProof, len(d.Blobs))

	for i, blob := range d.Blobs {
		cells, proofs, err := ckzg4844.ComputeCellsAndKZGProofs(blob)
		if err != nil {
			return err
		}
		d.Cells[i] = cells
		d.Proofs[i] = proofs
	}

	return nil
}
