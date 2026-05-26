package main

import (
	"context"
	"log/slog"
	"os"
	"time"

	"github.com/5nat/nft-auction-platform/backend/internal/config"
	ethinfra "github.com/5nat/nft-auction-platform/backend/internal/infra/blockchain/ethereum"
	"github.com/5nat/nft-auction-platform/backend/internal/infra/blockchain/ethereum/bindings"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cfg, err := config.Load()
	if err != nil {
		logger.Error("failed to load config", "error", err)
		os.Exit(1)
	}

	chainClient, err := ethinfra.NewClient(ctx, cfg.Chain, logger)
	if err != nil {
		logger.Error("failed to create chain client", "error", err)
		os.Exit(1)
	}
	defer chainClient.Close()

	auctionAddress := common.HexToAddress(cfg.Chain.AuctionContract)
	market, err := bindings.NewNFTAuctionMarket(auctionAddress, chainClient.EthClient())
	if err != nil {
		logger.Error("failed to create market client", "error", err)
		os.Exit(1)
	}

	nextAuctionID, err := market.NextAuctionId(&bind.CallOpts{
		Context: ctx,
	})
	if err != nil {
		logger.Error("failed to get next auction id", "error", err)
		os.Exit(1)
	}

	ethToken := common.Address{}
	ethFeed, err := market.PriceFeeds(&bind.CallOpts{
		Context: ctx,
	}, ethToken)
	if err != nil {
		logger.Error("failed to create market priceFeeds", "error", err)
		os.Exit(1)
	}

	logger.Info(
		"binding call success",
		"auction_contract", auctionAddress.Hex(),
		"next_auction_id", nextAuctionID.String(),
		//"max_price_delay", maxPriceDelay.String(),
		"eth_usd_feed", ethFeed.Hex(),
	)
}
