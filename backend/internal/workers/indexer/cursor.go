package indexer

import (
	"context"
	"fmt"
)

func (idx *Indexer) nextFromBlock(ctx context.Context) (uint64, error) {
	fromBlock, err := idx.repo.NextFromBlock(ctx, idx.startBlock)

	if err != nil {
		return 0, err
	}

	return fromBlock, nil
}

func (idx *Indexer) updateCursor(ctx context.Context, blockNumber uint64) error {
	blockHash, err := idx.scanner.BlockHash(ctx, blockNumber)
	if err != nil {
		return fmt.Errorf("get block hash for cursor: %w", err)
	}

	return idx.repo.UpsertCursor(ctx, blockNumber, normalizeHash(blockHash))
}
