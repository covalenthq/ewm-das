package internal

// DataBlock is an interface for data blocks.
type DataBlock interface {
	Describe() (size uint64, rows uint64, cols uint64)
	Commitment(row uint64) ([]byte, error)
	Proof(row uint64, col uint64) ([]byte, error)
	Cell(row uint64, col uint64) ([]byte, error)
	Verify() error
}

// DataBlockEncoder is an interface for encoding and decoding data blocks.
type DataBlockEncoder interface {
	Encode(data []byte) error
	Decode() ([]byte, error)
}
