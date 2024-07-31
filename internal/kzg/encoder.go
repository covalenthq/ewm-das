package kzg

import (
	"errors"
	"fmt"
	"sync"

	"github.com/protolambda/go-kzg/bls"
)

// Custom errors for encoding and decoding
var (
	ErrNoPolynomials         = errors.New("no polynomials found for decoding")
	ErrInvalidPolynomial     = errors.New("invalid polynomial data")
	ErrEncodingFailed        = errors.New("failed to encode data into polynomials")
	ErrDataLengthMismatch    = errors.New("data length does not match expected size")
	ErrInsufficientFieldSize = errors.New("insufficient field size for data")
	ErrInvalidDegree         = errors.New("invalid degree for polynomial encoding")
	ErrDataExceedsCapacity   = errors.New("data exceeds the capacity of the polynomial array")
	ErrCopyingData           = errors.New("error while copying data to bytes32 array")
)

// DataBlock represents encoded data with polynomials and related cryptographic proofs.
type DataBlock struct {
	Degree      uint64          // Degree of the polynomials
	FieldSize   uint64          // Size of the field elements
	TotalSize   uint64          // Total size of the encoded data
	Polynomials [][]bls.Fr      // Encoded polynomial cells
	Extended    [][]bls.Fr      // Extended polynomial cells
	Proofs      [][]bls.G1Point // Proofs associated with the polynomials
	Commitments []bls.G1Point   // Commitments for the polynomials
}

// NewDataBlock creates a new instance of DataBlock with the specified degree and field size.
func NewDataBlock(degree, fieldSize uint64) (*DataBlock, error) {
	if degree == 0 {
		return nil, ErrInvalidDegree
	}
	return &DataBlock{
		Degree:    degree,
		FieldSize: fieldSize,
	}, nil
}

// Encode converts raw data into a series of polynomials.
func (db *DataBlock) Encode(data []byte) error {
	if db.Degree == 0 {
		return ErrInvalidDegree
	}

	dataLen := uint64(len(data))
	if dataLen == 0 {
		return ErrDataLengthMismatch
	}

	// Calculate the number of polynomials required to encode the data
	polynomialCount := (dataLen + db.Degree*db.FieldSize - 1) / (db.Degree * db.FieldSize)

	// Initialize the polynomials slice
	db.Polynomials = make([][]bls.Fr, polynomialCount)
	for i := range db.Polynomials {
		db.Polynomials[i] = make([]bls.Fr, db.Degree)
	}

	// Encode the data into polynomials
	for i, offset := uint64(0), uint64(0); offset < dataLen; offset += db.FieldSize {
		availableBytes := dataLen - offset
		copySize := db.FieldSize
		if availableBytes < db.FieldSize {
			copySize = availableBytes
		}

		var bytes32 [32]byte
		if copySize > 32 {
			return ErrInsufficientFieldSize
		}

		copy(bytes32[:copySize], data[offset:offset+copySize])

		if !bls.FrFrom32(&db.Polynomials[offset/(db.Degree*db.FieldSize)][i%db.Degree], bytes32) {
			return fmt.Errorf("%w at index %d", ErrEncodingFailed, offset)
		}

		i++
	}

	db.TotalSize = dataLen
	return nil
}

// Decode converts the stored polynomials back into the original data.
func (db *DataBlock) Decode() ([]byte, error) {
	if len(db.Polynomials) == 0 {
		return nil, ErrNoPolynomials
	}

	polynomialCount := uint64(len(db.Polynomials))
	degree := uint64(len(db.Polynomials[0]))

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

			bytes32 := bls.FrTo32(&db.Polynomials[polyIndex][degIndex])

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

// Commit generates the commitments and proofs for the polynomials.
func (db *DataBlock) Commit(setup *TrustedSetup) error {
	var wg sync.WaitGroup
	resultCh := make(chan struct {
		index      int
		commitment bls.G1Point
		proofs     []bls.G1Point
		err        error
	}, len(db.Polynomials))

	commitments := make([]bls.G1Point, len(db.Polynomials))
	proofs := make([][]bls.G1Point, len(db.Polynomials))

	for i, poly := range db.Polynomials {
		wg.Add(1)
		go func(index int, data []bls.Fr) {
			defer wg.Done()

			coefficients, err := setup.Fk20Settings.KZGSettings.FFTSettings.FFT(data, true)
			if err != nil {
				resultCh <- struct {
					index      int
					commitment bls.G1Point
					proofs     []bls.G1Point
					err        error
				}{index, bls.G1Point{}, nil, err}
				return
			}

			commitment := setup.Fk20Settings.KZGSettings.CommitToPoly(coefficients)
			proof := setup.Fk20Settings.DAUsingFK20Multi(coefficients)

			resultCh <- struct {
				index      int
				commitment bls.G1Point
				proofs     []bls.G1Point
				err        error
			}{index, *commitment, proof, nil}
		}(i, poly)
	}

	wg.Wait()
	close(resultCh)

	for result := range resultCh {
		if result.err != nil {
			return result.err
		}
		commitments[result.index] = result.commitment
		proofs[result.index] = result.proofs
	}

	db.Commitments = commitments
	db.Proofs = proofs

	return nil
}

// Extend extends the polynomials by zero-padding and performing forward and inverse FFT.
func (db *DataBlock) Extend(setup *TrustedSetup) error {
	if db.Polynomials == nil || len(db.Polynomials) == 0 {
		return ErrNoPolynomials
	}

	var wg sync.WaitGroup
	resultCh := make(chan struct {
		index    int
		extended []bls.Fr
		err      error
	}, len(db.Polynomials))

	for i, poly := range db.Polynomials {
		wg.Add(1)
		go func(index int, data []bls.Fr) {
			defer wg.Done()

			// Perform forward FFT on the polynomial data
			coefficients, err := setup.Fk20Settings.KZGSettings.FFTSettings.FFT(data, true)
			if err != nil {
				resultCh <- struct {
					index    int
					extended []bls.Fr
					err      error
				}{index, nil, err}
				return
			}

			// Prepare extended coefficients with zero padding
			extendedCoeffs := make([]bls.Fr, setup.EvalLen*2)
			copy(extendedCoeffs[:setup.EvalLen], coefficients)
			// Zero padding is automatically done as the rest of extendedCoeffs is zero-initialized

			// Perform inverse FFT to obtain the extended polynomial
			extended, err := setup.Fk20Settings.KZGSettings.FFTSettings.FFT(extendedCoeffs, false)
			if err != nil {
				resultCh <- struct {
					index    int
					extended []bls.Fr
					err      error
				}{index, nil, err}
				return
			}
			reverseBitOrderFr(extended)

			resultCh <- struct {
				index    int
				extended []bls.Fr
				err      error
			}{index, extended, nil}
		}(i, poly)
	}

	wg.Wait()
	close(resultCh)

	db.Extended = make([][]bls.Fr, len(db.Polynomials))

	for result := range resultCh {
		if result.err != nil {
			return result.err
		}
		db.Extended[result.index] = result.extended
	}

	return nil
}
