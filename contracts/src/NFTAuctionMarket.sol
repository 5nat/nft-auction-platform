// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.13;

import {IERC721} from "@openzeppelin/contracts/token/ERC721/IERC721.sol";
import {ReentrancyGuard} from "@openzeppelin/contracts/utils/ReentrancyGuard.sol";
//import {ReentrancyGuardUpgradeable} from "@openzeppelin/contracts/utils/ReentrancyGuardUpgradeable.sol";
import {IERC20} from "@openzeppelin/contracts/token/ERC20/IERC20.sol";
import {SafeERC20} from "@openzeppelin/contracts/token/ERC20/utils/SafeERC20.sol";
// 接入 Chainlink 价格预言机
import {IERC20Metadata} from "@openzeppelin/contracts/token/ERC20/extensions/IERC20Metadata.sol";
import {AggregatorV3Interface} from "@chainlink/contracts/src/v0.8/shared/interfaces/AggregatorV3Interface.sol";
// import { Ownable } from "@openzeppelin/contracts/access/Ownable.sol";

// UUPS 升级
// Initializable 主要解决的是：可升级合约没有正常 constructor 初始化的问题
// UUPSUpgradeable 提供升级逻辑。也就是说，它让你的实现合约具备这个函数：upgradeToAndCall(address newImplementation, bytes memory data) 调用它可以让 proxy 指向新的 implementation。
// 普通 Ownable 依赖 constructor 设置 owner：constructor(address initialOwner) . 但代理模式不能依赖 constructor，所以要换成可升级版本：OwnableUpgradeable
// ERC1967Proxy 是什么？这个是在测试和部署脚本里用的，不是业务合约继承的。它是真正的代理合约。
import {Initializable} from "@openzeppelin/contracts-upgradeable/proxy/utils/Initializable.sol";
import {UUPSUpgradeable} from "@openzeppelin/contracts-upgradeable/proxy/utils/UUPSUpgradeable.sol";
import {OwnableUpgradeable} from "@openzeppelin/contracts-upgradeable/access/OwnableUpgradeable.sol";
// Initializable        让代理合约可以安全初始化一次
// OwnableUpgradeable   让 owner 写到代理存储里
// UUPSUpgradeable      给合约提供升级能力
// ERC1967Proxy         真正保存状态、接收用户调用、转发到实现合约

