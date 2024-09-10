package ckzgencoder

import "errors"

var (
	// ErrOutOfRange is returned when the index is out of range.
	ErrOutOfRange = errors.New("out of range")

	// ErrNotEnoughCells is returned when there are not enough cells to recover the data.
	ErrNotEnoughCells = errors.New("not enough cells")

	// ErrBadArgument is returned when the argument is invalid.
	ErrBadArgument = errors.New("bad argument")

	// ErrCellsOrProofsMissing is returned when cells or proofs are missing.
	ErrCellsOrProofsMissing = errors.New("cells or proofs missing")

	// ErrVerificationFailed is returned when verification fails.
	ErrVerificationFailed = errors.New("verification failed")
)
