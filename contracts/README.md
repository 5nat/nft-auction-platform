## Foundry

**Foundry is a blazing fast, portable and modular toolkit for Ethereum application development written in Rust.**

Foundry consists of:

- **Forge**: Ethereum testing framework (like Truffle, Hardhat and DappTools).
- **Cast**: Swiss army knife for interacting with EVM smart contracts, sending transactions and getting chain data.
- **Anvil**: Local Ethereum node, akin to Ganache, Hardhat Network.
- **Chisel**: Fast, utilitarian, and verbose solidity REPL.

## Documentation

https://book.getfoundry.sh/

## Usage

### Build

```shell
$ forge build
```

### Test

```shell
$ forge test
```

### Format

```shell
$ forge fmt
```

### Gas Snapshots

```shell
$ forge snapshot
```

### Anvil

```shell
$ anvil
```

### Deploy

```shell
$ forge script script/Counter.s.sol:CounterScript --rpc-url <your_rpc_url> --private-key <your_private_key>
```

### Cast

```shell
$ cast <subcommand>
```

### Help

```shell
$ forge --help
$ anvil --help
$ cast --help
```



## 部署&演示

### 一：部署前准备

```shell
export PRIVATE_KEY=0x部署者测试钱包私钥
export SEPOLIA_RPC_URL=你的Sepolia RPC
export ETH_USD_PRICE_FEED=0x694AA1769357215DE4FAC081bf1f309aDC325306
```

检查：

```shell
echo $PRIVATE_KEY
echo $SEPOLIA_RPC_URL
echo $ETH_USD_PRICE_FEED
```

然后确认本地代码没问题：

```shell
forge test
```

### 二：部署到 Sepolia

运行：

```shell
forge script script/DeployAll.s.sol:DeployAll \
  --rpc-url $SEPOLIA_RPC_URL \
  --broadcast
```

部署成功后，终端会打印类似：

```shell
AuctionNFT: 0x...
NFTAuctionMarket implementation: 0x...
NFTAuctionMarket proxy: 0x...
Use this market address: 0x...
```

保存两个地址：

```shell
export NFT=0x你的AuctionNFT地址
export MARKET=0x你的Use_this_market_address地址
```

注意：

```shell
MARKET 用 proxy 地址，也就是 Use this market address。
不要用 implementation 地址。
```

### 三、准备卖家和买家

准备两个测试钱包，都要有 Sepolia ETH。

```shell
export SELLER=0x卖家地址
export SELLER_PRIVATE_KEY=0x卖家私钥

export BIDDER=0x买家地址
export BIDDER_PRIVATE_KEY=0x买家私钥
不要用 implementation 地址。
```

### 四、演示业务流程

部署者 mint NFT 给卖家：

```shell
cast send $NFT \
 "mint(address,string)" \
 $SELLER \
 "ipfs://demo-token-0" \
 --private-key $PRIVATE_KEY \
 --rpc-url $SEPOLIA_RPC_URL
```

查看 NFT owner，确认是卖家：

```shell
cast call $NFT \
 "ownerOf(uint256)(address)" \
 0 \
 --rpc-url $SEPOLIA_RPC_URL
```

卖家授权市场合约转移 NFT：

```shell
cast send $NFT \
 "approve(address,uint256)" \
 $MARKET \
 0 \
 --private-key $SELLER_PRIVATE_KEY \
 --rpc-url $SEPOLIA_RPC_URL
```

卖家创建拍卖，起拍价 0.001 ETH，持续 60 秒：

```shell
cast send $MARKET \
 "createAuction(address,uint256,uint256,uint256)" \
 $NFT \
 0 \
 1000000000000000 \
 60 \
 --private-key $SELLER_PRIVATE_KEY \
 --rpc-url $SEPOLIA_RPC_URL
```

第一场拍卖的 auctionId 是 0。

买家出价 0.002 ETH：

```shell
cast send $MARKET \
 "bidEth(uint256)" \
 0 \
 --value 0.002ether \
 --private-key $BIDDER_PRIVATE_KEY \
 --rpc-url $SEPOLIA_RPC_URL
```

等 60 秒后结束拍卖：

```shell
cast send $MARKET \
 "endAuction(uint256)" \
 0 \
 --private-key $SELLER_PRIVATE_KEY \
 --rpc-url $SEPOLIA_RPC_URL
```

再查 NFT owner：

```shell
cast call $NFT \
 "ownerOf(uint256)(address)" \
 2 \
 --rpc-url $SEPOLIA_RPC_URL
```

如果返回的是 BIDDER 地址，说明流程成功：
