package kzg

import (
	"bytes"
	"testing"
)

// TestEncodeDecode tests the encoding and decoding process.
func TestEncodeDecode(t *testing.T) {
	// Initialize the test data
	originalData := []byte("test data for encoding and decoding")

	// Define the parameters for the DataBlock
	degree := uint64(4)
	fieldSize := uint64(32)

	// Create a new DataBlock
	db := NewDataBlock(degree, fieldSize)

	// Encode the data
	err := db.Encode(originalData)
	if err != nil {
		t.Fatalf("Encoding failed: %v", err)
	}

	// Decode the data
	decodedData, err := db.Decode()
	if err != nil {
		t.Fatalf("Decoding failed: %v", err)
	}

	// Verify the decoded data matches the original data
	if !bytes.Equal(originalData, decodedData) {
		t.Errorf("Decoded data does not match original. Got: %v, Want: %v", decodedData, originalData)
	}
}

// Additional tests for edge cases and error handling can be added here.
