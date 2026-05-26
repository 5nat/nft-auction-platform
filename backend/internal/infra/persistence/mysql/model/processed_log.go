package model

import "time"

type ProcessedLog struct {
	ID uint64 `gorm:"column:id;primaryKey;autoIncrement" json:"id"`

	// 数据来源：哪条链、哪个合约。
	ChainID int64 `gorm:"column:chain_id;not null;uniqueIndex:uk_processed_log,priority:1;index:idx_processed_logs_chain_id" json:"chain_id"`

	ContractAddress string `gorm:"column:contract_address;type:char(42);not null;uniqueIndex:uk_processed_log,priority:2;index:idx_processed_logs_contract_address" json:"contract_address"`

	// 已处理 log 所在交易。
	TxHash string `gorm:"column:tx_hash;type:char(66);not null;uniqueIndex:uk_processed_log,priority:3" json:"tx_hash"`

	// 已处理 log 在交易 receipt 中的位置。
	LogIndex uint64 `gorm:"column:log_index;not null;uniqueIndex:uk_processed_log,priority:4" json:"log_index"`

	// 已处理 log 所在区块号。
	BlockNumber uint64 `gorm:"column:block_number;not null;index:idx_processed_logs_block_number" json:"block_number"`

	// 已处理 log 所在区块 hash。
	BlockHash string `gorm:"column:block_hash;type:char(66);not null" json:"block_hash"`

	// 事件名称，例如：
	// AuctionCreated、BidPlaced、AuctionEnded、AuctionCancelled。
	EventName string `gorm:"column:event_name;type:varchar(64);not null;index:idx_processed_logs_event_name" json:"event_name"`

	// Indexer 成功处理该 log 并写入数据库的时间。
	CreatedAt time.Time `gorm:"column:created_at" json:"created_at"`
}

func (ProcessedLog) TableName() string {
	return "processed_logs"
}
