package kzg

import (
	"log"
	"os"
	"path/filepath"
	"strconv"
)

// Config holds the configuration for the KZG trusted setup.
type Config struct {
	SubGroupLen   uint64
	SubGroupCount uint64
	EvalLen       uint64
	Seed          string
	StorageDir    string
}

// DefaultConfig provides default values for the KZG setup.
var DefaultConfig = Config{
	SubGroupLen:   1 << 6,
	SubGroupCount: 1 << 8,
	Seed:          "1927409816240961209460912649124",
	StorageDir:    filepath.Join(os.Getenv("HOME"), ".pinner"), // Default storage directory
}

// LoadConfig loads the configuration from environment variables or uses default values.
func LoadConfig() Config {
	config := DefaultConfig

	if envSubGroupLen := os.Getenv("SUB_GROUP_LEN"); envSubGroupLen != "" {
		if val, err := strconv.ParseUint(envSubGroupLen, 10, 64); err == nil {
			config.SubGroupLen = val
		} else {
			log.Printf("Invalid SUB_GROUP_LEN value, using default: %v\n", DefaultConfig.SubGroupLen)
		}
	}

	if envSubGroupCount := os.Getenv("SUB_GROUP_COUNT"); envSubGroupCount != "" {
		if val, err := strconv.ParseUint(envSubGroupCount, 10, 64); err == nil {
			config.SubGroupCount = val
		} else {
			log.Printf("Invalid SUB_GROUP_COUNT value, using default: %v\n", DefaultConfig.SubGroupCount)
		}
	}

	if envSeed := os.Getenv("SEED"); envSeed != "" {
		config.Seed = envSeed
	}

	if envStorageDir := os.Getenv("STORAGE_DIR"); envStorageDir != "" {
		config.StorageDir = envStorageDir
	}

	config.EvalLen = config.SubGroupLen * config.SubGroupCount

	return config
}
