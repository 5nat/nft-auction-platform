package indexer

import (
	"context"

	"github.com/5nat/nft-auction-platform/backend/internal/model"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

// ChainScanner 定义 Indexer 需要的链上读取能力。
//
// 注意：
// 这个接口不是为了“抽象而抽象”，
// 而是为了让 Indexer 不直接依赖具体的 *Scanner 实现。
//
// 真实运行时：
//
//	使用 Scanner，通过 ethclient 访问链。
//
// 单元测试时：
//
//	可以使用 fakeScanner，不需要启动 Anvil，也不需要真实 RPC。
type ChainScanner interface {
	// LatestBlockNumber 返回当前链上的最新区块高度
	LatestBlockNumber(ctx context.Context) (uint64, error)

	// FilterLogs 扫描指定区块范围内的合约事件日志
	FilterLogs(ctx context.Context, fromBlock uint64, toBlock uint64) ([]types.Log, error)

	// BlockHash 返回指定区块的 block hash
	// cursor 写入 sync_cursors 时需要保存 block hash，为后续 reorg 处理预留
	BlockHash(ctx context.Context, blockNumber uint64) (common.Hash, error)
}

// EventRepository 定义 Indexer 需要的数据库读写能力。
//
// 这个接口的目的不是让所有数据库操作都抽象化，
// 而是让 Indexer 不直接依赖具体的 GORM Repository 实现。
//
// 真实运行时：
//
//	使用 Repository，通过 GORM 写 MySQL。
//
// 单元测试时：
//
//	可以使用 fakeRepository，不需要真实 MySQL。
type EventRepository interface {
	NextFromBlock(ctx context.Context, startBlock uint64) (uint64, error)

	UpsertCursor(ctx context.Context, blockNumber uint64, blockHash string) error

	WithTx(ctx context.Context, fn func(repo EventRepository) error) error

	InsertProcessedLog(ctx context.Context, lg types.Log, eventName string) (bool, error)

	CreateAuction(ctx context.Context, auction model.Auction) error

	CreateBid(ctx context.Context, bid model.Bid) error

	UpdateAuctionHighestBid(ctx context.Context, input UpdateAuctionHighestBidInput) error

	MarkAuctionEnded(ctx context.Context, input MarkAuctionEndedInput) error

	MarkAuctionCancelled(ctx context.Context, input MarkAuctionCancelledInput) error
}
