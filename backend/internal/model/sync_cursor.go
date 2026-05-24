package model

import "time"

type SyncCursor struct {
	ID uint64 `gorm:"primaryKey;autoIncrement" json:"id"`

	ChainID         int64  `gorm:"not null;uniqueIndex:uk_cursor" json:"chain_id"`
	ContractAddress string `gorm:"type:char(42);not null;uniqueIndex:uk_cursor" json:"contract_address"`

	LastProcessedBlock     uint64 `gorm:"column:last_processed_block;not null" json:"last_processed_block"`
	LastProcessedBlockHash string `gorm:"column:last_processed_block_hash;type:char(66)" json:"last_processed_block_hash"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (SyncCursor) TableName() string {
	return "sync_cursors"
}

/*
CREATE TABLE sync_cursors (
  id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,

  chain_id BIGINT NOT NULL,
  contract_address CHAR(42) NOT NULL,

  last_processed_block BIGINT UNSIGNED NOT NULL,
  last_processed_block_hash CHAR(66),

  created_at DATETIME,
  updated_at DATETIME,

  UNIQUE KEY uk_cursor (
    chain_id,
    contract_address
  )
);
*/
