package model

import "time"

type Bid struct {
	ID uint64 `gorm:"column:id;primaryKey;autoIncrement" json:"id"`

	// 数据来源：哪条链、哪个拍卖合约。
	ChainID int64 `gorm:"column:chain_id;not null;uniqueIndex:uk_bid_log,priority:1;index:idx_bids_auction_order,priority:1" json:"chain_id"`

	ContractAddress string `gorm:"column:contract_address;type:char(42);not null;uniqueIndex:uk_bid_log,priority:2;index:idx_bids_auction_order,priority:2" json:"contract_address"`

	// 业务归属：这笔出价属于哪个拍卖。
	AuctionID uint64 `gorm:"column:auction_id;not null;index:idx_bids_auction_order,priority:3" json:"auction_id"`

	// 出价人地址。
	Bidder string `gorm:"column:bidder;type:char(42);not null;index:idx_bids_bidder" json:"bidder"`

	// 出价资产。
	// 如果 bid_token 是零地址，表示使用 ETH 出价。
	// 如果不是零地址，表示使用 ERC20 出价。
	BidToken string `gorm:"column:bid_token;type:char(42);not null;index:idx_bids_bid_token" json:"bid_token"`

	// 原始出价金额，按链上最小单位保存。
	// ETH 是 wei，ERC20 是 token 最小单位。
	Amount string `gorm:"column:amount;type:varchar(78);not null" json:"amount"`

	// 折算后的 USD 金额，通常是 18 位精度。
	AmountUSD string `gorm:"column:amount_usd;type:varchar(78);not null" json:"amount_usd"`

	// 这条 BidPlaced 事件所在交易。
	TxHash string `gorm:"column:tx_hash;type:char(66);not null;uniqueIndex:uk_bid_log,priority:3" json:"tx_hash"`

	// 这条 log 在交易 receipt 中的位置。
	// 同一个 tx 里可能有多条 log，所以只用 tx_hash 不够，必须加 log_index。
	LogIndex uint64 `gorm:"column:log_index;not null;uniqueIndex:uk_bid_log,priority:4;index:idx_bids_auction_order,priority:5" json:"log_index"`

	// 这条事件所在区块号。
	BlockNumber uint64 `gorm:"column:block_number;not null;index:idx_bids_block_number;index:idx_bids_auction_order,priority:4" json:"block_number"`

	// 这条事件所在区块 hash。
	// 后续做 reorg 检测时很重要。
	BlockHash string `gorm:"column:block_hash;type:char(66);not null" json:"block_hash"`

	// 数据写入数据库的时间。
	// 注意：这不是链上出价时间，只是后端索引到这条事件的时间。
	CreatedAt time.Time `gorm:"column:created_at" json:"created_at"`
}

func (Bid) TableName() string {
	return "bids"
}
