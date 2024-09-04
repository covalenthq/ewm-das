package ckzgencoder

import (
	"fmt"
	"os"
	"path/filepath"

	ckzg4844 "github.com/ethereum/c-kzg-4844/v2/bindings/go"
)

// TrustedSetup is the trusted setup.
type TrustedSetup struct{}

// NewTrustedSetup creates a new trusted setup.
func NewTrustedSetup() *TrustedSetup {
	return &TrustedSetup{}
}

// GenerateTrustedSetup generates a trusted setup.
func (t *TrustedSetup) GenerateTrustedSetup() error {
	return nil
}

// LoadTrustedSetup loads a trusted setup.
func (t *TrustedSetup) LoadTrustedSetup(config Config) error {
	trustedSetupFile := filepath.Join(config.TrustedDir, "trusted_setup.txt")
	// check if the trusted setup file exists
	if _, err := os.Stat(trustedSetupFile); os.IsNotExist(err) {
		return fmt.Errorf("trusted setup file does not exist: %v", trustedSetupFile)
	}
	return ckzg4844.LoadTrustedSetupFile(trustedSetupFile, 0)
}

// FreeTrustedSetup frees a trusted setup.
func (t *TrustedSetup) FreeTrustedSetup() error {
	ckzg4844.FreeTrustedSetup()
	return nil
}
