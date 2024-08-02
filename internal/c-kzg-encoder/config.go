package ckzgencoder

import (
	"os"
	"path/filepath"
)

type Config struct {
	TrustedDir string
}

var DefaultConfig = Config{
	TrustedDir: filepath.Join(os.Getenv("HOME"), ".pinner"),
}

// LoadConfig loads the configuration from environment variables or uses default values.
func LoadConfig() Config {
	config := DefaultConfig

	if envTrustedDir := os.Getenv("TRUSTED_DIR"); envTrustedDir != "" {
		config.TrustedDir = envTrustedDir
	}

	return config
}
