# NFT Auction Indexer 设计说明

## 1. 背景

本项目中的 Indexer 用于监听 NFT 拍卖合约产生的链上事件，并将事件同步到 MySQL 中，形成便于后端 API 查询的 read model。

链上合约是最终事实来源，数据库不是权威状态，而是为了提高查询性能、分页能力和业务展示能力而构建的查询模型。

当前 Indexer 主要处理以下事件：

- AuctionCreated
- BidPlaced
- AuctionEnded
- AuctionCancelled

对应数据库表包括：

- auctions
- bids
- processed_logs
- sync_cursors

---

## 2. 整体流程

Indexer 的整体链路如下：

```text
NFTAuctionMarket 合约
        ↓
emit event
        ↓
Indexer 扫描链上 logs
        ↓
根据 topics[0] 判断事件类型
        ↓
ABI binding 解析事件
        ↓
写入 processed_logs，保证幂等
        ↓
写入或更新 auctions / bids
        ↓
更新 sync_cursors
```

当前核心调用链为：

```text
Start
  ↓
RunOnce
  ↓
processRange
  ↓
processLog
  ↓
processAuctionCreated / processBidPlaced / processAuctionEnded / processAuctionCancelled
  ↓
repo.WithTx
  ↓
repo.InsertProcessedLog / CreateAuction / CreateBid / UpdateAuction...
```

---

## 3. 目录结构

当前 Indexer 代码位于：

```text
backend/internal/indexer/
├── indexer.go       # 主流程：Indexer struct / New / Start / RunOnce / processRange / processLog
├── scanner.go       # 链上读取：latest block / filter logs / block hash
├── handlers.go      # 事件处理：AuctionCreated / BidPlaced / AuctionEnded / AuctionCancelled
├── repository.go    # 数据库读写：processed_logs / auctions / bids / sync_cursors
├── cursor.go        # 同步进度协调：nextFromBlock / updateCursor
├── decoder.go       # event topics / 地址和 hash 规范化
├── interfaces.go    # ChainScanner / EventRepository 接口
└── indexer_test.go  # RunOnce 单元测试
```

职责划分如下：

```text
Indexer：
负责同步流程控制。

Scanner：
负责读取链上数据。

Repository：
负责数据库读写和事务。

Handlers：
负责具体事件业务处理。

Cursor：
负责同步进度协调。

Decoder：
负责事件 topic 定义和地址格式化。

Interfaces：
定义 Indexer 所需依赖能力，便于测试替换。
```

---

## 4. Cursor 设计

Indexer 使用 `sync_cursors` 表记录同步进度。

核心字段包括：

- chain_id
- contract_address
- last_processed_block
- last_processed_block_hash

`RunOnce` 执行时，会先读取 cursor：

```text
如果没有 cursor：
    从配置中的 START_BLOCK 开始同步。

如果已有 cursor：
    从 last_processed_block + 1 开始同步。
```

处理完一个 batch 后，会更新 cursor：

```text
last_processed_block = 当前 batch 的 toBlock
last_processed_block_hash = toBlock 对应的 block hash
```

这样可以支持：

- 服务重启后继续同步
- 避免从头重复扫描
- 为后续 reorg 处理预留 block_hash

---

## 5. Confirmations 设计

Indexer 不直接处理最新区块，而是使用 confirmations 机制。

例如：

```text
latest_block = 100
confirmations = 1
target_block = 99
```

这表示 Indexer 当前最多处理到第 99 块，第 100 块暂时等待下一轮确认。

这样做的原因是：

```text
区块刚产生时仍可能发生短暂重组；
等待一定确认数可以降低同步到不稳定区块的风险。
```

当前本地 Anvil 环境中使用：

```text
confirmations = 1
```

生产环境中可以根据链的稳定性调整为更大的数值。

---

## 6. Batch 设计

Indexer 不一次性扫描所有区块，而是按 batchSize 分段处理。

例如：

```text
fromBlock = 20
targetBlock = 30
batchSize = 5
```

则会分成：

```text
20 - 24
25 - 29
30 - 30
```

每个 batch 处理成功后，都会更新 cursor。

这样做的好处是：

