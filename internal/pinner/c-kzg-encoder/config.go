package ckzgencoder

import (
	"os"
	"path/filepath"
)

// Config represents the configuration for the c-kzg encoder.
type Config struct {
	TrustedDir string
}

// DefaultConfig represents the default configuration for the c-kzg encoder.
var DefaultConfig = Config{
	TrustedDir: filepath.Join(os.Getenv("HOME"), ".pinner"),
}

// LoadConfig loads the configuration from environment variables or uses default values.
func LoadConfig() Config {
	config := DefaultConfig

	if envTrustedDir := os.Getenv("PINNER_DIR"); envTrustedDir != "" {
		config.TrustedDir = envTrustedDir
	}

	return config
}
