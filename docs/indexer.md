# Indexer 设计说明

## 1. 背景

本项目是一个基于 Go + Solidity 的 NFT 拍卖平台后端。链上合约负责处理拍卖创建、出价、结束和取消等核心逻辑；后端 Indexer 负责扫描链上事件，并将事件转换为 MySQL 中的查询模型。

Indexer 的核心目标不是简单“监听事件”，而是可靠地构建链上数据的 read model。

整体链路如下：

```text
用户发起链上交易
    ↓
合约执行并 emit event
    ↓
Indexer 扫描 confirmed blocks
    ↓
解析链上 logs
    ↓
processed_logs 保证幂等
    ↓
写入或更新 auctions / bids
    ↓
sync_cursors 记录同步进度
    ↓
REST API 查询 MySQL read model
```

## 2. Indexer 的定位

Indexer 是链上数据和后端查询系统之间的桥梁。

链上合约是最终事实来源，数据库是面向查询的 read model。

因此，本项目不会把数据库作为业务事实来源，而是通过 Indexer 根据链上事件持续更新数据库。

当前 Indexer 处理四类事件：

```text
AuctionCreated
BidPlaced
AuctionEnded
AuctionCancelled
```

事件到数据表的映射如下：

```text
AuctionCreated
    → 插入 auctions
    → status = active
    → 初始化 created_* 字段
    → 初始化 last_event_* 字段

BidPlaced
    → 插入 bids
    → 更新 auctions 当前最高出价
    → 更新 auctions.last_event_* 字段

AuctionEnded
    → 更新 auctions.status = ended
    → 更新最终 winner / amount / amountUsd
    → 更新 auctions.last_event_* 字段

AuctionCancelled
    → 更新 auctions.status = cancelled
    → 更新 auctions.last_event_* 字段
```

## 3. 当前目录结构

Indexer 代码位于：

```text
internal/indexer/
├── indexer.go       主流程：Start、RunOnce、processRange、processLog
├── scanner.go       链上读取：latest block、FilterLogs、BlockHash
├── handlers.go      事件处理：AuctionCreated、BidPlaced、AuctionEnded、AuctionCancelled
├── repository.go    数据库读写：cursor、processed_logs、auctions、bids
├── cursor.go        cursor 协调
├── decoder.go       event topic、地址和 hash 规范化
├── event_meta.go    事件元信息封装
├── interfaces.go    ChainScanner、EventRepository 接口
└── indexer_test.go  RunOnce 主流程测试
```

该结构采用单个 `indexer` package 多文件组织，而不是过早拆成多个子 package。

这样做的原因是：

```text
1. 当前 Indexer 仍然是一个内聚模块；
2. scanner、repository、handler、cursor 都围绕同一个同步流程；
3. 过早拆 package 容易制造循环依赖；
4. 多文件拆分已经能保证职责清晰；
5. 后续如果复杂度继续提高，再考虑拆成更细的 package。
```

## 4. 核心组件职责

### 4.1 Indexer

`Indexer` 是主流程控制器。

主要职责：

```text
读取最新区块
计算 confirmed target block
读取同步 cursor
按 batch 扫描 logs
按事件 topic 分发处理
处理成功后推进 cursor
长期循环运行
```

主要方法：

```text
Start(ctx)
RunOnce(ctx)
processRange(ctx, fromBlock, toBlock)
processLog(ctx, log)
```

### 4.2 Scanner

`Scanner` 负责链上读取。

主要职责：

```text
查询 latest block
按区块范围读取 logs
查询指定区块的 block hash
```

接口为：

```go
type ChainScanner interface {
    LatestBlockNumber(ctx context.Context) (uint64, error)
    FilterLogs(ctx context.Context, fromBlock uint64, toBlock uint64) ([]types.Log, error)
    BlockHash(ctx context.Context, blockNumber uint64) (common.Hash, error)
}
```

抽象 Scanner 的目的是让 Indexer 主流程不直接依赖 `ethclient`，方便测试和后续替换 RPC 实现。

### 4.3 Repository

`Repository` 负责数据库读写。

主要职责：

```text
读取和更新 sync_cursors
插入 processed_logs
插入 auctions
插入 bids
更新 auctions 最高出价
标记 auction ended
标记 auction cancelled
```

接口为：

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

