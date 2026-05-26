package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/5nat/nft-auction-platform/backend/internal/app"
	"github.com/5nat/nft-auction-platform/backend/internal/config"
)

/*
  main 很薄；
	run 负责启动流程；
	App 负责组件生命周期；
	os.Exit 只在最外层调用；
	错误通过 fmt.Errorf("%w") 向上包装。
*/

func main() {
	// 创建结构化 JSON logger。
	// 后端服务中建议统一使用结构化日志，方便后续接入 ELK、Loki、Datadog 等日志系统。
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))

	// 创建根 context。
	// signal.NotifyContext 会监听系统退出信号：
	// - Ctrl+C 对应 SIGINT
	// - kill / 容器停止通常对应 SIGTERM
	//
	// 一旦收到信号，ctx.Done() 会被触发。
	// 这个 ctx 会向下传递给整个应用，用于控制 HTTP server、Indexer 等组件的退出。
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// main 函数只负责进程入口和退出码。
	// 真正的启动逻辑放到 run 里，避免 main 过于臃肿。
	if runErr := run(ctx, logger); runErr != nil {
		logger.Error("api-server command failed", "error", runErr)
		os.Exit(1)
	}

	logger.Info("api-server command stopped")
}

func run(ctx context.Context, logger *slog.Logger) error {
	// 加载配置。
	// 配置包括：
	// - HTTP 服务配置
	// - MySQL 配置
	// - Chain RPC 配置
	// - Indexer 配置
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}

	// 创建 App。
	// App 是整个后端应用的组件容器，负责组装：
	// - HTTP Server
	// - DB
	// - Chain Client
	// - Indexer
	application, err := app.New(ctx, cfg, logger)
	if err != nil {
		return fmt.Errorf("initialize application: %w", err)
	}

	// 运行整个应用。
	// application.Run 会同时启动 HTTP server 和 Indexer，
	// 并负责监听 ctx.Done()，收到退出信号后执行优雅关闭。
	if runErr := application.Run(ctx); runErr != nil {
		return fmt.Errorf("run application: %w", runErr)
	}

	return nil
}
