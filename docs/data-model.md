# Data Model 设计说明

## 1. 背景

本项目是一个基于 Go + Solidity 的 NFT 拍卖平台后端。链上合约负责处理拍卖创建、出价、结束和取消等核心业务逻辑；后端 Indexer 负责扫描链上事件，并将链上事件转换为 MySQL 中的查询模型，供 REST API 查询。

本项目的数据模型不是传统 CRUD 模型，而是基于链上事件构建的 read model。

整体链路如下：

```text
用户发起链上交易
    ↓
合约执行并 emit event
    ↓
Indexer 扫描 confirmed blocks
    ↓
解析 AuctionCreated / BidPlaced / AuctionEnded / AuctionCancelled
    ↓
写入 processed_logs 保证幂等
    ↓
写入或更新 auctions / bids
    ↓
更新 sync_cursors
    ↓
REST API 查询 MySQL read model
```

## 2. 整体设计原则

### 2.1 链上事件是事实来源

本项目中，链上合约是最终事实来源，数据库不是最终事实来源。

数据库中的 `auctions` 和 `bids` 是由链上事件派生出来的 read model，用于支持高效查询、分页、筛选和前端展示。

因此，后端不直接修改拍卖业务状态，而是通过 Indexer 处理链上事件来更新数据库。

### 2.2 支持多链、多合约

链上业务 ID 不能单独作为全局唯一标识。

例如，`auction_id = 1` 只在某一条链、某一个拍卖合约内部唯一。不同链、不同合约都可能出现相同的 `auction_id`。

因此拍卖的业务唯一键设计为：

```text
chain_id + contract_address + auction_id
```

该设计可以支持后续多链部署、多合约部署和测试网络切换。

### 2.3 链上 log 的唯一定位

一笔交易中可能产生多条 log，因此不能只用 `tx_hash` 唯一定位一条链上事件。

一条链上 log 的唯一身份应为：

```text
chain_id + contract_address + tx_hash + log_index
```

该设计用于 `processed_logs` 和 `bids` 表，保证同一条链上 log 不会被重复处理或重复写入。

### 2.4 uint256 使用 string 保存

链上金额、tokenId 等字段通常是 `uint256` 类型，可能超过 Go 的 `uint64` 范围，也不能使用 `float64`，否则会出现精度损失。

因此项目中链上大整数统一使用字符串保存，数据库字段使用：

```text
varchar(78)
```

原因是 `uint256` 最大十进制长度约为 78 位。

涉及字段包括：

```text
token_id
min_bid_usd
amount
amount_usd
highest_bid_amount
highest_bid_usd
```

### 2.5 地址和 hash 使用固定长度字段

以太坊地址长度固定为：

```text
0x + 40 个十六进制字符 = 42
```

因此地址字段使用：

```text
char(42)
```

交易 hash 和区块 hash 长度固定为：

```text
0x + 64 个十六进制字符 = 66
```

因此 hash 字段使用：

```text
char(66)
```

### 2.6 保存 block_hash，为 reorg 预留

Indexer 不仅保存 `block_number`，还保存 `block_hash`。

原因是区块号在链重组后可能对应不同的区块 hash。保存 block hash 后，后续可以实现 reorg 检测、回滚和数据修复。

### 2.7 区分状态表、事件明细表和控制表

本项目主要包含四张表：

```text
auctions        拍卖当前状态表
bids            出价历史明细表
processed_logs  Indexer 幂等控制表
sync_cursors    Indexer 同步进度表
```

其中：

```text
auctions 是 read model，会被事件持续更新；
bids 是事件明细表，基本只插入不更新；
processed_logs 用于保证事件处理幂等；
sync_cursors 用于支持断点续扫。
```

## 3. auctions 表设计

### 3.1 表含义

`auctions` 表保存拍卖的当前状态。

它不是简单保存 `AuctionCreated` 事件，而是由多个链上事件共同维护的 read model。

事件与表更新关系如下：

```text
AuctionCreated    → 创建 auctions 记录，status = active
BidPlaced         → 更新最高出价信息
AuctionEnded      → 更新 status = ended
AuctionCancelled  → 更新 status = cancelled
```

因此，`auctions` 表面向的是查询 API，例如：

```text
GET /api/v1/auctions
GET /api/v1/auctions/:auctionId
```

### 3.2 字段说明

