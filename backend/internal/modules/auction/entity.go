package auction

import "time"

type Auction struct {
	ChainID         int64
	ContractAddress string
	AuctionID       uint64

	Seller      string
	NFTContract string
	TokenID     string

	MinBidUSD string

	HighestBidder    string
	HighestBidToken  string
	HighestBidAmount string
	HighestBidUSD    string

	Status  string
	EndTime uint64

	CreatedTxHash      string
	CreatedBlockNumber uint64
	CreatedBlockHash   string
	CreatedLogIndex    uint64

	LastEventName        string
	LastEventTxHash      string
	LastEventBlockNumber uint64
	LastEventBlockHash   string
	LastEventLogIndex    uint64

	CreatedAt time.Time
	UpdatedAt time.Time
}

type Bid struct {
	ChainID         int64
	ContractAddress string
	AuctionID       uint64

	Bidder    string
	BidToken  string
	Amount    string
	AmountUSD string

	TxHash      string
	LogIndex    uint64
	BlockNumber uint64
	BlockHash   string

	CreatedAt time.Time
}
