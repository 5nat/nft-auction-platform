# NFT Auction Platform
本项目是一个 Web3 NFT 拍卖系统，包含：

- contracts：Solidity 智能合约，使用 Foundry 开发、测试、部署
- backend：Go 后端服务，负责链上事件索引、数据库存储、REST API、交易发送
- deployments：记录不同网络的合约部署地址
- docs：架构设计、API 文档、索引器设计和面试复盘

项目目录
```aiignore
nft-auction-platform/
├── contracts/                         # 链上合约工程，Foundry 项目
│   ├── src/
│   │   ├── Auction.sol                # NFT 拍卖主合约
│   │   └── MockNFT.sol                # 本地测试用 ERC721 NFT 合约
│   ├── script/
│   │   └── Deploy.s.sol               # Foundry 部署脚本
│   ├── test/
│   │   └── Auction.t.sol              # Foundry 合约测试
│   ├── lib/                           # Foundry 依赖，如 OpenZeppelin
│   ├── out/                           # forge build 后自动生成
│   ├── cache/                         # forge 缓存
│   ├── foundry.toml
│   └── README.md
│
├── backend/                           # 链下 Go 后端工程
│   ├── cmd/
│   │   └── server/
│   │       └── main.go                # 后端启动入口
│   ├── configs/
│   │   └── config.example.yaml        # 示例配置
│   ├── internal/
│   │   ├── app/
│   │   │   └── app.go                 # 应用启动、依赖组装、优雅退出
│   │   ├── config/
│   │   │   └── config.go              # 配置加载
│   │   ├── api/
│   │   │   ├── routes.go              # Gin 路由注册
│   │   │   ├── health_handler.go
│   │   │   ├── auction_handler.go
│   │   │   └── bid_handler.go
│   │   ├── middleware/
│   │   │   ├── logger.go
│   │   │   └── recovery.go
│   │   ├── model/
│   │   │   ├── auction.go
│   │   │   ├── bid.go
│   │   │   ├── sync_cursor.go
│   │   │   └── processed_log.go
│   │   ├── store/
│   │   │   ├── db.go                  # GORM/MySQL 初始化
│   │   │   ├── auction_store.go
│   │   │   ├── bid_store.go
│   │   │   └── cursor_store.go
│   │   ├── service/
│   │   │   └── auction_service.go
│   │   ├── chain/
│   │   │   ├── client.go              # ethclient 初始化
│   │   │   ├── bindings/              # abigen 生成的 Go 合约绑定代码
│   │   │   │   ├── auction.go
│   │   │   │   └── mock_nft.go
│   │   │   └── tx.go                  # 交易发送、receipt 查询
│   │   ├── indexer/
│   │   │   ├── indexer.go             # 区块扫描主循环
│   │   │   ├── processor.go           # 事件解析处理
│   │   │   ├── events.go              # 事件结构定义
│   │   │   └── cursor.go              # 同步游标逻辑
│   │   └── metrics/
│   │       └── metrics.go             # 后续 Prometheus 指标
│   ├── migrations/
│   │   └── 001_init.sql               # 数据库表结构
│   ├── scripts/
│   │   ├── gen_bindings.sh            # 从合约 ABI 生成 Go binding
│   │   └── run_local.sh               # 本地启动脚本
│   ├── .env.example
│   ├── docker-compose.yml             # MySQL / Redis 等服务
│   ├── go.mod
│   ├── go.sum
│   └── README.md
│
├── deployments/                       # 合约部署结果记录
│   ├── anvil.json                     # 本地链部署地址
│   └── sepolia.json                   # Sepolia 部署地址
│
├── docs/                              # 项目文档
│   ├── architecture.md                # 架构说明
│   ├── api.md                         # API 文档
│   ├── indexer.md                     # 索引器设计
│   └── interview-notes.md             # 面试复盘笔记
│
├── scripts/                           # 根目录通用脚本
│   ├── dev.sh                         # 一键启动本地开发环境
│   └── clean.sh                       # 清理缓存
│
├── .gitignore
├── README.md                          # 总项目说明
└── Makefile                           # 常用命令封装
```

# 启动本地链
cd contracts
anvil

# 部署合约
forge script script/Deploy.s.sol --rpc-url http://127.0.0.1:8545 --broadcast

# 启动后端
cd ../backend
docker compose up -d
go run ./cmd/server


```aiignore
contracts/：
    专门管理 Solidity 合约、Foundry 测试、部署脚本。

backend/：
    专门管理 Go 后端、API、数据库、链上事件索引器。

deployments/：
    记录不同网络上的合约地址，后端可以读取这些地址。

docs/：
    记录架构、API、索引器设计，方便面试时展示。

scripts/：
    放一键启动、生成 binding、清理缓存等辅助脚本。
```