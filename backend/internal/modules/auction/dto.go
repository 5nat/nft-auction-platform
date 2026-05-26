package auction

import "time"

type AuctionDTO struct {
	ChainID         int64  `json:"chain_id"`
	ContractAddress string `json:"contract_address"`
	AuctionID       uint64 `json:"auction_id"`

	Seller      string `json:"seller"`
	NFTContract string `json:"nft_contract"`
	TokenID     string `json:"token_id"`

	MinBidUSD string `json:"min_bid_usd"`

	HighestBidder    string `json:"highest_bidder"`
	HighestBidToken  string `json:"highest_bid_token"`
	HighestBidAmount string `json:"highest_bid_amount"`
	HighestBidUSD    string `json:"highest_bid_usd"`

	Status  string `json:"status"`
	EndTime uint64 `json:"end_time"`

	CreatedTxHash      string `json:"created_tx_hash"`
	CreatedBlockNumber uint64 `json:"created_block_number"`
	CreatedBlockHash   string `json:"created_block_hash"`
	CreatedLogIndex    uint64 `json:"created_log_index"`

	LastEventName        string `json:"last_event_name"`
	LastEventTxHash      string `json:"last_event_tx_hash"`
	LastEventBlockNumber uint64 `json:"last_event_block_number"`
	LastEventBlockHash   string `json:"last_event_block_hash"`
	LastEventLogIndex    uint64 `json:"last_event_log_index"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type BidDTO struct {
	ChainID         int64  `json:"chain_id"`
	ContractAddress string `json:"contract_address"`
	AuctionID       uint64 `json:"auction_id"`

	Bidder   string `json:"bidder"`
	BidToken string `json:"bid_token"`

	Amount    string `json:"amount"`
	AmountUSD string `json:"amount_usd"`

	TxHash      string `json:"tx_hash"`
	LogIndex    uint64 `json:"log_index"`
	BlockNumber uint64 `json:"block_number"`
	BlockHash   string `json:"block_hash"`

	CreatedAt time.Time `json:"created_at"`
}

type PageMeta struct {
	Page       int   `json:"page"`
	PageSize   int   `json:"page_size"`
	Total      int64 `json:"total"`
	TotalPages int   `json:"total_pages"`
}

type PageResult[T any] struct {
	Items []T      `json:"items"`
	Meta  PageMeta `json:"meta"`
}
