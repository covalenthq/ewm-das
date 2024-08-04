package ckzgencoder

import (
	"path/filepath"

	ckzg4844 "github.com/ethereum/c-kzg-4844/bindings/go"
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
	return ckzg4844.LoadTrustedSetupFile(trustedSetupFile)
}

// FreeTrustedSetup frees a trusted setup.
func (t *TrustedSetup) FreeTrustedSetup() error {
	return nil
}
