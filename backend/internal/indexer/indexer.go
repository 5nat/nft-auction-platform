package indexer

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"math/big"
	"sort"
	"strings"
	"time"

	"github.com/5nat/nft-auction-platform/backend/internal/chain"
	"github.com/5nat/nft-auction-platform/backend/internal/chain/bindings"
	"github.com/5nat/nft-auction-platform/backend/internal/config"
	"github.com/5nat/nft-auction-platform/backend/internal/model"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var (
	auctionCreatedTopic   = crypto.Keccak256Hash([]byte("AuctionCreated(uint256,address,address,uint256,uint256,uint256)"))
	bidPlacedTopic        = crypto.Keccak256Hash([]byte("BidPlaced(uint256,address,address,uint256,uint256)"))
	auctionEndedTopic     = crypto.Keccak256Hash([]byte("AuctionEnded(uint256,address,address,uint256,uint256)"))
	auctionCancelledTopic = crypto.Keccak256Hash([]byte("AuctionCancelled(uint256)"))
)

type Indexer struct {
	db    *gorm.DB
	chain *chain.Client

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
	chainClient *chain.Client,
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

	return &Indexer{
		db:            db,
		chain:         chainClient,
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
		if err := idx.RunOne(ctx); err != nil {
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

func (idx *Indexer) RunOne(ctx context.Context) error {
	latestBlock, err := idx.chain.LatestBlockNumber(ctx)
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

func (idx *Indexer) nextFromBlock(ctx context.Context) (uint64, error) {
	var cursor model.SyncCursor

	err := idx.db.WithContext(ctx).
		Where("chain_id = ? AND contract_address =?", idx.chainID, normalizeAddress(idx.marketAddress)).
		First(&cursor).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return idx.startBlock, nil
	}

	if err != nil {
		return 0, err
	}

	return cursor.LastProcessedBlock + 1, nil
}

func (idx *Indexer) processRange(ctx context.Context, fromBlock uint64, toBlock uint64) error {
	logs, err := idx.chain.EthClient().FilterLogs(ctx, ethereum.FilterQuery{
		FromBlock: new(big.Int).SetUint64(fromBlock),
		ToBlock:   new(big.Int).SetUint64(toBlock),
		Addresses: []common.Address{idx.marketAddress},
	})
	if err != nil {
		return fmt.Errorf("failed to filter logs: %w", err)
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

func (idx *Indexer) processAuctionCreated(ctx context.Context, lg types.Log) (bool, error) {
	event, err := idx.market.ParseAuctionCreated(lg)
	if err != nil {
		return false, fmt.Errorf("failed to parse AuctionCreated: %w", err)
	}

	inserted := false

	err = idx.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var insertProcessedLogErr error
		inserted, insertProcessedLogErr = idx.insertProcessedLog(tx, lg, "AuctionCreated")
		if insertProcessedLogErr != nil {
			return insertProcessedLogErr
		}
		if !inserted {
			return nil
		}
		if !event.AuctionId.IsUint64() {
			return fmt.Errorf("auction id overflow: %s", event.AuctionId.String())
		}
		if !event.EndTime.IsUint64() {
			return fmt.Errorf("end time overflow: %s", event.EndTime.String())
		}

		auction := model.Auction{
			ChainID:            idx.chainID,
			ContractAddress:    normalizeAddress(idx.marketAddress),
			AuctionID:          event.AuctionId.Uint64(),
			Seller:             normalizeAddress(event.Seller),
			NFTContract:        normalizeAddress(event.Nft),
			TokenID:            event.TokenId.String(),
			MinBidUSD:          event.MinBidUsd.String(),
			HighestBidder:      "",
			HighestBidToken:    "",
			HighestBidAmount:   "0",
			HighestBidUSD:      "0",
			EndTime:            event.EndTime.Uint64(),
			Status:             model.AuctionStatusActive,
			CreatedTxHash:      normalizeHash(lg.TxHash),
			CreatedBlockNumber: lg.BlockNumber,
		}

		if insertAuctionErr := tx.Clauses(clause.OnConflict{
			Columns: []clause.Column{
				{Name: "chain_id"},
				{Name: "contract_address"},
				{Name: "auction_id"},
			},
			DoNothing: true,
		}).Create(&auction).Error; insertAuctionErr != nil {
			return fmt.Errorf("failed to insert auction: %w", insertAuctionErr)
		}

		return nil
	})

	if err != nil {
		return false, err
	}

	if !inserted {
		idx.logger.Debug(
			"AuctionCreated already processed",
			"auction_id", event.AuctionId.String(),
			"tx_hash", lg.TxHash.Hex(),
			"log_index", lg.Index,
		)
		return false, nil
	}

	idx.logger.Info(
		"AuctionCreated indexed",
		"auction_id", event.AuctionId.String(),
		"seller", event.Seller.Hex(),
		"nft", event.Nft.Hex(),
		"token_id", event.TokenId.String(),
		"block_number", lg.BlockNumber,
		"tx_hash", lg.TxHash.Hex(),
		"log_index", lg.Index,
	)

	return true, nil
}

func (idx *Indexer) processBidPlaced(ctx context.Context, lg types.Log) (bool, error) {
	event, err := idx.market.ParseBidPlaced(lg)
	if err != nil {
		return false, fmt.Errorf("failed to parse BidPlaced: %w", err)
	}

	inserted := false

	err = idx.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var insertProcessedLogErr error
		inserted, insertProcessedLogErr = idx.insertProcessedLog(tx, lg, "BidPlaced")

		if insertProcessedLogErr != nil {
			return insertProcessedLogErr
		}

		if !inserted {
			return nil
		}

		if !event.AuctionId.IsUint64() {
			return fmt.Errorf("auction id overflow: %s", event.AuctionId.String())
		}

		auctionID := event.AuctionId.Uint64()

		bid := model.Bid{
			ChainID:         idx.chainID,
			ContractAddress: normalizeAddress(idx.marketAddress),
			AuctionID:       auctionID,
			Bidder:          normalizeAddress(event.Bidder),
			BidToken:        normalizeAddress(event.BidToken),
			Amount:          event.Amount.String(),
			AmountUSD:       event.AmountUsd.String(),
			TxHash:          normalizeHash(lg.TxHash),
			LogIndex:        uint64(lg.Index),
			BlockNumber:     lg.BlockNumber,
			BlockHash:       normalizeHash(lg.BlockHash),
		}

		if createBidLogErr := tx.Clauses(clause.OnConflict{
			Columns: []clause.Column{
				{Name: "chain_id"},
				{Name: "contract_address"},
				{Name: "tx_hash"},
				{Name: "log_index"},
			},
		}).Create(&bid).Error; createBidLogErr != nil {
			return fmt.Errorf("failed to insert bid: %w", createBidLogErr)
		}

		res := tx.Model(&model.Auction{}).
			Where(
				"chain_id = ? AND contract_address = ? AND auction_id = ?",
				idx.chainID,
				normalizeAddress(idx.marketAddress),
				auctionID,
			).
			Updates(map[string]any{
				"highest_bidder":     normalizeAddress(event.Bidder),
				"highest_bid_token":  normalizeAddress(event.BidToken),
				"highest_bid_amount": event.Amount.String(),
				"highest_bid_usd":    event.AmountUsd.String(),
			})

		if res.Error != nil {
			return fmt.Errorf("failed to update auction: %w", res.Error)
		}

		if res.RowsAffected == 0 {
			return fmt.Errorf("auction not found when applying BidPlaced: auction_id=%d", auctionID)
		}
		return nil
	})

	if err != nil {
		return false, err
	}

	if !inserted {
		idx.logger.Debug(
			"BIdPlaced already processed",
			"auction_id", event.AuctionId.String(),
			"tx_hash", lg.TxHash.Hex(),
			"log_index", lg.Index,
		)
	}

	idx.logger.Info(
		"BidPlaced indexed",
		"auction_id", event.AuctionId.String(),
		"bidder", event.Bidder.Hex(),
		"bid_token", event.BidToken.Hex(),
		"amount", event.Amount.String(),
		"amount_usd", event.AmountUsd.String(),
		"block_number", lg.BlockNumber,
		"tx_hash", lg.TxHash.Hex(),
		"log_index", lg.Index,
	)

	return true, nil
}

func (idx *Indexer) processAuctionEnded(ctx context.Context, lg types.Log) (bool, error) {
	event, err := idx.market.ParseAuctionEnded(lg)
	if err != nil {
		return false, fmt.Errorf("failed to parse AuctionEnded: %w", err)
	}

	inserted := false

	err = idx.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var insertProcessedLogErr error
		inserted, insertProcessedLogErr = idx.insertProcessedLog(tx, lg, "AuctionEnded")
		if insertProcessedLogErr != nil {
			return insertProcessedLogErr
		}

		if !inserted {
			return nil
		}

		if !event.AuctionId.IsUint64() {
			return fmt.Errorf("auction id overflow: %s", event.AuctionId.String())
		}

		auctionID := event.AuctionId.Uint64()

		res := tx.Model(&model.Auction{}).
			Where(
				"chain_id = ? AND contract_address = ? AND auction_id = ?",
				idx.chainID,
				normalizeAddress(idx.marketAddress),
				auctionID,
			).Updates(map[string]any{
			"status":             model.AuctionStatusEnded,
			"highest_bidder":     normalizeAddress(event.Winner),
			"highest_bid_token":  normalizeAddress(event.BidToken),
			"highest_bid_amount": event.Amount.String(),
			"highest_bid_usd":    event.AmountUsd.String(),
		})

		if res.Error != nil {
			return fmt.Errorf("failed to update auction: %w", res.Error)
		}

		if res.RowsAffected == 0 {
			return fmt.Errorf("auction not found when applying AuctionEnded: auction_id=%d", auctionID)
		}

		return nil
	})

	if err != nil {
		return false, err
	}

	if !inserted {
		idx.logger.Debug(
			"AuctionEnded already processed",
			"auction_id", event.AuctionId.String(),
			"tx_hash", lg.TxHash.Hex(),
			"log_index", lg.Index,
		)
		return false, nil
	}

	idx.logger.Info(
		"AuctionEnded indexed",
		"auction_id", event.AuctionId.String(),
		"winner", event.Winner.Hex(),
		"bid_token", event.BidToken.Hex(),
		"amount", event.Amount.String(),
		"amount_usd", event.AmountUsd.String(),
		"block_number", lg.BlockNumber,
		"tx_hash", lg.TxHash.Hex(),
		"log_index", lg.Index,
	)

	return true, nil
}

