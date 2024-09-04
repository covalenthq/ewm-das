package ckzgencoder

import (
	"errors"

	ckzg4844 "github.com/ethereum/c-kzg-4844/v2/bindings/go"
)

var (
	// ErrOutOfRange is returned when the index is out of range.
	ErrOutOfRange = errors.New("out of range")
)

// DataBlockImpl is the data block implementation.
type DataBlockImpl struct {
	blobs       []*ckzg4844.Blob
	commitments []ckzg4844.KZGCommitment
	cells       [][ckzg4844.CellsPerExtBlob]ckzg4844.Cell
	proofs      [][ckzg4844.CellsPerExtBlob]ckzg4844.KZGProof
	size        uint64
}

// Describe returns the size and number of rows in the data block.
func (d *DataBlockImpl) Describe() (size uint64, rows uint64, cols uint64) {
	return d.size, uint64(len(d.cells)), ckzg4844.CellsPerExtBlob
}

// Commitment returns the commitment for the given row.
func (d *DataBlockImpl) Commitment(row uint64) ([]byte, error) {
	if row >= uint64(len(d.commitments)) {
		return nil, ErrOutOfRange
	}
	return d.commitments[row][:], nil
}

// Proof returns the KZG proof for the given row and column.
func (d *DataBlockImpl) Proof(row uint64, col uint64) ([]byte, error) {
	if row >= uint64(len(d.proofs)) || col >= ckzg4844.CellsPerExtBlob {
		return nil, ErrOutOfRange
	}
	return d.proofs[row][col][:], nil
}

// Cell returns the cell for the given row and column.
func (d *DataBlockImpl) Cell(row uint64, col uint64) ([]byte, error) {
	if row >= uint64(len(d.cells)) || col >= ckzg4844.CellsPerExtBlob {
		return nil, ErrOutOfRange
	}
	return d.cells[row][col][:], nil
}

// Verify verifies the data block.
func (d *DataBlockImpl) Verify() error {
	if d.blobs == nil {
		return nil
	}
	if d.cells == nil || d.proofs == nil {
		return ErrCellsOrProofsMissing
	}
	return d.verifyCommitmentsAndProofs()
}

// verifyCommitmentsAndProofs verifies the KZG commitments and proofs.
func (d *DataBlockImpl) verifyCommitmentsAndProofs() error {
	for i, commitment := range d.commitments {
		var commitments [1]ckzg4844.Bytes48
		copy(commitments[0][:], commitment[:])

		for j, cell := range d.cells[i] {
			var subCells [1]ckzg4844.Cell
			copy(subCells[0][:], cell[:])

			var proofs [1]ckzg4844.Bytes48
			copy(proofs[0][:], d.proofs[i][j][:])

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
