package ckzgencoder

import (
	"sync"

	ckzg4844 "github.com/ethereum/c-kzg-4844/v2/bindings/go"
	"golang.org/x/sync/errgroup"
)

const MinRequiredCells = ckzg4844.CellsPerExtBlob / 2

// DataBlockImpl is the data block implementation.
type DataBlockImpl struct {
	blobs       []*ckzg4844.Blob
	commitments []ckzg4844.KZGCommitment
	cells       [][ckzg4844.CellsPerExtBlob]ckzg4844.Cell
	proofs      [][ckzg4844.CellsPerExtBlob]ckzg4844.KZGProof
	size        uint64
}

// NewDataBlock creates a new data block.
func NewDataBlock() *DataBlockImpl {
	return &DataBlockImpl{}
}

// Describe returns the size and number of blobs in the data block.
func (d *DataBlockImpl) Describe() (size uint64, nBlobs uint64, nCells uint64) {
	return d.size, uint64(len(d.cells)), ckzg4844.CellsPerExtBlob
}

// Commitment returns the commitment to the given blob index.
func (d *DataBlockImpl) Commitment(nBlob uint64) ([]byte, error) {
	if nBlob >= uint64(len(d.commitments)) {
		return nil, ErrOutOfRange
	}
	return d.commitments[nBlob][:], nil
}

// ProofAndCell returns the KZG proof for the given blob and cell index.
func (d *DataBlockImpl) ProofAndCell(nBlob uint64, nCell uint64) ([]byte, []byte, error) {
	if nBlob >= uint64(len(d.proofs)) || nCell >= ckzg4844.CellsPerExtBlob {
		return nil, nil, ErrOutOfRange
	}
	return d.proofs[nBlob][nCell][:], d.cells[nBlob][nCell][:], nil
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

// Init initializes the data block.
func (d *DataBlockImpl) Init(size uint64, nBlobs uint64) {
	d.blobs = make([]*ckzg4844.Blob, nBlobs)
	d.commitments = make([]ckzg4844.KZGCommitment, nBlobs)
	d.cells = make([][ckzg4844.CellsPerExtBlob]ckzg4844.Cell, nBlobs)
	d.proofs = make([][ckzg4844.CellsPerExtBlob]ckzg4844.KZGProof, nBlobs)
	d.size = size
}

// RecoverData recovers the data from the KZG cells and proofs.
// RecoverData recovers the data from the KZG cells and proofs.
func (d *DataBlockImpl) RecoverData(bCells [][][]byte) error {
	if bCells == nil {
		return ErrBadArgument
	}

	if d.cells == nil || d.proofs == nil {
		return ErrCellsOrProofsMissing
	}

	// Using errgroup for parallel processing and error handling
	var g errgroup.Group
	mu := sync.Mutex{} // Protect access to shared memory (d.cells and d.proofs)

	for i, bBlob := range bCells {
		i, bBlob := i, bBlob // capture loop variables for the goroutine
		g.Go(func() error {
			// Pre-allocate to avoid multiple reallocations
			var validCells []ckzg4844.Cell
			var validIndexes []uint64

			// Iterate through the blob to collect valid cells and their indexes
			for k, byteCell := range bBlob {
				if byteCell == nil {
					continue // Skip nil cells
				}
				var cell ckzg4844.Cell
				copy(cell[:], byteCell)

				validCells = append(validCells, cell)
				validIndexes = append(validIndexes, uint64(k))
			}

			// Ensure we have at least the minimum number of valid cells
			if len(validCells) < MinRequiredCells {
				return ErrNotEnoughCells
			}

			// Recover cells and proofs using the valid indexes and cells
			rCells, rProofs, err := ckzg4844.RecoverCellsAndKZGProofs(validIndexes, validCells)
			if err != nil {
				return err
			}

			// Copy recovered cells and proofs into the DataBlockImpl (protected by mutex)
			mu.Lock()
			defer mu.Unlock()
			copy(d.cells[i][:], rCells[:])
			copy(d.proofs[i][:], rProofs[:])

			return nil
		})
	}

	// Wait for all goroutines to finish
	if err := g.Wait(); err != nil {
		return err
	}

	return nil
}
