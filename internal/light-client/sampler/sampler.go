package sampler

import (
	"context"
	"encoding/base64"
	"fmt"
	"math/rand"
	"time"

	"github.com/covalenthq/das-ipfs-pinner/common"
	"github.com/covalenthq/das-ipfs-pinner/internal"
	"github.com/covalenthq/das-ipfs-pinner/internal/gateway"
	verifier "github.com/covalenthq/das-ipfs-pinner/internal/light-client/c-kzg-verifier"
	publisher "github.com/covalenthq/das-ipfs-pinner/internal/light-client/publisher"
	ckzg4844 "github.com/ethereum/c-kzg-4844/v2/bindings/go"
	"github.com/ipfs/go-cid"
	ipfs "github.com/ipfs/go-ipfs-api"
	logging "github.com/ipfs/go-log/v2"
)

var log = logging.Logger("light-client")

// Sampler is a struct that samples data from IPFS and verifies it.
type Sampler struct {
	ipfsShell     *ipfs.Shell
	gh            *gateway.Handler
	pub           *publisher.Publisher
	samplingDelay uint
	samplingFn    func(int, int, float64) int
}

// NewSampler creates a new Sampler instance and checks the connection to the IPFS daemon.
func NewSampler(ipfsAddr string, samplingDelay uint, pub *publisher.Publisher) (*Sampler, error) {
	shell := ipfs.NewShell(ipfsAddr)

	if _, _, err := shell.Version(); err != nil {
		return nil, fmt.Errorf("failed to connect to IPFS daemon: %w", err)
	}

	gh := gateway.NewHandler(gateway.DefaultGateways, 3)

	return &Sampler{
		ipfsShell:     shell,
		gh:            gh,
		pub:           pub,
		samplingDelay: samplingDelay,
		samplingFn:    CalculateSamplesNeeded,
	}, nil
}

// ProcessEvent handles events asynchronously by processing the provided CID.
func (s *Sampler) ProcessEvent(request internal.SamplingRequest, signature []byte) {
	go func(request internal.SamplingRequest) {
		rawCid, err := cid.Decode(request.Cid)
		if err != nil {
			log.Errorf("Invalid CID: %v", err)
			return
		}

		if rawCid.Prefix().Codec != cid.DagCBOR {
			log.Debugf("Unsupported CID codec: %v. Skipping", rawCid.Prefix().Codec)
			return
		}

		log.Debugf("Processing event for CID [%s] is deferred for %d min", request.Cid, s.samplingDelay/60)
		time.Sleep(time.Duration(s.samplingDelay) * time.Second)
		log.Debugf("Processing event for CID [%s] ...", request.Cid)

		var rootNode internal.RootNode
		if err := s.GetData(request.Cid, &rootNode); err != nil {
			log.Errorf("Failed to fetch root DAG data: %v", err)
			return
		}

		sampleIterations := s.samplingFn(rootNode.Length, rootNode.Length/2, 0.95)
		stackSize := ckzg4844.CellsPerExtBlob / rootNode.Length

		log.Debugf("Sampling %d cells from %d blobs with stack size %d", sampleIterations, rootNode.Length, stackSize)

		for blobIndex, blobLink := range rootNode.Links {
			var links []internal.Link
			if err := s.GetData(blobLink.CID, &links); err != nil {
				log.Errorf("Failed to fetch link data: %v", err)
				return
			}

			// Track sampled column indices to avoid duplicates
			sampledCols := make(map[int]bool)

			for i := 0; i < sampleIterations; i++ {
				// Find a unique column index that hasn't been sampled yet
				var colIndex int
				for {
					colIndex = rand.Intn(len(links))
					if !sampledCols[colIndex] {
						sampledCols[colIndex] = true
						break
					}
				}

				var data internal.DataMap
				if err := s.GetData(links[colIndex].CID, &data); err != nil {
					log.Errorf("Failed to fetch data node: %v", err)
					return
				}

				commitment := rootNode.Commitments[blobIndex].Nested.Bytes
				proof := data.Proof.Nested.Bytes
				cell := data.Cell.Nested.Bytes
				res, err := verifier.NewKZGVerifier(commitment, proof, cell, uint64(colIndex), uint64(stackSize)).VerifyBatch()
				if err != nil {
					log.Errorf("Failed to verify proof and cell: %v", err)
					return
				}

				log.Infof("cell=[%2d,%3d], verified=%-5v, cid=%-40v", blobIndex, colIndex, res, request.Cid)

				storeReq := internal.StoreRequest{
					SamplingReqest:    request,
					SamplingSignature: fmt.Sprintf("%x", signature),
					Status:            res,
					Commitment:        base64.StdEncoding.EncodeToString(commitment),
					Proof:             base64.StdEncoding.EncodeToString(proof),
					Cell:              base64.StdEncoding.EncodeToString(cell),
					Version:           fmt.Sprintf("%s-%s", common.Version, common.GitCommit),
					BlobIndex:         blobIndex,
					CellIndex:         colIndex,
				}

				if err := s.pub.SendStoreRequest(&storeReq); err != nil {
					log.Errorf("Failed to publish to Cloud Storage: %v", err)
					return
				}
			}
		}
	}(request)
}

// GetData tries to fetch data from both the IPFS node and gateways simultaneously.
// It cancels the other request as soon as one returns successfully.
func (s *Sampler) GetData(cidStr string, data interface{}) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	results := make(chan error, 1) // Buffer of 1 to capture first result

	go s.fetchDataFromIPFS(ctx, cidStr, data, results)
	go s.fetchDataFromGateways(ctx, cidStr, data, results)

	// Return immediately after the first successful result
	select {
	case <-ctx.Done():
		return fmt.Errorf("operation canceled")
	case err := <-results:
		return err
	}
}

// fetchDataFromIPFS starts a concurrent fetch from the IPFS node.
func (s *Sampler) fetchDataFromIPFS(ctx context.Context, cidStr string, data interface{}, results chan<- error) {
	if err := s.ipfsShell.DagGet(cidStr, &data); err != nil {
		select {
		case results <- fmt.Errorf("IPFS node: %v", err):
		case <-ctx.Done():
		}
		return
	}
	select {
	case results <- nil:
	case <-ctx.Done():
	}
}

// fetchDataFromGateways starts a concurrent fetch from the gateways.
func (s *Sampler) fetchDataFromGateways(ctx context.Context, cidStr string, data interface{}, results chan<- error) {
	if err := s.gh.FetchFromGateways(ctx, cidStr, data); err != nil {
		select {
		case results <- fmt.Errorf("gateways: %v", err):
		case <-ctx.Done():
		}
		return
	}
	select {
	case results <- nil:
	case <-ctx.Done():
	}
}
