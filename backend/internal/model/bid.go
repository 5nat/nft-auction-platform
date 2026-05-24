package model

import "time"

type Bid struct {
	ID uint64 `gorm:"primaryKey;autoIncrement" json:"id"`

	ChainID         int64  `gorm:"not null;uniqueIndex:uk_bid_log;index:idx_bids_auction_order" json:"chain_id"`
	ContractAddress string `gorm:"type:char(42);not null;uniqueIndex:uk_bid_log;index:idx_bids_auction_order" json:"contract_address"`

	AuctionID uint64 `gorm:"not null;index:idx_bids_auction_order" json:"auction_id"`
	Bidder    string `gorm:"type:char(42);not null;index" json:"bidder"`

	BidToken  string `gorm:"column:bid_token;type:char(42);not null;index" json:"bid_token"`
	Amount    string `gorm:"type:varchar(78);not null" json:"amount"`
	AmountUSD string `gorm:"column:amount_usd;type:varchar(78);not null" json:"amount_usd"`

	TxHash      string `gorm:"column:tx_hash;type:char(66);not null;uniqueIndex:uk_bid_log" json:"tx_hash"`
	LogIndex    uint64 `gorm:"column:log_index;not null;uniqueIndex:uk_bid_log;index:idx_bids_auction_order" json:"log_index"`
	BlockNumber uint64 `gorm:"column:block_number;not null;index;index:idx_bids_auction_order" json:"block_number"`
	BlockHash   string `gorm:"column:block_hash;type:char(66)" json:"block_hash"`

	CreatedAt time.Time `json:"created_at"`
}

func (Bid) TableName() string {
	return "bids"
}

/*
CREATE TABLE bids (
  id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,

  chain_id BIGINT NOT NULL,
  contract_address CHAR(42) NOT NULL,

  auction_id BIGINT UNSIGNED NOT NULL,
  bidder CHAR(42) NOT NULL,

  bid_token CHAR(42) NOT NULL,
  amount VARCHAR(78) NOT NULL,
  amount_usd VARCHAR(78) NOT NULL,

  tx_hash CHAR(66) NOT NULL,
  log_index BIGINT UNSIGNED NOT NULL,
  block_number BIGINT UNSIGNED NOT NULL,
  block_hash CHAR(66),

  created_at DATETIME,

  UNIQUE KEY uk_bid_log (
    chain_id,
    contract_address,
    tx_hash,
    log_index
  ),

  INDEX idx_bids_auction_order (
    chain_id,
    contract_address,
    auction_id,
    block_number,
    log_index
  ),

  INDEX idx_bids_bidder (bidder),
  INDEX idx_bids_bid_token (bid_token),
  INDEX idx_bids_block_number (block_number)
);
*/