| 字段 | 含义 |
|---|---|
| id | 数据库内部自增主键，只用于数据库内部定位记录，不代表链上的 auctionId。 |
| chain_id | 数据来源链 ID，例如 Anvil 为 31337，Ethereum Mainnet 为 1，Sepolia 为 11155111。 |
| contract_address | 拍卖市场合约地址。用于区分同一条链上的不同拍卖合约。 |
| auction_id | 链上合约中的拍卖 ID。它只在 `chain_id + contract_address` 范围内唯一。 |
| seller | 拍卖发起人的钱包地址。 |
| nft_contract | 被拍卖的 NFT 合约地址。 |
| token_id | 被拍卖的 NFT tokenId。链上通常是 uint256，因此使用 string 保存。 |
| min_bid_usd | 最低起拍价格，按合约中的 USD 精度保存。 |
| highest_bidder | 当前最高出价人地址。拍卖刚创建时可以为空，发生 BidPlaced 后更新。 |
| highest_bid_token | 当前最高出价使用的资产地址。零地址表示 ETH，非零地址表示 ERC20 token。 |
| highest_bid_amount | 当前最高出价原始金额。ETH 使用 wei，ERC20 使用 token 最小单位。 |
| highest_bid_usd | 当前最高出价折算后的 USD 金额，用于支持 ETH 和 ERC20 混合出价比较。 |
| status | 拍卖当前状态，可选值为 active、ended、cancelled。 |
| end_time | 拍卖结束时间，通常是 Unix timestamp 秒级时间戳。 |
| created_at | 数据库记录创建时间，不等于链上事件时间。 |
| updated_at | 数据库记录最近更新时间。 |

### 3.3 created_* 字段

`created_*` 字段记录 `AuctionCreated` 事件的位置。

包括：

```text
created_tx_hash
created_block_number
created_block_hash
created_log_index
```

它回答的问题是：

```text
这个拍卖最初是由哪条链上事件创建的？
```

这些字段在创建后通常不再变化。

字段含义如下：

| 字段 | 含义 |
|---|---|
| created_tx_hash | 创建该拍卖的交易 hash。 |
| created_block_number | AuctionCreated 事件所在区块号。 |
| created_block_hash | AuctionCreated 事件所在区块 hash，为后续 reorg 预留。 |
| created_log_index | AuctionCreated 事件在交易 receipt 中的 log 序号。 |

### 3.4 last_event_* 字段

`last_event_*` 字段记录最近一次改变该 auction read model 的链上事件位置。

包括：

```text
last_event_name
last_event_tx_hash
last_event_block_number
last_event_block_hash
last_event_log_index
```

它回答的问题是：

```text
当前这条 auctions read model 最近一次是被哪条链上事件更新的？
```

可能的事件包括：

```text
AuctionCreated
BidPlaced
AuctionEnded
AuctionCancelled
```

示例：

```text
拍卖刚创建：
last_event_name = AuctionCreated

有人出价：
last_event_name = BidPlaced

拍卖结束：
last_event_name = AuctionEnded

拍卖取消：
last_event_name = AuctionCancelled
```

该设计有利于：

```text
排查数据来源
定位最后一次状态变更
后续实现 reorg 检测
后续实现数据修复
面试时解释事件驱动 read model
```

### 3.5 auctions 表索引设计

#### uk_auction_chain_contract_id

```text
chain_id + contract_address + auction_id
```

作用：

```text
保证同一条链、同一个合约内的 auction_id 唯一。
```

#### idx_auctions_chain_status_end

```text
chain_id + status + end_time
```

作用：

```text
支持按链、状态和结束时间查询拍卖列表。
例如查询 active 拍卖，并按 end_time 排序。
```

#### idx_auctions_seller

作用：

```text
支持查询某个用户创建的拍卖。
```

#### idx_auctions_nft_contract

作用：

```text
支持查询某个 NFT 合约下的拍卖。
```

#### idx_auctions_last_event_block_number

作用：

```text
支持按照最近一次事件所在区块进行排查、回滚或数据修复。
```

## 4. bids 表设计

### 4.1 表含义

`bids` 表保存出价历史。

一条 `BidPlaced` 链上事件对应 `bids` 表中的一条记录。

它是事件明细表，通常只插入，不更新。

### 4.2 字段说明

| 字段 | 含义 |
|---|---|
| id | 数据库内部自增主键。 |
| chain_id | 出价事件来自哪条链。 |
| contract_address | 出价事件来自哪个拍卖合约。 |
| auction_id | 这条出价属于哪个拍卖。 |
| bidder | 出价人钱包地址。 |
| bid_token | 出价使用的资产地址。零地址表示 ETH，非零地址表示 ERC20 token。 |
| amount | 原始出价金额。ETH 使用 wei，ERC20 使用 token 最小单位。 |
| amount_usd | 合约中折算后的 USD 金额，用于比较不同 token 的出价大小。 |
| tx_hash | 触发 BidPlaced 事件的交易 hash。 |
| log_index | 该 log 在交易 receipt 中的位置。 |
| block_number | 事件所在区块号，用于按链上发生顺序排序。 |
| block_hash | 事件所在区块 hash，后续 reorg 检查会用到。 |
| created_at | 数据库记录创建时间，不等于链上出价时间。 |

### 4.3 bids 表索引设计

