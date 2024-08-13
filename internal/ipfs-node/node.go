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
	"github.com/ipld/go-ipld-prime/codec/dagcbor"
	mh "github.com/multiformats/go-multihash"
)

// IPFSNode struct encapsulates the IPFS node and CoreAPI.
type IPFSNode struct {
	Node *core.IpfsNode
	API  iface.CoreAPI
}

// NewIPFSNode initializes and returns a new IPFSNode instance.
func NewIPFSNode() (*IPFSNode, error) {
	cfg := core.BuildCfg{
		Online:    true, // networking
		Permanent: true, // data persists across restarts?
	}
	cfg.Routing = libp2p.DHTOption
	cfg.Host = libp2p.DefaultHostOption

	repoPath, err := initIpfsRepo()
	if err != nil {
		return nil, err
	}

	if err := setupPlugins(repoPath); err != nil {
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

	cfg.Repo, err = fsrepo.Open(repoPath)
	if err != nil {
		return nil, err
	}

	node, err := core.NewNode(context.Background(), &cfg)
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
func (ipfsNode *IPFSNode) PublishBlock(block *ipldencoder.IPLDDataBlock) error {
	for _, subNodes := range block.DataNodes {
		for _, dataNode := range subNodes {
			var buf bytes.Buffer
			if err := dagcbor.Encode(dataNode, &buf); err != nil {
				return err
			}

			cidPrefix := cid.Prefix{
				Version:  1,
				Codec:    uint64(cid.DagCBOR),
				MhType:   uint64(mh.SHA2_256),
				MhLength: -1,
			}

			blockCid, err := cidPrefix.Sum(buf.Bytes())
			if err != nil {
				return err
			}

			blk, err := blocks.NewBlockWithCid(buf.Bytes(), blockCid)
			if err != nil {
				return err
			}

			// Use blockstore to store the block
			blockStore := blockstore.NewBlockstore(ipfsNode.Node.Repo.Datastore())
			if err := blockStore.Put(context.Background(), blk); err != nil {
				return err
			}

			// Log the CID of the stored CBOR DAG node
			cborCID := blk.Cid()
			log.Printf("CBOR DAG object added to IPFS with CID: %s", cborCID.String())

			// Example: Retrieve the CBOR DAG object back
			retrievedNode, err := ipfsNode.API.Dag().Get(context.Background(), cborCID)
			if err != nil {
				return err
			}

			log.Printf("Retrieved CBOR DAG object with CID: %s", retrievedNode.Cid().String())
		}
	}

	return nil
}

// setupPlugins loads and initializes any external plugins.
func setupPlugins(externalPluginsPath string) error {
	plugins, err := loader.NewPluginLoader(filepath.Join(externalPluginsPath, "plugins"))
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

// initIpfsRepo initializes the IPFS repository.
func initIpfsRepo() (string, error) {
	pathRoot, err := config.PathRoot() // IFPS path root, can be changed via env variable too
	if err != nil {
		return "", err
	}
	if err = os.MkdirAll(pathRoot, fs.ModeDir); err != nil {
		return "", err
	}

	return pathRoot, nil
}
