package tx

type BuildApproveNFTTxRequest struct {
	ChainID int64 `json:"chain_id"`

	// NFT 合约地址，交易会发到这个地址。
	NFTContract string `json:"nft_contract"`

	// 要授权的 NFT tokenId。
	TokenID string `json:"token_id"`

	// 被授权的地址，默认是拍卖合约地址。
	// 前端可以不传，后端会使用 config 里的 AuctionContract。
	Operator string `json:"operator"`
}

type BuildCreateAuctionTxRequest struct {
	ChainID         int64  `json:"chain_id"`
	ContractAddress string `json:"contract_address"`

	NFTContract string `json:"nft_contract"`
	TokenID     string `json:"token_id"`
	MinBidUSD   string `json:"min_bid_usd"`
	Duration    uint64 `json:"duration"`
}

type BuildPlaceBidTxRequest struct {
	ChainID         int64  `json:"chain_id"`
	ContractAddress string `json:"contract_address"`

	AuctionID uint64 `json:"auction_id"`

	// 如果 bid_token 为空或 0x000...000，表示 ETH 出价，后端会构造 bidEth。
	// 如果 bid_token 是 ERC20 地址，表示 ERC20 出价，后端会构造 bidERC20。
	BidToken string `json:"bid_token"`
	Amount   string `json:"amount"`
}

type BuildCancelAuctionTxRequest struct {
	ChainID         int64  `json:"chain_id"`
	ContractAddress string `json:"contract_address"`
	AuctionID       uint64 `json:"auction_id"`
}

type BuildEndAuctionTxRequest struct {
	ChainID         int64  `json:"chain_id"`
	ContractAddress string `json:"contract_address"`
	AuctionID       uint64 `json:"auction_id"`
}

type TransactionRequestDTO struct {
	ChainID int64  `json:"chain_id"`
	To      string `json:"to"`
	Data    string `json:"data"`
	Value   string `json:"value"`
}
