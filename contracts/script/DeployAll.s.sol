// SPDX-License-Identifier: MIT
pragma solidity ^0.8.22;

import {Script} from "forge-std/Script.sol";
import {console2} from "forge-std/console2.sol";
import {ERC1967Proxy} from "@openzeppelin/contracts/proxy/ERC1967/ERC1967Proxy.sol";

import {AuctionNFT} from "../src/AuctionNFT.sol";
import {NFTAuctionMarket} from "../src/NFTAuctionMarket.sol";

contract DeployAll is Script {
    function run()
        external
        returns (AuctionNFT nft, NFTAuctionMarket implementation, ERC1967Proxy proxy, NFTAuctionMarket market)
    {
        uint256 deployerPrivateKey = vm.envOr("PRIVATE_KEY", uint256(0));
        address finalOwner = vm.envOr("OWNER", address(0));
        address ethUsdPriceFeed = vm.envOr("ETH_USD_PRICE_FEED", address(0));
        address erc20Token = vm.envOr("ERC20_TOKEN", address(0));
        address erc20UsdPriceFeed = vm.envOr("ERC20_USD_PRICE_FEED", address(0));

        address deployer;
        if (deployerPrivateKey == 0) {
            deployer = msg.sender;
            if (finalOwner == address(0)) {
                finalOwner = deployer;
            }
            vm.startBroadcast();
        } else {
            deployer = vm.addr(deployerPrivateKey);
            if (finalOwner == address(0)) {
                finalOwner = deployer;
            }
            vm.startBroadcast(deployerPrivateKey);
        }

        nft = new AuctionNFT();

        implementation = new NFTAuctionMarket();

        bytes memory initData = abi.encodeCall(NFTAuctionMarket.initialize, (deployer));

        proxy = new ERC1967Proxy(address(implementation), initData);
        market = NFTAuctionMarket(address(proxy));

        if (ethUsdPriceFeed != address(0)) {
            market.setPriceFeed(address(0), ethUsdPriceFeed);
        }

        if (erc20Token != address(0) && erc20UsdPriceFeed != address(0)) {
            market.setPriceFeed(erc20Token, erc20UsdPriceFeed);
        }

        if (finalOwner != deployer) {
            market.transferOwnership(finalOwner);
        }

        vm.stopBroadcast();

        console2.log("Deployer:", deployer);
        console2.log("Owner:", finalOwner);
        console2.log("AuctionNFT:", address(nft));
        console2.log("NFTAuctionMarket implementation:", address(implementation));
        console2.log("NFTAuctionMarket proxy:", address(proxy));
        console2.log("Use this market address:", address(market));
        console2.log("ETH/USD price feed:", ethUsdPriceFeed);
        console2.log("ERC20 token:", erc20Token);
        console2.log("ERC20/USD price feed:", erc20UsdPriceFeed);
    }
}
