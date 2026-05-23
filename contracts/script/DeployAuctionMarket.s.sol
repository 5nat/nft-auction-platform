// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.13;

import {Script} from "forge-std/Script.sol";
import {ERC1967Proxy} from "@openzeppelin/contracts/proxy/ERC1967/ERC1967Proxy.sol";

import {NFTAuctionMarket} from "../src/NFTAuctionMarket.sol";

contract DeployAuctionMarket is Script {
    function run() external returns (NFTAuctionMarket implementation, ERC1967Proxy proxy, NFTAuctionMarket market) {
        uint256 deployPrivateKey = vm.envUint("PRIVATE_KEY");
        address owner = vm.addr(deployPrivateKey);

        vm.startBroadcast(deployPrivateKey);

        implementation = new NFTAuctionMarket();

        bytes memory initData = abi.encodeCall(NFTAuctionMarket.initialize, (owner));

        proxy = new ERC1967Proxy(address(implementation), initData);

        market = NFTAuctionMarket(address(proxy));

        vm.stopBroadcast();
    }
}
