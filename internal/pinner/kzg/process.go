package kzg

// Encode converts raw data into a series of polynomials.
func Encode(data []byte) (*DataBlock, error) {

	// Encode the data into polynomials
	encoded, err := NewDataBlock(setup.EvalLen, 31)
	if err != nil {
		return nil, err
	}

	err = encoded.Encode(data)
	if err != nil {
		return nil, err
	}

	decoded, err := encoded.Decode()
	if err != nil {
		return nil, err
	}

	println("decoded", string(decoded))

	// Commit to the polynomials
	return nil, nil
}
