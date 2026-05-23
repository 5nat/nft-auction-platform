// SPDX-License-Identifier: MIT
pragma solidity ^0.8.13;

import {Test} from "forge-std/Test.sol";
import {ERC1967Proxy} from "@openzeppelin/contracts/proxy/ERC1967/ERC1967Proxy.sol";

import {NFTAuctionMarket} from "../src/NFTAuctionMarket.sol";
import {NFTAuctionMarketV2} from "../src/NFTAuctionMarketV2.sol";

import {AuctionNFT} from "../src/AuctionNFT.sol";

// 代理部署成功
// owner 能升级
// 非 owner 不能升级
// 升级后状态不丢

contract NFTAuctionUpgradeTest is Test {
    NFTAuctionMarket public market;
    AuctionNFT public nft;
    address public seller = address(2);

    address public owner = address(this);
    address public notOwner = address(1);

    function setUp() public {
        NFTAuctionMarket implementation = new NFTAuctionMarket();

        ERC1967Proxy proxy =
            new ERC1967Proxy(address(implementation), abi.encodeCall(NFTAuctionMarket.initialize, (owner)));

        market = NFTAuctionMarket(address(proxy));

        nft = new AuctionNFT();
        nft.mint(seller, "ipfs://token-0");
    }

    function testOwnerCanUpgradeToV2() public {
        NFTAuctionMarketV2 newImplementation = new NFTAuctionMarketV2();

        market.upgradeToAndCall(address(newImplementation), "");

        NFTAuctionMarketV2 upgradedMarket = NFTAuctionMarketV2(address(market));

        assertEq(upgradedMarket.version(), "v2");
    }

    function testNonOwnerCannotUpgrade() public {
        NFTAuctionMarketV2 newImplementation = new NFTAuctionMarketV2();

        vm.prank(notOwner);
        vm.expectRevert();

        market.upgradeToAndCall(address(newImplementation), "");
    }

    // 1. 升级前创建的 auctionId 还存在
    // 2. NFT 仍然托管在 proxy 地址
    // 3. 升级后可以使用 V2 的新函数
    function testUpgradeKeepsAuctionState() public {
        vm.startPrank(seller);
        nft.approve(address(market), 0);

        uint256 auctionId = market.createAuction(address(nft), 0, 1 ether, 1 days);

        vm.stopPrank();

        NFTAuctionMarketV2 newImplementation = new NFTAuctionMarketV2();

        market.upgradeToAndCall(address(newImplementation), "");

        NFTAuctionMarketV2 upgradedMarket = NFTAuctionMarketV2(address(market));

        assertEq(upgradedMarket.version(), "v2");
        assertEq(nft.ownerOf(0), address(market));

        (
            address auctionSeller,
            address auctionNft,
            uint256 tokenId,
            uint256 minBidUsd,
            uint256 highestBid,
            address highestBidder,
            uint256 highestBidUsd,
            address bidToken,
            uint256 endTime,
            bool ended,
            bool cancelled
        ) = upgradedMarket.auctions(auctionId);

        assertEq(auctionSeller, seller);
        assertEq(auctionNft, address(nft));
        assertEq(tokenId, 0);
        assertEq(minBidUsd, 1 ether);
        assertEq(highestBid, 0);
        assertEq(highestBidder, address(0));
        assertEq(highestBidUsd, 0);
        assertEq(bidToken, address(0));
        assertGt(endTime, block.timestamp);
        assertFalse(ended);
        assertFalse(cancelled);
    }
}