contract NFTAuctionMarket is Initializable, ReentrancyGuard, OwnableUpgradeable, UUPSUpgradeable {
    using SafeERC20 for IERC20;
    address public constant ETH_TOKEN = address(0);
    uint256 public constant USD_PRECISION = 1e18;

    struct Auction {
        address seller; // 拍卖发起人，也就是 NFT 原持有人
        address nft; // 被拍卖 NFT 所属的 NFT 合约地址
        uint256 tokenId; // 被拍卖 NFT 的编号
        uint256 minBidUsd; // 起拍价。目前是原始数量单位
        uint256 highestBid; // 当前最高出价的原始数量。如果是 ETH，就是 wei；如果是 ERC20，就是 token 最小单位。
        address highestBidder; // 当前最高出价人
        uint256 highestBidUsd; // 当前最高出价折算成 USD 后的价值
        address bidToken; // 当前最高出价使用的资产。address(0) 表示 ETH；非零地址表示 ERC20 token。
        uint256 endTime; //  拍卖结束时间
        bool ended; // 拍卖是否已经结束
        bool cancelled; // 拍卖已取消
    }

    uint256 public nextAuctionId;
    mapping(uint256 => Auction) public auctions;

    //  uint256 public maxPriceDelay;
    mapping(address => address) public priceFeeds;

    event AuctionCreated(
        uint256 indexed auctionId,
        address indexed seller,
        address indexed nft,
        uint256 tokenId,
        uint256 minBidUsd,
        uint256 endTime
    );

    event BidPlaced(
        uint256 indexed auctionId, address indexed bidder, address indexed bidToken, uint256 amount, uint256 amountUsd
    );

    event AuctionEnded(
        uint256 indexed auctionId, address indexed winner, address indexed bidToken, uint256 amount, uint256 amountUsd
    );

    event AuctionCancelled(uint256 indexed auctionId);

    event PriceFeedSet(address indexed token, address indexed priceFeed);

    //  event MaxPriceDelayUpdated(
    //    uint256 maxPriceDelay
    //  );

    // constructor() Ownable(msg.sender) {}
    // -----  UUPS 升级逻辑 -------
    constructor() {
        // _disableInitializers() 可以锁住 implementation，防止别人直接调用 initialize()。
        _disableInitializers();
    }

    // initialize() 才是真正给代理合约初始化 owner 的函数。
    function initialize(address initialOwner) public initializer {
        __Ownable_init(initialOwner);
        //  __UUPSUpgradeable_init();
        //  __ReentrancyGuard_init();
        //  maxPriceDelay = 1 days;
    }
    // _authorizeUpgrade() 决定谁能升级，这里限制为 onlyOwner。
    function _authorizeUpgrade(address newImplementation) internal override onlyOwner {}

    function createAuction(address nft, uint256 tokenId, uint256 minBidUsd, uint256 duration)
        external
        returns (uint256)
    {
        require(nft != address(0), "nft is zero");
        require(duration > 0, "duration is zero");
        require(minBidUsd > 0, "min bid usd is zero");
        require(IERC721(nft).ownerOf(tokenId) == msg.sender, "not token owner");

        // 创建拍卖前，seller 必须先 approve market 合约, 否则 transferFrom 会失败。
        // 这一步表示 NFT 被托管到拍卖市场合约中。
        IERC721(nft).transferFrom(msg.sender, address(this), tokenId);

        uint256 auctionId = nextAuctionId;
        nextAuctionId++;

        auctions[auctionId] = Auction({
            seller: msg.sender,
            nft: nft,
            tokenId: tokenId,
            minBidUsd: minBidUsd,
            highestBid: 0,
            highestBidder: address(0),
            highestBidUsd: 0,
            bidToken: ETH_TOKEN,
            endTime: block.timestamp + duration,
            ended: false,
            cancelled: false
        });

        emit AuctionCreated(auctionId, msg.sender, nft, tokenId, minBidUsd, block.timestamp + duration);

        return auctionId;
    }

    function bidEth(uint256 auctionId) external payable nonReentrant {
        require(msg.value > 0, "value is zero");

        Auction storage auction = auctions[auctionId];

        _validateAuctionForBid(auction);

        //  require(msg.sender != auction.seller, "seller cannot bid");

        uint256 bidUsd = getBidUsdValue(ETH_TOKEN, msg.value);

        _placeBid(auctionId, auction, ETH_TOKEN, msg.value, bidUsd, msg.sender);
    }

    function bidERC20(uint256 auctionId, address token, uint256 amount) external nonReentrant {
        require(token != address(0), "token is zero");
        require(amount > 0, "amount is zero");

        Auction storage auction = auctions[auctionId];

        _validateAuctionForBid(auction);

        //  require(msg.sender != auction.seller, "seller cannot bid");

        uint256 bidUsd = getBidUsdValue(token, amount);

        _placeBid(auctionId, auction, token, amount, bidUsd, msg.sender);

        IERC20(token).safeTransferFrom(msg.sender, address(this), amount);
    }

    function endAuction(uint256 auctionId) external nonReentrant {
        Auction storage auction = auctions[auctionId];

        require(auction.seller != address(0), "auction not found");
        require(!auction.ended, "already ended");
        require(!auction.cancelled, "auction cancelled");
        require(block.timestamp >= auction.endTime, "auction not expired");

        auction.ended = true;

        if (auction.highestBidder == address(0)) {
            IERC721(auction.nft).transferFrom(address(this), auction.seller, auction.tokenId);
            emit AuctionEnded(auctionId, address(0), ETH_TOKEN, 0, 0);
            return;
        }

        IERC721(auction.nft).transferFrom(address(this), auction.highestBidder, auction.tokenId);

        if (auction.bidToken == ETH_TOKEN) {
            (bool paid,) = auction.seller.call{value: auction.highestBid}("");
            require(paid, "pay seller failed");
        } else {
            IERC20(auction.bidToken).safeTransfer(auction.seller, auction.highestBid);
        }

        emit AuctionEnded(auctionId, auction.highestBidder, auction.bidToken, auction.highestBid, auction.highestBidUsd);
    }

    function cancelAuction(uint256 auctionId) external nonReentrant {
        Auction storage auction = auctions[auctionId];

        require(auction.seller != address(0), "auction not found");
        require(msg.sender == auction.seller, "not seller");
        require(!auction.ended, "already ended");
        require(!auction.cancelled, "already cancelled");
        require(auction.highestBidder == address(0), "already has bid");

        auction.cancelled = true;
        auction.ended = true;

        IERC721(auction.nft).transferFrom(address(this), auction.seller, auction.tokenId);

        emit AuctionCancelled(auctionId);
    }

    function setPriceFeed(address token, address priceFeed) external onlyOwner {
        require(priceFeed != address(0), "price feed is zero");
        priceFeeds[token] = priceFeed;
        emit PriceFeedSet(token, priceFeed);
    }

    //  function setMaxPriceDelay(uint256 newMaxPriceDelay) external onlyOwner {
    //    require(newMaxPriceDelay > 0, "delay is zero");
    //
    //    maxPriceDelay = newMaxPriceDelay;
    //
    //    emit MaxPriceDelayUpdated(newMaxPriceDelay);
    //  }

    function getBidUsdValue(address token, uint256 amount) public view returns (uint256) {
        address priceFeed = priceFeeds[token];
        require(priceFeed != address(0), "price feed not set");

        (uint80 roundId, int256 price,, uint256 updatedAt, uint80 answeredInRound) =
            AggregatorV3Interface(priceFeed).latestRoundData();
        require(price > 0, "invalid price");
        require(updatedAt > 0, "round not complete");
        require(answeredInRound >= roundId, "stale round");
        require(block.timestamp >= updatedAt, "invalid updatedAt");
        //    require(block.timestamp - updatedAt <= maxPriceDelay, "stale price");

        uint8 feedDecimals = AggregatorV3Interface(priceFeed).decimals();

        uint8 tokenDecimals = 18;
        if (token != ETH_TOKEN) {
            tokenDecimals = IERC20Metadata(token).decimals();
        }
        // price is safe to cast because it is checked to be greater than zero above.
        // forge-lint: disable-next-line(unsafe-typecast)
        uint256 unsignedPrice = uint256(price);
        return amount * unsignedPrice * USD_PRECISION / (10 ** tokenDecimals) / (10 ** feedDecimals);
    }

    function _validateAuctionForBid(Auction storage auction) internal view {
        require(auction.seller != address(0), "auction not found");
        require(!auction.ended, "auction ended");
        require(!auction.cancelled, "auction cancelled");
        require(block.timestamp < auction.endTime, "auction expired");
    }

    function _placeBid(
        uint256 auctionId,
        Auction storage auction,
        address bidToken,
        uint256 amount,
        uint256 amountUsd,
        address bidder
    ) internal {
        require(amountUsd >= auction.minBidUsd, "bid below min usd");
        require(amountUsd > auction.highestBidUsd, "bid not high enough");

        address previousBidder = auction.highestBidder;
        uint256 previousBid = auction.highestBid;
        address previousBidToken = auction.bidToken;

        auction.highestBidder = bidder;
        auction.highestBid = amount;
        auction.bidToken = bidToken;
        auction.highestBidUsd = amountUsd;

        if (previousBidder != address(0)) {
            _refundPreviousBid(previousBidder, previousBidToken, previousBid);
        }

        emit BidPlaced(auctionId, bidder, bidToken, amount, amountUsd);
    }

    function _refundPreviousBid(address previousBidder, address previousBidToken, uint256 previousBid) internal {
        if (previousBidToken == ETH_TOKEN) {
            (bool refunded,) = previousBidder.call{value: previousBid}("");
            require(refunded, "refund failed");
        } else {
            IERC20(previousBidToken).safeTransfer(previousBidder, previousBid);
        }
    }
}