Repository 中封装 GORM 细节，避免 handler 直接写 SQL。

### 4.4 EventMeta

`EventMeta` 用于统一封装链上事件位置信息。

字段包括：

```text
event_name
tx_hash
block_number
block_hash
log_index
```

它的作用是避免在每个 handler 中重复提取：

```go
normalizeHash(lg.TxHash)
normalizeHash(lg.BlockHash)
uint64(lg.Index)
```

同时，`EventMeta` 会写入：

```text
auctions.created_*
auctions.last_event_*
bids.tx_hash / block_number / block_hash / log_index
processed_logs
```

## 5. 同步流程

### 5.1 Start 长期运行

`Start(ctx)` 是长期运行入口。

它会循环调用 `RunOnce(ctx)`：

```text
Start
    ↓
RunOnce
    ↓
等待 poll_interval
    ↓
RunOnce
    ↓
等待 poll_interval
    ↓
...
```

当收到 `ctx.Done()` 时，Indexer 会优雅退出。

该设计使 Indexer 可以作为：

```text
cmd/indexer 独立 worker
cmd/server 中的后台 goroutine
```

两种方式运行。

### 5.2 RunOnce 单次同步

`RunOnce(ctx)` 表示执行一次同步。

核心流程：

```text
1. 查询 latest block
2. 根据 confirmations 计算 target block
3. 如果 latest block 不足 confirmations，直接返回
4. 从 sync_cursors 读取 fromBlock
5. 如果 fromBlock > targetBlock，说明没有新区块
6. 按 batchSize 分段扫描 logs
7. 每个 batch 处理完成后推进 cursor
```

伪流程：

```text
latestBlock = scanner.LatestBlockNumber()
targetBlock = latestBlock - confirmations
fromBlock = repo.NextFromBlock(startBlock)

for fromBlock <= targetBlock:
    toBlock = min(fromBlock + batchSize - 1, targetBlock)
    processRange(fromBlock, toBlock)
    updateCursor(toBlock)
    fromBlock = toBlock + 1
```

## 6. Cursor 机制

`sync_cursors` 表记录同步进度。

唯一键：

```text
chain_id + contract_address
```

核心字段：

```text
last_processed_block
last_processed_block_hash
```

当 Indexer 第一次运行时，如果 cursor 不存在，则从配置中的 `START_BLOCK` 开始。

如果 cursor 存在，则从：

```text
last_processed_block + 1
```

继续同步。

每个 batch 成功处理完成后，Indexer 更新 cursor：

```text
last_processed_block = toBlock
last_processed_block_hash = blockHash(toBlock)
```

这样可以支持服务重启后的断点续扫。

## 7. Confirmations 机制

Indexer 不直接处理最新区块，而是等待一定确认数。

计算方式：

```text
targetBlock = latestBlock - confirmations
```

例如：

```text
latestBlock = 100
confirmations = 3
targetBlock = 97
```

此时 Indexer 只处理到 97 区块。

这样做的原因是：

```text
最新区块可能发生 reorg
等待 confirmations 可以降低处理不稳定区块的风险
```

当前本地 Anvil 测试中 confirmations 可以设置为 1。

在真实网络中，可以根据链的稳定性设置更高确认数。

## 8. Batch 扫描机制

Indexer 不一次性扫描从 startBlock 到 latestBlock 的全部区块，而是按 batch 分段扫描。

例如：

```text
fromBlock = 20
targetBlock = 30
batchSize = 5
```

扫描范围为：

```text
20 - 24
25 - 29
30 - 30
```

这样做的原因是：

```text
1. 避免一次 RPC 请求区块范围过大；
2. 降低 RPC 超时风险；
3. 每个 batch 成功后都能推进 cursor；
4. 失败时只需要重试当前 batch。
```

## 9. 幂等机制

Indexer 可能会重复扫描同一区块范围。

例如：

```text
第一次扫描 100 ~ 200 成功；
服务重启；
为了安全回扫，又扫描 180 ~ 220。
```

如果没有幂等控制，同一条事件可能被重复写入业务表。

本项目使用 `processed_logs` 保证事件级幂等。

唯一键为：

```text
chain_id + contract_address + tx_hash + log_index
```

处理每条 log 时，先尝试插入 `processed_logs`：

