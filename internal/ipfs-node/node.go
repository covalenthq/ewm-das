package ipfsnode

import (
	"context"
	"io/fs"
	"os"
	"path/filepath"

	config "github.com/ipfs/kubo/config"
	"github.com/ipfs/kubo/core"
	"github.com/ipfs/kubo/core/coreapi"
	iface "github.com/ipfs/kubo/core/coreiface"
	"github.com/ipfs/kubo/core/node/libp2p"
	"github.com/ipfs/kubo/plugin/loader"
	"github.com/ipfs/kubo/repo/fsrepo"
)

// IPFSNode struct encapsulates the IPFS node and CoreAPI.
type IPFSNode struct {
	Node *core.IpfsNode
	API  iface.CoreAPI
	W3   *W3Storage
}

// NewIPFSNode initializes and returns a new IPFSNode instance.
func NewIPFSNode(w3Key, w3DelegationProofPath string) (*IPFSNode, error) {
	w3, err := NewW3Storage(w3Key, w3DelegationProofPath)
	if err != nil {
		return nil, err
	}

	if err := w3.Initialize(); err != nil {
		return nil, err
	}

	buildConfig, err := initializeIPFSConfig()
	if err != nil {
		return nil, err
	}

	node, err := core.NewNode(context.Background(), buildConfig)
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
		W3:   w3,
	}, nil
}

// initializeIPFSConfig sets up the IPFS configuration.
func initializeIPFSConfig() (*core.BuildCfg, error) {
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

	repo, err := fsrepo.Open(repoPath)
	if err != nil {
		return nil, err
	}

	return &core.BuildCfg{
		Online:    true,
		Permanent: true,
		Routing:   libp2p.DHTOption,
		Host:      libp2p.DefaultHostOption,
		Repo:      repo,
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
