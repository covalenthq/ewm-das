package ipfsnode

import (
	"context"
	"os"
	"path/filepath"

	"github.com/covalenthq/das-ipfs-pinner/internal/gateway"
	logging "github.com/ipfs/go-log/v2"
	config "github.com/ipfs/kubo/config"
	"github.com/ipfs/kubo/core"
	"github.com/ipfs/kubo/core/coreapi"
	iface "github.com/ipfs/kubo/core/coreiface"
	"github.com/ipfs/kubo/core/node/libp2p"
	"github.com/ipfs/kubo/plugin/loader"
	"github.com/ipfs/kubo/repo/fsrepo"
)

var log = logging.Logger("das-pinner") // Initialize the logger

// IPFSNode struct encapsulates the IPFS node and CoreAPI.
type IPFSNode struct {
	node *core.IpfsNode
	api  iface.CoreAPI
	w3   *W3Storage
	gh   *gateway.Handler
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

	gh := gateway.NewHandler(gateway.DefaultGateways)

	// Stuf from Kubo client to consider
	// err = cctx.Plugins.Start(node)
	// if err != nil {
	// 	return err
	// }
	// node.Process.AddChild(goprocess.WithTeardown(cctx.Plugins.Close))

	return &IPFSNode{
		node: node,
		api:  api,
		w3:   w3,
		gh:   gh,
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
		Routing:   libp2p.ConstructDefaultRouting(ipfsConfig, libp2p.DHTOption),
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
	if err = os.MkdirAll(repoPath, 0755); err != nil {
		return "", err
	}

	return repoPath, nil
}
