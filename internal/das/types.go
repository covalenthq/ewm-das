package das

type DataBlock interface {
	Encode(data []byte) error
	Decode() ([]byte, error)
}
