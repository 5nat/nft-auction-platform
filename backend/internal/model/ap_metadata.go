package model

import "time"

type AppMetadata struct {
	MetaKey   string    `gorm:"column:meta_key;primaryKey;type:varchar(128)" json:"key"`
	Value     string    `gorm:"type:varchar(255);not null" json:"value"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (AppMetadata) TableName() string {
	return "app_metadata"
}

// 整体设计原则
//
//auctions：拍卖当前状态表
//bids：出价历史明细表
//processed_logs：事件幂等表
//sync_cursors：同步进度表

//1. chain_id + contract_address 表示数据来源。
//2. tx_hash + log_index 表示一条链上 log 的唯一位置。
//3. block_number + block_hash 为排序、排查和 reorg 预留。
//4. uint256 金额统一用 string 保存，数据库用 varchar(78)。
//5. address 统一用 char(42)，hash 统一用 char(66)。
//6. bids / processed_logs 是事件级数据，基本不可变。
//7. auctions 是 read model 状态表，会被事件持续更新。