```text
插入成功：
    说明该 log 第一次处理，继续写业务表

插入失败且是唯一键冲突：
    说明该 log 已经处理过，直接跳过
```

这样可以确保同一条链上 log 只被处理一次。

## 10. 事务一致性

每条事件的处理必须满足：

```text
processed_logs 写入
+
业务表写入或更新
=
同一个数据库事务
```

否则可能出现严重不一致。

错误场景：

```text
1. processed_logs 插入成功；
2. bids 或 auctions 写入失败；
3. 服务退出；
4. 下次重扫时发现 processed_logs 已存在；
5. 业务表永远不会补上。
```

因此本项目在处理每条 log 时使用：

```go
repo.WithTx(ctx, func(repo EventRepository) error {
    inserted, err := repo.InsertProcessedLog(...)
    if !inserted {
        return nil
    }

    // 写入或更新业务表
    return nil
})
```

确保幂等记录和业务数据在同一个事务里提交或回滚。

## 11. 事件处理逻辑

### 11.1 AuctionCreated

处理逻辑：

```text
1. ParseAuctionCreated
2. 检查 auctionId / endTime 是否能安全转换为 uint64
3. 构造 EventMeta
4. 事务内插入 processed_logs
5. 插入 auctions
6. 初始化 status = active
7. 初始化 created_* 字段
8. 初始化 last_event_* 字段
```

写入 `auctions` 的核心字段：

```text
chain_id
contract_address
auction_id
seller
nft_contract
token_id
min_bid_usd
end_time
status = active

created_tx_hash
created_block_number
created_block_hash
created_log_index

last_event_name = AuctionCreated
last_event_tx_hash
last_event_block_number
last_event_block_hash
last_event_log_index
```

### 11.2 BidPlaced

处理逻辑：

```text
1. ParseBidPlaced
2. 检查 auctionId 是否能安全转换为 uint64
3. 构造 EventMeta
4. 事务内插入 processed_logs
5. 插入 bids
6. 更新 auctions.highest_bidder
7. 更新 auctions.highest_bid_token
8. 更新 auctions.highest_bid_amount
9. 更新 auctions.highest_bid_usd
10. 更新 auctions.last_event_* 字段
```

`bids` 表保存历史出价记录。

`auctions` 表保存当前最高出价。

### 11.3 AuctionEnded

处理逻辑：

```text
1. ParseAuctionEnded
2. 检查 auctionId 是否能安全转换为 uint64
3. 构造 EventMeta
4. 事务内插入 processed_logs
5. 更新 auctions.status = ended
6. 更新最终 winner / bidToken / amount / amountUsd
7. 更新 auctions.last_event_* 字段
```

### 11.4 AuctionCancelled

处理逻辑：

```text
1. ParseAuctionCancelled
2. 检查 auctionId 是否能安全转换为 uint64
3. 构造 EventMeta
4. 事务内插入 processed_logs
5. 更新 auctions.status = cancelled
6. 更新 auctions.last_event_* 字段
```

## 12. Polling 与实时监听

当前 Indexer 使用 polling 模式，而不是只依赖 WebSocket 实时监听。

当前流程是：

```text
定时查询 latest block
根据 confirmations 计算 target block
根据 sync_cursors 决定 fromBlock
扫描 logs
写库
更新 cursor
```

这样设计的优点是：

```text
1. 服务重启后可以继续同步；
2. RPC 短暂异常后不会漏事件；
3. WebSocket 断线不会导致数据永久丢失；
4. 支持人工回放历史区块；
5. 配合 processed_logs 可以安全重复扫描。
```

未来可以增加 WebSocket newHeads 作为触发器：

```text
收到新区块 header
    ↓
触发 RunOnce
```

但真正的数据同步仍然依赖：

```text
cursor + confirmations + processed_logs
```

也就是说，WebSocket 只用于降低延迟，不作为唯一可靠来源。

## 13. Reorg 预留设计

当前模型已经保存了多个 block hash：

```text
auctions.created_block_hash
auctions.last_event_block_hash
bids.block_hash
processed_logs.block_hash
sync_cursors.last_processed_block_hash
```

后续可以实现 reorg 检测：

```text
1. 每次 RunOnce 前读取 sync_cursors.last_processed_block
2. 查询链上该 block 的当前 block_hash
3. 与 sync_cursors.last_processed_block_hash 对比
4. 如果不一致，说明发生 reorg
5. 停止 Indexer 或执行回滚
```