func (idx *Indexer) processAuctionCancelled(ctx context.Context, lg types.Log) (bool, error) {
	event, err := idx.market.ParseAuctionCancelled(lg)
	if err != nil {
		return false, fmt.Errorf("failed to parse AuctionCancelled: %w", err)
	}
	inserted := false
	err = idx.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var insertProcessedLogErr error
		inserted, insertProcessedLogErr = idx.insertProcessedLog(tx, lg, "AuctionCancelled")
		if insertProcessedLogErr != nil {
			return insertProcessedLogErr
		}
		if !inserted {
			return nil
		}
		if !event.AuctionId.IsUint64() {
			return fmt.Errorf("auction id overflow: %s", event.AuctionId.String())
		}

		auctionID := event.AuctionId.Uint64()

		res := tx.Model(&model.Auction{}).
			Where(
				"chain_id = ? AND contract_address = ? AND auction_id = ?",
				idx.chainID,
				normalizeAddress(idx.marketAddress),
				auctionID,
			).Updates(map[string]any{
			"status": model.AuctionStatusCancelled,
		})

		if res.Error != nil {
			return fmt.Errorf("failed to update auction: %w", res.Error)
		}
		if res.RowsAffected == 0 {
			return fmt.Errorf("auction not found when processAuctionCancelled: auction_id=%d", auctionID)
		}
		return nil
	})

	if err != nil {
		return false, err
	}
	if !inserted {
		idx.logger.Debug(
			"AuctionCancelled already processed",
			"auction_id", event.AuctionId.String(),
			"tx_hash", lg.TxHash.Hex(),
			"log_index", lg.Index,
		)
	}

	idx.logger.Info(
		"AuctionCancelled indexed",
		"auction_id", event.AuctionId.String(),
		"block_number", lg.BlockNumber,
		"tx_hash", lg.TxHash.Hex(),
		"log_index", lg.Index,
	)

	return true, nil
}

