package model

import "time"

type Bid struct {
	ID uint64 `gorm:"primary_key;auto_increment" json:"id"`

	ChainID         int64  `gorm:"not null;uniqueIndex:uk_bid_log;index" json:"chain_id"`
	ContractAddress string `gorm:"type:char(42);not null;uniqueIndex:uk_bid_log;index" json:"contract_address"`

	AuctionID uint64 `gorm:"not null;index" json:"auction_id"`
	Bidder    string `gorm:"type:char(42);not null;index" json:"bidder"`

	BidToken  string `gorm:"column:bid_token;type:char(42);not null;index" json:"bid_token"`
	Amount    string `gorm:"type:varchar(78);not null" json:"amount"`
	AmountUSD string `gorm:"column:amount_usd;type:varchar(78);not null" json:"amount_usd"`

	TxHash      string `gorm:"column:tx_hash;type:char(66);not null;uniqueIndex:uk_bid_log" json:"tx_hash"`
	LogIndex    uint64 `gorm:"column:log_index;not null;uniqueIndex:uk_bid_log" json:"log_index"`
	BlockNumber uint64 `gorm:"column:block_number;not null;index" json:"block_number"`
	BlockHash   string `gorm:"column:block_hash;type:char(66)" json:"block_hash"`

	CreatedAt time.Time `json:"created_at"`
}

func (Bid) TableName() string {
	return "bids"
}
