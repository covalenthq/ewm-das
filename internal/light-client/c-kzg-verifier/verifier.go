package verifier

import (
	"fmt"

	ckzg4844 "github.com/ethereum/c-kzg-4844/v2/bindings/go"
)

// KZGVerifier represents a structure to verify KZG commitments.
type KZGVerifier struct {
	commitment  ckzg4844.Bytes48
	proofStack  []byte
	cellStack   []byte
	indexOffset uint64
	stackSize   uint64
}

// NewKZGVerifier creates a new instance of KZGVerifier with the given commitment, proof, and cell.
func NewKZGVerifier(commitment, proofs, cells []byte, index uint64, stackSize uint64) *KZGVerifier {
	temp := ckzg4844.Bytes48{}
	copy(temp[:], commitment)

	return &KZGVerifier{
		commitment:  temp,
		proofStack:  proofs,
		cellStack:   cells,
		indexOffset: index,
		stackSize:   stackSize,
	}
}

// VerifyBatch verifies the KZG commitment in batch.
func (v *KZGVerifier) VerifyBatch() (bool, error) {
	var (
		proof ckzg4844.Bytes48
		cell  ckzg4844.Cell
	)

	if len(v.commitment) != ckzg4844.BytesPerCommitment {
		return false, fmt.Errorf("invalid length of commitment")
	}

	var results []bool
	commitments := [1]ckzg4844.Bytes48{v.commitment}
	for i := 0; i < int(v.stackSize); i++ {

		copy(proof[:], v.proofStack[i*ckzg4844.BytesPerProof:(i+1)*ckzg4844.BytesPerProof])
		copy(cell[:], v.cellStack[i*ckzg4844.BytesPerCell:(i+1)*ckzg4844.BytesPerCell])

		indexes := [1]uint64{v.indexOffset*v.stackSize + uint64(i)}
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

func (v *KZGVerifier) VerifyBatch2() (bool, error) {
	proofs := make([]ckzg4844.Bytes48, v.stackSize)
	cells := make([]ckzg4844.Cell, v.stackSize)
	commitments := make([]ckzg4844.Bytes48, v.stackSize)
	indeces := make([]uint64, v.stackSize)

	for i := 0; i < int(v.stackSize); i++ {
		copy(proofs[i][:], v.proofStack[i*ckzg4844.BytesPerProof:(i+1)*ckzg4844.BytesPerProof])
		copy(cells[i][:], v.cellStack[i*ckzg4844.BytesPerCell:(i+1)*ckzg4844.BytesPerCell])
		commitments[i] = v.commitment
		indeces[i] = v.indexOffset*v.stackSize + uint64(i)
	}

	results, err := ckzg4844.VerifyCellKZGProofBatch(commitments, indeces, cells, proofs)
	if err != nil {
		return false, err
	}

	return results, nil
}

func all(results []bool) bool {
	for _, res := range results {
		if !res {
			return false
		}
	}
	return true
}
