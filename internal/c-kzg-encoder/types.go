package ckzgencoder

import ckzg4844 "github.com/ethereum/c-kzg-4844/bindings/go"

type DataBlockImpl struct {
	Size        uint64
	Blobs       []*ckzg4844.Blob
	Commitments []ckzg4844.KZGCommitment
}
