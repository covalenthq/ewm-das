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
	format "github.com/ipfs/go-ipld-format"
	merkledag "github.com/ipfs/go-merkledag"
	iface "github.com/ipfs/kubo/core/coreiface"
	gocar "github.com/ipld/go-car"
	carutil "github.com/ipld/go-car/util"
	parse "github.com/ipld/go-ipld-prime/traversal/selector/parse"
	"github.com/multiformats/go-multihash"
)

type dagStore struct {
	dag iface.APIDagService
	ctx context.Context
}

func (ds dagStore) Get(ctx context.Context, c cid.Cid) (blocks.Block, error) {
	return ds.dag.Get(ctx, c)
}

// Pin builds a CAR of the local DAG rooted at root and uploads it to Pinata.
// The CAR is wrapped in a one-Link dag-pb root before upload because Pinata's
// `car=true` async validator silently drops CARs whose root is dag-cbor (and
// our DAS DAGs are always dag-cbor). The wrapper lets the upload succeed while
// keeping every inner block — including root — indexed by its original CID, so
// downstream consumers fetch the same dag-cbor bytes from any IPFS gateway.
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

	wrappedCAR, err := os.CreateTemp(os.TempDir(), "*-wrapped.car")
	if err != nil {
		return cid.Undef, err
	}
	defer cleanupCAR(wrappedCAR)

	wrapperCid, err := writeWrappedCAR(innerCAR, wrappedCAR, root)
	if err != nil {
		return cid.Undef, fmt.Errorf("wrap CAR for Pinata: %w", err)
	}

	if err := ipfsNode.pin.Pin(wrappedCAR, wrapperCid); err != nil {
		return cid.Undef, err
	}
	return root, nil
}

// writeWrappedCAR reads a single-root CAR from in, builds a dag-pb node whose
// only Link points to that root, and writes a new CAR to out whose root is the
// wrapper and whose body is the wrapper block followed by every original
// block. Returns the wrapper CID.
func writeWrappedCAR(in *os.File, out *os.File, expectedInnerRoot cid.Cid) (cid.Cid, error) {
	if _, err := in.Seek(0, io.SeekStart); err != nil {
		return cid.Undef, err
	}
	r, err := gocar.NewCarReader(in)
	if err != nil {
		return cid.Undef, err
	}
	if len(r.Header.Roots) != 1 || !r.Header.Roots[0].Equals(expectedInnerRoot) {
		return cid.Undef, fmt.Errorf("inner CAR has roots %v, expected single root %s",
			r.Header.Roots, expectedInnerRoot)
	}

	type rawBlock struct {
		c   cid.Cid
		raw []byte
	}
	var collected []rawBlock
	var totalSize uint64
	for {
		b, err := r.Next()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return cid.Undef, fmt.Errorf("read inner CAR: %w", err)
		}
		collected = append(collected, rawBlock{c: b.Cid(), raw: b.RawData()})
		totalSize += uint64(len(b.RawData()))
	}

	wrapper := merkledag.NodeWithData(nil)
	wrapper.SetCidBuilder(cid.V1Builder{Codec: cid.DagProtobuf, MhType: multihash.SHA2_256})
	if err := wrapper.AddRawLink("das-root", &format.Link{
		Name: "das-root",
		Size: totalSize,
		Cid:  expectedInnerRoot,
	}); err != nil {
		return cid.Undef, fmt.Errorf("build wrapper link: %w", err)
	}
	wrapperCid := wrapper.Cid()
	wrapperRaw := wrapper.RawData()

	if err := gocar.WriteHeader(&gocar.CarHeader{Roots: []cid.Cid{wrapperCid}, Version: 1}, out); err != nil {
		return cid.Undef, err
	}
	if err := carutil.LdWrite(out, wrapperCid.Bytes(), wrapperRaw); err != nil {
		return cid.Undef, err
	}
	for _, b := range collected {
		if err := carutil.LdWrite(out, b.c.Bytes(), b.raw); err != nil {
			return cid.Undef, err
		}
	}
	return wrapperCid, nil
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
