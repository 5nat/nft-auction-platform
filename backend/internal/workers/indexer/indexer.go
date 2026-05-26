package indexer

import (
	"context"
	"fmt"
	"log/slog"
	"sort"
	"time"

	"github.com/5nat/nft-auction-platform/backend/internal/config"
	ethinfra "github.com/5nat/nft-auction-platform/backend/internal/infra/blockchain/ethereum"
	"github.com/5nat/nft-auction-platform/backend/internal/infra/blockchain/ethereum/bindings"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"gorm.io/gorm"
)

type Indexer struct {
	// scanner 是链上读取接口。
	// Indexer 只依赖 ChainScanner 行为，不依赖具体的 *Scanner 实现。
	// 这样后续单元测试可以注入 fakeScanner。
	scanner ChainScanner

	// repo 是数据库读写接口。
	// Indexer 只依赖 EventRepository 行为，不依赖具体 GORM 实现。
	repo EventRepository

	marketAddress common.Address
	market        *bindings.NFTAuctionMarket

	chainID       int64
	startBlock    uint64
	confirmations uint64
	batchSize     uint64
	pollInterval  time.Duration

	logger *slog.Logger
}

func New(
	db *gorm.DB,
	chainClient *ethinfra.Client,
	cfg config.Config,
	logger *slog.Logger,
) (*Indexer, error) {
	if db == nil {
		return nil, fmt.Errorf("db is required")
	}
	if chainClient == nil {
		return nil, fmt.Errorf("chain client is required")
	}
	if logger == nil {
		logger = slog.Default()
	}
	if !common.IsHexAddress(cfg.Chain.AuctionContract) {
		return nil, fmt.Errorf("invalid chain contract address: %s", cfg.Chain.AuctionContract)
	}

	marketAddress := common.HexToAddress(cfg.Chain.AuctionContract)

	market, err := bindings.NewNFTAuctionMarket(
		marketAddress,
		chainClient.EthClient(),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create market: %w", err)
	}

	batchSize := cfg.Indexer.BatchSize
	if batchSize == 0 {
		batchSize = 500
	}

	pollInterval := cfg.Indexer.PollInterval
	if pollInterval <= 0 {
		pollInterval = 3 * time.Second
	}

	repo := NewRepository(db, cfg.Chain.ChainID, marketAddress)
	scanner := NewScanner(chainClient, marketAddress)

	return &Indexer{
		scanner:       scanner,
		repo:          repo,
		marketAddress: marketAddress,
		market:        market,
		chainID:       cfg.Chain.ChainID,
		startBlock:    cfg.Chain.StartBlock,
		confirmations: cfg.Indexer.Confirmations,
		batchSize:     batchSize,
		pollInterval:  pollInterval,
		logger:        logger,
	}, nil
}

// Start 启动长期运行的 Indexer 循环
// 它会立即执行一次 RunOnce(), 然后按照 pollInterval 周期性轮询新区块，收到 ctx cancellation 时会优雅退出。
func (idx *Indexer) Start(ctx context.Context) error {
	idx.logger.Info(
		"Starting indexer",
		"chain_id", idx.chainID,
		"contract_address", normalizeAddress(idx.marketAddress),
		"start_block", idx.startBlock,
		"confirmations", idx.confirmations,
		"batch_size", idx.batchSize,
		"poll_interval", idx.pollInterval.String(),
	)

	for {
		if err := idx.RunOnce(ctx); err != nil {
			if ctx.Err() != nil {
				idx.logger.Info("Indexer context cancelled")
				return nil
			}

			// 真实 Indexer 运行中，RPC 偶尔超时、节点短暂不可用都可能发生。
			idx.logger.Error("indexer run once failed", "error", err)
		}

		timer := time.NewTimer(idx.pollInterval)

		select {
		case <-ctx.Done():
			if !timer.Stop() {
				<-timer.C
			}
			idx.logger.Info("Indexer context cancelled")
			return nil
		case <-timer.C:
		}
	}
}

func (idx *Indexer) RunOnce(ctx context.Context) error {
	latestBlock, err := idx.scanner.LatestBlockNumber(ctx)
	if err != nil {
		return fmt.Errorf("failed to get latest block number: %w", err)
	}

	if latestBlock < idx.confirmations {
		idx.logger.Info(
			"not enough block for confirmations",
			"latest_block", latestBlock,
			"confirmations", idx.confirmations,
		)
		return nil
	}

	targetBlock := latestBlock - idx.confirmations

	fromBlock, err := idx.nextFromBlock(ctx)
	if err != nil {
		return fmt.Errorf("failed to get next block: %w", err)
	}
	if fromBlock > targetBlock {
		idx.logger.Info(
			"no new blocks to index",
			"from_block", fromBlock,
			"target_block", targetBlock,
			"latest_block", latestBlock,
		)
		return nil
	}

	idx.logger.Info(
		"indexing new block",
		"from_block", fromBlock,
		"target_block", targetBlock,
		"latest_block", latestBlock,
		"confirmations", idx.confirmations,
		"batch_size", idx.batchSize,
	)

	for fromBlock <= targetBlock {
		toBlock := fromBlock + idx.batchSize - 1
		if toBlock > targetBlock {
			toBlock = targetBlock
		}

		if processRangeErr := idx.processRange(ctx, fromBlock, toBlock); processRangeErr != nil {
			return fmt.Errorf("failed to process range %d-%d: %w", fromBlock, toBlock, processRangeErr)
		}

		if updateCursorErr := idx.updateCursor(ctx, toBlock); updateCursorErr != nil {
			return fmt.Errorf("update cursor to block %d: %w", toBlock, updateCursorErr)
		}

		fromBlock = toBlock + 1
	}

	idx.logger.Info("indexer run once completed", "target_block", targetBlock)

	return nil
}

func (idx *Indexer) processRange(ctx context.Context, fromBlock uint64, toBlock uint64) error {
	logs, err := idx.scanner.FilterLogs(ctx, fromBlock, toBlock)
	if err != nil {
		return err
	}

	sort.Slice(logs, func(i, j int) bool {
		if logs[i].BlockNumber != logs[j].BlockNumber {
			return logs[i].BlockNumber < logs[j].BlockNumber
		}
		if logs[i].TxIndex != logs[j].TxIndex {
			return logs[i].TxIndex < logs[j].TxIndex
		}
		return logs[i].Index < logs[j].Index
	})

	processed := 0
	skipped := 0

	for _, log := range logs {
		ok, processLogErr := idx.processLog(ctx, log)
		if processLogErr != nil {
			return processLogErr
		}
		if ok {
			processed++
		} else {
			skipped++
		}
	}

	idx.logger.Info(
		"block range processed",
		"from_block", fromBlock,
		"to_block", toBlock,
		"logs", len(logs),
		"processed", processed,
		"skipped", skipped,
	)

	return nil
}

func (idx *Indexer) processLog(ctx context.Context, lg types.Log) (bool, error) {
	if len(lg.Topics) == 0 {
		return false, nil
	}

	switch lg.Topics[0] {
	case auctionCreatedTopic:
		return idx.processAuctionCreated(ctx, lg)
	case bidPlacedTopic:
		return idx.processBidPlaced(ctx, lg)
	case auctionEndedTopic:
		return idx.processAuctionEnded(ctx, lg)
	case auctionCancelledTopic:
		return idx.processAuctionCancelled(ctx, lg)
	default:
		idx.logger.Debug(
			"unsupported event skipped",
			"topic", lg.Topics[0].Hex(),
			"tx_hash", lg.TxHash.Hex(),
			"log_index", lg.Index,
		)
		return false, nil
	}
}
