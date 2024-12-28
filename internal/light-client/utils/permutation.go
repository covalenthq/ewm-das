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

// Pseudo-random number generator based on a polynomial permutation
// `f(x) = (a*x + b) mod maxElements“
// where `a` and `maxElements` are coprime, guarantees unique indices
// and `b` is a random integer - ofset
type PolynomialPermutation struct {
	a           int
	b           int
	maxElements int
}

// NewPolynomialPermutation creates a new PolynomialPermutation
func NewPolynomialPermutation(publicSeed []byte, maxElements int) *PolynomialPermutation {
	// Hash the public seed to derive coefficients
	hash := sha256.Sum256(publicSeed)
	seed := new(big.Int).SetBytes(hash[:]).Int64()

	// Derive coefficients a and b
	a := int(seed%int64(maxElements)) + 1
	b := int((seed / int64(maxElements)) % int64(maxElements))

	// Ensure a is coprime with maxElements
	for gcd(a, maxElements) != 1 {
		a = (a + 1) % maxElements
	}

	return &PolynomialPermutation{
		a:           a,
		b:           b,
		maxElements: maxElements,
	}
}

// Generate generates a unique set of integers between 0 and maxElements
func (p *PolynomialPermutation) Generate(numSamples int) []int {
	indices := make([]int, numSamples)
	for i := 0; i < numSamples; i++ {
		indices[i] = (p.a*i + p.b) % p.maxElements
	}

	return indices
}
