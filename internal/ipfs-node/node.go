package ipfsnode

import (
	"bytes"
	"context"
	"io/fs"
	"log"
	"os"
	"path/filepath"

	ipldencoder "github.com/covalenthq/das-ipfs-pinner/internal/ipld-encoder"
	"github.com/ipfs/boxo/blockstore"
	blocks "github.com/ipfs/go-block-format"
	"github.com/ipfs/go-cid"
	config "github.com/ipfs/kubo/config"
	"github.com/ipfs/kubo/core"
	"github.com/ipfs/kubo/core/coreapi"
	iface "github.com/ipfs/kubo/core/coreiface"
	"github.com/ipfs/kubo/core/node/libp2p"
	"github.com/ipfs/kubo/plugin/loader"
	"github.com/ipfs/kubo/repo/fsrepo"
	"github.com/ipld/go-ipld-prime"
	"github.com/ipld/go-ipld-prime/codec"
	"github.com/ipld/go-ipld-prime/fluent/qp"
	cidlink "github.com/ipld/go-ipld-prime/linking/cid"
	"github.com/ipld/go-ipld-prime/multicodec"
	"github.com/ipld/go-ipld-prime/node/basicnode"
	mc "github.com/multiformats/go-multicodec"
)

// CodecConfig struct encapsulates codec encoder, decoder, and CID prefix.
type CodecConfig struct {
	Encoder codec.Encoder
	Prefix  cid.Prefix
}

// IPFSNode struct encapsulates the IPFS node and CoreAPI.
type IPFSNode struct {
	Node *core.IpfsNode
	API  iface.CoreAPI
}

// NewIPFSNode initializes and returns a new IPFSNode instance.
func NewIPFSNode() (*IPFSNode, error) {
	buildConfig := core.BuildCfg{
		Online:    true,
		Permanent: true,
		Routing:   libp2p.DHTOption,
		Host:      libp2p.DefaultHostOption,
	}

	repoPath, err := initializeRepo()
	if err != nil {
		return nil, err
	}

	if err := loadPlugins(repoPath); err != nil {
		return nil, err
	}

	ipfsConfig, err := config.Init(os.Stdout, 2048)
	if err != nil {
		return nil, err
	}

	ipfsConfig.Datastore = config.DefaultDatastoreConfig()
	if err = fsrepo.Init(repoPath, ipfsConfig); err != nil {
		return nil, err
	}

	buildConfig.Repo, err = fsrepo.Open(repoPath)
	if err != nil {
		return nil, err
	}

	node, err := core.NewNode(context.Background(), &buildConfig)
	if err != nil {
		return nil, err
	}

	api, err := coreapi.NewCoreAPI(node)
	if err != nil {
		return nil, err
	}

	return &IPFSNode{
		Node: node,
		API:  api,
	}, nil
}

// PublishBlock publishes a block to IPFS.
func (ipfsNode *IPFSNode) PublishBlock(dataBlock *ipldencoder.IPLDDataBlock) error {
	codecConfig, err := prepareCodec("dag-cbor", "sha2-256")
	if err != nil {
		return err
	}

	// Encode the block data to CBOR DAG
	for _, nodeGroup := range dataBlock.DataNodes {
		for _, node := range nodeGroup {
			var buffer bytes.Buffer
			if err := codecConfig.Encoder(node, &buffer); err != nil {
				return err
			}

			blockCid, err := codecConfig.Prefix.Sum(buffer.Bytes())
			if err != nil {
				return err
			}

			block, err := blocks.NewBlockWithCid(buffer.Bytes(), blockCid)
			if err != nil {
				return err
			}

			// TODO: maybe use kubo's approach to store the block
			// Use blockstore to store the block
			blockStore := blockstore.NewBlockstore(ipfsNode.Node.Repo.Datastore())
			if err := blockStore.Put(context.Background(), block); err != nil {
				return err
			}

			// Log the CID of the stored CBOR DAG node
			log.Printf("CBOR DAG object added to IPFS with CID: %s", blockCid.String())

			// Retrieve and log the CBOR DAG object
			retrievedNode, err := ipfsNode.API.Dag().Get(context.Background(), blockCid)
			if err != nil {
				return err
			}
			log.Printf("Retrieved CBOR DAG object with CID: %s", retrievedNode.Cid().String())
		}
	}

	for _, subLinks := range dataBlock.Links {
		node, err := qp.BuildList(basicnode.Prototype.List, -1, func(la ipld.ListAssembler) {
			for _, link := range subLinks {
				newLink := cidlink.Link{Cid: link.(cidlink.Link).Cid}
				qp.ListEntry(la, qp.Link(newLink))
			}
		})
		if err != nil {
			return err
		}

		var buffer bytes.Buffer
		if err := codecConfig.Encoder(node, &buffer); err != nil {
			return err
		}

		blockCid, err := codecConfig.Prefix.Sum(buffer.Bytes())
		if err != nil {
			return err
		}

		block, err := blocks.NewBlockWithCid(buffer.Bytes(), blockCid)
		if err != nil {
			return err
		}

		// TODO: maybe use kubo's approach to store the block
		// Use blockstore to store the block
		blockStore := blockstore.NewBlockstore(ipfsNode.Node.Repo.Datastore())
		if err := blockStore.Put(context.Background(), block); err != nil {
			return err
		}

		// Log the CID of the stored CBOR DAG node
		log.Printf("CBOR DAG object added to IPFS with CID: %s", blockCid.String())
	}

	// Root

	var buffer bytes.Buffer
	if err := codecConfig.Encoder(dataBlock.Root, &buffer); err != nil {
		return err
	}

	blockCid, err := codecConfig.Prefix.Sum(buffer.Bytes())
	if err != nil {
		return err
	}

	block, err := blocks.NewBlockWithCid(buffer.Bytes(), blockCid)
	if err != nil {
		return err
	}

	// TODO: maybe use kubo's approach to store the block
	// Use blockstore to store the block
	blockStore := blockstore.NewBlockstore(ipfsNode.Node.Repo.Datastore())
	if err := blockStore.Put(context.Background(), block); err != nil {
		return err
	}

	// Log the CID of the stored CBOR DAG node
	log.Printf("CBOR DAG object added to IPFS with CID: %s", blockCid.String())

	return nil
}

// prepareCodec sets up the encoder and CID prefix.
func prepareCodec(storageFormat, hashAlgorithm string) (*CodecConfig, error) {
	var storageCodec mc.Code
	if err := storageCodec.Set(storageFormat); err != nil {
		return nil, err
	}
	var multihashType mc.Code
	if err := multihashType.Set(hashAlgorithm); err != nil {
		return nil, err
	}

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

	return &CodecConfig{
		Encoder: encoder,
		Prefix:  cidPrefix,
	}, nil
}

// loadPlugins loads and initializes any external plugins.
func loadPlugins(pluginPath string) error {
	plugins, err := loader.NewPluginLoader(filepath.Join(pluginPath, "plugins"))
	if err != nil {
		return err
	}

	if err := plugins.Initialize(); err != nil {
		return err
	}

	if err := plugins.Inject(); err != nil {
		return err
	}

	return nil
}

// initializeRepo initializes the IPFS repository.
func initializeRepo() (string, error) {
	repoPath, err := config.PathRoot() // IPFS path root, can be changed via env variable
	if err != nil {
		return "", err
	}
	if err = os.MkdirAll(repoPath, fs.ModeDir); err != nil {
		return "", err
	}

	return repoPath, nil
}