第一阶段可以先实现：

```text
检测到 reorg 后停止并报警
```

后续再实现自动回滚：

```text
1. 回退 N 个区块
2. 删除该区块范围内的 processed_logs
3. 删除该区块范围内的 bids
4. 重建受影响的 auctions read model
5. 从回退点重新同步
```

## 14. 测试策略

当前已经为 `RunOnce` 主流程设计了单元测试。

测试重点包括：

```text
没有新 confirmed block 时不扫描
有新区块但没有 logs 时仍推进 cursor
confirmations 不足时不处理
batchSize 分段逻辑正确
latest block 查询错误能正确返回
cursor 查询错误能正确返回
filter logs 错误能正确返回
block hash 查询错误能正确返回
cursor 更新错误能正确返回
```

测试使用 fake scanner 和 fake repository，而不是连接真实 RPC 和 MySQL。

这样做的好处是：

```text
测试速度快
不依赖外部服务
可以稳定覆盖边界条件
便于验证主流程逻辑
```

后续可以补充：

```text
Repository 集成测试
事件 handler 单元测试
Indexer + MySQL 集成测试
reorg 检测测试
```

## 15. 当前运行方式

### 15.1 独立运行 Indexer

```bash
go run ./cmd/indexer
```

适合把 Indexer 作为独立 worker 运行。

### 15.2 随 server 一起运行

通过配置开启：

```env
INDEXER_ENABLED=true
```

然后启动：

```bash
go run ./cmd/server
```

此时 `cmd/server` 会同时启动：

```text
HTTP server
Indexer goroutine
```

当前阶段推荐开发时分开运行：

```text
终端 1：go run ./cmd/server
终端 2：go run ./cmd/indexer
```

这样日志更清楚，问题更容易定位。

## 16. 与 REST API 的关系

Indexer 负责构建 MySQL read model。

REST API 不直接查询链上合约，而是查询数据库：

```text
GET /api/v1/auctions
    ↓
读取 auctions

GET /api/v1/auctions/:auctionId
    ↓
读取 auctions

GET /api/v1/auctions/:auctionId/bids
    ↓
读取 bids
```

这样可以带来：

```text
查询速度更快
支持分页和筛选
降低 RPC 依赖
前端体验更稳定
便于后续做缓存、搜索和统计
```

## 17. 当前完成度

当前 Indexer 已完成：

```text
四类事件处理
MySQL read model 写入
processed_logs 幂等
sync_cursors 断点续扫
confirmations 机制
batch 扫描
created_* 字段写入
last_event_* 字段更新
Scanner / Repository 接口抽象
RunOnce 单元测试
cmd/indexer 独立运行
cmd/server 集成运行
```

尚未完成：

```text
reorg 自动检测
reorg 自动回滚
Repository 集成测试
事件 handler 单元测试
Prometheus metrics
RPC retry / backoff
WebSocket newHeads 触发
```

## 18. 面试表述

可以这样介绍 Indexer 设计：

```text
我没有把 WebSocket 订阅作为唯一的数据同步来源，而是实现了一个基于 cursor 的 polling indexer。Indexer 会根据 confirmations 只处理稳定区块，通过 sync_cursors 支持断点续扫，通过 processed_logs 保证事件级幂等，并且把 processed_logs 和业务表写入放在同一个事务中，避免处理记录和业务数据不一致。auctions 表是由链上事件构建的 read model，AuctionCreated 初始化记录，BidPlaced、AuctionEnded、AuctionCancelled 持续更新状态。为了提高可追溯性和为后续 reorg 做准备，我保存了 created_* 和 last_event_* 事件位置信息，以及 block_hash。
```

## 19. 后续演进方向

后续可以继续增强：

```text
1. 实现 reorg 检测，先检测后报警
2. 实现 reorg 自动回滚和重建 read model
3. 给 RPC 调用增加 retry / backoff
4. 增加 Prometheus metrics
5. 增加 readiness 检查
6. 增加 Repository 集成测试
7. 增加事件 handler 单元测试
8. 增加 WebSocket newHeads 触发 RunOnce
9. 与 REST API 查询层联调
10. 与 Tx Service 交易参数构造联调
```