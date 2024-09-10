package internal

// DataBlock is an interface for data blocks.
type DataBlock interface {
	Describe() (size uint64, rows uint64, cols uint64)
	Commitment(row uint64) ([]byte, error)
	ProofAndCell(row uint64, col uint64) ([]byte, []byte, error)
	Verify() error

	Init(size uint64, rows uint64)
	SetProofAndCell(row uint64, col uint64, proof []byte, cell []byte) error
	SetCellBytes(row uint64, col uint64, cell []byte) error
	RecoverData(byteCells [][][]byte) error
}

// DataBlockEncoder is an interface for encoding and decoding data blocks.
type DataBlockEncoder interface {
	Encode(data []byte) error
	Decode() ([]byte, error)
}
