package verifier

import (
	"fmt"

	ckzg4844 "github.com/ethereum/c-kzg-4844/bindings/go"
)

// KZGVerifier represents a structure to verify KZG commitments.
type KZGVerifier struct {
	Commitment []byte
	Proof      []byte
	Cell       []byte
	Index      uint64
}

// NewKZGVerifier creates a new instance of KZGVerifier with the given commitment, proof, and cell.
func NewKZGVerifier(commitment, proof, cell []byte, index uint64) *KZGVerifier {
	return &KZGVerifier{
		Commitment: commitment,
		Proof:      proof,
		Cell:       cell,
		Index:      index,
	}
}

// Verify checks the validity of the proof against the commitment and cell.
func (v *KZGVerifier) Verify() (bool, error) {
	var (
		commitment ckzg4844.Bytes48
		proof      ckzg4844.Bytes48
		cell       ckzg4844.Cell
	)

	// verify the length of the commitment, proof, and cell
	if len(v.Commitment) != ckzg4844.BytesPerCommitment ||
		len(v.Proof) != ckzg4844.BytesPerProof ||
		len(v.Cell) != ckzg4844.BytesPerCell {
		return false, fmt.Errorf("invalid length of commitment, proof, or cell")
	}

	copy(commitment[:], v.Commitment)
	copy(proof[:], v.Proof)
	copy(cell[:], v.Cell)

	commitments := [1]ckzg4844.Bytes48{commitment}
	proofs := [1]ckzg4844.Bytes48{proof}
	cells := [1]ckzg4844.Cell{cell}
	indexes := [1]uint64{v.Index}

	return ckzg4844.VerifyCellKZGProofBatch(commitments[:], indexes[:], cells[:], proofs[:])
}
