// SPDX-License-Identifier: MIT
pragma solidity ^0.8.13;

import {NFTAuctionMarket} from "./NFTAuctionMarket.sol";

contract NFTAuctionMarketV2 is NFTAuctionMarket {
    function version() external pure returns (string memory) {
        return "v2";
    }
}
