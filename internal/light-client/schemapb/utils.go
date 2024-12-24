package schemapb

import (
	"fmt"

	"github.com/ipfs/go-cid"
)

func (w *Workload) ReadableString() string {
	_, cid, _ := cid.CidFromBytes(w.IpfsCid)
	return fmt.Sprintf(
		"Workload:\n  ChainID: %d\n  BlockHeight: %d\n  BlobIndex: %d\n  ExpirationTimestamp: %d\n  Hash: %x\n  BlockHash: %x\n  SpecimenHash: %x\n  Commitment: %x\n  IPFSCID: %s\n  Challenge: %x",
		w.ChainId,
		w.BlockHeight,
		w.BlobIndex,
		w.ExpirationTimestamp,
		w.Hash,
		w.BlockHash,
		w.SpecimenHash,
		w.Commitment,
		cid.String(),
		w.Challenge,
	)
}
