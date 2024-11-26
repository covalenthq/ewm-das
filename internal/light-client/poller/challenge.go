package poller

import (
	"bytes"
	"crypto/sha256"
	"encoding/base32"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/covalenthq/das-ipfs-pinner/internal/light-client/utils"
)

// ClauseType defines the various clause types
type ClauseType struct {
	Type    string
	M       uint64
	K       uint64
	H       uint64
	T       uint64
	Delta   uint64
	Prefix  uint64
	Unknown bool
}

// Challenge represents the decoded challenge
type Challenge struct {
	Version      uint8
	HashFunction uint8
	ClauseType   ClauseType
}

// Decode a Base32 encoded string into a Challenge
func Decode(encoded string) (*Challenge, error) {
	// Validate prefix
	if !strings.HasPrefix(encoded, "ewm") {
		return nil, errors.New("invalid prefix: must start with 'ewm'")
	}

	// Remove the "ewm" prefix
	encodedBody := encoded[3:]

	// Decode Base32
	decoded, err := base32.StdEncoding.WithPadding(base32.NoPadding).DecodeString(strings.ToUpper(encodedBody))
	if err != nil {
		return nil, fmt.Errorf("failed to decode Base32: %w", err)
	}

	reader := bytes.NewReader(decoded)

	// Read version
	var version uint8
	if err := binary.Read(reader, binary.BigEndian, &version); err != nil {
		return nil, fmt.Errorf("failed to read version: %w", err)
	}

	// Read hash_function
	var hashFunction uint8
	if err := binary.Read(reader, binary.BigEndian, &hashFunction); err != nil {
		return nil, fmt.Errorf("failed to read hash_function: %w", err)
	}

	// Read clause_type
	var clauseTypeByte uint8
	if err := binary.Read(reader, binary.BigEndian, &clauseTypeByte); err != nil {
		return nil, fmt.Errorf("failed to read clause_type: %w", err)
	}

	// Decode clause_type and its parameters
	var clauseType ClauseType
	switch clauseTypeByte {
	case 1:
		// Modulo
		var m, k uint64
		if err := binary.Read(reader, binary.BigEndian, &m); err != nil {
			return nil, fmt.Errorf("failed to read modulo m: %w", err)
		}
		if err := binary.Read(reader, binary.BigEndian, &k); err != nil {
			return nil, fmt.Errorf("failed to read modulo k: %w", err)
		}
		clauseType = ClauseType{Type: "Modulo", M: m, K: k}

	case 2:
		// XOR
		var h, t, delta uint64
		if err := binary.Read(reader, binary.BigEndian, &h); err != nil {
			return nil, fmt.Errorf("failed to read XOR h: %w", err)
		}
		if err := binary.Read(reader, binary.BigEndian, &t); err != nil {
			return nil, fmt.Errorf("failed to read XOR t: %w", err)
		}
		if err := binary.Read(reader, binary.BigEndian, &delta); err != nil {
			return nil, fmt.Errorf("failed to read XOR delta: %w", err)
		}
		clauseType = ClauseType{Type: "Xor", H: h, T: t, Delta: delta}

	case 3:
		// Hash Prefix
		var prefix uint64
		if err := binary.Read(reader, binary.BigEndian, &prefix); err != nil {
			return nil, fmt.Errorf("failed to read hash prefix: %w", err)
		}
		clauseType = ClauseType{Type: "HashPrefix", Prefix: prefix}

	default:
		clauseType = ClauseType{Unknown: true}
	}

	return &Challenge{
		Version:      version,
		HashFunction: hashFunction,
		ClauseType:   clauseType,
	}, nil
}

func (c *Challenge) Solve(workload *Workload, identity *utils.Identity) (bool, error) {
	// Calculate the target
	target, err := c.computeTarget(workload, identity)
	if err != nil {
		return false, fmt.Errorf("failed to calculate target: %w", err)
	}

	// Compare the target
	switch c.ClauseType.Type {
	case "Modulo":
		return c.solveModulo(target, c.ClauseType.M, c.ClauseType.K)
	default:
		return false, fmt.Errorf("unsupported clause type: %s", c.ClauseType.Type)
	}
}

func (c *Challenge) computeTarget(workload *Workload, identity *utils.Identity) ([]byte, error) {
	switch c.HashFunction {
	case 1:
		// SHA256
		data, err := json.Marshal(workload)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal workload: %w", err)
		}

		// Compute the hash
		hash := sha256.New()
		hash.Write(data)
		hash.Write(identity.GetAddress().Bytes())

		// Compare the hash
		return hash.Sum(nil), nil
	default:
		return nil, fmt.Errorf("unsupported hash function: %d", c.HashFunction)
	}
}

func (c *Challenge) solveModulo(target []byte, m, k uint64) (bool, error) {
	// Compute the modulo
	var value uint64
	if err := binary.Read(bytes.NewReader(target), binary.BigEndian, &value); err != nil {
		return false, fmt.Errorf("failed to read target: %w", err)
	}

	return value%m < k, nil
}
