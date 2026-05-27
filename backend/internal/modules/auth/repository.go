package auth

import (
	"context"
	"time"
)

// UserRepository 定义 Auth 模块需要的用户身份持久化能力。
type UserRepository interface {
	// FindWalletByAddress 根据钱包地址和 chain_id 查找钱包绑定关系
	// address + chain_id 在 wallets 表中是唯一约束
	FindWalletByAddress(ctx context.Context, address string, chainID int64) (*Wallet, error)

	// CreateUserWithWallet 创建平台用户并绑定钱包
	// 具体实现必须使用 数据库事务，保证 users 和 wallets 同时创建成功或同时回滚 ；否则可能出现 user 已创建但 wallet 创建失败的脏数据
	CreateUserWithWallet(ctx context.Context, address string, chainID int64, now time.Time) (*User, *Wallet, error)

	// FindUserByID 根据平台 user_id 查询用户
	FindUserByID(ctx context.Context, userID uint64) (*User, error)

	// ListWalletsByUserID 查询用户绑定的所有钱包
	ListWalletsByUserID(ctx context.Context, userID uint64) ([]*Wallet, error)
}

// NonceStore 定义钱包登录挑战 nonce 的存取能力。
type NonceStore interface {
	// Save 保存一次登录挑战
	Save(ctx context.Context, n *AuthNonce) error

	// FindByMessageHash 根据 message_hash 查找登录挑战
	FindByMessageHash(ctx context.Context, messageHash string) (*AuthNonce, error)

	// Consume 原子消费 nonce
	// nonce 是防重放攻击的核心状态，不能先查询 used_at 再普通更新
	// 否则两个并发请求可能同时通过 used_at 检查
	//
	// 具体实现应使用类似：
	// UPDATE auth_nonces
	// SET used_at = ?
	// WHERE id = ? AND used_at IS NULL AND expires_at > ?
	//
	// 如果 affected rows = 1，表示消费成功；
	// 如果 affected rows = 0，表示 nonce 已使用、已过期或不存在。
	Consume(ctx context.Context, id uint64, now time.Time) error
}
