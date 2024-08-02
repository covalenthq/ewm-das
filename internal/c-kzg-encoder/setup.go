package ckzgencoder

import (
	"path/filepath"

	ckzg4844 "github.com/ethereum/c-kzg-4844/bindings/go"
)

type trustedSetup struct{}

func NewTrustedSetup() *trustedSetup {
	return &trustedSetup{}
}

func (t *trustedSetup) GenerateTrustedSetup() error {
	return nil
}

func (t *trustedSetup) LoadTrustedSetup(config Config) error {
	trustedSetupFile := filepath.Join(config.TrustedDir, "trusted_setup.txt")
	return ckzg4844.LoadTrustedSetupFile(trustedSetupFile)
}

func (t *trustedSetup) FreeTrustedSetup() error {
	return nil
}
