package das

// DataBlock is an interface for encoding and decoding data.
type DataBlock interface {
	Encode(data []byte) error
	Decode() ([]byte, error)
}
