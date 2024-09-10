package internal

// DataBlock is an interface for data blocks.
type DataBlock interface {
	Describe() (size uint64, nBlobs uint64, nCells uint64)
	Commitment(nBlob uint64) ([]byte, error)
	ProofAndCell(nBlob uint64, nCell uint64) ([]byte, []byte, error)
	Verify() error

	Init(size uint64, nBlobs uint64)
	RecoverData(byteCells [][][]byte) error
}

// DataBlockEncoder is an interface for encoding and decoding data blocks.
type DataBlockEncoder interface {
	Encode(data []byte) error
	Decode() ([]byte, error)
}
