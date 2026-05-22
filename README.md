# NFT Auction Platform

这是一个用于学习和实践 Web3 后端开发的 NFT 拍卖系统，采用 monorepo 结构，同时包含链上合约工程和链下后端工程。

## 项目结构

```text
nft-auction-platform/
├── backend/       # Go 后端服务，负责 API、数据库、链上事件索引
├── contracts/     # Solidity 智能合约，使用 Foundry 开发和测试
├── deployments/   # 不同网络的合约部署地址记录
├── docs/          # 架构设计、API 文档、索引器设计和面试复盘
├── scripts/       # 根目录通用脚本
├── Makefile       # 常用命令封装
└── README.md
```

