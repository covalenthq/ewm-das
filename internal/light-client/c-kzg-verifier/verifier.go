package verifier

import (
	"fmt"

	ckzg4844 "github.com/ethereum/c-kzg-4844/v2/bindings/go"
)

// KZGVerifier represents a structure to verify KZG commitments.
type KZGVerifier struct {
	Commitment ckzg4844.Bytes48
	Proof      []byte
	Cell       []byte
	Index      uint64
}

// NewKZGVerifier creates a new instance of KZGVerifier with the given commitment, proof, and cell.
func NewKZGVerifier(commitment, proof, cell []byte, index uint64) *KZGVerifier {
	temp := ckzg4844.Bytes48{}
	copy(temp[:], commitment)

	return &KZGVerifier{
		Commitment: temp,
		Proof:      proof,
		Cell:       cell,
		Index:      index,
	}
}

// Verify checks the validity of the proof against the commitment and cell.
func (v *KZGVerifier) Verify() (bool, error) {
	var (
		proof ckzg4844.Bytes48
		cell  ckzg4844.Cell
	)

	// verify the length of the commitment, proof, and cell
	if len(v.Commitment) != ckzg4844.BytesPerCommitment ||
		len(v.Proof) != ckzg4844.BytesPerProof ||
		len(v.Cell) != ckzg4844.BytesPerCell {
		return false, fmt.Errorf("invalid length of commitment, proof, or cell")
	}

	copy(proof[:], v.Proof)
	copy(cell[:], v.Cell)

	commitments := [1]ckzg4844.Bytes48{v.Commitment}
	proofs := [1]ckzg4844.Bytes48{proof}
	cells := [1]ckzg4844.Cell{cell}
	indexes := [1]uint64{v.Index}

	return ckzg4844.VerifyCellKZGProofBatch(commitments[:], indexes[:], cells[:], proofs[:])
}

func (v *KZGVerifier) VerifyBatch() (bool, error) {
	var (
		proof   ckzg4844.Bytes48
		cell    ckzg4844.Cell
		indexes [1]uint64
	)

	if len(v.Commitment) != ckzg4844.BytesPerCommitment {
		return false, fmt.Errorf("invalid length of commitment")
	}

	var results []bool
	commitments := [1]ckzg4844.Bytes48{v.Commitment}
	for i := 0; i < 64; i++ {

		copy(proof[:], v.Proof[i*ckzg4844.BytesPerProof:(i+1)*ckzg4844.BytesPerProof])
		copy(cell[:], v.Cell[i*ckzg4844.BytesPerCell:(i+1)*ckzg4844.BytesPerCell])

		indexes[0] = v.Index*64 + uint64(i)
		proofs := [1]ckzg4844.Bytes48{proof}
		cells := [1]ckzg4844.Cell{cell}

		res, err := ckzg4844.VerifyCellKZGProofBatch(commitments[:], indexes[:], cells[:], proofs[:])
		if err != nil {
			return false, err
		}
		results = append(results, res)
	}

	return all(results), nil
}

func all(results []bool) bool {
	for _, res := range results {
		if !res {
			return false
		}
	}
	return true
}
