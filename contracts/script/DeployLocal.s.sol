// SPDX-License-Identifier: MIT
pragma solidity ^0.8.13;

import {Script} from "forge-std/Script.sol";
import {console2} from "forge-std/console2.sol";
import {ERC1967Proxy} from "@openzeppelin/contracts/proxy/ERC1967/ERC1967Proxy.sol";

import {AuctionNFT} from "../src/AuctionNFT.sol";
import {NFTAuctionMarket} from "../src/NFTAuctionMarket.sol";
import {MockERC20} from "../src/MockERC20.sol";
import {MockV3Aggregator} from "../src/MockV3Aggregator.sol";

contract DeployLocal is Script {
    function run()
        external
        returns (
            AuctionNFT nft,
            MockERC20 token,
            MockV3Aggregator ethUsdFeed,
            MockV3Aggregator tokenUsdFeed,
            NFTAuctionMarket implementation,
            ERC1967Proxy proxy,
            NFTAuctionMarket market
        )
    {
        uint256 deployerPrivateKey = vm.envUint("PRIVATE_KEY");
        address deployer = vm.addr(deployerPrivateKey);

        vm.startBroadcast(deployerPrivateKey);

        nft = new AuctionNFT();
        token = new MockERC20();

        ethUsdFeed = new MockV3Aggregator(8, 2000e8);
        tokenUsdFeed = new MockV3Aggregator(8, 1e8);

        implementation = new NFTAuctionMarket();

        bytes memory initData = abi.encodeCall(NFTAuctionMarket.initialize, (deployer));

        proxy = new ERC1967Proxy(address(implementation), initData);
        market = NFTAuctionMarket(address(proxy));

        market.setPriceFeed(address(0), address(ethUsdFeed));
        market.setPriceFeed(address(token), address(tokenUsdFeed));

        vm.stopBroadcast();

        console2.log("Deployer:", deployer);
        console2.log("AuctionNFT:", address(nft));
        console2.log("MockERC20:", address(token));
        console2.log("ETH/USD Feed:", address(ethUsdFeed));
        console2.log("Token/USD Feed:", address(tokenUsdFeed));
        console2.log("Market Implementation:", address(implementation));
        console2.log("Market Proxy:", address(proxy));
        console2.log("Use this market address:", address(market));
    }
}
