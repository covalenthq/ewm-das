package ipldencoder

import "github.com/ipld/go-ipld-prime/datamodel"

type IPLDDataBlock struct {
	Root      datamodel.Node
	DataNodes [][]datamodel.Node
	Links     [][]datamodel.Link
}
