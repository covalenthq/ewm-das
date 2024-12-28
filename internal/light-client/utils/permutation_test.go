package utils

import (
	"testing"
	"time"
)

// TestPermutationAlgorihm tests the GenerateIndices function.
func TestPermutationAlgorihm(t *testing.T) {
	now := time.Now().UnixNano()
	publicSeed := []byte{byte(now), byte(now >> 8), byte(now >> 16), byte(now >> 24)}
	numSamples := 128
	totalIndices := 128

	rand := NewPolynomialPermutation(publicSeed, totalIndices)
	indices := rand.Generate(numSamples)

	if len(indices) != numSamples {
		t.Fatalf("Expected %d indices, got %d", numSamples, len(indices))
	}

	// Check for uniqueness
	uniqueIndices := make(map[int]struct{})
	for _, idx := range indices {
		if _, ok := uniqueIndices[idx]; ok {
			t.Fatalf("Index %d is not unique", idx)
		}
		uniqueIndices[idx] = struct{}{}
	}
}
