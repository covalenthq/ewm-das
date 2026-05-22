package ipfsnode

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"syscall"

	blocks "github.com/ipfs/go-block-format"
	"github.com/ipfs/go-cid"
	iface "github.com/ipfs/kubo/core/coreiface"
	gocar "github.com/ipld/go-car"
	carutil "github.com/ipld/go-car/util"
	parse "github.com/ipld/go-ipld-prime/traversal/selector/parse"
)

type dagStore struct {
	dag iface.APIDagService
	ctx context.Context
}

func (ds dagStore) Get(ctx context.Context, c cid.Cid) (blocks.Block, error) {
	return ds.dag.Get(ctx, c)
}

// Pin uploads the DAG rooted at root to Filebase and returns root on success.
//
// Two CAR passes: the first walks the DAG into a temp CAR via SelectiveCar,
// the second repacks that CAR with every block CID declared as a root in the
// header. Filebase's /dag/import?pin-roots=true then pins each block
// individually, which is the only way the dedicated gateway will serve inner
// dag-cbor CIDs (Filebase's CAR processor does not recursively pin children
// of a single-root dag-cbor CAR — see filebase.go for the underlying detail).
func (ipfsNode *IPFSNode) Pin(ctx context.Context, root cid.Cid) (cid.Cid, error) {
	innerCAR, err := os.CreateTemp(os.TempDir(), "*-inner.car")
	if err != nil {
		return cid.Undef, err
	}
	defer cleanupCAR(innerCAR)

	store := dagStore{dag: ipfsNode.api.Dag(), ctx: ctx}
	dag := gocar.Dag{Root: root, Selector: parse.CommonSelector_ExploreAllRecursively}
	scar := gocar.NewSelectiveCar(ctx, store, []gocar.Dag{dag}, gocar.TraverseLinksOnlyOnce())
	if err := scar.Write(innerCAR); err != nil {
		return cid.Undef, err
	}

	multiCAR, err := os.CreateTemp(os.TempDir(), "*-multi.car")
	if err != nil {
		return cid.Undef, err
	}
	defer cleanupCAR(multiCAR)

	allBlockCIDs, err := writeMultiRootCAR(innerCAR, multiCAR, root)
	if err != nil {
		return cid.Undef, fmt.Errorf("repack CAR with all roots: %w", err)
	}

	if err := ipfsNode.fb.Pin(multiCAR, allBlockCIDs); err != nil {
		return cid.Undef, err
	}
	return root, nil
}

// writeMultiRootCAR reads a single-root CAR from in and writes a CAR to out
// whose header lists every block CID as a root. Block bytes and order are
// preserved verbatim. Returns the full list of block CIDs (the DAS root is
// guaranteed to be first because go-car SelectiveCar emits the root before
// descendants); the caller passes this list to FilebaseStorage.Pin so every
// block is asserted to be pinned.
func writeMultiRootCAR(in *os.File, out *os.File, expectedInnerRoot cid.Cid) ([]cid.Cid, error) {
	if _, err := in.Seek(0, io.SeekStart); err != nil {
		return nil, err
	}
	r, err := gocar.NewCarReader(in)
	if err != nil {
		return nil, err
	}
	if len(r.Header.Roots) != 1 || !r.Header.Roots[0].Equals(expectedInnerRoot) {
		return nil, fmt.Errorf("inner CAR has roots %v, expected single root %s",
			r.Header.Roots, expectedInnerRoot)
	}

	type rawBlock struct {
		c   cid.Cid
		raw []byte
	}
	var blocks []rawBlock
	var roots []cid.Cid
	seen := make(map[string]struct{})
	for {
		b, err := r.Next()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("read inner CAR: %w", err)
		}
		k := b.Cid().KeyString()
		if _, dup := seen[k]; dup {
			continue
		}
		seen[k] = struct{}{}
		blocks = append(blocks, rawBlock{c: b.Cid(), raw: b.RawData()})
		roots = append(roots, b.Cid())
	}

	if err := gocar.WriteHeader(&gocar.CarHeader{Roots: roots, Version: 1}, out); err != nil {
		return nil, err
	}
	for _, b := range blocks {
		if err := carutil.LdWrite(out, b.c.Bytes(), b.raw); err != nil {
			return nil, err
		}
	}
	return roots, nil
}

func cleanupCAR(f *os.File) {
	if f == nil {
		return
	}
	_ = f.Close()
	if err := syscall.Unlink(f.Name()); err != nil {
		log.Errorf("error unlinking %s: %v", f.Name(), err)
	}
}
