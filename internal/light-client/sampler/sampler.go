package sampler

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/covalenthq/das-ipfs-pinner/internal"
	"github.com/covalenthq/das-ipfs-pinner/internal/gateway"
	verifier "github.com/covalenthq/das-ipfs-pinner/internal/light-client/c-kzg-verifier"
	publisher "github.com/covalenthq/das-ipfs-pinner/internal/light-client/publisher"
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

	gh := gateway.NewHandler(gateway.DefaultGateways)

	return &Sampler{
		ipfsShell:     shell,
		gh:            gh,
		pub:           pub,
		samplingDelay: samplingDelay,
		samplingFn:    CalculateSamplesNeeded,
	}, nil
}

// ProcessEvent handles events asynchronously by processing the provided CID.
func (s *Sampler) ProcessEvent(cidStr string, blockHeight uint64) {
	go func(cidStr string) {
		rawCid, err := cid.Decode(cidStr)
		if err != nil {
			log.Errorf("Invalid CID: %v", err)
			return
		}

		if rawCid.Prefix().Codec != cid.DagCBOR {
			log.Debugf("Unsupported CID codec: %v. Skipping", rawCid.Prefix().Codec)
			return
		}

		log.Debugf("Processing event for CID [%s] is deferred for %d min", cidStr, s.samplingDelay/60)
		time.Sleep(time.Duration(s.samplingDelay) * time.Second)
		log.Debugf("Processing event for CID [%s] ...", cidStr)

		var rootNode internal.RootNode
		if err := s.GetData(cidStr, &rootNode); err != nil {
			log.Errorf("Failed to fetch root DAG data: %v", err)
			return
		}

		sampleIterations := s.samplingFn(rootNode.Size, rootNode.Size/2, 0.95)

		for rowIndex, rowLink := range rootNode.Links {
			var links []internal.Link
			if err := s.GetData(rowLink.CID, &links); err != nil {
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

				commitment := rootNode.Commitments[rowIndex].Nested.Bytes
				proof := data.Proof.Nested.Bytes
				cell := data.Cell.Nested.Bytes
				res, err := verifier.NewKZGVerifier(commitment, proof, cell, uint64(colIndex)).Verify()
				if err != nil {
					log.Errorf("Failed to verify proof and cell: %v", err)
					return
				}

				log.Infof("cell=[%2d,%3d], verified=%-5v, cid=%-40v", rowIndex, colIndex, res, cidStr)

				if err := s.pub.PublishToCS(cidStr, rowIndex, colIndex, res, commitment, proof, cell, blockHeight); err != nil {
					log.Errorf("Failed to publish to Cloud Storage: %v", err)
					return
				}
			}
		}
	}(cidStr)
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
