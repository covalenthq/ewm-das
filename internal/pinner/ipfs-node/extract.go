package ipfsnode

import (
	"context"
	"errors"
	"sync"

	"github.com/covalenthq/das-ipfs-pinner/internal"
	ckzgencoder "github.com/covalenthq/das-ipfs-pinner/internal/pinner/c-kzg-encoder"
)

// ExtractData extracts the block from IPFS and downloads all cells.
func (ipfsNode *IPFSNode) ExtractData(ctx context.Context, cidStr string) ([]byte, error) {
	var root internal.RootNode
	if err := ipfsNode.GetData(ctx, cidStr, &root); err != nil {
		return nil, err
	}

	byteCells := make([][][]byte, len(root.Links))

	// Channel for handling errors
	errorChan := make(chan error, 1)

	// WaitGroup to synchronize the completion of all downloads
	var wg sync.WaitGroup

	// Context with cancelation to stop downloads if an error occurs
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// Start processing each root link in parallel
	for i, link := range root.Links {
		wg.Add(1)
		go func(i int, link internal.Link) {
			defer wg.Done()

			// Fetch the next set of links (128 cells per link)
			var blobLinks []internal.Link
			if err := ipfsNode.GetData(ctx, link.CID, &blobLinks); err != nil {
				select {
				case errorChan <- err:
				default:
				}
				return
			}

			// Download up to 64 cells from the blob links
			err := downloadCells(ctx, byteCells, ipfsNode, i, blobLinks, errorChan, 64)
			if err != nil {
				return
			}
		}(i, link)
	}

	// Goroutine to close error channel when all downloads are done
	go func() {
		wg.Wait()
		close(errorChan)
	}()

	// Wait for either all downloads to complete or an error to occur
	select {
	case err := <-errorChan:
		// If an error occurs, cancel the remaining operations
		cancel()
		return nil, err
	case <-ctx.Done():
		// If the context is canceled, return an error
		return nil, errors.New("context canceled")
	case <-errorChan:
		// All downloads completed successfully, combine them into a block
		return combineDownloadedCells(root, byteCells)
	}
}

// downloadCells downloads up to the specified limit of cells from the provided blob links.
// It preserves the order of the cells.
func downloadCells(ctx context.Context, byteCells [][][]byte, ipfsNode *IPFSNode, blobIndex int, blobLinks []internal.Link, errorChan chan<- error, limit int) error {
	var wg sync.WaitGroup
	mu := sync.Mutex{} // Mutex to ensure safe access to shared state
	count := 0         // Track number of downloaded cells

	byteCells[blobIndex] = make([][]byte, 128)

	for i, link := range blobLinks {
		if count >= limit {
			break
		}

		wg.Add(1)
		go func(i int, link internal.Link) {
			defer wg.Done()

			var cell internal.DataMap
			if err := ipfsNode.GetData(ctx, link.CID, &cell); err != nil {
				select {
				case errorChan <- err:
				default:
				}
				return
			}

			// Allocate space for each byte slice within the cell
			cellBytes := make([]byte, len(cell.Cell.Nested.Bytes))

			mu.Lock()
			defer mu.Unlock()

			// Insert the cell at the correct index and increment the count
			if count < limit {
				copy(cellBytes, cell.Cell.Nested.Bytes)
				byteCells[blobIndex][i] = cellBytes
				count++
			}
		}(i, link)
	}

	// Wait for all downloads to finish
	wg.Wait()

	return nil
}

// combineDownloadedCells combines the downloaded cells into a block.
func combineDownloadedCells(root internal.RootNode, byteCells [][][]byte) ([]byte, error) {
	dataBlock := ckzgencoder.NewDataBlock()
	dataBlock.Init(uint64(root.Size), uint64(len(root.Links)))

	if err := dataBlock.RecoverData(byteCells); err != nil {
		return nil, err
	}
	return dataBlock.Decode()
}
