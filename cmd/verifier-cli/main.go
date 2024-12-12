package main

import (
	"github.com/covalenthq/das-ipfs-pinner/internal"
	verifier "github.com/covalenthq/das-ipfs-pinner/internal/light-client/c-kzg-verifier"
	"github.com/covalenthq/das-ipfs-pinner/internal/light-client/sampler"
	"github.com/covalenthq/das-ipfs-pinner/internal/pinner/das"
	"github.com/ipfs/go-cid"
	logging "github.com/ipfs/go-log/v2"
	"github.com/spf13/cobra"

	ckzg4844 "github.com/ethereum/c-kzg-4844/v2/bindings/go"
)

var log = logging.Logger("ewm.verifier")

func main() {
	var (
		cid       string
		blobIndex uint
		cellIndex uint
		ipfsAddr  string
	)
	var rootCmd = &cobra.Command{
		Use:   "verifier-cli",
		Short: "A lightweight client to verify data",
		Run: func(cmd *cobra.Command, args []string) {
			logging.SetAllLoggers(logging.LevelDebug)

			config := das.LoadConfig()
			// Initialize the KZG trusted setup
			if err := das.InitializeTrustedSetup(config); err != nil {
				log.Fatalf("Failed to initialize trusted setup: %v", err)
			}

			sample(cid, blobIndex, cellIndex, ipfsAddr)
		},
	}

	rootCmd.PersistentFlags().StringVar(&cid, "cid", "", "CID of the data to verify")
	rootCmd.PersistentFlags().UintVar(&blobIndex, "blob-index", 0, "Index of the blob to verify")
	rootCmd.PersistentFlags().UintVar(&cellIndex, "cell-index", 0, "Index of the cell to verify")
	rootCmd.PersistentFlags().StringVar(&ipfsAddr, "ipfs-addr", ":5001", "IPFS node address")

	rootCmd.MarkPersistentFlagRequired("cid")
	rootCmd.MarkPersistentFlagRequired("blob-index")
	rootCmd.MarkPersistentFlagRequired("cell-index")

	if err := rootCmd.Execute(); err != nil {
		log.Error(err)
	}
}

func sample(cidStr string, blobIndex uint, cellIndex uint, ipfsAddr string) {
	s, err := sampler.NewSampler(ipfsAddr, 0, nil)
	if err != nil {
		log.Errorf("Failed to create sampler: %v", err)
		return
	}

	rawCid, err := cid.Decode(cidStr)
	if err != nil {
		log.Errorf("Invalid CID: %v", err)
		return
	}

	if rawCid.Prefix().Codec != cid.DagCBOR {
		log.Debugf("Unsupported CID codec: %v. Skipping", rawCid.Prefix().Codec)
		return
	}

	log.Debugf("Processing event for CID [%s] Blob [%d] ...", cidStr, blobIndex)

	var rootNode internal.RootNode
	if err := s.GetData(cidStr, &rootNode); err != nil {
		log.Errorf("Failed to fetch root DAG data: %v", err)
		return
	}

	stackSize := ckzg4844.CellsPerExtBlob / rootNode.Length

	log.Infof("Root CID=%s version=%s, length=%d, size=%d, links=%d", cidStr, rootNode.Version, rootNode.Length, rootNode.Size, len(rootNode.Links))

	blobLink := rootNode.Links[blobIndex]

	var links []internal.Link
	if err := s.GetData(blobLink.CID, &links); err != nil {
		log.Errorf("Failed to fetch link data: %v", err)
		return
	}

	var data internal.DataMap
	if err := s.GetData(links[cellIndex].CID, &data); err != nil {
		log.Errorf("Failed to fetch data node: %v", err)
		return
	}

	commitment := rootNode.Commitments[blobIndex].Nested.Bytes
	proof := data.Proof.Nested.Bytes
	cell := data.Cell.Nested.Bytes
	res, err := verifier.NewKZGVerifier(commitment, proof, cell, uint64(cellIndex), uint64(stackSize)).VerifyBatch()
	if err != nil {
		log.Errorf("Failed to verify proof and cell: %v", err)
		return
	}

	log.Infof("cell=[%2d,%3d], verified=%-5v, cid=%-40v", blobIndex, cellIndex, res, cidStr)
}
