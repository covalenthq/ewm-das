package poller

import (
	"bytes"
	"crypto/sha256"
	"encoding/base32"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"strings"

	"github.com/covalenthq/das-ipfs-pinner/internal"
	pb "github.com/covalenthq/das-ipfs-pinner/internal/light-client/schemapb"
	"github.com/covalenthq/das-ipfs-pinner/internal/light-client/utils"
	"google.golang.org/protobuf/proto"
)

// ClauseType defines the various clause types
type ClauseType struct {
	Type    string
	M       *big.Int
	K       *big.Int
	H       *big.Int
	T       *big.Int
	Delta   *big.Int
	Prefix  *big.Int
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
	// Validate the prefix
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

	return DecodeRaw(decoded)
}

func DecodeRaw(data []byte) (*Challenge, error) {
	if len(data) < 3 {
		return nil, errors.New("invalid challenge length")
	}

	// Check if the prefix "ewm" exists and remove it
	if bytes.HasPrefix(data, []byte("ewm")) {
		data = data[3:] // Remove the first 3 bytes
	}

	reader := bytes.NewReader(data)

	// Read version
	var version uint8
	if err := binary.Read(reader, binary.BigEndian, &version); err != nil {
		return nil, fmt.Errorf("failed to read version: %w", err)
	}

	// Read hash function
	var hashFunction uint8
	if err := binary.Read(reader, binary.BigEndian, &hashFunction); err != nil {
		return nil, fmt.Errorf("failed to read hash function: %w", err)
	}

	// Read clause type
	var clauseTypeByte uint8
	if err := binary.Read(reader, binary.BigEndian, &clauseTypeByte); err != nil {
		return nil, fmt.Errorf("failed to read clause type: %w", err)
	}

	// Decode clause type and parameters
	var clauseType ClauseType
	switch clauseTypeByte {
	case 1:
		// Modulo Clause
		m, err := readBigInt(reader)
		if err != nil {
			return nil, fmt.Errorf("failed to read Modulo parameter m: %w", err)
		}
		k, err := readBigInt(reader)
		if err != nil {
			return nil, fmt.Errorf("failed to read Modulo parameter k: %w", err)
		}
		clauseType = ClauseType{Type: "Modulo", M: m, K: k}

	case 2:
		// XOR Clause
		h, err := readBigInt(reader)
		if err != nil {
			return nil, fmt.Errorf("failed to read XOR parameter h: %w", err)
		}
		t, err := readBigInt(reader)
		if err != nil {
			return nil, fmt.Errorf("failed to read XOR parameter t: %w", err)
		}
		delta, err := readBigInt(reader)
		if err != nil {
			return nil, fmt.Errorf("failed to read XOR parameter delta: %w", err)
		}
		clauseType = ClauseType{Type: "Xor", H: h, T: t, Delta: delta}

	case 3:
		// Hash Prefix Clause
		prefix, err := readBigInt(reader)
		if err != nil {
			return nil, fmt.Errorf("failed to read Hash Prefix parameter: %w", err)
		}
		clauseType = ClauseType{Type: "HashPrefix", Prefix: prefix}

	default:
		// Unknown Clause
		clauseType = ClauseType{Unknown: true}
	}

	// Ensure no leftover bytes in the stream
	if reader.Len() > 0 {
		return nil, errors.New("extra unexpected bytes in the encoded challenge")
	}

	return &Challenge{
		Version:      version,
		HashFunction: hashFunction,
		ClauseType:   clauseType,
	}, nil
}

// Helper to read a big.Int from a binary reader
func readBigInt(reader *bytes.Reader) (*big.Int, error) {
	buf := make([]byte, 32) // U256 is 32 bytes
	if _, err := reader.Read(buf); err != nil {
		return nil, err
	}
	return new(big.Int).SetBytes(buf), nil
}

// Solve solves a challenge for a given workload and identity
func (c *Challenge) Solve(workload *internal.Workload, identity *utils.Identity) (bool, error) {
	// Calculate the target
	target, err := c.computeTarget(workload, identity)
	if err != nil {
		return false, fmt.Errorf("failed to calculate target: %w", err)
	}

	// Compare the target
	switch c.ClauseType.Type {
	case "Modulo":
		log.Infof("Solving Modulo challenge M=%s, K=%s", c.ClauseType.M, c.ClauseType.K)
		return c.solveModulo(target, c.ClauseType.M, c.ClauseType.K)
	default:
		return false, fmt.Errorf("unsupported clause type: %s", c.ClauseType.Type)
	}
}

func (c *Challenge) SolveProt(workload *pb.Workload, identity *utils.Identity) (bool, error) {
	// Calculate the target
	target, err := c.computeTargetProt(workload, identity)
	if err != nil {
		return false, fmt.Errorf("failed to calculate target: %w", err)
	}

	// Compare the target
	switch c.ClauseType.Type {
	case "Modulo":
		log.Infof("Solving Modulo challenge M=%s, K=%s", c.ClauseType.M, c.ClauseType.K)
		return c.solveModulo(target, c.ClauseType.M, c.ClauseType.K)
	default:
		return false, fmt.Errorf("unsupported clause type: %s", c.ClauseType.Type)
	}
}

func (c *Challenge) computeTargetProt(workload *pb.Workload, identity *utils.Identity) ([]byte, error) {
	switch c.HashFunction {
	case 1:
		// SHA256
		data, err := proto.Marshal(workload)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal workload: %w", err)
		}

		// Compute the hash
		hash := sha256.New()
		hash.Write(data)
		hash.Write(identity.GetAddress().Bytes())

		return hash.Sum(nil), nil
	default:
		return nil, fmt.Errorf("unsupported hash function: %d", c.HashFunction)
	}
}

func (c *Challenge) computeTarget(workload *internal.Workload, identity *utils.Identity) ([]byte, error) {
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

		return hash.Sum(nil), nil
	default:
		return nil, fmt.Errorf("unsupported hash function: %d", c.HashFunction)
	}
}

func (c *Challenge) solveModulo(target []byte, m, k *big.Int) (bool, error) {
	// Compute the modulo
	targetInt := new(big.Int).SetBytes(target)
	remainder := new(big.Int).Mod(targetInt, m)

	return remainder.Cmp(k) < 0, nil
}
