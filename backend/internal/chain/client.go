package chain

import (
	"context"
	"fmt"
	"log/slog"
	"math/big"

	"github.com/5nat/nft-auction-platform/backend/internal/config"
	"github.com/ethereum/go-ethereum/ethclient"
)

// Client 封装 go-ethereum 的 ethClient
// 以后所有的链上查询、事件扫描、交易发送，都会基于这个 Client 扩展
type Client struct {
	eth    *ethclient.Client
	cfg    config.ChainConfig
	logger *slog.Logger
}

// NewClient 连接以太坊节点，并校验当前 RPC 的 chainID
func NewClient(ctx context.Context, cfg config.ChainConfig, logger *slog.Logger) (*Client, error) {
	if cfg.PRCURL == "" {
		return nil, fmt.Errorf("chain config has no prc url")
	}

	eth, err := ethclient.DialContext(ctx, cfg.PRCURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to eth client: %w", err)
	}

	c := &Client{
		eth:    eth,
		cfg:    cfg,
		logger: logger,
	}

	if err := c.checkChainId(ctx); err != nil {
		eth.Close()
		return nil, err
	}

	return c, nil
}

// EthClient 暴露底层 ethClient, 后面调用合约、扫 logs 会用到
func (c *Client) EthClient() *ethclient.Client {
	return c.eth
}

// Close 关闭 PRC 连接
func (c *Client) Close() {
	if c.eth != nil {
		c.eth.Close()
	}
}

// ChainID 查询当前 PRC 对应的链 ID
func (c *Client) ChainID(ctx context.Context) (*big.Int, error) {
	chainID, err := c.eth.ChainID(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch chain ID: %w", err)
	}
	return chainID, nil
}

// LatestBlockNumber 查询最新区块高度。
func (c *Client) LatestBlockNumber(ctx context.Context) (uint64, error) {
	blockNumber, err := c.eth.BlockNumber(ctx)
	if err != nil {
		return 0, fmt.Errorf("failed to fetch block number: %w", err)
	}
	return blockNumber, nil
}

// checkChainId 校验配置里的 CHAIN_ID 和 RPC 实际 chainID 是否一致
func (c *Client) checkChainId(ctx context.Context) error {
	actualChainID, err := c.ChainID(ctx)
	if err != nil {
		return err
	}

	expectedChainID := big.NewInt(c.cfg.ChainID)

	if actualChainID.Cmp(expectedChainID) != 0 {
		return fmt.Errorf(
			"chain id mismatch: expected  %s, got %s",
			expectedChainID.String(),
			actualChainID.String(),
		)
	}

	c.logger.Info(
		"chain id verified",
		"chain_id", actualChainID.String(),
		"rpc_url", c.cfg.PRCURL,
	)
	return nil
}
