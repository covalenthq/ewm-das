package ipfsnode

import (
	"bytes"
	"context"
	"fmt"

	"github.com/ipfs/go-cid"
	"github.com/ipld/go-ipld-prime/datamodel"
	"github.com/ipld/go-ipld-prime/multicodec"
	"github.com/ipld/go-ipld-prime/node/basicnode"

	ipldencoder "github.com/covalenthq/das-ipfs-pinner/internal/ipld-encoder"
)

func (ipfsNode *IPFSNode) ExtractBlock(ctx context.Context, cidStr string) (*ipldencoder.IPLDDataBlock, error) {
	// Parse the CID
	cid, err := cid.Parse(cidStr)
	if err != nil {
		return nil, err
	}

	// Retrieve the block from blockstore
	node, err := ipfsNode.API.Dag().Get(context.Background(), cid)
	if err != nil {
		return nil, err
	}

	block, err := ipfsNode.decodeNode(ctx, cid, node.RawData())
	if err != nil {
		return nil, err
	}

	// Decode the block
	return block, nil
}

func (ipfsNode *IPFSNode) decodeNode(ctx context.Context, rootCid cid.Cid, data []byte) (*ipldencoder.IPLDDataBlock, error) {
	decoder, err := multicodec.LookupDecoder(uint64(rootCid.Type()))
	if err != nil {
		return nil, err
	}

	// Wrap the data in a bytes.Reader to satisfy io.Reader interface
	reader := bytes.NewReader(data)

	// Decode the data using the DAG-CBOR decoder
	nb := basicnode.Prototype.Any.NewBuilder()

	if err := decoder(nb, reader); err != nil {
		return nil, fmt.Errorf("failed to decode CBOR data: %w", err)
	}

	root := nb.Build()

	links, err := root.LookupByString("links")
	if err != nil {
		return nil, fmt.Errorf("failed to lookup links: %w", err)
	}

	var dataNodes [][]datamodel.Node
	var linksNodes [][]datamodel.Link

	iter := links.ListIterator()
	for !iter.Done() {
		_, linkNode, err := iter.Next()
		if err != nil {
			return nil, fmt.Errorf("failed to get link node: %w", err)
		}

		// Get the link CID

		link, err := linkNode.AsLink()
		if err != nil {
			return nil, fmt.Errorf("failed to lookup link CID: %w", err)
		}

		linkCid, err := cid.Parse(link.String())
		if err != nil {
			return nil, err
		}

		// Get the link node
		legacyNode, err := ipfsNode.API.Dag().Get(ctx, linkCid)
		if err != nil {
			return nil, fmt.Errorf("failed to get link node: %w", err)
		}

		links, datanodes, err := ipfsNode.decodeNode2(ctx, linkCid, legacyNode.RawData())
		if err != nil {
			return nil, err
		}

		linksNodes = append(linksNodes, links)
		dataNodes = append(dataNodes, datanodes)
	}

	return &ipldencoder.IPLDDataBlock{
		Version:   rootCid.Version(),
		Codec:     rootCid.Type(),
		Links:     linksNodes,
		DataNodes: dataNodes,
		Root:      root,
	}, nil
}

func (ipfsNode *IPFSNode) decodeNode2(ctx context.Context, rootCid cid.Cid, data []byte) ([]datamodel.Link, []datamodel.Node, error) {
	decoder, err := multicodec.LookupDecoder(uint64(rootCid.Type()))
	if err != nil {
		return nil, nil, err
	}

	// Wrap the data in a bytes.Reader to satisfy io.Reader interface
	reader := bytes.NewReader(data)

	// Decode the data using the DAG-CBOR decoder
	nb := basicnode.Prototype.Any.NewBuilder()

	if err := decoder(nb, reader); err != nil {
		return nil, nil, fmt.Errorf("failed to decode CBOR data: %w", err)
	}

	root := nb.Build()

	var dataNodes []datamodel.Node
	var links []datamodel.Link
	iter := root.ListIterator()
	for !iter.Done() {
		_, linkNode, err := iter.Next()
		if err != nil {
			return nil, nil, fmt.Errorf("failed to get link node: %w", err)
		}

		// Get the link CID

		link, err := linkNode.AsLink()
		if err != nil {
			return nil, nil, fmt.Errorf("failed to lookup link CID: %w", err)
		}

		links = append(links, link)

		linkCid, err := cid.Parse(link.String())
		if err != nil {
			return nil, nil, err
		}

		// Get the link node
		legacyNode, err := ipfsNode.API.Dag().Get(ctx, linkCid)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to get link node: %w", err)
		}

		node, err := ipfsNode.decodeNode3(linkCid, legacyNode.RawData())
		if err != nil {
			return nil, nil, err
		}

		dataNodes = append(dataNodes, node)
	}

	return links, dataNodes, nil
}

func (ipfsNode *IPFSNode) decodeNode3(rootCid cid.Cid, data []byte) (datamodel.Node, error) {
	decoder, err := multicodec.LookupDecoder(uint64(rootCid.Type()))
	if err != nil {
		return nil, err
	}

	// Wrap the data in a bytes.Reader to satisfy io.Reader interface
	reader := bytes.NewReader(data)

	// Decode the data using the DAG-CBOR decoder
	nb := basicnode.Prototype.Any.NewBuilder()

	if err := decoder(nb, reader); err != nil {
		return nil, fmt.Errorf("failed to decode CBOR data: %w", err)
	}

	return nb.Build(), nil
}
