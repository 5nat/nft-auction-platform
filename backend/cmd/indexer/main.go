package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/5nat/nft-auction-platform/backend/internal/chain"
	"github.com/5nat/nft-auction-platform/backend/internal/config"
	"github.com/5nat/nft-auction-platform/backend/internal/indexer"
	"github.com/5nat/nft-auction-platform/backend/internal/store"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	if runErr := run(ctx, logger); runErr != nil {
		logger.Error("indexer command failed", "error", runErr)
		os.Exit(1)
	}

	logger.Info("indexer command stopped")
}

func run(ctx context.Context, logger *slog.Logger) error {
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}

	db, err := store.NewMySQL(ctx, cfg.Database.MySQLDSN)
	if err != nil {
		return fmt.Errorf("connect mysql: %w", err)
	}

	defer func() {
		if dbCloseErr := db.Close(); dbCloseErr != nil {
			logger.Error("close mysql failed", "error", dbCloseErr)
		}
	}()

	chainClient, err := chain.NewClient(ctx, cfg.Chain, logger)
	if err != nil {
		return fmt.Errorf("create chain client: %w", err)
	}
	defer chainClient.Close()

	idx, err := indexer.New(db.Gorm, chainClient, cfg, logger)
	if err != nil {
		return fmt.Errorf("create indexer: %w", err)
	}

	if startErr := idx.Start(ctx); startErr != nil {
		return fmt.Errorf("start indexer: %w", startErr)
	}

	return nil
}
