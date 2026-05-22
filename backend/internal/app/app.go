package app

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/5nat/nft-auction-platform/backend/internal/api"
	"github.com/5nat/nft-auction-platform/backend/internal/config"
	"github.com/5nat/nft-auction-platform/backend/internal/store"
	"github.com/gin-gonic/gin"
)

type App struct {
	cfg       config.Config
	logger    *slog.Logger
	server    *http.Server
	db        *store.DB
	startedAt time.Time
}

func New(ctx context.Context, cfg config.Config, logger *slog.Logger) (*App, error) {
	gin.SetMode(cfg.Server.GinMode)

	db, err := store.NewMySQL(ctx, cfg.Database.MySQLDSN)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	startedAt := time.Now()

	router := api.NewRouter(api.Dependencies{
		Logger:    logger,
		DB:        db,
		StartedAt: startedAt,
	})

	server := &http.Server{
		Addr:         cfg.Server.Addr,
		Handler:      router,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
		IdleTimeout:  cfg.Server.IdleTimeout,
	}

	return &App{
		cfg:       cfg,
		logger:    logger,
		server:    server,
		db:        db,
		startedAt: startedAt,
	}, nil
}

func (a *App) Start() error {
	a.logger.Info("http server starting", "addr", a.server.Addr)

	err := a.server.ListenAndServe()
	if errors.Is(err, http.ErrServerClosed) {
		return nil
	}

	return err
}

func (a *App) Stop(ctx context.Context) error {
	a.logger.Info("application shutting down")

	var shutdownErr error

	if err := a.server.Shutdown(ctx); err != nil {
		shutdownErr = fmt.Errorf("failed to gracefully shutdown: %w", err)
	}

	if a.db != nil {
		if err := a.db.Close(); err != nil {
			a.logger.Error("failed to close database", "error", err)
		}
	}

	a.logger.Info("application shut down")

	return shutdownErr
}
