package ipfsnode

import (
	"bytes"
	"context"
	"fmt"

	ipldencoder "github.com/covalenthq/das-ipfs-pinner/internal/pinner/ipld-encoder"
	"github.com/ipfs/boxo/blockstore"
	blocks "github.com/ipfs/go-block-format"
	cid "github.com/ipfs/go-cid"
	"github.com/ipld/go-ipld-prime"
	"github.com/ipld/go-ipld-prime/codec"
	"github.com/ipld/go-ipld-prime/fluent/qp"
	cidlink "github.com/ipld/go-ipld-prime/linking/cid"
	"github.com/ipld/go-ipld-prime/multicodec"
	"github.com/ipld/go-ipld-prime/node/basicnode"
	mc "github.com/multiformats/go-multicodec"
)

// codecConfig struct encapsulates codec encoder, decoder, and CID prefix.
type codecConfig struct {
	Encoder codec.Encoder
	Prefix  cid.Prefix
}

// PublishBlock publishes a block to IPFS.
func (ipfsNode *IPFSNode) PublishBlock(dataBlock *ipldencoder.IPLDDataBlock, pin bool) (cid.Cid, error) {
	codecConfig, err := prepareCodec(mc.Code(dataBlock.Codec), mc.Code(dataBlock.MhType))
	if err != nil {
		return cid.Undef, err
	}

	// Process data nodes
	for _, nodeGroup := range dataBlock.DataNodes {
		for _, node := range nodeGroup {
			if _, err := ipfsNode.processAndStoreNode(codecConfig, node); err != nil {
				return cid.Undef, err
			}
		}
	}

	// Process links
	for _, subLinks := range dataBlock.Links {
		node, err := buildLinkNode(subLinks)
		if err != nil {
			return cid.Undef, err
		}
		if _, err := ipfsNode.processAndStoreNode(codecConfig, node); err != nil {
			return cid.Undef, err
		}
	}

	// Process root
	rootCid, err := ipfsNode.processAndStoreNode(codecConfig, dataBlock.Root)
	if err != nil {
		return rootCid, err
	}

	if pin {
		pinnedCid, err := ipfsNode.Pin(context.Background(), rootCid)
		if err != nil {
			return cid.Undef, err
		}

		if pinnedCid != rootCid {
			return cid.Undef, fmt.Errorf("pinned CID %s does not match root CID %s", pinnedCid, rootCid)
		}
	}

	return rootCid, nil
}

// processAndStoreNode encodes, creates a block, stores it, and logs the CID.
func (ipfsNode *IPFSNode) processAndStoreNode(codecConfig *codecConfig, node ipld.Node) (cid.Cid, error) {
	buffer, blockCid, err := encodeAndCreateCID(codecConfig, node)
	if err != nil {
		return cid.Undef, err
	}

	block, err := blocks.NewBlockWithCid(buffer.Bytes(), blockCid)
	if err != nil {
		return cid.Undef, err
	}

	blockStore := blockstore.NewBlockstore(ipfsNode.node.Repo.Datastore())
	if err := blockStore.Put(context.Background(), block); err != nil {
		return cid.Undef, nil
	}

	retrievedNode, err := ipfsNode.api.Dag().Get(context.Background(), blockCid)
	if err != nil {
		return cid.Undef, err
	}
	if retrievedNode.Cid() != blockCid {
		return cid.Undef, fmt.Errorf("stored CID %s does not match retrieved CID %s", blockCid, retrievedNode.Cid())
	}

	return retrievedNode.Cid(), nil
}

// prepareCodec sets up the encoder and CID prefix.
func prepareCodec(storageCodec, multihashType mc.Code) (*codecConfig, error) {
	cidPrefix := cid.Prefix{
		Version:  1,
		Codec:    uint64(storageCodec),
		MhType:   uint64(multihashType),
		MhLength: -1,
	}

	encoder, err := multicodec.LookupEncoder(uint64(storageCodec))
	if err != nil {
		return nil, err
	}

	return &codecConfig{
		Encoder: encoder,
		Prefix:  cidPrefix,
	}, nil
}

// encodeAndCreateCID encodes the node and creates a CID.
func encodeAndCreateCID(codecConfig *codecConfig, node ipld.Node) (*bytes.Buffer, cid.Cid, error) {
	var buffer bytes.Buffer
	if err := codecConfig.Encoder(node, &buffer); err != nil {
		return nil, cid.Undef, err
	}

	blockCid, err := codecConfig.Prefix.Sum(buffer.Bytes())
	if err != nil {
		return nil, cid.Undef, err
	}

	return &buffer, blockCid, nil
}

// buildLinkNode creates an IPLD node for the given links.
func buildLinkNode(subLinks []ipld.Link) (ipld.Node, error) {
	return qp.BuildList(basicnode.Prototype.List, -1, func(la ipld.ListAssembler) {
		for _, link := range subLinks {
			newLink := cidlink.Link{Cid: link.(cidlink.Link).Cid}
			qp.ListEntry(la, qp.Link(newLink))
		}
	})
}
