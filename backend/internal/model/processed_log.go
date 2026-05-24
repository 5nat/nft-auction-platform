package model

import "time"

type ProcessedLog struct {
	ID uint64 `gorm:"primaryKey;autoIncrement" json:"id"`

	ChainID         int64  `gorm:"not null;uniqueIndex:uk_processed_log;index" json:"chain_id"`
	ContractAddress string `gorm:"type:char(42);not null;uniqueIndex:uk_processed_log;index" json:"contract_address"`

	TxHash   string `gorm:"column:tx_hash;type:char(66);not null;uniqueIndex:uk_processed_log" json:"tx_hash"`
	LogIndex uint64 `gorm:"column:log_index;not null;uniqueIndex:uk_processed_log" json:"log_index"`

	BlockNumber uint64 `gorm:"column:block_number;not null;index" json:"block_number"`
	BlockHash   string `gorm:"column:block_hash;type:char(66)" json:"block_hash"`

	EventName string `gorm:"column:event_name;type:varchar(64);not null;index" json:"event_name"`

	CreatedAt time.Time `json:"created_at"`
}

func (ProcessedLog) TableName() string {
	return "processed_logs"
}

/*
CREATE TABLE processed_logs (
  id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,

  chain_id BIGINT NOT NULL,
  contract_address CHAR(42) NOT NULL,

  tx_hash CHAR(66) NOT NULL,
  log_index BIGINT UNSIGNED NOT NULL,

  block_number BIGINT UNSIGNED NOT NULL,
  block_hash CHAR(66),

  event_name VARCHAR(64) NOT NULL,

  created_at DATETIME,

  UNIQUE KEY uk_processed_log (
    chain_id,
    contract_address,
    tx_hash,
    log_index
  ),

  INDEX idx_processed_logs_chain_id (chain_id),
  INDEX idx_processed_logs_contract_address (contract_address),
  INDEX idx_processed_logs_block_number (block_number),
  INDEX idx_processed_logs_event_name (event_name)
);
*/
