package auction

import (
	"context"
	"errors"
	"fmt"

	"github.com/5nat/nft-auction-platform/backend/internal/infra/persistence/mysql/model"
	"gorm.io/gorm"
)

type Repository interface {
	ListAuctions(ctx context.Context, query ListAuctionsQuery) ([]model.Auction, int64, error)
	GetAuction(ctx context.Context, query GetAuctionQuery) (*model.Auction, error)
	ListBids(ctx context.Context, query ListBidsQuery) ([]model.Bid, int64, error)
}

type GormRepository struct {
	db *gorm.DB
}

func NewGormRepository(db *gorm.DB) *GormRepository {
	return &GormRepository{db: db}
}

func (r *GormRepository) ListAuctions(ctx context.Context, query ListAuctionsQuery) ([]model.Auction, int64, error) {
	db := r.db.WithContext(ctx).Model(&model.Auction{})
	db = applyAuctionFilters(db, query)

	var total int64
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("count auctions: %w", err)
	}

	offset := (query.Page - 1) * query.PageSize

	var auctions []model.Auction
	if err := db.
		Scopes(applyAuctionSort(query.Sort)).
		Offset(offset).
		Limit(query.PageSize).
		Find(&auctions).Error; err != nil {
		return nil, 0, fmt.Errorf("list auctions: %w", err)
	}

	return auctions, total, nil
}

func (r *GormRepository) GetAuction(ctx context.Context, query GetAuctionQuery) (*model.Auction, error) {
	db := r.db.WithContext(ctx).Model(&model.Auction{})

	if query.ChainID != 0 {
		db = db.Where("chain_id = ?", query.ChainID)
	}

	if query.ContractAddress != "" {
		db = db.Where("contract_address = ?", query.ContractAddress)
	}

	var auction model.Auction
	err := db.
		Where("auction_id = ?", query.AuctionID).
		First(&auction).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrAuctionNotFound
		}
		return nil, fmt.Errorf("get auction: %w", err)
	}

	return &auction, nil
}

func (r *GormRepository) ListBids(ctx context.Context, query ListBidsQuery) ([]model.Bid, int64, error) {
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

	var bids []model.Bid
	if err := db.
		Order("block_number ASC").
		Order("log_index ASC").
		Limit(query.PageSize).
		Offset(offset).
		Find(&bids).Error; err != nil {
		return nil, 0, fmt.Errorf("list bids: %w", err)
	}

	return bids, total, nil
}

func applyAuctionFilters(db *gorm.DB, query ListAuctionsQuery) *gorm.DB {
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
