package auction

type ListAuctionsQuery struct {
	ChainID         int64
	ContractAddress string

	Status string
	Seller string
	NFT    string

	Page     int
	PageSize int
	Sort     string
}

type GetAuctionQuery struct {
	ChainID         int64
	ContractAddress string
	AuctionID       uint64
}

type ListBidsQuery struct {
	ChainID         int64
	ContractAddress string
	AuctionID       uint64

	Page     int
	PageSize int
}
