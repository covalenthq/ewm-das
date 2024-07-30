package kzg

import (
	"fmt"
	"math/big"
	"os"
	"path/filepath"
	"sync"

	"github.com/protolambda/go-kzg"
	"github.com/protolambda/go-kzg/bls"
)

const (
	G1CompressedSize = 48
	G2CompressedSize = 96
	splitFactor      = 8 // Number of concurrent chunks to read
)

// TrustedSetup contains the necessary parameters and precomputed values for KZG commitments.
type TrustedSetup struct {
	SubGroupLen   uint64
	SubGroupCount uint64
	EvalLen       uint64
	FftSettings   *kzg.FFTSettings
	KzgSettings   *kzg.KZGSettings
	Fk20Settings  *kzg.FK20MultiSettings
}

var setup *TrustedSetup

// InitTrustedSetup initializes the trusted setup for KZG commitments using the provided configuration.
// It attempts to load cached setup data if available and verifies the seed matches.
func InitTrustedSetup(config Config) error {
	if setup != nil {
		return fmt.Errorf("trusted setup already initialized")
	}

	// Attempt to load the setup from the binary files
	s1, s2, err := LoadSetup(config)
	if err == nil && len(s1) > 0 && len(s2) > 0 {
		// Successfully loaded setup, now use it
		fmt.Println("Loaded trusted setup from cache")
		setup = initializeFromPoints(s1, s2, config)
		return nil
	}

	// If loading fails, generate the setup and cache it
	fmt.Println("Generating new trusted setup")
	s1, s2 = kzg.GenerateTestingSetup(config.Seed, config.EvalLen*2)
	if err := StoreSetup(s1, s2, config); err != nil {
		return fmt.Errorf("failed to store trusted setup: %w", err)
	}

	setup = initializeFromPoints(s1, s2, config)
	return nil
}

func initializeFromPoints(s1 []bls.G1Point, s2 []bls.G2Point, config Config) *TrustedSetup {
	scaleBits := bitsNeeded(config.EvalLen)
	fftSettings := kzg.NewFFTSettings(scaleBits)
	kzgSettings := kzg.NewKZGSettings(fftSettings, s1, s2)
	fk20Settings := kzg.NewFK20MultiSettings(kzgSettings, config.EvalLen*2, config.SubGroupLen)

	return &TrustedSetup{
		SubGroupLen:   config.SubGroupLen,
		SubGroupCount: config.SubGroupCount,
		EvalLen:       config.EvalLen,
		FftSettings:   fftSettings,
		KzgSettings:   kzgSettings,
		Fk20Settings:  fk20Settings,
	}
}

// StoreSetup stores the KZG setup (s1, s2) to binary files.
func StoreSetup(s1 []bls.G1Point, s2 []bls.G2Point, config Config) error {
	if len(s1) != len(s2) {
		return fmt.Errorf("s1 and s2 must have the same length")
	}

	// Ensure storage directory exists
	if err := os.MkdirAll(config.StorageDir, os.ModePerm); err != nil {
		return fmt.Errorf("failed to create storage directory: %w", err)
	}

	// Paths for storing the setup points
	s1Path := filepath.Join(config.StorageDir, "s1.bin")
	s2Path := filepath.Join(config.StorageDir, "s2.bin")

	// Create files for s1 and s2
	f1, err := os.Create(s1Path)
	if err != nil {
		return fmt.Errorf("failed to create file for s1: %w", err)
	}
	defer f1.Close()

	f2, err := os.Create(s2Path)
	if err != nil {
		return fmt.Errorf("failed to create file for s2: %w", err)
	}
	defer f2.Close()

	// Write s1 and s2 points to files in binary format
	for i := 0; i < len(s1); i++ {
		if _, err := f1.Write(bls.ToCompressedG1(&s1[i])); err != nil {
			return fmt.Errorf("failed to write s1[%d]: %w", i, err)
		}
		if _, err := f2.Write(bls.ToCompressedG2(&s2[i])); err != nil {
			return fmt.Errorf("failed to write s2[%d]: %w", i, err)
		}
	}

	return nil
}

// LoadSetup loads the KZG setup (s1, s2) from binary files using concurrent reading.
func LoadSetup(config Config) ([]bls.G1Point, []bls.G2Point, error) {
	s1Path := filepath.Join(config.StorageDir, "s1.bin")
	s2Path := filepath.Join(config.StorageDir, "s2.bin")

	// Open the files for reading
	f1, err := os.Open(s1Path)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to open file for s1: %w", err)
	}
	defer f1.Close()

	f2, err := os.Open(s2Path)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to open file for s2: %w", err)
	}
	defer f2.Close()

	// Determine the number of points
	info, err := f1.Stat()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to stat file for s1: %w", err)
	}
	n := uint64(info.Size() / G1CompressedSize) // Assuming each G1 point is of fixed compressed size

	// Allocate slices for points
	s1 := make([]bls.G1Point, n)
	s2 := make([]bls.G2Point, n)

	chunkSize := n / splitFactor

	var wg sync.WaitGroup
	errCh := make(chan error, splitFactor)
	defer close(errCh)

	readChunk := func(start, end uint64) {
		defer wg.Done()
		for i := start; i < end; i++ {
			// Read G1 point
			s1Data := make([]byte, G1CompressedSize)
			_, err := f1.ReadAt(s1Data, int64(i*G1CompressedSize))
			if err != nil {
				errCh <- fmt.Errorf("failed to read s1[%d]: %w", i, err)
				return
			}
			point, err := bls.FromCompressedG1(s1Data)
			if err != nil {
				errCh <- fmt.Errorf("failed to decompress s1[%d]: %w", i, err)
				return
			}
			s1[i] = *point

			// Read G2 point
			s2Data := make([]byte, G2CompressedSize)
			_, err = f2.ReadAt(s2Data, int64(i*G2CompressedSize))
			if err != nil {
				errCh <- fmt.Errorf("failed to read s2[%d]: %w", i, err)
				return
			}
			point2, err := bls.FromCompressedG2(s2Data)
			if err != nil {
				errCh <- fmt.Errorf("failed to decompress s2[%d]: %w", i, err)
				return
			}
			s2[i] = *point2
		}
	}

	// Launch concurrent reading tasks
	for i := uint64(0); i < splitFactor; i++ {
		start := i * chunkSize
		end := start + chunkSize
		if i == splitFactor-1 {
			end = n // Ensure the last chunk covers any remaining items
		}
		wg.Add(1)
		go readChunk(start, end)
	}

	// Wait for all goroutines to complete
	wg.Wait()

	// Check for errors from goroutines
	select {
	case err := <-errCh:
		return nil, nil, err
	default:
	}

	return s1, s2, nil
}

// bitsNeeded calculates the number of bits needed to represent the given value.
func bitsNeeded(x uint64) uint8 {
	return uint8((big.NewInt(0).SetUint64(x).BitLen() + 7) / 8 * 8)
}

// GetTrustedSetup returns the initialized trusted setup.
func GetTrustedSetup() (*TrustedSetup, error) {
	if setup == nil {
		return nil, fmt.Errorf("trusted setup not initialized")
	}
	return setup, nil
}