- 避免单次 RPC 查询范围过大
- 避免一次事务或一次循环处理过多日志
- 某个 batch 失败时，可以从上一个成功 cursor 继续恢复
- 更适合后续扩展 retry / backoff / metrics

---

## 7. 幂等设计

链上事件同步必须保证幂等。

原因是 Indexer 可能会因为以下情况重复处理同一段区块：

- 服务重启
- RPC 请求失败后重试
- cursor 未及时更新
- 人工回放历史区块
- 后续 reorg 修复

本项目使用 `processed_logs` 表保证事件级幂等。

唯一键为：

```text
chain_id + contract_address + tx_hash + log_index
```

每个事件处理时，都会先插入 processed_logs。

如果插入成功：

```text
说明这条 log 第一次处理，继续执行业务写入。
```

如果唯一键冲突：

```text
说明这条 log 已处理过，本次跳过。
```

该逻辑通过 GORM 的 `OnConflict DoNothing` 实现。

---

## 8. 事务设计

事件处理时，`processed_logs` 和业务表写入必须在同一个数据库事务中完成。

例如处理 `BidPlaced` 时：

```text
插入 processed_logs
插入 bids
更新 auctions.highest_bidder / highest_bid_amount / highest_bid_usd
```

这几步必须原子执行。

如果 processed_logs 写入成功，但 bids 或 auctions 更新失败，会造成严重问题：

```text
下一次重试时 processed_logs 已存在，事件会被跳过；
但业务表其实没有正确更新。
```

因此当前设计中，每个事件 handler 都通过：

```go
repo.WithTx(...)
```

统一包裹事务。

---

## 9. Read Model 设计

数据库中的 `auctions` 和 `bids` 是链上事件派生出来的 read model。

### 9.1 auctions

`auctions` 保存拍卖的当前状态，服务于列表和详情查询。

主要字段包括：

- auction_id
- seller
- nft_contract
- token_id
- min_bid_usd
- highest_bidder
- highest_bid_token
- highest_bid_amount
- highest_bid_usd
- end_time
- status
- block_number
- tx_hash
- log_index

### 9.2 bids

`bids` 保存出价明细，一条 `BidPlaced` 事件对应一条 bid 记录。

主要字段包括：

- auction_id
- bidder
- bid_token
- amount
- amount_usd
- block_number
- tx_hash
- log_index

---

## 10. 事件处理逻辑

### 10.1 AuctionCreated

链上事件：

```text
AuctionCreated(auctionId, seller, nft, tokenId, minBidUsd, endTime)
```

处理逻辑：

```text
插入 processed_logs
创建 auctions 记录
状态设置为 active
```

---

### 10.2 BidPlaced

链上事件：

```text
BidPlaced(auctionId, bidder, bidToken, amount, amountUsd)
```

处理逻辑：

```text
插入 processed_logs
插入 bids 记录
更新 auctions 的最高出价信息
```

合约已经保证新的出价高于当前最高价，因此后端不重新比较金额，只根据链上事件更新 read model。

---

### 10.3 AuctionEnded

链上事件：

```text
AuctionEnded(auctionId, winner, bidToken, amount, amountUsd)
```

处理逻辑：

```text
插入 processed_logs
更新 auctions.status = ended
刷新最终成交信息
```

---

### 10.4 AuctionCancelled

链上事件：

```text
AuctionCancelled(auctionId)
```

处理逻辑：

```text
插入 processed_logs
更新 auctions.status = cancelled
```

---

## 11. 接口抽象

当前 Indexer 依赖两个接口：

```go
type ChainScanner interface {
	LatestBlockNumber(ctx context.Context) (uint64, error)
	FilterLogs(ctx context.Context, fromBlock uint64, toBlock uint64) ([]types.Log, error)
	BlockHash(ctx context.Context, blockNumber uint64) (common.Hash, error)
}
```

```go
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
```

接口定义在 `indexer` 包中，原因是：

```text
接口应该由使用方定义，而不是由实现方定义。
```

`Indexer` 只关心自己需要什么能力，不关心具体实现是：

- 真实 RPC Scanner
- fakeScanner
- GORM Repository
- fakeRepository

当前 `EventRepository` 是一个中间态接口。后续如果方法继续增多，可以进一步拆分为：