#### uk_bid_log

```text
chain_id + contract_address + tx_hash + log_index
```

作用：

```text
保证同一条链上 log 不会重复写入 bids 表。
```

#### idx_bids_auction_order

```text
chain_id + contract_address + auction_id + block_number + log_index
```

作用：

```text
支持查询某个 auction 的出价记录，并按链上顺序排序。
```

#### idx_bids_bidder

作用：

```text
支持查询某个用户的所有出价记录。
```

#### idx_bids_bid_token

作用：

```text
支持按 ETH / ERC20 出价资产筛选。
```

#### idx_bids_block_number

作用：

```text
支持后续 reorg 回滚时按区块范围查找和删除。
```

## 5. processed_logs 表设计

### 5.1 表含义

`processed_logs` 表是 Indexer 的幂等控制表。

它不是业务表，不直接面向前端查询。

它的作用是记录：

```text
哪些链上 log 已经被 Indexer 成功处理过
```

### 5.2 为什么需要 processed_logs？

Indexer 可能重复扫描同一区块范围。

例如：

```text
第一次扫描 100 ~ 200 区块，处理成功。
服务重启。
为了安全回扫，第二次扫描 180 ~ 220 区块。
```

如果没有 `processed_logs`，同一条 `BidPlaced` 事件可能被重复写入。

有了 `processed_logs`，处理流程变为：

```text
1. 开启数据库事务
2. 尝试插入 processed_logs
3. 如果插入成功，说明该 log 第一次处理
4. 继续写业务表
5. 如果唯一键冲突，说明已经处理过，直接跳过
6. 提交事务
```

### 5.3 字段说明

| 字段 | 含义 |
|---|---|
| id | 数据库内部自增主键。 |
| chain_id | 已处理 log 来自哪条链。 |
| contract_address | 已处理 log 来自哪个合约。 |
| tx_hash | 已处理 log 所在交易 hash。 |
| log_index | 已处理 log 在交易 receipt 中的位置。 |
| block_number | 已处理 log 所在区块号。 |
| block_hash | 已处理 log 所在区块 hash，为 reorg 检测预留。 |
| event_name | 事件名称，例如 AuctionCreated、BidPlaced、AuctionEnded、AuctionCancelled。 |
| created_at | Indexer 成功处理该 log 并写入数据库的时间。 |

### 5.4 processed_logs 表索引设计

#### uk_processed_log

```text
chain_id + contract_address + tx_hash + log_index
```

作用：

```text
唯一标识一条已处理 log，防止重复处理。
```

#### idx_processed_logs_block_number

作用：

```text
后续做 reorg 回滚时，可以按区块号范围查找和删除。
```

#### idx_processed_logs_event_name

作用：

```text
方便调试和统计不同事件处理数量。
```

## 6. sync_cursors 表设计

### 6.1 表含义

`sync_cursors` 表记录 Indexer 的同步进度。

它回答的问题是：

```text
某条链上的某个合约已经同步到哪个区块？
```

### 6.2 字段说明

| 字段 | 含义 |
|---|---|
| id | 数据库内部自增主键。 |
| chain_id | 当前 cursor 属于哪条链。 |
| contract_address | 当前 cursor 属于哪个合约。 |
| last_processed_block | Indexer 已经成功处理完成的最后一个区块号。 |
| last_processed_block_hash | last_processed_block 对应的区块 hash，用于后续 reorg 检测。 |
| created_at | cursor 第一次创建时间。 |
| updated_at | cursor 最近一次推进时间。 |

### 6.3 sync_cursors 表索引设计

#### uk_cursor

```text
chain_id + contract_address
```

作用：

```text
一条链上的一个合约只有一条同步进度记录。
```

### 6.4 为什么 cursor 只记录到 block？

当前 Indexer 设计是按区块范围批量处理：

```text
fromBlock → toBlock
```

每个 batch 成功处理完后，才把 cursor 更新到：

```text
last_processed_block = toBlock
```

这种设计简单、可靠，适合当前阶段。

如果未来要做更细粒度的恢复，也可以演进为：

```text
block_number + log_index
```

但复杂度会明显增加。当前项目阶段，block 级 cursor 已经足够。

## 7. Indexer 处理流程

### 7.1 启动时读取 cursor

Indexer 启动后，先读取 `sync_cursors`。

如果 cursor 存在：

```text
fromBlock = last_processed_block + 1
```

如果 cursor 不存在：

```text
fromBlock = configured_start_block
```

### 7.2 计算 confirmed target block

为了避免处理最新不稳定区块，Indexer 使用 confirmations 机制：

```text
targetBlock = latestBlock - confirmations
```

只有达到确认数的区块才会被处理。

### 7.3 分批扫描 logs

Indexer 根据 batch size 分段扫描：

```text
fromBlock ~ toBlock
```

