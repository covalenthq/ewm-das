package ckzgencoder

import (
	"errors"

	ckzg4844 "github.com/ethereum/c-kzg-4844/bindings/go"
)

var (
	// ErrOutOfRange is returned when the index is out of range.
	ErrOutOfRange = errors.New("out of range")
)

// DataBlockImpl is the data block implementation.
type DataBlockImpl struct {
	Blobs       []*ckzg4844.Blob
	Commitments []ckzg4844.KZGCommitment
	Cells       [][ckzg4844.CellsPerExtBlob]ckzg4844.Cell
	Proofs      [][ckzg4844.CellsPerExtBlob]ckzg4844.KZGProof
	Size        uint64
}

// Describe returns the size and number of rows in the data block.
func (d *DataBlockImpl) Describe() (size uint64, rows uint64, cols uint64) {
	return d.Size, uint64(len(d.Cells)), ckzg4844.CellsPerExtBlob
}

// Commitment returns the commitment for the given row.
func (d *DataBlockImpl) Commitment(row uint64) ([]byte, error) {
	if row >= uint64(len(d.Commitments)) {
		return nil, ErrOutOfRange
	}
	return d.Commitments[row][:], nil
}

// Proof returns the KZG proof for the given row and column.
func (d *DataBlockImpl) Proof(row uint64, col uint64) ([]byte, error) {
	if row >= uint64(len(d.Proofs)) || col >= ckzg4844.CellsPerExtBlob {
		return nil, ErrOutOfRange
	}
	return d.Proofs[row][col][:], nil
}

// Cell returns the cell for the given row and column.
func (d *DataBlockImpl) Cell(row uint64, col uint64) ([]byte, error) {
	if row >= uint64(len(d.Cells)) || col >= ckzg4844.CellsPerExtBlob {
		return nil, ErrOutOfRange
	}
	return d.Cells[row][col][:], nil
}

// Verify verifies the data block.
func (d *DataBlockImpl) Verify() error {
	if d.Blobs == nil {
		return nil
	}
	if d.Cells == nil || d.Proofs == nil {
		return ErrCellsOrProofsMissing
	}
	return d.verifyCommitmentsAndProofs()
}

// verifyCommitmentsAndProofs verifies the KZG commitments and proofs.
func (d *DataBlockImpl) verifyCommitmentsAndProofs() error {
	for i, commitment := range d.Commitments {
		var commitments [1]ckzg4844.Bytes48
		copy(commitments[0][:], commitment[:])

		for j, cell := range d.Cells[i] {
			var subCells [1]ckzg4844.Cell
			copy(subCells[0][:], cell[:])

			var proofs [1]ckzg4844.Bytes48
			copy(proofs[0][:], d.Proofs[i][j][:])

			indexes := [1]uint64{uint64(j)}

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
