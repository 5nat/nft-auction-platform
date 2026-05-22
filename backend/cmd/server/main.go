package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/5nat/nft-auction-platform/backend/internal/app"
	"github.com/5nat/nft-auction-platform/backend/internal/config"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))

	cfg, configLoadErr := config.Load()
	if configLoadErr != nil {
		logger.Error("failed to load config", "error", configLoadErr)
		os.Exit(1)
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	application, appInitErr := app.New(ctx, cfg, logger)
	if appInitErr != nil {
		logger.Error("failed to initialize application", "error", appInitErr)
		os.Exit(1)
	}

	errCh := make(chan error, 1)

	go func() {
		errCh <- application.Start()
	}()

	select {
	case err := <-errCh:
		if err != nil {
			logger.Error("failed to start application", "error", err)
		}
	case <-ctx.Done():
		logger.Info("shutting down server")
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), cfg.Server.ShutdownTimeout)
	defer cancel()

	if shutdownErr := application.Stop(shutdownCtx); shutdownErr != nil {
		logger.Error("failed to stop application", "error", shutdownErr)
		os.Exit(1)
	}
}
