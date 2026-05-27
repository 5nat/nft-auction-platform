package auth

import "time"

type CreateNonceCommand struct {
	Address string
	ChainID int64
}

type VerifyCommand struct {
	Address   string
	ChainID   int64
	Message   string
	Signature string
}

type NonceChallenge struct {
	Nonce     string
	Message   string
	ExpiresAt time.Time
}

type LoginResult struct {
	AccessToken string
	TokenType   string
	ExpiresIn   int64
	User        *LoginUser
}

type LoginUser struct {
	ID            uint64
	WalletID      uint64
	WalletAddress string
	ChainID       int64
}

type UserProfile struct {
	ID      uint64
	Status  string
	Wallets []*Wallet
}
