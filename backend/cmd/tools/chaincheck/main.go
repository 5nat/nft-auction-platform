package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/5nat/nft-auction-platform/backend/internal/config"
	ethinfra "github.com/5nat/nft-auction-platform/backend/internal/infra/blockchain/ethereum"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
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
		logger.Error("error loading config", "error", err)
		os.Exit(1)
	}

	chainClient, err := ethinfra.NewClient(ctx, cfg.Chain, logger)
	if err != nil {
		logger.Error("error creating chain client", "error", err)
		os.Exit(1)
	}
	defer chainClient.Close()

	latestBlock, err := chainClient.LatestBlockNumber(ctx)
	if err != nil {
		logger.Error("error getting latest block", "error", err)
		os.Exit(1)
	}

	logger.Info("connected to chain", "latest_block", latestBlock)

	if !common.IsHexAddress(cfg.Chain.AuctionContract) {
		logger.Error("invalid chain contract address", "address", cfg.Chain.AuctionContract)
		os.Exit(1)
	}

	auctionContract := common.HexToAddress(cfg.Chain.AuctionContract)

	nextAuctionID, err := callNextAuctionID(ctx, chainClient, auctionContract)
	if err != nil {
		logger.Error("error calling next auction ID", "error", err)
		os.Exit(1)
	}

	logger.Info(
		"contract call success",
		"auction_contract", auctionContract.Hex(),
		"next_auction_id", nextAuctionID,
	)
}

// callNextAuctionID 演示手动 ABI 调用合约 view 方法。 这里先不用 abigen，目的是理解底层 eth_call 是怎么工作的。
func callNextAuctionID(ctx context.Context, chainClient *ethinfra.Client, contract common.Address) (string, error) {
	abiFile, err := os.Open("abi/NFTAuctionMarket.abi.json")
	if err != nil {
		return "", fmt.Errorf("error opening ABI file: %w", err)
	}
	defer abiFile.Close()

	parsedABI, err := abi.JSON(abiFile)
	if err != nil {
		return "", fmt.Errorf("error parsing ABI: %w", err)
	}

	callData, err := parsedABI.Pack("nextAuctionId")
	if err != nil {
		return "", fmt.Errorf("pack nextAuctionId calldata: %w", err)
	}

	result, err := chainClient.EthClient().CallContract(ctx, ethereum.CallMsg{
		To:   &contract,
		Data: callData,
	}, nil)
	if err != nil {
		return "", fmt.Errorf("eth_call nextAuctionId: %w", err)
	}

	values, err := parsedABI.Unpack("nextAuctionId", result)
	if err != nil {
		return "", fmt.Errorf("unpack nextAuctionId result: %w", err)
	}

	if len(values) != 1 {
		return "", fmt.Errorf("unexpected nextAuctionId return values: %d", len(values))
	}

	return fmt.Sprintf("%v", values[0]), nil
}
