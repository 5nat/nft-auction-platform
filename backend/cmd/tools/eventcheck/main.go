package main

import (
	"context"
	"log/slog"
	"math/big"
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

	latestBlock, err := chainClient.LatestBlockNumber(ctx)
	if err != nil {
		logger.Error("failed to fetch latest block", "error", err)
		os.Exit(1)
	}

	auctionAddress := common.HexToAddress(cfg.Chain.AuctionContract)

	market, err := bindings.NewNFTAuctionMarket(
		auctionAddress,
		chainClient.EthClient(),
	)
	if err != nil {
		logger.Error("failed to create market client", "error", err)
		os.Exit(1)
	}

	startBlock := uint64(cfg.Chain.StartBlock)
	endBlock := latestBlock

	logger.Info(
		"filter auction created events",
		"contract", auctionAddress.Hex(),
		"from_block", startBlock,
		"to_block", endBlock,
	)

	iter, err := market.FilterAuctionCreated(&bind.FilterOpts{
		Start:   startBlock,
		End:     &endBlock,
		Context: ctx,
	}, nil, nil, nil)
	if err != nil {
		logger.Error("failed to filter auction created", "error", err)
		os.Exit(1)
	}
	defer iter.Close()

	count := 0

	for iter.Next() {
		event := iter.Event
		count++

		logger.Info(
			"auction created event found",
			"auction_id", event.AuctionId.String(),
			"seller", event.Seller.Hex(),
			"nft", event.Nft.Hex(),
			"token_id", event.TokenId.String(),
			"min_bid_usd", formatWei(event.MinBidUsd),
			"end_time", event.EndTime.String(),
			"block_number", event.Raw.BlockNumber,
			"tx_hash", event.Raw.TxHash.Hex(),
			"log_index", event.Raw.Index,
		)
	}

	if err := iter.Error(); err != nil {
		logger.Error("failed to iterate auction created", "error", err)
		os.Exit(1)
	}

	logger.Info("event check completed", "auction_created_count", count)
}

func formatWei(v *big.Int) string {
	if v == nil {
		return "0"
	}

	base := new(big.Int).Exp(big.NewInt(10), big.NewInt(18), nil)

	intPart := new(big.Int).Div(v, base)
	fracPart := new(big.Int).Mod(v, base)

	if fracPart.Sign() == 0 {
		return intPart.String()
	}

	return v.String()
}
