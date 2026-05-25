package auction

type ListAuctionsQuery struct {
	Status   string
	Seller   string
	NFT      string
	Page     int
	PageSize int
}

type AuctionDTO struct {
	AuctionID        uint64 `json:"auction_id"`
	ContractAddress  string `json:"contract_address"`
	Seller           string `json:"seller"`
	NFTContract      string `json:"nft_contract"`
	TokenID          string `json:"token_id"`
	MinBidUSD        uint64 `json:"min_bid_usd"`
	HighestBidder    uint64 `json:"highest_bidder"`
	HighestBIdToken  uint64 `json:"highest_bid_token"`
	HighestBidAmount uint64 `json:"highest_bid_amount"`
	HighestBidUSD    uint64 `json:"highest_bid_usd"`
	EndTime          uint64 `json:"end_time"`
	Status           string `json:"status"`
	BlockNumber      uint64 `json:"block_number"`
	BlockHash        string `json:"block_hash"`
	TxHash           string `json:"tx_hash"`
	LogIndex         int    `json:"log_index"`
	CreatedAt        uint64 `json:"created_at"`
	UpdatedAt        uint64 `json:"updated_at"`
}

type BidDTO struct {
	AuctionID   uint64 `json:"auction_id"`
	Bidder      string `json:"bidder"`
	BidToken    string `json:"bid_token"`
	Amount      uint64 `json:"amount"`
	AmountUSD   uint64 `json:"amount_usd"`
	BlockNumber uint64 `json:"block_number"`
	BlockHash   string `json:"block_hash"`
	TxHash      string `json:"tx_hash"`
	LogIndex    int    `json:"log_index"`
	CreatedAt   uint64 `json:"created_at"`
}

type PageMeta struct {
	Page      int   `json:"page"`
	PageSize  int   `json:"page_size"`
	Total     int64 `json:"total"`
	TotalPage int   `json:"total_page"`
}

type PageResult[T any] struct {
	Items []T      `json:"items"`
	Meta  PageMeta `json:"meta"`
}
