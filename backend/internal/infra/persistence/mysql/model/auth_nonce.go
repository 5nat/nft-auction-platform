package model

import "time"

// AuthNonce 保存一次钱包登录挑战。
// Web3 登录不能直接相信前端传来的 address，
// 因此服务端需要生成 nonce 和完整 SIWE 风格 message，
// 前端钱包对 message 签名后，后端再通过签名恢复地址来证明用户确实控制该钱包。
type AuthNonce struct {
	// ID 是登录挑战记录的主键。
	ID uint64 `gorm:"column:id;primaryKey;autoIncrement"`

	// Address 是本次请求登录的钱包地址。
	// 这里不使用 wallet_id，因为用户第一次登录时，wallets 表中可能还没有这条记录。
	// Ethereum 地址长度固定为 42，因此使用 char(42)。
	Address string `gorm:"column:address;type:char(42);not null;index:idx_auth_nonce_address_chain,priority:1"`

	// ChainID 表示本次登录挑战所属链。
	// SIWE message 中也包含 Chain ID，用于避免跨链或错误网络下的登录混淆。
	ChainID uint64 `gorm:"column:chain_id;not null;index:idx_auth_nonce_address_chain,priority:2"`

	// Nonce 是服务端生成的一次性随机字符串。
	// 它用于防止重放攻击：同一份 message + signature 只能成功登录一次。
	// varchar(64) 足够保存常见的 hex / base64url 随机 nonce。
	Nonce string `gorm:"column:nonce;type:varchar(64);not null;uniqueIndex:uk_auth_nonce_nonce"`

	// MessageHash 是 sha256(message) 的十六进制字符串，不带 0x，固定长度为 64。
	// message 本身是 TEXT，不适合直接做唯一索引和高效查询，
	// 因此 Verify 阶段通过 message_hash 找回对应登录挑战。
	MessageHash string `gorm:"column:message_hash;type:char(64);not null;uniqueIndex:uk_auth_nonce_message_hash"`

	// Message 是完整的 SIWE 风格登录消息原文。
	// 钱包实际签名的是这份完整 message，保存原文便于后续验签、审计和排查问题。
	Message string `gorm:"column:message;type:text;not null"`

	// Domain 表示发起登录的网站域名，例如 localhost:8080 或 nft-auction.example.com。
	// SIWE message 中包含 domain，可以降低签名被跨站复用的风险。
	Domain string `gorm:"column:domain;type:varchar(255);not null"`

	// URI 表示本次登录对应的站点 URI，例如 http://localhost:8080。
	// 它和 Domain 一起用于描述签名适用的站点上下文。
	URI string `gorm:"column:uri;type:varchar(512);not null"`

	// IssuedAt 表示登录挑战生成时间。
	// 它会写入 SIWE message，用于说明这条签名请求是什么时候发起的。
	IssuedAt time.Time `gorm:"column:issued_at;not null"`

	// ExpiresAt 表示登录挑战过期时间。
	// Verify 阶段必须检查该字段，过期 nonce 不能再用于登录。
	// 该字段加索引，方便后续清理过期 nonce。
	ExpiresAt time.Time `gorm:"column:expires_at;not null;index:idx_auth_nonce_expires_at"`

	// UsedAt 表示 nonce 被成功消费的时间。
	// NULL 表示未使用；非 NULL 表示已使用。
	// 后续实现 Consume 时必须使用原子 UPDATE，
	// 防止两个并发请求同时使用同一个 nonce 登录。
	UsedAt *time.Time `gorm:"column:used_at"`

	// CreatedAt 表示记录创建时间。
	CreatedAt time.Time `gorm:"column:created_at;not null"`
}

func (AuthNonce) TableName() string {
	return "auth_nonces"
}
