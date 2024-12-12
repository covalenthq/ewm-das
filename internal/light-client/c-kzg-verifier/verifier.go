package verifier

import (
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

func (v *KZGVerifier) VerifyBatch() (bool, error) {
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
