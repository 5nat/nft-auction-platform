package model

import "time"

const (
	AuctionStatusActive    = "active"
	AuctionStatusEnded     = "ended"
	AuctionStatusCancelled = "cancelled"
)

type Auction struct {
	ID uint64 `gorm:"primaryKey;autoIncrement" json:"id"`

	ChainID         int64  `gorm:"not null;uniqueIndex:uk_auction_chain_contract_id;index:idx_auctions_chain_status_end" json:"chain_id"`
	ContractAddress string `gorm:"type:char(42);not null;uniqueIndex:uk_auction_chain_contract_id" json:"contract_address"`
	AuctionID       uint64 `gorm:"not null;uniqueIndex:uk_auction_chain_contract_id" json:"auction_id"`

	Seller      string `gorm:"type:char(42);not null;index" json:"seller"`
	NFTContract string `gorm:"column:nft_contract;type:char(42);not null;index" json:"nft_contract"`
	TokenID     string `gorm:"column:token_id;type:varchar(78);not null" json:"token_id"`

	MinBidUSD string `gorm:"column:min_bid_usd;type:varchar(78);not null" json:"min_bid_usd"`

	HighestBidder    string `gorm:"column:highest_bidder;type:char(42)" json:"highest_bidder"`
	HighestBidToken  string `gorm:"column:highest_bid_token;type:char(42)" json:"highest_bid_token"`
	HighestBidAmount string `gorm:"column:highest_bid_amount;type:varchar(78)" json:"highest_bid_amount"`
	HighestBidUSD    string `gorm:"column:highest_bid_usd;type:varchar(78)" json:"highest_bid_usd"`

	EndTime uint64 `gorm:"column:end_time;not null;index:idx_auctions_chain_status_end" json:"end_time"`
	Status  string `gorm:"type:varchar(32);not null;index:idx_auctions_chain_status_end" json:"status"`

	CreatedTxHash      string `gorm:"column:created_tx_hash;type:char(66)" json:"created_tx_hash"`
	CreatedBlockNumber uint64 `gorm:"column:created_block_number;index" json:"created_block_number"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (Auction) TableName() string {
	return "auctions"
}

/*
CREATE TABLE auctions (
  id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,

  chain_id BIGINT NOT NULL,
  contract_address CHAR(42) NOT NULL,
  auction_id BIGINT UNSIGNED NOT NULL,

  seller CHAR(42) NOT NULL,
  nft_contract CHAR(42) NOT NULL,
  token_id VARCHAR(78) NOT NULL,

  min_bid_usd VARCHAR(78) NOT NULL,

  highest_bidder CHAR(42),
  highest_bid_token CHAR(42),
  highest_bid_amount VARCHAR(78),
  highest_bid_usd VARCHAR(78),

  end_time BIGINT UNSIGNED NOT NULL,
  status VARCHAR(32) NOT NULL,

  created_tx_hash CHAR(66),
  created_block_number BIGINT UNSIGNED,

  created_at DATETIME,
  updated_at DATETIME,

  UNIQUE KEY uk_auction_chain_contract_id (
    chain_id,
    contract_address,
    auction_id
  ),

  INDEX idx_auctions_chain_status_end (
    chain_id,
    status,
    end_time
  ),

  INDEX idx_auctions_seller (seller),
  INDEX idx_auctions_nft_contract (nft_contract),
  INDEX idx_auctions_created_block_number (created_block_number)
);


*/