每一段处理完成后，推进 cursor。

### 7.4 每条 log 的处理顺序

每条 log 都按以下顺序处理：

```text
1. 开启数据库事务
2. 插入 processed_logs
3. 如果该 log 已经处理过，跳过
4. 如果是新 log，解析事件
5. 写入或更新 auctions / bids
6. 提交事务
```

关键原则：

```text
processed_logs 和业务表写入必须在同一个事务里。
```

否则可能出现：

```text
processed_logs 已经写入；
业务表写入失败；
下次重扫时因为 processed_logs 已存在；
业务数据永远不会补上。
```

因此，中高级 Indexer 写法必须保证：

```text
幂等记录 + 业务写入 = 同事务
```

## 8. 事件到数据表的映射

### 8.1 AuctionCreated

处理结果：

```text
插入 auctions
status = active
初始化 created_* 字段
初始化 last_event_* 字段
插入 processed_logs
```

### 8.2 BidPlaced

处理结果：

```text
插入 bids
更新 auctions.highest_bidder
更新 auctions.highest_bid_token
更新 auctions.highest_bid_amount
更新 auctions.highest_bid_usd
更新 auctions.last_event_* 字段
插入 processed_logs
```

### 8.3 AuctionEnded

处理结果：

```text
更新 auctions.status = ended
更新最终 winner / bidToken / amount / amountUsd
更新 auctions.last_event_* 字段
插入 processed_logs
```

### 8.4 AuctionCancelled

处理结果：

```text
更新 auctions.status = cancelled
更新 auctions.last_event_* 字段
插入 processed_logs
```

## 9. 查询 API 与数据模型关系

REST API 不直接读取链上合约，而是读取 MySQL read model。

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

这样做的优点是：

```text
查询速度快
支持分页和筛选
避免频繁 RPC 调用
前端体验更稳定
便于后续做缓存、搜索和统计
```

## 10. 后续 reorg 设计预留

当前模型已经为 reorg 预留字段：

```text
created_block_hash
last_event_block_hash
block_hash
last_processed_block_hash
```

后续可以实现：

```text
1. 每次同步前检查 last_processed_block_hash 是否仍然匹配链上 block hash
2. 如果不匹配，说明发生 reorg
3. 回退一定区块范围
4. 删除该范围内的 processed_logs / bids
5. 重建受影响的 auctions read model
6. 重新同步
```

第一阶段也可以先实现为：

```text
检测到 reorg 后停止 Indexer 并报警
```

## 11. 表与业务流程关系

完整业务流程如下：

```text
用户创建拍卖
    ↓
合约 emit AuctionCreated
    ↓
Indexer 插入 auctions
    ↓
REST API 可以查询拍卖列表

用户出价
    ↓
合约 emit BidPlaced
    ↓
Indexer 插入 bids，并更新 auctions 当前最高出价
    ↓
REST API 可以查询最新拍卖状态和出价历史

用户结束拍卖
    ↓
合约 emit AuctionEnded
    ↓
Indexer 更新 auctions.status = ended
    ↓
REST API 查询到拍卖已结束

用户取消拍卖
    ↓
合约 emit AuctionCancelled
    ↓
Indexer 更新 auctions.status = cancelled
    ↓
REST API 查询到拍卖已取消
```

## 12. 当前模型体现的工程能力

该数据模型体现了以下中高级 Web3 后端设计点：

```text
多链、多合约数据建模
链上 uint256 精度处理
event-driven read model
Indexer 幂等机制
断点续扫
confirmed block 同步
为 reorg 预留 block_hash
事件溯源字段 created_* / last_event_*
按查询模式设计索引
REST API 与链上事件解耦
```

## 13. 面试表述

可以这样介绍本项目的数据模型：

```text
本项目的数据模型不是传统 CRUD，而是基于链上事件构建 read model。auctions 表保存拍卖当前状态，bids 表保存 BidPlaced 事件历史，processed_logs 表用于保证 Indexer 幂等，sync_cursors 表用于断点续扫。为了支持多链和多合约，我在业务唯一键中加入 chain_id 和 contract_address。为了后续处理 reorg，我在事件表和 cursor 中保存 block_hash。同时，auctions 表区分了 created_* 和 last_event_* 字段，分别记录拍卖创建事件和最近一次改变 read model 的链上事件位置，从而提高数据可追溯性和可维护性。
```

## 14. 后续演进方向

后续可以在当前数据模型基础上继续增强：

```text
1. 实现 reorg 检测与回滚
2. 增加 repository 集成测试
3. 增加 Indexer handler 单元测试
4. 增加 REST API 查询层
5. 增加 Tx Service 交易参数构造
6. 增加钱包登录 SIWE / JWT
7. 增加 metrics、health、readiness
8. 增加 Docker Compose 一键启动和部署文档
```