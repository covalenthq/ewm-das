package ckzgencoder

import (
	"bytes"
	"os"
	"testing"

	ckzg4844 "github.com/ethereum/c-kzg-4844/bindings/go"
	logging "github.com/ipfs/go-log/v2"
	"golang.org/x/exp/rand"
)

var log = logging.Logger("das-pinner") // Initialize the logger

func fillRandomBytes(size int) []byte {
	data := make([]byte, size)
	rand.Read(data)
	return data
}

// TestMain sets up the testing environment
func TestMain(m *testing.M) {
	err := ckzg4844.LoadTrustedSetupFile("../../../test/data/trusted_setup.txt", 0)
	if err != nil {
		log.Fatalf("Failed to load trusted setup: %v", err)
	}

	// Run the tests
	code := m.Run()

	// Cleanup code if needed
	os.Exit(code)
}

func TestEncodeDecode(t *testing.T) {
	tests := []struct {
		name string
		data []byte
	}{
		{
			name: "Empty data",
			data: []byte{},
		},
		{
			name: "Small data",
			data: []byte("hello world"),
		},
		{
			name: "Exact multiple of 31 bytes",
			data: bytes.Repeat([]byte("a"), 31*3),
		},
		{
			name: "Non-exact multiple of 31 bytes",
			data: bytes.Repeat([]byte("a"), 31*3+10),
		},
		{
			name: "Large data",
			data: fillRandomBytes(1 << 20),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var db DataBlockImpl

			// Test encoding
			err := db.Encode(tt.data)
			if err != nil {
				t.Fatalf("Encode() error = %v", err)
			}

			// Test if Size matches
			if db.size != uint64(len(tt.data)) {
				t.Fatalf("Size mismatch after Encode, got = %v, want = %v", db.size, len(tt.data))
			}

			// Test decoding
			decodedData, err := db.Decode()
			if err != nil {
				t.Fatalf("Decode() error = %v", err)
			}

			// Test if original data matches decoded data
			if !bytes.Equal(tt.data, decodedData) {
				t.Fatalf("Data mismatch after Decode, got = %v, want = %v", decodedData, tt.data)
			}
		})
	}
}
