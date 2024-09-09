package ipfsnode

import (
	"context"
	"errors"
	"sync"

	"github.com/covalenthq/das-ipfs-pinner/internal"
)

type Block struct {
	Cells [][]*internal.DataMap
}

// ExtractBlock extracts the block from IPFS and downloads all cells.
func (ipfsNode *IPFSNode) ExtractBlock(ctx context.Context, cidStr string) (*Block, error) {
	var root internal.RootNode
	if err := ipfsNode.GetData(ctx, cidStr, &root); err != nil {
		return nil, err
	}

	// Pre-allocate space for cells from each of the root links (up to 13)
	downloadedCells := make([][]*internal.DataMap, len(root.Links))

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
			var rowLinks []internal.Link
			if err := ipfsNode.GetData(ctx, link.CID, &rowLinks); err != nil {
				select {
				case errorChan <- err:
				default:
				}
				return
			}

			// Download up to 64 cells from the row links
			cells, err := downloadCells(ctx, ipfsNode, rowLinks, errorChan, 64)
			if err != nil {
				return
			}

			// Store the downloaded cells in the correct index of the root link
			downloadedCells[i] = cells
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
		return combineDownloadedCells(downloadedCells), nil
	}
}

// downloadCells downloads up to the specified limit of cells from the provided row links.
// It preserves the order of the cells.
func downloadCells(ctx context.Context, ipfsNode *IPFSNode, rowLinks []internal.Link, errorChan chan<- error, limit int) ([]*internal.DataMap, error) {
	var wg sync.WaitGroup
	cells := make([]*internal.DataMap, len(rowLinks)) // Pre-allocate with the full length to maintain order
	mu := sync.Mutex{}                                // Mutex to ensure safe access to shared state
	count := 0                                        // Track number of downloaded cells

	for i, link := range rowLinks {
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

			mu.Lock()
			defer mu.Unlock()

			// Insert the cell at the correct index and increment the count
			if count < limit {
				cells[i] = &cell
				count++
			}
		}(i, link)
	}

	// Wait for all downloads to finish
	wg.Wait()

	return cells, nil
}

// combineDownloadedCells combines the downloaded cells into a block.
func combineDownloadedCells(cells [][]*internal.DataMap) *Block {
	// Combine the downloaded cells into a block.
	return &Block{
		Cells: cells,
	}
}
