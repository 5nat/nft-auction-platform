// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.13;

import {Test} from "forge-std/Test.sol";
import {Ownable} from "@openzeppelin/contracts/access/Ownable.sol";
import {AuctionNFT} from "../src/AuctionNFT.sol";

contract AuctionNFTTest is Test {
    AuctionNFT public nft;

    address public nathan = address(1);

    function setUp() public {
        nft = new AuctionNFT();
    }

    function test_OwnerCanMint() public {
        uint256 tokenId = nft.mint(nathan, "ipfs://test-token-uri");

        assertEq(tokenId, 0);
        assertEq(nft.ownerOf(tokenId), nathan);
        assertEq(nft.tokenURI(tokenId), "ipfs://test-token-uri");
    }

    function test_NonOwnerCannotMint() public {
        assertEq(nft.owner(), address(this));
        assertTrue(nathan != nft.owner());
        vm.prank(nathan);

        vm.expectRevert(abi.encodeWithSelector(Ownable.OwnableUnauthorizedAccount.selector, nathan));

        nft.mint(nathan, "ipfs://test-token-uri");
    }

    function test_MintTokenIdIncreases() public {
        uint256 firstTokenId = nft.mint(nathan, "ipfs://token-1");
        uint256 secondTokenId = nft.mint(nathan, "ipfs://token-2");

        assertEq(firstTokenId, 0);
        assertEq(secondTokenId, 1);
        assertEq(nft.ownerOf(firstTokenId), nathan);
        assertEq(nft.ownerOf(secondTokenId), nathan);
    }
}
