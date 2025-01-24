package ipfsnode

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"time"

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

var (
	IPFS_HTTP_GATEWAYS = []string{"https://w3s.link/ipfs/%s", "https://dweb.link/ipfs/%s", "https://ipfs.io/ipfs/%s"}
)

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
