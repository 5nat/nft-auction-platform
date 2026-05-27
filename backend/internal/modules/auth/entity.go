package auth

import "time"

type User struct {
	ID        uint64
	Status    string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Wallet struct {
	ID        uint64
	UserID    uint64
	Address   string
	ChainID   int64
	IsPrimary bool
	CreatedAt time.Time
	UpdatedAt time.Time
}

type AuthNonce struct {
	ID          uint64
	Address     string
	ChainID     int64
	Nonce       string
	MessageHash string
	Message     string
	Domain      string
	URI         string
	IssuedAt    time.Time
	ExpiresAt   time.Time
	UsedAt      *time.Time
	CreatedAt   time.Time
}

type CurrentUser struct {
	UserID        uint64
	WalletID      uint64
	WalletAddress string
	ChainID       int64
}
