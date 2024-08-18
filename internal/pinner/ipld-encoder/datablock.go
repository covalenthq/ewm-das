package ipldencoder

import "github.com/ipld/go-ipld-prime/datamodel"

type IPLDDataBlock struct {
	Version   uint64
	Codec     uint64
	MhType    uint64
	Root      datamodel.Node
	DataNodes [][]datamodel.Node
	Links     [][]datamodel.Link
}
