package ipldencoder

import (
	"bytes"

	"github.com/ipld/go-ipld-prime/codec/dagcbor"
	"github.com/ipld/go-ipld-prime/datamodel"
	"github.com/ipld/go-ipld-prime/node/basicnode"
)

func DecodeNode(data []byte) (datamodel.Node, error) {
	// Create a NodeAssembler to build the decoded node
	nb := basicnode.Prototype.Any.NewBuilder()

	// Wrap the data in a bytes.Reader to satisfy io.Reader interface
	reader := bytes.NewReader(data)

	// Decode the data using the DAG-CBOR decoder
	if err := dagcbor.Decode(nb, reader); err != nil {
		return nil, err
	}

	// Return the decoded node
	return nb.Build(), nil
}
