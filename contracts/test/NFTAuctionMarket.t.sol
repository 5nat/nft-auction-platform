// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.13;

import {Test} from "forge-std/Test.sol";
import {AuctionNFT} from "../src/AuctionNFT.sol";
import {NFTAuctionMarket} from "../src/NFTAuctionMarket.sol";
import {MockERC20} from "../src/MockERC20.sol";
import {MockV3Aggregator} from "../src/MockV3Aggregator.sol";
import {ERC1967Proxy} from "@openzeppelin/contracts/proxy/ERC1967/ERC1967Proxy.sol";

contract NFTAuctionMarketTest is Test {
    AuctionNFT public nft;
    NFTAuctionMarket public market;
    MockERC20 public token;
    MockV3Aggregator public ethUsdFeed;
    MockV3Aggregator public tokenUsdFeed;

    address public seller = address(1);
    address public bidder1 = address(2);
    address public bidder2 = address(3);

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

    function setUp() public {
        nft = new AuctionNFT();
        // market = new NFTAuctionMarket();
        NFTAuctionMarket implementation = new NFTAuctionMarket();
        ERC1967Proxy proxy =
            new ERC1967Proxy(address(implementation), abi.encodeCall(NFTAuctionMarket.initialize, (address(this))));
        market = NFTAuctionMarket(address(proxy));
        token = new MockERC20();

        ethUsdFeed = new MockV3Aggregator(8, 2000e8);
        tokenUsdFeed = new MockV3Aggregator(8, 1e8);

        market.setPriceFeed(address(0), address(ethUsdFeed));
        market.setPriceFeed(address(token), address(tokenUsdFeed));

        nft.mint(seller, "ipfs://token-0");
        token.mint(bidder1, 10_000 ether);
        token.mint(bidder2, 10_000 ether);

        vm.deal(bidder1, 1_000 ether);
        vm.deal(bidder2, 1_000 ether);
    }

    function testCreateAuction() public {
        vm.startPrank(seller);
        nft.approve(address(market), 0);

        uint256 expectedEndTime = block.timestamp + 1 days;
        vm.expectEmit(true, true, true, true, address(market));
        emit AuctionCreated(0, seller, address(nft), 0, 100 ether, expectedEndTime);

        uint256 auctionId = market.createAuction(address(nft), 0, 100 ether, 1 days);
        vm.stopPrank();

        assertEq(auctionId, 0);
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
        ) = market.auctions(auctionId);

        assertEq(auctionSeller, seller);
        assertEq(auctionNft, address(nft));
        assertEq(tokenId, 0);
        assertEq(minBidUsd, 100 ether);
        assertEq(highestBid, 0);
        assertEq(highestBidder, address(0));
        assertEq(highestBidUsd, 0);
        assertEq(bidToken, address(0));
        assertGt(endTime, block.timestamp);
        assertFalse(ended);
        assertFalse(cancelled);
    }

    function testBidEth() public {
        _createDefaultAuction();

        vm.expectEmit(true, true, true, true, address(market));
        emit BidPlaced(0, bidder1, address(0), 1 ether, 2000 ether);

        vm.prank(bidder1);
        market.bidEth{value: 1 ether}(0);
        (,,,, uint256 highestBid, address highestBidder,, address bidToken,,,) = market.auctions(0);

        assertEq(highestBid, 1 ether);
        assertEq(highestBidder, bidder1);
        assertEq(bidToken, address(0));
    }

    function testHighestBidRefundsPreviousBidder() public {
        _createDefaultAuction();

        vm.prank(bidder1);
        market.bidEth{value: 1 ether}(0);

        uint256 bidder1BalanceBefore = bidder1.balance;

        vm.prank(bidder2);
        market.bidEth{value: 2 ether}(0);

        assertEq(bidder1.balance, bidder1BalanceBefore + 1 ether);

        (,,,, uint256 highestBid, address highestBidder,, address bidToken,,,) = market.auctions(0);

        assertEq(highestBid, 2 ether);
        assertEq(highestBidder, bidder2);
        assertEq(bidToken, address(0));
    }

    function testCancelAuction() public {
        _createDefaultAuction();

        vm.prank(seller);

        vm.expectEmit(true, false, false, false, address(market));
        emit AuctionCancelled(0);

        market.cancelAuction(0);

        // NFT 应该退回 seller。
        assertEq(nft.ownerOf(0), seller);

        (,,,,,,,,, bool ended, bool cancelled) = market.auctions(0);

        assertTrue(ended);
        assertTrue(cancelled);
    }

    function testCannotCancelAuctionWithBid() public {
        _createDefaultAuction();

        vm.prank(bidder1);
        market.bidEth{value: 1 ether}(0);

        vm.prank(seller);
        vm.expectRevert("already has bid");
        market.cancelAuction(0);
    }

    //    function testNonSellerCannotCancelAuction() public {
    //        _createDefaultAuction();
    //
    //        vm.prank(bidder1);
    //        vm.expectRevert("not seller");
    //        market.cancelAuction(0);
    //    }
    //
    //    function testPriceFeedStaleReverts() public {
    //        _createDefaultAuction();
    //
    //        market.setMaxPriceDelay(1 hours);
    //
    //        vm.warp(3 hours);
    //
    //        ethUsdFeed.updateAnswerWithTimestamp(
    //            2000e8,
    //            block.timestamp - 2 hours
    //        );
    //
    //        vm.prank(bidder1);
    //        vm.expectRevert("stale price");
    //        market.bidEth{value: 1 ether}(0);
    //    }

    function testEndAuctionWithWinner() public {
        _createDefaultAuction();

        vm.prank(bidder1);
        market.bidEth{value: 2 ether}(0);

        uint256 sellerBalanceBefore = seller.balance;

        vm.warp(block.timestamp + 1 days + 1);

        vm.expectEmit(true, true, true, true, address(market));
        emit AuctionEnded(0, bidder1, address(0), 2 ether, 4000 ether);

        market.endAuction(0);

        assertEq(nft.ownerOf(0), bidder1);
        assertEq(seller.balance, sellerBalanceBefore + 2 ether);

        (,,,,,,,,, bool ended, bool cancelled) = market.auctions(0);

        assertTrue(ended);
        assertFalse(cancelled);
    }

    function testEndAuctionWithoutBidReturnsNFTToSeller() public {
        _createDefaultAuction();

        vm.warp(block.timestamp + 1 days + 1);

        vm.expectEmit(true, true, true, true, address(market));
        emit AuctionEnded(0, address(0), address(0), 0, 0);

        market.endAuction(0);

        assertEq(nft.ownerOf(0), seller);

        (,,,,,,,,, bool ended, bool cancelled) = market.auctions(0);

        assertTrue(ended);
        assertFalse(cancelled);
    }

    function testCannotBidBelowMinUsd() public {
        _createDefaultAuction();

        vm.prank(bidder1);
        vm.expectRevert("bid below min usd");
        market.bidEth{value: 0.04 ether}(0);
    }

    function testCannotEndBeforeExpired() public {
        _createDefaultAuction();

        vm.expectRevert("auction not expired");
        market.endAuction(0);
    }

    function testBidERC20() public {
        _createDefaultAuction();

        vm.startPrank(bidder1);
        token.approve(address(market), 100 ether);
        market.bidERC20(0, address(token), 100 ether);
        vm.stopPrank();

        (,,,, uint256 highestBid, address highestBidder,, address bidToken,,,) = market.auctions(0);

        assertEq(highestBid, 100 ether);
        assertEq(highestBidder, bidder1);
        assertEq(bidToken, address(token));
        assertEq(token.balanceOf(address(market)), 100 ether);
    }

    function testHigherERC20BidRefundsPreviousERC20Bidder() public {
        _createDefaultAuction();

        vm.startPrank(bidder1);
        token.approve(address(market), 100 ether);
        market.bidERC20(0, address(token), 100 ether);
        vm.stopPrank();

        uint256 bidder1BalanceBefore = token.balanceOf(bidder1);

        vm.startPrank(bidder2);
        token.approve(address(market), 200 ether);
        market.bidERC20(0, address(token), 200 ether);
        vm.stopPrank();

        assertEq(token.balanceOf(bidder1), bidder1BalanceBefore + 100 ether);
        assertEq(token.balanceOf(address(market)), 200 ether);

        (,,,, uint256 highestBid, address highestBidder,, address bidToken,,,) = market.auctions(0);

        assertEq(highestBid, 200 ether);
        assertEq(highestBidder, bidder2);
        assertEq(bidToken, address(token));
    }

    function testERC20BidRefundsPreviousETHBidder() public {
        _createDefaultAuction();

        vm.prank(bidder1);
        market.bidEth{value: 1 ether}(0);

        uint256 bidder1BalanceBefore = bidder1.balance;

        vm.startPrank(bidder2);
        token.approve(address(market), 2_001 ether);
        market.bidERC20(0, address(token), 2_001 ether);
        vm.stopPrank();

        assertEq(bidder1.balance, bidder1BalanceBefore + 1 ether);

        (,,,, uint256 highestBid, address highestBidder,, address bidToken,,,) = market.auctions(0);

        assertEq(highestBid, 2_001 ether);
        assertEq(highestBidder, bidder2);
        assertEq(bidToken, address(token));
    }

    function testETHBidRefundsPreviousERC20Bidder() public {
        _createDefaultAuction();

        vm.startPrank(bidder1);
        token.approve(address(market), 100 ether);
        market.bidERC20(0, address(token), 100 ether);
        vm.stopPrank();

        uint256 bidder1TokenBalanceBefore = token.balanceOf(bidder1);

        vm.prank(bidder2);
        market.bidEth{value: 200 ether}(0);

        assertEq(token.balanceOf(bidder1), bidder1TokenBalanceBefore + 100 ether);

        (,,,, uint256 highestBid, address highestBidder,, address bidToken,,,) = market.auctions(0);

        assertEq(highestBid, 200 ether);
        assertEq(highestBidder, bidder2);
        assertEq(bidToken, address(0));
    }

    function testGetBidUsdValue() public view {
        uint256 ethUsdValue = market.getBidUsdValue(address(0), 1 ether);
        uint256 tokenUsdValue = market.getBidUsdValue(address(token), 100 ether);

        assertEq(ethUsdValue, 2000 ether);
        assertEq(tokenUsdValue, 100 ether);
    }

    function _createDefaultAuction() internal returns (uint256) {
        vm.startPrank(seller);
        nft.approve(address(market), 0);

        uint256 auctionId = market.createAuction(address(nft), 0, 100 ether, 1 days);

        vm.stopPrank();

        return auctionId;
    }
}
