package ipfsnode

import (
	"bytes"
	"context"

	"github.com/ipfs/go-cid"
	"github.com/ipld/go-ipld-prime/datamodel"
	"github.com/ipld/go-ipld-prime/multicodec"
	"github.com/ipld/go-ipld-prime/node/basicnode"

	ipldencoder "github.com/covalenthq/das-ipfs-pinner/internal/pinner/ipld-encoder"
)

func (ipfsNode *IPFSNode) ExtractBlock(ctx context.Context, cidStr string) (*ipldencoder.IPLDDataBlock, error) {
	// Parse the CID
	cid, err := cid.Parse(cidStr)
	if err != nil {
		return nil, err
	}

	// Retrieve the block from blockstore
	node, err := ipfsNode.API.Dag().Get(ctx, cid)
	if err != nil {
		return nil, err
	}

	// Decode the block
	return ipfsNode.decodeNode(ctx, cid, node.RawData())
}

func (ipfsNode *IPFSNode) decodeNode(ctx context.Context, rootCid cid.Cid, data []byte) (*ipldencoder.IPLDDataBlock, error) {
	root, err := ipfsNode.decodeData(rootCid, data)
	if err != nil {
		return nil, err
	}

	links, err := root.LookupByString("links")
	if err != nil {
		return nil, err
	}

	linkNodes, dataNodes, err := ipfsNode.processLinks(ctx, links)
	if err != nil {
		return nil, err
	}

	return &ipldencoder.IPLDDataBlock{
		Version:   rootCid.Version(),
		Codec:     rootCid.Type(),
		Links:     linkNodes,
		DataNodes: dataNodes,
		Root:      root,
	}, nil
}

func (ipfsNode *IPFSNode) processLinks(ctx context.Context, links datamodel.Node) ([][]datamodel.Link, [][]datamodel.Node, error) {
	var linkNodes [][]datamodel.Link
	var dataNodes [][]datamodel.Node

	iter := links.ListIterator()
	for !iter.Done() {
		_, linkNode, err := iter.Next()
		if err != nil {
			return nil, nil, err
		}

		link, err := linkNode.AsLink()
		if err != nil {
			return nil, nil, err
		}

		linkCid, err := cid.Parse(link.String())
		if err != nil {
			return nil, nil, err
		}

		// Get the linked node
		legacyNode, err := ipfsNode.API.Dag().Get(ctx, linkCid)
		if err != nil {
			return nil, nil, err
		}

		links, data, err := ipfsNode.decodeLinkedNode(ctx, linkCid, legacyNode.RawData())
		if err != nil {
			return nil, nil, err
		}

		linkNodes = append(linkNodes, links)
		dataNodes = append(dataNodes, data)
	}

	return linkNodes, dataNodes, nil
}

func (ipfsNode *IPFSNode) decodeLinkedNode(ctx context.Context, rootCid cid.Cid, data []byte) ([]datamodel.Link, []datamodel.Node, error) {
	root, err := ipfsNode.decodeData(rootCid, data)
	if err != nil {
		return nil, nil, err
	}

	var dataNodes []datamodel.Node
	var links []datamodel.Link

	iter := root.ListIterator()
	for !iter.Done() {
		_, linkNode, err := iter.Next()
		if err != nil {
			return nil, nil, err
		}

		link, err := linkNode.AsLink()
		if err != nil {
			return nil, nil, err
		}

		links = append(links, link)

		linkCid, err := cid.Parse(link.String())
		if err != nil {
			return nil, nil, err
		}

		// Get the linked node
		legacyNode, err := ipfsNode.API.Dag().Get(ctx, linkCid)
		if err != nil {
			return nil, nil, err
		}

		node, err := ipfsNode.decodeData(linkCid, legacyNode.RawData())
		if err != nil {
			return nil, nil, err
		}

		dataNodes = append(dataNodes, node)
	}

	return links, dataNodes, nil
}

func (ipfsNode *IPFSNode) decodeData(rootCid cid.Cid, data []byte) (datamodel.Node, error) {
	decoder, err := multicodec.LookupDecoder(uint64(rootCid.Type()))
	if err != nil {
		return nil, err
	}

	// Wrap the data in a bytes.Reader to satisfy io.Reader interface
	reader := bytes.NewReader(data)

	// Decode the data using the appropriate decoder
	nb := basicnode.Prototype.Any.NewBuilder()
	if err := decoder(nb, reader); err != nil {
		return nil, err
	}

	return nb.Build(), nil
}
