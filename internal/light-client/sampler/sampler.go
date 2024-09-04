package sampler

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/covalenthq/das-ipfs-pinner/internal"
	gatewayhandler "github.com/covalenthq/das-ipfs-pinner/internal/gateway-handler"
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
	gh            *gatewayhandler.GatewayHandler
	pub           *publisher.Publisher
	samplingDelay uint
	samplingFn    func(int, int, float64) int
}

type errorContext struct {
	Err     error
	Context string
}

type resultContext struct {
	Result  interface{}
	Context string
}

// NewSampler creates a new Sampler instance and checks the connection to the IPFS daemon.
func NewSampler(ipfsAddr string, samplingDelay uint, pub *publisher.Publisher) (*Sampler, error) {
	shell := ipfs.NewShell(ipfsAddr)

	if _, _, err := shell.Version(); err != nil {
		return nil, fmt.Errorf("failed to connect to IPFS daemon: %w", err)
	}

	gatewayhandler := gatewayhandler.NewGatewayHandler(gatewayhandler.DefaultGateways)

	return &Sampler{
		ipfsShell:     shell,
		gh:            gatewayhandler,
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
	// Create a context with cancel to stop the other fetch when one completes.
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel() // Ensure all routines are stopped when done.

	// Channels to capture results and errors
	resultChan := make(chan resultContext)
	errorChan := make(chan errorContext)

	// Run IPFS node fetch concurrently
	go func() {
		if err := s.ipfsShell.DagGet(cidStr, &data); err != nil {
			errorChan <- errorContext{Err: err, Context: "IPFS node"}
		} else {
			resultChan <- resultContext{Result: data, Context: "IPFS node"}
		}
	}()

	// Run gateway fetch concurrently
	go func() {
		if err := s.gh.FetchFromGateways(ctx, cidStr, data); err != nil {
			errorChan <- errorContext{Err: err, Context: "Gateways"}
		} else {
			resultChan <- resultContext{Result: data, Context: "Gateways"}
		}
	}()

	// Wait for the first successful result
	for {
		select {
		case result := <-resultChan:
			// Successful fetch, cancel the other operation
			cancel()
			log.Debugf("Data fetched from %s", result.Context)
			return nil

		case errCxt := <-errorChan:
			log.Debugf("Error getting data from %s: %v", errCxt.Context, errCxt.Err)
			// If we receive errors from both sources, return the last error.
			if errCxt.Context == "Gateways" {
				return errCxt.Err
			}

		case <-ctx.Done():
			// If the context is canceled, return.
			return nil
		}
	}
}
