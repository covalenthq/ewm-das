package ckzgencoder

import ckzg4844 "github.com/ethereum/c-kzg-4844/bindings/go"

// DataBlockImpl is the data block implementation.
type DataBlockImpl struct {
	Blobs       []*ckzg4844.Blob
	Commitments []ckzg4844.KZGCommitment
	Cells       [][ckzg4844.CellsPerExtBlob]ckzg4844.Cell
	Proofs      [][ckzg4844.CellsPerExtBlob]ckzg4844.KZGProof
	Size        uint64
}
