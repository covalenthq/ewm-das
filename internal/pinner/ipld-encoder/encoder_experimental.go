//go:build experimental

package ipldencoder

import (
	"github.com/covalenthq/das-ipfs-pinner/internal"
	"github.com/ipfs/go-cid"
	"github.com/ipld/go-ipld-prime"

	// Import the DAG-CBOR codec
	_ "github.com/ipld/go-ipld-prime/codec/dagcbor"
	"github.com/ipld/go-ipld-prime/datamodel"
	"github.com/ipld/go-ipld-prime/fluent/qp"
	"github.com/ipld/go-ipld-prime/linking"
	cidlink "github.com/ipld/go-ipld-prime/linking/cid"
	"github.com/ipld/go-ipld-prime/node/basicnode"
	"github.com/ipld/go-ipld-prime/storage/memstore"
	mh "github.com/multiformats/go-multihash"
)

// EncodeDatablock encodes the given DataBlock into an IPLDDataBlock.
func EncodeDatablock(block internal.DataBlock) (*IPLDDataBlock, error) {
	datablock := &IPLDDataBlock{
		Version: 1,
		Codec:   cid.DagCBOR,
		MhType:  mh.SHA2_256,
	}
	err := datablock.Encode(block)
	if err != nil {
		return nil, err
	}

	return datablock, nil
}

// Encode encodes the data from the given DataBlock.
func (b *IPLDDataBlock) Encode(block internal.DataBlock) error {
	// Encode data nodes
	if err := b.encodeDataNodes(block); err != nil {
		return err
	}

	// Create the LinkSystem
	lsys := createLinkSystem()

	// Encode links
	if err := b.encodeLinks(lsys, block); err != nil {
		return err
	}

	// Encode the root node
	if err := b.encodeRoot(lsys, block); err != nil {
		return err
	}

	// Return the root node
	return nil
}

func (b *IPLDDataBlock) encodeDataNodes(block internal.DataBlock) error {
	_, nBlobs, nCells := block.Describe()
	cellsInNode := nCells / internal.StackSize

	b.DataNodes = make([][]datamodel.Node, nBlobs)
	for nBlob := uint64(0); nBlob < nBlobs; nBlob++ {
		b.DataNodes[nBlob] = make([]datamodel.Node, cellsInNode)

		for nCellStack := uint64(0); nCellStack < cellsInNode; nCellStack++ {
			var stackedProof []byte
			var stackedCell []byte

			for i := uint64(0); i < internal.StackSize; i++ {
				proof, cell, err := block.ProofAndCell(nBlob, nCellStack*internal.StackSize+i)
				if err != nil {
					return err
				}

				stackedProof = append(stackedProof, proof...)
				stackedCell = append(stackedCell, cell...)
			}

			node, err := createCellNode(stackedProof, stackedCell)
			if err != nil {
				return err
			}

			b.DataNodes[nBlob][nCellStack] = node
		}
	}

	return nil
}

func (b *IPLDDataBlock) encodeLinks(lsys *linking.LinkSystem, block internal.DataBlock) error {
	_, nBlobs, _ := block.Describe()

	b.Links = make([][]datamodel.Link, nBlobs)
	for i := uint64(0); i < nBlobs; i++ {
		nCells := len(b.DataNodes[i])

		b.Links[i] = make([]datamodel.Link, nCells)
		for j := 0; j < nCells; j++ {
			link, err := b.createLink(lsys, b.DataNodes[i][j])
			if err != nil {
				return err
			}

			b.Links[i][j] = link
		}
	}

	return nil
}

func (b *IPLDDataBlock) encodeRoot(lsys *linking.LinkSystem, block internal.DataBlock) error {
	size, nBlob, nCell := block.Describe()

	// Create an array of links for the root node
	listLinks := make([]datamodel.Link, nBlob)
	for i, subLinks := range b.Links {
		node, err := qp.BuildList(basicnode.Prototype.List, int64(len(subLinks)), func(la ipld.ListAssembler) {
			for _, link := range subLinks {
				qp.ListEntry(la, qp.Link(link))
			}
		})
		if err != nil {
			return err
		}
		link, err := b.createLink(lsys, node)
		if err != nil {
			return err
		}
		listLinks[i] = link
	}

	// Create the root DAG-CBOR object
	rootNode, err := qp.BuildMap(basicnode.Prototype.Map, -1, func(ma datamodel.MapAssembler) {
		qp.MapEntry(ma, "version", qp.String("v0.2.0"))
		qp.MapEntry(ma, "length", qp.Int(int64(nCell/internal.StackSize)))
		qp.MapEntry(ma, "size", qp.Int(int64(size)))
		qp.MapEntry(ma, "commitments", qp.List(int64(nBlob), func(la ipld.ListAssembler) {
			for i := uint64(0); i < nBlob; i++ {
				commitment, err := block.Commitment(i)
				if err != nil {
					return
				}
				qp.ListEntry(la, qp.Bytes(commitment[:]))
			}
		}))
		qp.MapEntry(ma, "links", qp.List(int64(len(listLinks)), func(la ipld.ListAssembler) {
			for _, link := range listLinks {
				qp.ListEntry(la, qp.Link(link))
			}
		}))
	})
	if err != nil {
		return err
	}

	b.Root = rootNode
	return nil
}

// Utility function to create a link from a node using a given LinkSystem
func (b *IPLDDataBlock) createLink(ls *linking.LinkSystem, node datamodel.Node) (datamodel.Link, error) {
	lp := cidlink.LinkPrototype{Prefix: cid.Prefix{
		Version:  b.Version,
		Codec:    b.Codec,
		MhType:   b.MhType,
		MhLength: 32,
	}}

	return ls.Store(linking.LinkContext{}, lp, node)
}

// Utility function to create a cell node from proof and cell data
func createCellNode(proof, cell []byte) (datamodel.Node, error) {
	return qp.BuildMap(basicnode.Prototype.Map, 2, func(ma datamodel.MapAssembler) {
		qp.MapEntry(ma, "proof", qp.Bytes(proof))
		qp.MapEntry(ma, "cell", qp.Bytes(cell))
	})
}

// Utility function to create a new LinkSystem with an in-memory store
func createLinkSystem() *linking.LinkSystem {
	store := memstore.Store{}
	lsys := cidlink.DefaultLinkSystem()
	lsys.SetWriteStorage(&store)
	lsys.SetReadStorage(&store)

	return &lsys
}
