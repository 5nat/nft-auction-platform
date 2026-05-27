package model

import "time"

// Wallet 是用户绑定的钱包表。
// users 与 wallets 是一对多关系：一个平台用户可以绑定多个钱包。
// 不把 address 直接放到 users 表，是为了支持多钱包、多链、主钱包、解绑和换绑等后续能力。
type Wallet struct {
	// ID 是钱包绑定记录的主键。
	// 后续 JWT 和 CurrentUser 中可以携带 WalletID，
	// 以明确当前登录态对应的是哪一个钱包。
	ID uint64 `gorm:"column:id;primaryKey;autoIncrement"`

	// UserID 关联 users.id。
	// 这里先不强制添加数据库外键，保持和当前项目其他 read model / indexer 表风格一致，
	// 由 Repository 和 Service 负责业务一致性。
	UserID uint64 `gorm:"column:user_id;not null;index:idx_wallet_user_id"`

	// Address 是钱包地址。
	// Ethereum 地址标准展示形式为 0x + 40 位十六进制字符，因此使用 char(42)。
	// Service 层应统一规范化地址，例如统一转小写或校验 checksum。
	Address string `gorm:"column:address;type:char(42);not null;uniqueIndex:uk_wallet_address_chain,priority:1"`

	// ChainID 表示钱包所属链。
	// address + chain_id 联合唯一，避免同一条链上的同一钱包被绑定到多个用户。
	// 这样也为后续多链扩展预留空间。
	ChainID int64 `gorm:"column:chain_id;not null;uniqueIndex:uk_wallet_address_chain,priority:2"`

	// IsPrimary 表示是否为用户主钱包。
	// 当前一个用户可能只有一个钱包，但后续支持多钱包时，
	// 可以用该字段标识默认展示或默认交易钱包。
	IsPrimary bool `gorm:"column:is_primary;not null;default:true"`

	// CreatedAt 表示钱包绑定时间。
	CreatedAt time.Time `gorm:"column:created_at;not null"`

	// UpdatedAt 表示钱包绑定记录最近更新时间。
	UpdatedAt time.Time `gorm:"column:updated_at;not null"`
}

func (Wallet) TableName() string {
	return "wallets"
}
