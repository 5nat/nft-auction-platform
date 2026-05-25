package indexer

import (
	"context"
	"errors"
	"io"
	"log/slog"
	"testing"
	"time"

	"github.com/5nat/nft-auction-platform/backend/internal/model"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

type blockRange struct {
	from uint64
	to   uint64
}

type fakeScanner struct {
	latestBlock uint64
	logs        []types.Log
	blockHash   common.Hash

	latestErr    error
	filterErr    error
	blockHashErr error

	filterCalled bool
	filterFrom   uint64
	filterTo     uint64
	filterRanges []blockRange
}

func (f *fakeScanner) LatestBlockNumber(ctx context.Context) (uint64, error) {
	if f.latestErr != nil {
		return 0, f.latestErr
	}

	return f.latestBlock, nil
}

func (f *fakeScanner) FilterLogs(ctx context.Context, fromBlock uint64, toBlock uint64) ([]types.Log, error) {
	f.filterCalled = true
	f.filterFrom = fromBlock
	f.filterTo = toBlock
	f.filterRanges = append(f.filterRanges, blockRange{
		from: fromBlock,
		to:   toBlock,
	})

	if f.filterErr != nil {
		return nil, f.filterErr
	}

	return f.logs, nil
}

func (f *fakeScanner) BlockHash(ctx context.Context, blockNumber uint64) (common.Hash, error) {
	if f.blockHashErr != nil {
		return common.Hash{}, f.blockHashErr
	}

	return f.blockHash, nil
}

type fakeRepository struct {
	nextFromBlock uint64

	nextFromErr error
	upsertErr   error

	nextFromCalled bool

	upsertCalled bool
	upsertBlock  uint64
	upsertHash   string
	upsertBlocks []uint64
}

func (f *fakeRepository) NextFromBlock(ctx context.Context, startBlock uint64) (uint64, error) {
	f.nextFromCalled = true

	if f.nextFromErr != nil {
		return 0, f.nextFromErr
	}

	return f.nextFromBlock, nil
}

func (f *fakeRepository) UpsertCursor(ctx context.Context, blockNumber uint64, blockHash string) error {
	f.upsertCalled = true
	f.upsertBlock = blockNumber
	f.upsertHash = blockHash
	f.upsertBlocks = append(f.upsertBlocks, blockNumber)

	if f.upsertErr != nil {
		return f.upsertErr
	}

	return nil
}

func (f *fakeRepository) WithTx(ctx context.Context, fn func(repo EventRepository) error) error {
	return fn(f)
}

func (f *fakeRepository) InsertProcessedLog(ctx context.Context, lg types.Log, eventName string) (bool, error) {
	return true, nil
}

func (f *fakeRepository) CreateAuction(ctx context.Context, auction model.Auction) error {
	return nil
}

func (f *fakeRepository) CreateBid(ctx context.Context, bid model.Bid) error {
	return nil
}

func (f *fakeRepository) UpdateAuctionHighestBid(ctx context.Context, input UpdateAuctionHighestBidInput) error {
	return nil
}

func (f *fakeRepository) MarkAuctionEnded(ctx context.Context, input MarkAuctionEndedInput) error {
	return nil
}

func (f *fakeRepository) MarkAuctionCancelled(ctx context.Context, input MarkAuctionCancelledInput) error {
	return nil
}

func newTestIndexer(scanner ChainScanner, repo EventRepository) *Indexer {
	return &Indexer{
		scanner:       scanner,
		repo:          repo,
		chainID:       31337,
		startBlock:    0,
		confirmations: 1,
		batchSize:     500,
		pollInterval:  3 * time.Second,
		logger:        slog.New(slog.NewTextHandler(io.Discard, nil)),
	}
}

func TestRunOnce_NoNewBlocks(t *testing.T) {
	ctx := context.Background()

	scanner := &fakeScanner{
		latestBlock: 20,
	}

	repo := &fakeRepository{
		nextFromBlock: 20,
	}

	idx := newTestIndexer(scanner, repo)

	if err := idx.RunOnce(ctx); err != nil {
		t.Fatalf("RunOnce returned error: %v", err)
	}

	if scanner.filterCalled {
		t.Fatal("expected FilterLogs not to be called when there are no new confirmed blocks")
	}

	if repo.upsertCalled {
		t.Fatal("expected cursor not to be updated when there are no new confirmed blocks")
	}
}

func TestRunOnce_EmptyRangeUpdatesCursor(t *testing.T) {
	ctx := context.Background()

	blockHash := common.HexToHash("0xabc")

	scanner := &fakeScanner{
		latestBlock: 24,
		logs:        nil,
		blockHash:   blockHash,
	}

	repo := &fakeRepository{
		nextFromBlock: 20,
	}

	idx := newTestIndexer(scanner, repo)

	if err := idx.RunOnce(ctx); err != nil {
		t.Fatalf("RunOnce returned error: %v", err)
	}

	if !scanner.filterCalled {
		t.Fatal("expected FilterLogs to be called")
	}

	if scanner.filterFrom != 20 {
		t.Fatalf("expected filter from block 20, got %d", scanner.filterFrom)
	}

	if scanner.filterTo != 23 {
		t.Fatalf("expected filter to block 23, got %d", scanner.filterTo)
	}

	if !repo.upsertCalled {
		t.Fatal("expected cursor to be updated")
	}

	if repo.upsertBlock != 23 {
		t.Fatalf("expected cursor block 23, got %d", repo.upsertBlock)
	}

	if repo.upsertHash != normalizeHash(blockHash) {
		t.Fatalf("expected cursor hash %s, got %s", normalizeHash(blockHash), repo.upsertHash)
	}
}

func TestRunOnce_NotEnoughConfirmations(t *testing.T) {
	ctx := context.Background()

	scanner := &fakeScanner{
		latestBlock: 0,
	}

	repo := &fakeRepository{
		nextFromBlock: 0,
	}

	idx := newTestIndexer(scanner, repo)
	idx.confirmations = 1

	if err := idx.RunOnce(ctx); err != nil {
		t.Fatalf("RunOnce returned error: %v", err)
	}

	if repo.nextFromCalled {
		t.Fatal("expected NextFromBlock not to be called when latest block is below confirmations")
	}

	if scanner.filterCalled {
		t.Fatal("expected FilterLogs not to be called when confirmations are insufficient")
	}

	if repo.upsertCalled {
		t.Fatal("expected cursor not to be updated when confirmations are insufficient")
	}
}

func TestRunOnce_BatchRanges(t *testing.T) {
	ctx := context.Background()

	blockHash := common.HexToHash("0xabc")

	scanner := &fakeScanner{
		latestBlock: 31,
		logs:        nil,
		blockHash:   blockHash,
	}

	repo := &fakeRepository{
		nextFromBlock: 20,
	}

	idx := newTestIndexer(scanner, repo)
	idx.confirmations = 1
	idx.batchSize = 5

	if err := idx.RunOnce(ctx); err != nil {
		t.Fatalf("RunOnce returned error: %v", err)
	}

	expectedRanges := []blockRange{
		{from: 20, to: 24},
		{from: 25, to: 29},
		{from: 30, to: 30},
	}

	if len(scanner.filterRanges) != len(expectedRanges) {
		t.Fatalf(
			"expected %d filter ranges, got %d: %+v",
			len(expectedRanges),
			len(scanner.filterRanges),
			scanner.filterRanges,
		)
	}

	for i := range expectedRanges {
		if scanner.filterRanges[i] != expectedRanges[i] {
			t.Fatalf(
				"range %d mismatch: expected %+v, got %+v",
				i,
				expectedRanges[i],
				scanner.filterRanges[i],
			)
		}
	}

	expectedCursorBlocks := []uint64{24, 29, 30}

	if len(repo.upsertBlocks) != len(expectedCursorBlocks) {
		t.Fatalf(
			"expected %d cursor updates, got %d: %+v",
			len(expectedCursorBlocks),
			len(repo.upsertBlocks),
			repo.upsertBlocks,
		)
	}

	for i := range expectedCursorBlocks {
		if repo.upsertBlocks[i] != expectedCursorBlocks[i] {
			t.Fatalf(
				"cursor update %d mismatch: expected block %d, got %d",
				i,
				expectedCursorBlocks[i],
				repo.upsertBlocks[i],
			)
		}
	}

	if repo.upsertBlock != 30 {
		t.Fatalf("expected final cursor block 30, got %d", repo.upsertBlock)
	}

	if repo.upsertHash != normalizeHash(blockHash) {
		t.Fatalf("expected cursor hash %s, got %s", normalizeHash(blockHash), repo.upsertHash)
	}
}

func TestRunOnce_ReturnsErrors(t *testing.T) {
	ctx := context.Background()

	t.Run("latest block error", func(t *testing.T) {
		scanner := &fakeScanner{
			latestErr: errors.New("latest block failed"),
		}

		repo := &fakeRepository{
			nextFromBlock: 20,
		}

		idx := newTestIndexer(scanner, repo)

		if err := idx.RunOnce(ctx); err == nil {
			t.Fatal("expected RunOnce to return latest block error")
		}
	})

	t.Run("next from block error", func(t *testing.T) {
		scanner := &fakeScanner{
			latestBlock: 24,
		}

		repo := &fakeRepository{
			nextFromErr: errors.New("cursor failed"),
		}

		idx := newTestIndexer(scanner, repo)

		if err := idx.RunOnce(ctx); err == nil {
			t.Fatal("expected RunOnce to return cursor error")
		}
	})

	t.Run("filter logs error", func(t *testing.T) {
		scanner := &fakeScanner{
			latestBlock: 24,
			filterErr:   errors.New("filter logs failed"),
		}

		repo := &fakeRepository{
			nextFromBlock: 20,
		}

		idx := newTestIndexer(scanner, repo)

		if err := idx.RunOnce(ctx); err == nil {
			t.Fatal("expected RunOnce to return filter logs error")
		}
	})

	t.Run("block hash error", func(t *testing.T) {
		scanner := &fakeScanner{
			latestBlock:  24,
			logs:         nil,
			blockHashErr: errors.New("block hash failed"),
		}

		repo := &fakeRepository{
			nextFromBlock: 20,
		}

		idx := newTestIndexer(scanner, repo)

		if err := idx.RunOnce(ctx); err == nil {
			t.Fatal("expected RunOnce to return block hash error")
		}
	})

	t.Run("upsert cursor error", func(t *testing.T) {
		scanner := &fakeScanner{
			latestBlock: 24,
			logs:        nil,
			blockHash:   common.HexToHash("0xabc"),
		}

		repo := &fakeRepository{
			nextFromBlock: 20,
			upsertErr:     errors.New("upsert cursor failed"),
		}

		idx := newTestIndexer(scanner, repo)

		if err := idx.RunOnce(ctx); err == nil {
			t.Fatal("expected RunOnce to return upsert cursor error")
		}
	})
}
