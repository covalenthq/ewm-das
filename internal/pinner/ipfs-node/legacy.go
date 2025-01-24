package ipfsnode

import (
	iface "github.com/ipfs/kubo/core/coreiface"
	"github.com/ipfs/kubo/core/coreiface/options"
	mh "github.com/multiformats/go-multihash"
)

func (ipfsNode *IPFSNode) UnixFs() iface.UnixfsAPI {
	return ipfsNode.api.Unixfs()
}

func (ipfsNode *IPFSNode) AddOptions(upload bool) []options.UnixfsAddOption {
	cidVersion := 1

	addOptions := []options.UnixfsAddOption{}
	addOptions = append(addOptions, options.Unixfs.CidVersion(cidVersion))
	addOptions = append(addOptions, options.Unixfs.HashOnly(!upload))
	addOptions = append(addOptions, options.Unixfs.Pin(upload))

	// default merkle dag creation options
	// we want to use the same options throughout, and provide these values explicitly
	// even if the default values by ipfs libs change in future
	addOptions = append(addOptions, options.Unixfs.Hash(mh.SHA2_256))
	addOptions = append(addOptions, options.Unixfs.Inline(false))
	addOptions = append(addOptions, options.Unixfs.InlineLimit(32))
	// for cid version, raw leaves is used by default
	addOptions = append(addOptions, options.Unixfs.RawLeaves(cidVersion == 1))
	addOptions = append(addOptions, options.Unixfs.Chunker("size-262144"))
	addOptions = append(addOptions, options.Unixfs.Layout(options.BalancedLayout))
	addOptions = append(addOptions, options.Unixfs.Nocopy(false))

	return addOptions
}
