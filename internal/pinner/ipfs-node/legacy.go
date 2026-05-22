package ipfsnode

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/covalenthq/das-ipfs-pinner/internal/gateway"
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

// IPFS_HTTP_GATEWAYS is the gateway pool used by the deprecated /get
// endpoint. It is derived from gateway.DefaultGateways (the same source of
// truth the /api/v1/get path uses) so future updates to that pool flow
// through automatically and the legacy path can never drift back to a
// retired gateway. The legacy fetcher injects the CID via fmt.Sprintf, so
// each base URL is suffixed with "/ipfs/%s".
//
// The optional DEDICATED_GATEWAY is intentionally not honoured here — the
// legacy endpoint is deprecated and callers should migrate to /api/v1/get
// to pick up dedicated-gateway support and the new path's parallel-fetch
// worker pool.
var IPFS_HTTP_GATEWAYS = func() []string {
	out := make([]string, 0, len(gateway.DefaultGateways))
	for _, g := range gateway.DefaultGateways {
		out = append(out, strings.TrimRight(g, "/")+"/ipfs/%s")
	}
	return out
}()

type httpContentFetcher struct {
	cursor        int
	ipfsFetchUrls []string
}

func NewHttpContentFetcher(ipfsFetchUrls []string) *httpContentFetcher {
	return &httpContentFetcher{cursor: 0, ipfsFetchUrls: ipfsFetchUrls}
}

func (fetcher *httpContentFetcher) FetchCidViaHttp(ctx context.Context, cid string) ([]byte, error) {
	previous := fetcher.cursor

	for {
		content, err := fetcher.tryFetch(ctx, cid, fetcher.ipfsFetchUrls[fetcher.cursor])
		if err != nil {
			log.Errorf("%s", err)
		} else {
			return content, nil
		}

		fetcher.cursor = (fetcher.cursor + 1) % len(fetcher.ipfsFetchUrls)
		log.Debugf("value of cursor: %d", fetcher.cursor)
		if fetcher.cursor == previous {
			return nil, fmt.Errorf("exhausted listed gateways, but content not found")
		}
	}
}

func (fetcher *httpContentFetcher) tryFetch(ctx context.Context, cid string, url string) ([]byte, error) {
	url = fmt.Sprintf(url, cid)
	log.Debugf("trying out %s", url)
	timeoutCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := fetcher.Get(timeoutCtx, url)
	if err != nil {
		return nil, err
	}

	defer func() {
		if err := resp.Body.Close(); err != nil {
			log.Errorf("failed to close response body: %v", err)
		}
	}()
	if resp.StatusCode == 200 {
		return io.ReadAll(resp.Body)
	} else {
		return nil, fmt.Errorf("status from GET %s is %d", url, resp.StatusCode)
	}
}

func (fetcher *httpContentFetcher) Get(ctx context.Context, url string) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	resp, err := http.DefaultClient.Do(req)
	return resp, err
}
