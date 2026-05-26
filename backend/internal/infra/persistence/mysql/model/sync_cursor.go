package model

import "time"

type SyncCursor struct {
	ID uint64 `gorm:"column:id;primaryKey;autoIncrement" json:"id"`

	// 一个 cursor 对应某条链上的某个合约。
	ChainID int64 `gorm:"column:chain_id;not null;uniqueIndex:uk_cursor,priority:1" json:"chain_id"`

	ContractAddress string `gorm:"column:contract_address;type:char(42);not null;uniqueIndex:uk_cursor,priority:2" json:"contract_address"`

	// 已经成功处理完成的最后一个区块号。
	LastProcessedBlock uint64 `gorm:"column:last_processed_block;not null" json:"last_processed_block"`

	// last_processed_block 对应的区块 hash。
	// 后续检测 reorg 时使用。
	LastProcessedBlockHash string `gorm:"column:last_processed_block_hash;type:char(66);not null" json:"last_processed_block_hash"`

	CreatedAt time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at" json:"updated_at"`
}

func (SyncCursor) TableName() string {
	return "sync_cursors"
}
