package utils

import (
	"crypto/sha256"
	"math/big"
)

func gcd(a, b int) int {
	if b == 0 {
		return a
	}
	return gcd(b, a%b)
}

// Generate unique indices using polynomial permutation
// f(x) = (a*x + b) mod totalIndices
// where a and totalIndices are coprime, guarantees unique indices
// and b is a random integer - ofset
func GenerateIndices(publicSeed []byte, numSamples, totalIndices int) []int {
	// Hash the public seed to derive coefficients
	hash := sha256.Sum256(publicSeed)
	seed := new(big.Int).SetBytes(hash[:]).Int64()

	// Derive coefficients a and b
	a := int(seed%int64(totalIndices)) + 1
	b := int((seed / int64(totalIndices)) % int64(totalIndices))

	// Ensure a is coprime with totalIndices
	for gcd(a, totalIndices) != 1 {
		a = (a + 1) % totalIndices
	}

	// Generate indices
	indices := make([]int, numSamples)
	for i := 0; i < numSamples; i++ {
		indices[i] = (a*i + b) % totalIndices
	}

	return indices
}
