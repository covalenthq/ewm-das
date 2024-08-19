package ipfsnode

import (
	"context"
	"os"
	"syscall"

	blocks "github.com/ipfs/go-block-format"
	"github.com/ipfs/go-cid"
	iface "github.com/ipfs/kubo/core/coreiface"
	gocar "github.com/ipld/go-car"
	parse "github.com/ipld/go-ipld-prime/traversal/selector/parse"
)

type dagStore struct {
	dag iface.APIDagService
	ctx context.Context
}

func (ds dagStore) Get(ctx context.Context, c cid.Cid) (blocks.Block, error) {
	return ds.dag.Get(ctx, c)
}

// Pin pins the given CID from local blockstore to W3.
func (ipfsNode *IPFSNode) Pin(ctx context.Context, root cid.Cid) (cid.Cid, error) {
	carFile, err := os.CreateTemp(os.TempDir(), "*.car")
	if err != nil {
		return cid.Undef, err
	}
	defer carFile.Close() // should delete the file due to unlink
	defer func() {
		err := syscall.Unlink(carFile.Name())
		if err != nil {
			log.Errorf("error in unlinking:%v", err)
		}
	}()

	store := dagStore{
		dag: ipfsNode.API.Dag(),
		ctx: ctx,
	}

	dag := gocar.Dag{Root: root, Selector: parse.CommonSelector_ExploreAllRecursively}
	scar := gocar.NewSelectiveCar(ctx, store, []gocar.Dag{dag}, gocar.TraverseLinksOnlyOnce())

	if err := scar.Write(carFile); err != nil {
		return cid.Undef, err
	}

	pinnedCid, err := ipfsNode.W3.Pin(carFile)
	if err != nil {
		return cid.Undef, err
	}

	return pinnedCid, nil
}
