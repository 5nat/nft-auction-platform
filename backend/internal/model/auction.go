package model

import "time"

const (
	AuctionStatusActive    = "active"
	AuctionStatusEnded     = "ended"
	AuctionStatusCancelled = "cancelled"
)

type Auction struct {
	ID uint64 `gorm:"primaryKey;autoIncrement" json:"id"`

	ChainID         int64  `gorm:"not null;uniqueIndex:uk_auction_chain_auction;index:idx_auctions_chain_status" json:"chain_id"`
	ContractAddress string `gorm:"type:char(42);not null;uniqueIndex:uk_auction_chain_auction" json:"contract_address"`
	AuctionID       uint64 `gorm:"not null;uniqueIndex:uk_auction_chain_auction" json:"auction_id"`

	Seller      string `gorm:"type:char(42);not null;index" json:"seller"`
	NFTContract string `gorm:"column:nft_contract;type:char(42);not null;index" json:"nft_contract"`
	TokenID     string `gorm:" column:token_id;type:varchar(78);not null" json:"token_id"`

	MinBidUSD string `gorm:"column:min_bid_usd;type:varchar(78);not null" json:"min_bid_usd"`

	HighestBidder    string `gorm:"type:char(42)" json:"highest_bidder"`
	HighestBidToken  string `gorm:"column:highest_bid_token;type:char(42)" json:"highest_bid_token"`
	HighestBidAmount string `gorm:"column:highest_bid_amount;type:varchar(78)" json:"highest_bid_amount"`
	HighestBidUSD    string `gorm:"column:highest_bid_amount;type:varchar(78)" json:"highest_bid_usd"`

	EndTime uint64 `gorm:"not null;index" json:"end_time"`
	Status  string `gorm:"type:varchar(32);not null;index:idx_auctions_chain_status" json:"status"`

	CreatedTxHash      string `gorm:"column:created_tx_hash;type:char(66)" json:"created_tx_hash"`
	CreatedBlockNumber uint64 `gorm:"column:created_block_number;index" json:"created_block_number"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (Auction) TableName() string {
	return "auctions"
}