- CursorStore
- ProcessedLogStore
- AuctionStore
- TxRunner

但当前阶段先保持稳定，不做过度抽象。

---

## 12. 测试策略

当前已经为 `RunOnce` 编写单元测试。

测试不依赖：

- Anvil
- MySQL
- 真实 RPC
- 真实合约事件

而是通过：

- fakeScanner
- fakeRepository

测试同步流程。

已覆盖场景包括：

```text
1. 没有新 confirmed block 时，不扫描、不更新 cursor。
2. 有新区块但没有 logs 时，正常推进 cursor。
3. confirmations 不足时，不读取 cursor、不扫描 logs。
4. batchSize 分段逻辑正确。
5. latest block、cursor、filter logs、block hash、upsert cursor 出错时，RunOnce 正确返回 error。
```

运行命令：

```bash
go test ./internal/indexer -v
go test ./...
go vet ./...
```

---

## 13. 当前运行方式

单独运行 Indexer：

```bash
go run ./cmd/indexer
```

运行 API 服务并同时启动 Indexer：

```bash
go run ./cmd/server
```

其中是否在 server 中启动 Indexer 由配置项控制：

```text
INDEXER_ENABLED=true / false
```

本地开发时可以让 server 同时启动 API 和 Indexer。

未来生产环境可以拆成两个独立进程：

```text
api-server
indexer-worker
```

两个进程可以使用同一个镜像，但启动命令不同。

例如：

```text
api-server:
    ./server

indexer-worker:
    ./indexer
```

---

## 14. 当前项目阶段

当前 Indexer 已完成以下能力：

```text
链上事件扫描
事件解析
四类核心事件处理
processed_logs 幂等
sync_cursors 断点续扫
confirmations 确认块机制
batch 分段处理
数据库事务一致性
长期轮询 Start(ctx)
集成到 App runtime
Scanner / Repository 解耦
ChainScanner / EventRepository 接口抽象
RunOnce 单元测试
```

当前阶段可以视为：

```text
中级 Web3 Go 后端 Indexer 雏形完成。
```

但它还不是完整生产级 Indexer，主要还缺：

```text
Reorg 检测与修复
Repository 集成测试
事件 handler 单元测试
Prometheus metrics
health / readiness API
RPC retry / backoff
多链 / 多合约支持
历史区块重放工具
```

---

## 15. 后续计划

Indexer 后续还需要继续增强：

1. Reorg 检测与修复
2. Repository 集成测试
3. 事件 handler 单元测试
4. Prometheus metrics
5. health / readiness API
6. RPC retry / backoff
7. 多链 / 多合约支持
8. 历史区块重放工具

其中最重要的是 reorg 处理。

当前已经保存：

```text
last_processed_block_hash
block_hash
```

为后续检测 block hash 是否变化做准备。

---

## 16. 面试表述

可以这样描述本项目的 Indexer 设计：

```text
我实现了一个 NFT 拍卖合约的链上事件 Indexer。链上合约是最终事实来源，Indexer 负责扫描 AuctionCreated、BidPlaced、AuctionEnded、AuctionCancelled 等事件，并同步到 MySQL read model 中，供后端 API 快速查询。

为了保证可靠性，我设计了 processed_logs 表做事件级幂等，使用 sync_cursors 表做断点续扫，并通过 confirmations 机制避免处理最新不稳定区块。事件处理过程中，processed_logs 和业务表写入在同一个数据库事务中完成，避免出现事件被标记已处理但业务表未更新的问题。

工程结构上，我将 Indexer 拆分为 Scanner、Repository、Handlers、Cursor、Decoder 等模块，并通过 ChainScanner 和 EventRepository 接口降低对真实 RPC 和 MySQL 的依赖，从而可以使用 fakeScanner 和 fakeRepository 对 RunOnce 的核心同步流程进行单元测试。
```

---

## 17. 提交建议

建议分两次提交：

```bash
git add internal/indexer
git commit -m "add indexer interfaces and run once tests"
```

```bash
git add docs/indexer.md
git commit -m "document indexer architecture"
```

如果前面的接口和测试还没有单独提交，也可以合并提交：

```bash
git add .
git commit -m "add indexer interfaces tests and architecture docs"
```