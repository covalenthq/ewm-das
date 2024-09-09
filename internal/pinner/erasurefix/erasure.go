package erasurefix

import (
	"github.com/covalenthq/das-ipfs-pinner/internal"
	ckzg4844 "github.com/ethereum/c-kzg-4844/v2/bindings/go"
)

type Erasurer struct {
}

func NewErasurer() *Erasurer {
	return &Erasurer{}
}

func (e *Erasurer) Fix(cells [][]*internal.DataMap) ([][128]ckzg4844.Cell, error) {
	block := make([][128]ckzg4844.Cell, len(cells))

	for i, row := range cells {
		kzgCells := make([]ckzg4844.Cell, len(row)/2)
		var indexes []uint64

		kzgIndex := 0
		for j, cell := range row {
			if cell == nil {
				continue
			}

			copy(kzgCells[kzgIndex][:], cell.Cell.Nested.Bytes[:])
			indexes = append(indexes, uint64(j))
			kzgIndex++
		}

		// Fix the row
		recovered, _, err := ckzg4844.RecoverCellsAndKZGProofs(indexes, kzgCells)
		if err != nil {
			return nil, err
		}

		block[i] = recovered
	}

	return block, nil
}
