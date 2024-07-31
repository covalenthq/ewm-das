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
func TestLargeData(t *testing.T) {
	originalData := make([]byte, 1024*1024*4+7) // 4MB of data
	for i := range originalData {
		originalData[i] = byte(i % 256)
	}

	// Define the parameters for the DataBlock
	degree := uint64(4)
	fieldSize := uint64(31)

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

func TestInvalidDegree(t *testing.T) {
	// Define the parameters for the DataBlock with degree 0
	degree := uint64(0)
	fieldSize := uint64(32)

	// Try to create a new DataBlock
	db := NewDataBlock(degree, fieldSize)

	if db != nil {
		t.Fatalf("Expected DataBlock to be nil when degree is invalid")
	}
}

func TestPartialFieldSizeData(t *testing.T) {
	originalData := []byte("incomplete field data")

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

func TestOversizedFieldSize(t *testing.T) {
	originalData := []byte("data that will be oversized when encoded")

	// Define the parameters for the DataBlock
	degree := uint64(4)
	fieldSize := uint64(64) // Oversized field size

	// Create a new DataBlock
	db := NewDataBlock(degree, fieldSize)

	// Encode the data and expect an error
	err := db.Encode(originalData)
	if err == nil {
		t.Fatalf("Expected encoding to fail due to oversized field size, but it succeeded")
	}
}
