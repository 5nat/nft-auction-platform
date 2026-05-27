package model

import "time"

// User 是平台内部用户表。
// 在 Web3 场景中，用户身份主要由钱包证明，
// 因此 users 表只保存平台用户的基础信息，钱包信息通过 wallets 表关联。
// 后续如果需要昵称、头像、邮箱、简介等资料，可以继续扩展 users 表，
// 不影响现有 wallet 登录模型。
type User struct {
	// ID 是平台内部用户 ID。
	// 不直接使用钱包地址作为用户 ID，是为了支持一个用户绑定多个钱包、
	// 后续扩展用户资料、风控状态、通知设置和交易审计等能力。
	ID uint64 `gorm:"column:id;primaryKey;autoIncrement"`

	// Status 表示用户状态。
	// 当前可以使用 active，后续可扩展 disabled、banned、frozen 等状态。
	// 这样即使用户被禁用，也可以保留其历史拍卖、出价和交易意图记录。
	Status string `gorm:"column:status;type:varchar(32);not null;default:active"`

	// CreatedAt 表示用户创建时间。
	// 用于基础审计和后续用户生命周期分析。
	CreatedAt time.Time `gorm:"column:created_at;not null"`

	// UpdatedAt 表示用户最近更新时间。
	// GORM 会在更新记录时自动维护该字段。
	UpdatedAt time.Time `gorm:"column:updated_at;not null"`
}

func (User) TableName() string {
	return "users"
}
