package indexer

import (
	"context"
	"fmt"
	"math/big"

	ethinfra "github.com/5nat/nft-auction-platform/backend/internal/infra/blockchain/ethereum"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

type Scanner struct {
	chainClient     *ethinfra.Client
	contractAddress common.Address
}

// 如果 Scanner 没有实现 ChainScanner，编译期直接报错。
var _ ChainScanner = (*Scanner)(nil)

func NewScanner(chainClient *ethinfra.Client, contractAddress common.Address) *Scanner {
	return &Scanner{
		chainClient:     chainClient,
		contractAddress: contractAddress,
	}
}

func (s *Scanner) LatestBlockNumber(ctx context.Context) (uint64, error) {
	return s.chainClient.LatestBlockNumber(ctx)
}

func (s *Scanner) FilterLogs(ctx context.Context, fromBlock uint64, toBlock uint64) ([]types.Log, error) {
	query := ethereum.FilterQuery{
		FromBlock: new(big.Int).SetUint64(fromBlock),
		ToBlock:   new(big.Int).SetUint64(toBlock),
		Addresses: []common.Address{s.contractAddress},
	}

	logs, err := s.chainClient.EthClient().FilterLogs(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("filter logs: from_block=%d to_block=%d: %w", fromBlock, toBlock, err)
	}

	return logs, nil
}

func (s *Scanner) BlockHash(ctx context.Context, blockNumber uint64) (common.Hash, error) {
	header, err := s.chainClient.EthClient().HeaderByNumber(ctx, new(big.Int).SetUint64(blockNumber))
	if err != nil {
		return common.Hash{}, fmt.Errorf("get block header: block_number=%d: %w", blockNumber, err)
	}

	return header.Hash(), nil
}
