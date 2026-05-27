package tx

import "errors"

var (
	ErrInvalidChainID         = errors.New("invalid chain_id")
	ErrInvalidContractAddress = errors.New("invalid contract_address")
	ErrInvalidNFTContract     = errors.New("invalid nft_contract")
	ErrInvalidOperator        = errors.New("invalid operator")
	ErrInvalidTokenID         = errors.New("invalid token_id")
	ErrInvalidAuctionID       = errors.New("invalid auction_id")
	ErrInvalidBidToken        = errors.New("invalid bid_token")
	ErrInvalidAmount          = errors.New("invalid amount")
	ErrInvalidDuration        = errors.New("invalid duration")
	ErrUnauthorized           = errors.New("unauthorized tx actor")
	ErrForbidden              = errors.New("auction action forbidden")
	ErrSellerCannotBid        = errors.New("seller cannot bid on own auction")
	ErrInvalidActor           = errors.New("invalid auction actor")
	ErrChainMismatch          = errors.New("actor chain id does not match request chain id")
	ErrActorChainMismatch     = errors.New("actor chain id does not match request chain id")
)
