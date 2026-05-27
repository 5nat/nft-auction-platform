package ethereum

import (
	"context"
	"encoding/hex"
	"fmt"
	"math/big"
	"strings"

	"github.com/5nat/nft-auction-platform/backend/internal/infra/blockchain/ethereum/bindings"
	"github.com/5nat/nft-auction-platform/backend/internal/modules/tx"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
)

const erc721ApprovalABI = `[
  {
    "type": "function",
    "name": "approve",
    "inputs": [
      {"name": "to", "type": "address"},
      {"name": "tokenId", "type": "uint256"}
    ],
    "outputs": [],
    "stateMutability": "nonpayable"
  }
]`

type TxCalldataBuilder struct {
	auctionABI *abi.ABI
	erc721ABI  *abi.ABI
}

var _ tx.CalldataBuilder = (*TxCalldataBuilder)(nil)

func (b *TxCalldataBuilder) BuildApproveNFTCalldata(ctx context.Context, req tx.BuildApproveNFTTxRequest) (string, error) {
	_ = ctx

	tokenID, err := parseUint256(req.TokenID)
	if err != nil {
		return "", fmt.Errorf("parse token_id: %w", err)
	}

	data, err := b.erc721ABI.Pack(
		"approve",
		common.HexToAddress(req.Operator),
		tokenID,
	)
	if err != nil {
		return "", fmt.Errorf("pack approve calldata: %w", err)
	}

	return "0x" + hex.EncodeToString(data), nil
}

func NewTxCalldataBuilder() (*TxCalldataBuilder, error) {
	auctionABI, err := bindings.NFTAuctionMarketMetaData.GetAbi()
	if err != nil {
		return nil, fmt.Errorf("get auction market abi: %w", err)
	}

	erc721ABI, err := abi.JSON(strings.NewReader(erc721ApprovalABI))
	if err != nil {
		return nil, fmt.Errorf("parse erc721 approval abi: %w", err)
	}

	return &TxCalldataBuilder{
		auctionABI: auctionABI,
		erc721ABI:  &erc721ABI,
	}, nil
}

func (b *TxCalldataBuilder) BuildCreateAuctionCalldata(ctx context.Context, req tx.BuildCreateAuctionTxRequest) (string, error) {
	_ = ctx

	tokenID, err := parseUint256(req.TokenID)
	if err != nil {
		return "", fmt.Errorf("parse token_id: %w", err)
	}

	minBidUSD, err := parseUint256(req.MinBidUSD)
	if err != nil {
		return "", fmt.Errorf("parse min_bid_usd: %w", err)
	}

	duration := new(big.Int).SetUint64(req.Duration)

	data, err := b.auctionABI.Pack(
		"createAuction",
		common.HexToAddress(req.NFTContract),
		tokenID,
		minBidUSD,
		duration,
	)
	if err != nil {
		return "", fmt.Errorf("pack createAuction calldata: %w", err)
	}

	return "0x" + hex.EncodeToString(data), nil
}

func (b *TxCalldataBuilder) BuildPlaceBidCalldata(ctx context.Context, req tx.BuildPlaceBidTxRequest) (string, error) {
	_ = ctx

	auctionID := new(big.Int).SetUint64(req.AuctionID)

	if common.HexToAddress(req.BidToken) == (common.Address{}) {
		data, err := b.auctionABI.Pack("bidEth", auctionID)
		if err != nil {
			return "", fmt.Errorf("pack bidEth calldata: %w", err)
		}

		return "0x" + hex.EncodeToString(data), nil
	}

	amount, err := parseUint256(req.Amount)
	if err != nil {
		return "", fmt.Errorf("parse amount: %w", err)
	}

	data, err := b.auctionABI.Pack(
		"bidERC20",
		auctionID,
		common.HexToAddress(req.BidToken),
		amount,
	)
	if err != nil {
		return "", fmt.Errorf("pack bidERC20 calldata: %w", err)
	}

	return "0x" + hex.EncodeToString(data), nil
}

func (b *TxCalldataBuilder) BuildCancelAuctionCalldata(ctx context.Context, req tx.BuildCancelAuctionTxRequest) (string, error) {
	_ = ctx

	data, err := b.auctionABI.Pack(
		"cancelAuction",
		new(big.Int).SetUint64(req.AuctionID),
	)
	if err != nil {
		return "", fmt.Errorf("pack cancelAuction calldata: %w", err)
	}

	return "0x" + hex.EncodeToString(data), nil
}

func (b *TxCalldataBuilder) BuildEndAuctionCalldata(ctx context.Context, req tx.BuildEndAuctionTxRequest) (string, error) {
	_ = ctx

	data, err := b.auctionABI.Pack(
		"endAuction",
		new(big.Int).SetUint64(req.AuctionID),
	)
	if err != nil {
		return "", fmt.Errorf("pack endAuction calldata: %w", err)
	}

	return "0x" + hex.EncodeToString(data), nil
}

func parseUint256(value string) (*big.Int, error) {
	value = strings.TrimSpace(value)
	if value == "" {
		return nil, fmt.Errorf("empty uint256")
	}

	n, ok := new(big.Int).SetString(value, 10)
	if !ok {
		return nil, fmt.Errorf("invalid uint256")
	}

	if n.Sign() < 0 {
		return nil, fmt.Errorf("negative uint256")
	}

	return n, nil
}
