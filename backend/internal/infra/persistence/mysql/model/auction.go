package model

import "time"

const (
	// AuctionStatusActive 表示拍卖正在进行。
	AuctionStatusActive = "active"

	// AuctionStatusEnded 表示拍卖已经正常结束。
	AuctionStatusEnded = "ended"

	// AuctionStatusCancelled 表示拍卖已被取消。
	AuctionStatusCancelled = "cancelled"
)

type Auction struct {
	// ID 是数据库内部自增主键。
	// 它只用于数据库内部定位记录，不代表链上的 auctionId。
	ID uint64 `gorm:"column:id;primaryKey;autoIncrement" json:"id"`

	// ChainID 表示该拍卖事件来自哪条链。
	// 例如 Anvil 本地链通常是 31337，Ethereum Mainnet 是 1，Sepolia 是 11155111。
	// 在多链场景下，同一个 auction_id 可能在不同链上重复出现，因此必须保存 chain_id。
	ChainID int64 `gorm:"column:chain_id;not null;uniqueIndex:uk_auction_chain_contract_id,priority:1;index:idx_auctions_chain_status_end,priority:1" json:"chain_id"`

	// ContractAddress 表示该拍卖所属的拍卖市场合约地址。
	// 同一条链上也可能部署多个拍卖合约，因此 contract_address 也是业务唯一键的一部分。
	ContractAddress string `gorm:"column:contract_address;type:char(42);not null;uniqueIndex:uk_auction_chain_contract_id,priority:2;index:idx_auctions_contract_address" json:"contract_address"`

	// AuctionID 是链上合约里的拍卖 ID。
	// 它只在 chain_id + contract_address 这个范围内唯一。
	// 因此业务唯一键是：
	// chain_id + contract_address + auction_id
	AuctionID uint64 `gorm:"column:auction_id;not null;uniqueIndex:uk_auction_chain_contract_id,priority:3" json:"auction_id"`

	// Seller 是拍卖发起人的钱包地址。
	Seller string `gorm:"column:seller;type:char(42);not null;index:idx_auctions_seller" json:"seller"`

	// NFTContract 是被拍卖的 NFT 合约地址。
	NFTContract string `gorm:"column:nft_contract;type:char(42);not null;index:idx_auctions_nft_contract" json:"nft_contract"`

	// TokenID 是被拍卖的 NFT tokenId。
	// 链上 tokenId 通常是 uint256，为了避免溢出和精度问题，这里使用 string 保存。
	TokenID string `gorm:"column:token_id;type:varchar(78);not null" json:"token_id"`

	// MinBidUSD 是最低起拍价格，按合约中的 USD 精度保存。
	// uint256 金额统一使用 string，避免 uint64 溢出和 float64 精度丢失。
	MinBidUSD string `gorm:"column:min_bid_usd;type:varchar(78);not null" json:"min_bid_usd"`

	// HighestBidder 是当前最高出价人地址。
	// 拍卖刚创建时可以为空，发生 BidPlaced 后由 Indexer 更新。
	HighestBidder string `gorm:"column:highest_bidder;type:char(42)" json:"highest_bidder"`

	// HighestBidToken 是当前最高出价使用的资产地址。
	// 零地址表示 ETH，非零地址表示 ERC20 token。
	HighestBidToken string `gorm:"column:highest_bid_token;type:char(42)" json:"highest_bid_token"`

	// HighestBidAmount 是当前最高出价的原始金额。
	// ETH 使用 wei，ERC20 使用 token 的最小单位。
	HighestBidAmount string `gorm:"column:highest_bid_amount;type:varchar(78)" json:"highest_bid_amount"`

	// HighestBidUSD 是当前最高出价折算后的 USD 金额。
	// 用于支持 ETH 和 ERC20 混合出价时的统一比较。
	HighestBidUSD string `gorm:"column:highest_bid_usd;type:varchar(78)" json:"highest_bid_usd"`

	// Status 表示拍卖当前状态。
	// 由 AuctionCreated 初始化为 active；
	// AuctionEnded 更新为 ended；
	// AuctionCancelled 更新为 cancelled。
	Status string `gorm:"column:status;type:varchar(32);not null;index:idx_auctions_chain_status_end,priority:2" json:"status"`

	// EndTime 是拍卖结束时间，通常是 Unix timestamp 秒级时间戳。
	EndTime uint64 `gorm:"column:end_time;not null;index:idx_auctions_chain_status_end,priority:3" json:"end_time"`

	// CreatedTxHash 记录 AuctionCreated 事件所在交易 hash。
	// 这个字段表示“该拍卖是由哪笔链上交易创建的”。
	CreatedTxHash string `gorm:"column:created_tx_hash;type:char(66);not null;index:idx_auctions_created_tx_hash" json:"created_tx_hash"`

	// CreatedBlockNumber 记录 AuctionCreated 事件所在区块号。
	// 可用于按创建时间排序，也可用于排查链上同步问题。
	CreatedBlockNumber uint64 `gorm:"column:created_block_number;not null;index:idx_auctions_created_block_number" json:"created_block_number"`

	// CreatedBlockHash 记录 AuctionCreated 事件所在区块 hash。
	// 后续处理 reorg 时，可以用它判断创建事件所在区块是否发生变化。
	CreatedBlockHash string `gorm:"column:created_block_hash;type:char(66);not null" json:"created_block_hash"`

	// CreatedLogIndex 记录 AuctionCreated 事件在交易 receipt 中的 log 序号。
	// 同一个交易中可能产生多条 log，所以 log_index 是定位事件的重要字段。
	CreatedLogIndex uint64 `gorm:"column:created_log_index;not null" json:"created_log_index"`

	// LastEventName 记录最近一次改变该 auction read model 的事件名称。
	// 可能是 AuctionCreated、BidPlaced、AuctionEnded、AuctionCancelled。
	LastEventName string `gorm:"column:last_event_name;type:varchar(64);not null;index:idx_auctions_last_event_name" json:"last_event_name"`

	// LastEventTxHash 记录最近一次改变该拍卖状态的交易 hash。
	// 例如最近一次最高出价、结束拍卖或取消拍卖对应的交易。
	LastEventTxHash string `gorm:"column:last_event_tx_hash;type:char(66);not null" json:"last_event_tx_hash"`

	// LastEventBlockNumber 记录最近一次改变该拍卖状态的事件所在区块号。
	LastEventBlockNumber uint64 `gorm:"column:last_event_block_number;not null;index:idx_auctions_last_event_block_number" json:"last_event_block_number"`

	// LastEventBlockHash 记录最近一次改变该拍卖状态的事件所在区块 hash。
	// 后续做 reorg 回滚或数据修复时会用到。
	LastEventBlockHash string `gorm:"column:last_event_block_hash;type:char(66);not null" json:"last_event_block_hash"`

	// LastEventLogIndex 记录最近一次改变该拍卖状态的事件 log 序号。
	LastEventLogIndex uint64 `gorm:"column:last_event_log_index;not null" json:"last_event_log_index"`

	// CreatedAt 是数据库记录创建时间。
	// 注意：它不是链上拍卖创建时间，而是 Indexer 写入数据库的时间。
	CreatedAt time.Time `gorm:"column:created_at" json:"created_at"`

	// UpdatedAt 是数据库记录最近更新时间。
	// 当 BidPlaced、AuctionEnded、AuctionCancelled 更新该记录时，GORM 会自动更新这个字段。
	UpdatedAt time.Time `gorm:"column:updated_at" json:"updated_at"`
}

func (Auction) TableName() string {
	return "auctions"
}
