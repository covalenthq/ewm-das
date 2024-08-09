package ckzgencoder

import (
	"errors"

	ckzg4844 "github.com/ethereum/c-kzg-4844/bindings/go"
)

var (
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
