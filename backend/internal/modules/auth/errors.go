package auth

import "errors"

var (
	ErrInvalidAddress   = errors.New("invalid wallet address")
	ErrInvalidChainID   = errors.New("invalid chain id")
	ErrInvalidMessage   = errors.New("invalid auth message")
	ErrInvalidSignature = errors.New("invalid signature")

	ErrNonceNotFound    = errors.New("auth nonce not found")
	ErrNonceExpired     = errors.New("auth nonce expired")
	ErrNonceUsed        = errors.New("auth nonce already used")
	ErrNonceUnavailable = errors.New("auth nonce unavailable")

	ErrWalletNotFound = errors.New("wallet not found")
	ErrUserNotFound   = errors.New("user not found")

	ErrUnauthorized = errors.New("unauthorized")
	ErrInvalidToken = errors.New("invalid token")
)
