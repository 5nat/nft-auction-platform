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
)
