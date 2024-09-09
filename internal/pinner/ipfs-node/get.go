package ipfsnode

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"

	"github.com/ipfs/go-cid"
	"github.com/ipld/go-ipld-prime/codec/dagjson"
	"github.com/ipld/go-ipld-prime/multicodec"
	"github.com/ipld/go-ipld-prime/node/basicnode"
)

type fetchContext struct {
	Data    interface{}
	Context string
	Err     error
}

// GetData concurrently fetches data from both IPFS and gateways.
func (s *IPFSNode) GetData(ctx context.Context, cidStr string, data interface{}) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	results := make(chan fetchContext, 2)

	go s.fetchDataFromIPFS(ctx, cidStr, data, results)
	// go s.fetchDataFromGateways(ctx, cidStr, data, results)

	for i := 0; i < 2; i++ {
		select {
		case res := <-results:
			if res.Err == nil {
				log.Debugf("Data fetched from %s", res.Context)
				return nil
			}
			log.Debugf("Error getting data from %s: %v", res.Context, res.Err)
			if res.Context == "Gateways" {
				return res.Err
			}
		case <-ctx.Done():
			return fmt.Errorf("operation canceled")
		}
	}
	return fmt.Errorf("failed to fetch data")
}

// fetchDataFromIPFS starts a concurrent fetch from the IPFS node.
func (s *IPFSNode) fetchDataFromIPFS(ctx context.Context, cidStr string, data interface{}, results chan<- fetchContext) {
	err := s.FetchFromDagApi(ctx, cidStr, data)
	results <- fetchContext{Data: data, Context: "IPFS node", Err: err}
}

// fetchDataFromGateways starts a concurrent fetch from the gateways.
// func (s *IPFSNode) fetchDataFromGateways(ctx context.Context, cidStr string, data interface{}, results chan<- fetchContext) {
// 	err := s.gh.FetchFromGateways(ctx, cidStr, data)
// 	results <- fetchContext{Data: data, Context: "Gateways", Err: err}
// }

// FetchFromDagApi fetches data from IPFS using the Dag API and decodes it.
func (ipfsNode *IPFSNode) FetchFromDagApi(ctx context.Context, cidStr string, data interface{}) error {
	parsedCid, err := cid.Parse(cidStr)
	if err != nil {
		return err
	}

	node, err := ipfsNode.api.Dag().Get(ctx, parsedCid)
	if err != nil {
		return err
	}

	decoder, err := multicodec.LookupDecoder(uint64(parsedCid.Type()))
	if err != nil {
		return err
	}

	reader := bytes.NewReader(node.RawData())
	nb := basicnode.Prototype.Any.NewBuilder()
	if err := decoder(nb, reader); err != nil {
		return err
	}

	var jsonData bytes.Buffer
	if err := dagjson.Encode(nb.Build(), &jsonData); err != nil {
		return err
	}

	return json.Unmarshal(jsonData.Bytes(), data)
}
