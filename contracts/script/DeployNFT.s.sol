// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.13;

import {Script} from "forge-std/Script.sol";
import {AuctionNFT} from "../src/AuctionNFT.sol";

contract DeployNFT is Script {
    function run() external returns (AuctionNFT nft) {
        vm.startBroadcast();

        nft = new AuctionNFT();

        vm.stopBroadcast();
    }
}
