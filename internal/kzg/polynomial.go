package kzg

import (
	"fmt"

	"github.com/protolambda/go-kzg/bls"
)

// Custom errors
var (
	ErrNoPolynomials     = fmt.Errorf("no polynomials found for decoding")
	ErrInvalidPolynomial = fmt.Errorf("invalid polynomial data")
	ErrEncodingFailed    = fmt.Errorf("failed to encode data into polynomials")
)

// DataBlock represents encoded data with polynomials and related cryptographic proofs.
type DataBlock struct {
	Degree      uint64           // Degree of the polynomials
	FieldSize   uint64           // Size of the field elements
	TotalSize   uint64           // Total size of the encoded data
	Polynomials *[][]bls.Fr      // Encoded polynomial cells
	Proofs      *[][]bls.G1Point // Proofs associated with the polynomials
	Commitments *[]bls.G1Point   // Commitments for the polynomials
}

// NewDataBlock creates a new instance of DataBlock with specified degree and field size.
func NewDataBlock(degree, fieldSize uint64) *DataBlock {
	return &DataBlock{
		Degree:    degree,
		FieldSize: fieldSize,
		TotalSize: 0,
	}
}

// Encode converts raw data into a series of polynomials.
func (db *DataBlock) Encode(data []byte) error {
	dataLen := uint64(len(data))

	// Calculate the number of polynomials required to encode the data
	polynomialCount := dataLen / (db.Degree * db.FieldSize)
	if remainingBytes := dataLen % (db.Degree * db.FieldSize); remainingBytes > 0 {
		polynomialCount++
	}

	// Initialize the polynomials slice
	polynomials := make([][]bls.Fr, polynomialCount)
	for i := range polynomials {
		polynomials[i] = make([]bls.Fr, db.Degree)
	}

	fmt.Println("Data Length:", dataLen, "Polynomial Count:", polynomialCount)

	// Encode the data into polynomials
	for i, offset := uint64(0), uint64(0); offset < dataLen; offset += db.FieldSize {
		var bytes32 [32]byte
		copy(bytes32[:db.FieldSize], data[offset:offset+db.FieldSize])

		if !bls.FrFrom32(&polynomials[offset/(db.Degree*db.FieldSize)][i%db.Degree], bytes32) {
			return fmt.Errorf("%w at index %d", ErrEncodingFailed, offset)
		}

		i++

		if offset+db.FieldSize >= dataLen {
			remainingBytes := dataLen - offset
			if remainingBytes > 0 {
				copy(bytes32[:remainingBytes], data[offset:offset+remainingBytes])
				if !bls.FrFrom32(&polynomials[offset/(db.Degree*db.FieldSize)][i%db.Degree], bytes32) {
					return fmt.Errorf("%w at index %d", ErrEncodingFailed, offset)
				}
			}
			break
		}
	}

	db.Polynomials = &polynomials
	db.TotalSize = dataLen

	return nil
}

// Decode converts the stored polynomials back into the original data.
func (db *DataBlock) Decode() ([]byte, error) {
	if db.Polynomials == nil || len(*db.Polynomials) == 0 {
		return nil, ErrNoPolynomials
	}

	polynomialCount := uint64(len(*db.Polynomials))
	degree := uint64(len((*db.Polynomials)[0]))

	if degree*polynomialCount*db.FieldSize < db.TotalSize {
		return nil, ErrInvalidPolynomial
	}

	data := make([]byte, db.TotalSize)
	dataIndex := uint64(0)

	// Decode the polynomials back into the original data
	for polyIndex := uint64(0); polyIndex < polynomialCount; polyIndex++ {
		for degIndex := uint64(0); degIndex < degree; degIndex++ {
			if dataIndex >= db.TotalSize {
				break
			}

			bytes32 := bls.FrTo32(&(*db.Polynomials)[polyIndex][degIndex])

			// Calculate the number of bytes to copy, considering the potential end of the data
			bytesToCopy := db.FieldSize
			if dataIndex+db.FieldSize > db.TotalSize {
				bytesToCopy = db.TotalSize - dataIndex
			}

			copy(data[dataIndex:dataIndex+bytesToCopy], bytes32[:bytesToCopy])
			dataIndex += bytesToCopy
		}
	}

	return data, nil
}
