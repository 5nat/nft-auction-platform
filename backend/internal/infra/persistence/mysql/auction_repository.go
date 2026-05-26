package mysql

import (
	"context"
	"errors"
	"fmt"

	"github.com/5nat/nft-auction-platform/backend/internal/infra/persistence/mysql/model"
	"github.com/5nat/nft-auction-platform/backend/internal/modules/auction"
	"gorm.io/gorm"
)

type AuctionRepository struct {
	db *gorm.DB
}

// 编译期检查：确保 AuctionRepository 实现了 auction.Repository 接口。
var _ auction.Repository = (*AuctionRepository)(nil)

func NewAuctionRepository(db *gorm.DB) *AuctionRepository {
	return &AuctionRepository{db: db}
}

func (r *AuctionRepository) ListAuctions(ctx context.Context, query auction.ListAuctionsQuery) ([]auction.Auction, int64, error) {
	db := r.db.WithContext(ctx).Model(&model.Auction{})
	db = applyAuctionFilters(db, query)

	var total int64
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("count auctions: %w", err)
	}

	offset := (query.Page - 1) * query.PageSize

	var rows []model.Auction
	if err := db.
		Scopes(applyAuctionSort(query.Sort)).
		Offset(offset).
		Limit(query.PageSize).
		Find(&rows).Error; err != nil {
		return nil, 0, fmt.Errorf("list auctions: %w", err)
	}

	return toAuctionEntities(rows), total, nil
}

func (r *AuctionRepository) GetAuction(ctx context.Context, query auction.GetAuctionQuery) (*auction.Auction, error) {
	db := r.db.WithContext(ctx).Model(&model.Auction{})

	if query.ChainID != 0 {
		db = db.Where("chain_id = ?", query.ChainID)
	}

	if query.ContractAddress != "" {
		db = db.Where("contract_address = ?", query.ContractAddress)
	}

	var row model.Auction
	err := db.
		Where("auction_id = ?", query.AuctionID).
		First(&row).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, auction.ErrAuctionNotFound
		}
		return nil, fmt.Errorf("get auction: %w", err)
	}

	entity := toAuctionEntity(row)

	return &entity, nil
}

func (r *AuctionRepository) ListBids(ctx context.Context, query auction.ListBidsQuery) ([]auction.Bid, int64, error) {
	db := r.db.WithContext(ctx).Model(&model.Bid{})

	if query.ChainID != 0 {
		db = db.Where("chain_id = ?", query.ChainID)
	}

	if query.ContractAddress != "" {
		db = db.Where("contract_address = ?", query.ContractAddress)
	}

	db = db.Where("auction_id = ?", query.AuctionID)

	var total int64
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("count bids: %w", err)
	}

	offset := (query.Page - 1) * query.PageSize

	var rows []model.Bid
	if err := db.
		Order("block_number ASC").
		Order("log_index ASC").
		Limit(query.PageSize).
		Offset(offset).
		Find(&rows).Error; err != nil {
		return nil, 0, fmt.Errorf("list bids: %w", err)
	}

	return toBidEntities(rows), total, nil
}

func applyAuctionFilters(db *gorm.DB, query auction.ListAuctionsQuery) *gorm.DB {
	if query.ChainID != 0 {
		db = db.Where("chain_id = ?", query.ChainID)
	}

	if query.ContractAddress != "" {
		db = db.Where("contract_address = ?", query.ContractAddress)
	}

	if query.Status != "" {
		db = db.Where("status = ?", query.Status)
	}

	if query.Seller != "" {
		db = db.Where("seller = ?", query.Seller)
	}

	if query.NFT != "" {
		db = db.Where("nft_contract = ?", query.NFT)
	}

	return db
}

func applyAuctionSort(sort string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		switch sort {
		case "end_time_asc":
			return db.
				Order("end_time ASC").
				Order("created_block_number DESC").
				Order("created_log_index DESC")
		case "last_event_desc":
			return db.
				Order("last_event_block_number DESC").
				Order("last_event_log_index DESC")
		case "created_desc":
			fallthrough
		default:
			return db.
				Order("created_block_number DESC").
				Order("created_log_index DESC")
		}
	}
}

func toAuctionEntities(rows []model.Auction) []auction.Auction {
	items := make([]auction.Auction, 0, len(rows))

	for _, row := range rows {
		items = append(items, toAuctionEntity(row))
	}

	return items
}

func toAuctionEntity(row model.Auction) auction.Auction {
	return auction.Auction{
		ChainID:         row.ChainID,
		ContractAddress: row.ContractAddress,
		AuctionID:       row.AuctionID,

		Seller:      row.Seller,
		NFTContract: row.NFTContract,
		TokenID:     row.TokenID,

		MinBidUSD: row.MinBidUSD,

		HighestBidder:    row.HighestBidder,
		HighestBidToken:  row.HighestBidToken,
		HighestBidAmount: row.HighestBidAmount,
		HighestBidUSD:    row.HighestBidUSD,

		Status:  row.Status,
		EndTime: row.EndTime,

		CreatedTxHash:      row.CreatedTxHash,
		CreatedBlockNumber: row.CreatedBlockNumber,
		CreatedBlockHash:   row.CreatedBlockHash,
		CreatedLogIndex:    row.CreatedLogIndex,

		LastEventName:        row.LastEventName,
		LastEventTxHash:      row.LastEventTxHash,
		LastEventBlockNumber: row.LastEventBlockNumber,
		LastEventBlockHash:   row.LastEventBlockHash,
		LastEventLogIndex:    row.LastEventLogIndex,

		CreatedAt: row.CreatedAt,
		UpdatedAt: row.UpdatedAt,
	}
}

func toBidEntities(rows []model.Bid) []auction.Bid {
	items := make([]auction.Bid, 0, len(rows))

	for _, row := range rows {
		items = append(items, toBidEntity(row))
	}

	return items
}

func toBidEntity(row model.Bid) auction.Bid {
	return auction.Bid{
		ChainID:         row.ChainID,
		ContractAddress: row.ContractAddress,
		AuctionID:       row.AuctionID,

		Bidder:    row.Bidder,
		BidToken:  row.BidToken,
		Amount:    row.Amount,
		AmountUSD: row.AmountUSD,

		TxHash:      row.TxHash,
		LogIndex:    row.LogIndex,
		BlockNumber: row.BlockNumber,
		BlockHash:   row.BlockHash,

		CreatedAt: row.CreatedAt,
	}
}
