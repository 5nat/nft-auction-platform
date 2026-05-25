package indexer

import (
	"context"
	"errors"
	"fmt"

	"github.com/5nat/nft-auction-platform/backend/internal/model"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Repository struct {
	db              *gorm.DB
	chainID         int64
	contractAddress string
}

// 如果 Repository 没有完整实现 EventRepository， 编译期直接报错。
var _ EventRepository = (*Repository)(nil)

func NewRepository(db *gorm.DB, chainID int64, contractAddress common.Address) *Repository {
	return &Repository{
		db:              db,
		chainID:         chainID,
		contractAddress: normalizeAddress(contractAddress),
	}
}

// WithTx 统一封装数据库事务。
//
// handler 不再直接调用 idx.db.Transaction。
// 这样后面如果要换事务策略、加 trace、加日志，都可以集中在这里处理。
func (r *Repository) WithTx(ctx context.Context, fn func(repo EventRepository) error) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		txRepo := &Repository{
			db:              tx,
			chainID:         r.chainID,
			contractAddress: r.contractAddress,
		}
		return fn(txRepo)
	})
}

// InsertProcessedLog 插入 processed_logs，作为事件幂等闸门。
//
// 返回值含义：
// true  ：这条 log 第一次处理。
// false ：这条 log 已经处理过，本次应该跳过。
func (r *Repository) InsertProcessedLog(ctx context.Context, lg types.Log, eventName string) (bool, error) {
	processedLog := model.ProcessedLog{
		ChainID:         r.chainID,
		ContractAddress: r.contractAddress,
		TxHash:          normalizeHash(lg.TxHash),
		LogIndex:        uint64(lg.Index),
		BlockNumber:     lg.BlockNumber,
		BlockHash:       normalizeHash(lg.BlockHash),
		EventName:       eventName,
	}

	res := r.db.WithContext(ctx).Clauses(clause.OnConflict{
		Columns: []clause.Column{
			{Name: "chain_id"},
			{Name: "contract_address"},
			{Name: "tx_hash"},
			{Name: "log_index"},
		},
		DoNothing: true,
	}).Create(&processedLog)

	if res.Error != nil {
		return false, fmt.Errorf("insert processed log: %w", res.Error)
	}

	return res.RowsAffected == 1, nil
}

func (r *Repository) CreateAuction(ctx context.Context, auction model.Auction) error {
	res := r.db.WithContext(ctx).Clauses(clause.OnConflict{
		Columns: []clause.Column{
			{Name: "chain_id"},
			{Name: "contract_address"},
			{Name: "auction_id"},
		},
		DoNothing: true,
	}).Create(&auction)

	if res.Error != nil {
		return fmt.Errorf("create auction: %w", res.Error)
	}

	return nil
}

func (r *Repository) CreateBid(ctx context.Context, bid model.Bid) error {
	res := r.db.WithContext(ctx).Clauses(clause.OnConflict{
		Columns: []clause.Column{
			{Name: "chain_id"},
			{Name: "contract_address"},
			{Name: "tx_hash"},
			{Name: "log_index"},
		},
		DoNothing: true,
	}).Create(&bid)

	if res.Error != nil {
		return fmt.Errorf("create bid: %w", res.Error)
	}

	return nil
}

type UpdateAuctionHighestBidInput struct {
	AuctionID uint64
	Bidder    string
	BidToken  string
	Amount    string
	AmountUSD string
}

func (r *Repository) UpdateAuctionHighestBid(ctx context.Context, input UpdateAuctionHighestBidInput) error {
	res := r.db.WithContext(ctx).
		Model(&model.Auction{}).
		Where(
			"chain_id = ? AND contract_address = ? AND auction_id = ?",
			r.chainID,
			r.contractAddress,
			input.AuctionID,
		).
		Updates(map[string]any{
			"highest_bidder":     input.Bidder,
			"highest_bid_token":  input.BidToken,
			"highest_bid_amount": input.Amount,
			"highest_bid_usd":    input.AmountUSD,
		})

	if res.Error != nil {
		return fmt.Errorf("update auction highest bid: %w", res.Error)
	}

	if res.RowsAffected == 0 {
		return fmt.Errorf("auction not found when updating highest bid: auction_id=%d", input.AuctionID)
	}

	return nil
}

type MarkAuctionEndedInput struct {
	AuctionID uint64
	Winner    string
	BidToken  string
	Amount    string
	AmountUSD string
}

func (r *Repository) MarkAuctionEnded(ctx context.Context, input MarkAuctionEndedInput) error {
	res := r.db.WithContext(ctx).
		Model(&model.Auction{}).
		Where(
			"chain_id = ? AND contract_address = ? AND auction_id = ?",
			r.chainID,
			r.contractAddress,
			input.AuctionID,
		).Updates(map[string]any{
		"status":             model.AuctionStatusEnded,
		"highest_bidder":     input.Winner,
		"highest_bid_token":  input.BidToken,
		"highest_bid_amount": input.Amount,
		"highest_bid_usd":    input.AmountUSD,
	})

	if res.Error != nil {
		return fmt.Errorf("mark auction ended: %w", res.Error)
	}

	if res.RowsAffected == 0 {
		return fmt.Errorf("auction not found when marking ended: auction_id=%d", input.AuctionID)
	}

	return nil
}

type MarkAuctionCancelledInput struct {
	AuctionID uint64
}

func (r *Repository) MarkAuctionCancelled(ctx context.Context, input MarkAuctionCancelledInput) error {
	res := r.db.WithContext(ctx).
		Model(&model.Auction{}).
		Where(
			"chain_id = ? AND contract_address = ? AND auction_id = ?",
			r.chainID,
			r.contractAddress,
			input.AuctionID,
		).
		Updates(map[string]any{
			"status": model.AuctionStatusCancelled,
		})

	if res.Error != nil {
		return fmt.Errorf("mark auction cancelled: %w", res.Error)
	}

	if res.RowsAffected == 0 {
		return fmt.Errorf("auction not found when marking cancelled: auction_id=%d", input.AuctionID)
	}

	return nil
}

func (r *Repository) NextFromBlock(ctx context.Context, startBlock uint64) (uint64, error) {
	var cursor model.SyncCursor

	err := r.db.WithContext(ctx).
		Where(
			"chain_id = ? AND contract_address = ?",
			r.chainID,
			r.contractAddress,
		).First(&cursor).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return startBlock, nil
		}

		return 0, fmt.Errorf("get sync cursor: %w", err)
	}

	return cursor.LastProcessedBlock + 1, nil
}

func (r *Repository) UpsertCursor(ctx context.Context, blockNumber uint64, blockHash string) error {
	cursor := model.SyncCursor{
		ChainID:                r.chainID,
		ContractAddress:        r.contractAddress,
		LastProcessedBlock:     blockNumber,
		LastProcessedBlockHash: blockHash,
	}

	res := r.db.WithContext(ctx).Clauses(clause.OnConflict{
		Columns: []clause.Column{
			{Name: "chain_id"},
			{Name: "contract_address"},
		},
		DoUpdates: clause.AssignmentColumns([]string{
			"last_processed_block",
			"last_processed_block_hash",
			"updated_at",
		}),
	}).Create(&cursor)

	if res.Error != nil {
		return fmt.Errorf("upsert sync cursor: %w", res.Error)
	}

	return nil
}