func (idx *Indexer) insertProcessedLog(tx *gorm.DB, lg types.Log, eventName string) (bool, error) {
	processedLog := model.ProcessedLog{
		ChainID:         idx.chainID,
		ContractAddress: normalizeAddress(idx.marketAddress),
		TxHash:          normalizeHash(lg.TxHash),
		LogIndex:        uint64(lg.Index),
		BlockNumber:     lg.BlockNumber,
		BlockHash:       normalizeHash(lg.BlockHash),
		EventName:       eventName,
	}

	res := tx.Clauses(clause.OnConflict{
		Columns: []clause.Column{
			{Name: "chain_id"},
			{Name: "contract_address"},
			{Name: "tx_hash"},
			{Name: "log_index"},
		},
		DoNothing: true,
	}).Create(&processedLog)

	if res.Error != nil {
		return false, fmt.Errorf("failed to insert processed log: %w", res.Error)
	}

	return res.RowsAffected == 1, nil
}

func (idx *Indexer) updateCursor(ctx context.Context, blockNumber uint64) error {
	header, err := idx.chain.EthClient().BlockByNumber(ctx, new(big.Int).SetUint64(blockNumber))
	if err != nil {
		return fmt.Errorf("failed to get block header by number: %w", err)
	}

	cursor := model.SyncCursor{
		ChainID:                idx.chainID,
		ContractAddress:        normalizeAddress(idx.marketAddress),
		LastProcessedBlock:     blockNumber,
		LastProcessedBlockHash: normalizeHash(header.Hash()),
	}

	return idx.db.WithContext(ctx).Clauses(clause.OnConflict{
		Columns: []clause.Column{
			{Name: "chain_id"},
			{Name: "contract_address"},
		},
		DoUpdates: clause.AssignmentColumns([]string{
			"last_processed_block",
			"last_processed_block_hash",
			"updated_at",
		}),
	}).Create(&cursor).Error
}

func normalizeAddress(addr common.Address) string {
	return strings.ToLower(addr.Hex())
}

func normalizeHash(hash common.Hash) string {
	return strings.ToLower(hash.Hex())
}
