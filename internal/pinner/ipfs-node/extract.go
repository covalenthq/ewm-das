package ipfsnode

import (
	"context"
	"errors"
	"sync"

	"github.com/covalenthq/das-ipfs-pinner/internal"
)

// ExtractBlock extracts the block from IPFS and downloads all cells.
func (ipfsNode *IPFSNode) ExtractBlock(ctx context.Context, cidStr string) ([]byte, error) {
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
		return ipfsNode.combineDownloadedCells(root, downloadedCells)
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
func (ipfsNode *IPFSNode) combineDownloadedCells(root internal.RootNode, cells [][]*internal.DataMap) ([]byte, error) {
	// Fix the cells using error correction, if needed
	block, err := ipfsNode.ef.Fix(cells)
	if err != nil {
		return nil, err
	}

	data := make([]byte, root.Size) // Preallocate the final data array with exact size
	dataOffset := 0                 // Keep track of the offset in the final data array

	for _, row := range block {
		for _, cell := range row[:64] {
			cellLen := len(cell)

			// Copy every 31 bytes from each 32-byte chunk, skipping the 0 padding
			for k := 0; k < cellLen; k += 32 {
				if k+32 <= cellLen {
					// Calculate the number of remaining bytes to copy
					bytesToCopy := 31
					if dataOffset+31 > root.Size {
						bytesToCopy = root.Size - dataOffset // Ensure we don't exceed root.Size
					}

					// Copy 31 bytes from position k+1 to k+32 (skip the first byte)
					copy(data[dataOffset:], cell[k+1:k+1+bytesToCopy])
					dataOffset += bytesToCopy

					// Stop if we've copied enough data
					if dataOffset >= root.Size {
						return data, nil
					}
				}
			}
		}
	}

	return data, nil
}
