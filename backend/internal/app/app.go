package app

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"sync"
	"time"

	"github.com/5nat/nft-auction-platform/backend/internal/config"
	ethinfra "github.com/5nat/nft-auction-platform/backend/internal/infra/blockchain/ethereum"
	mysql "github.com/5nat/nft-auction-platform/backend/internal/infra/persistence/mysql"
	httptransport "github.com/5nat/nft-auction-platform/backend/internal/transport/http"
	"github.com/5nat/nft-auction-platform/backend/internal/workers/indexer"
	"github.com/gin-gonic/gin"
)

// App 表示整个后端应用进程，而不只是 HTTP server。
// 现在这个 App 管理多个组件：
// 1. HTTP Server：对外提供 REST API。
// 2. MySQL DB：保存 read model、事件同步进度等。
// 3. Chain Client：连接 Ethereum RPC。
// 4. Indexer：后台扫描链上事件并写入数据库。
type App struct {
	cfg    config.Config
	logger *slog.Logger

	// HTTP 服务组件。
	server *http.Server

	// 数据库连接。
	db *mysql.DB

	// 链客户端。Indexer 需要用它访问 RPC。
	chainClient *ethinfra.Client

	// 链上事件同步器。
	// 当 cfg.Indexer.Enabled = false 时，这个字段可以为 nil。
	indexer *indexer.Indexer

	startedAt time.Time
}

// New 负责创建应用所需的所有组件，但不真正启动它们。
// 也就是说：
// - New 负责依赖组装。
// - Run 负责生命周期运行。
// - shutdown 负责资源释放。
func New(ctx context.Context, cfg config.Config, logger *slog.Logger) (*App, error) {
	// 设置 Gin 的运行模式，例如 debug / release。
	gin.SetMode(cfg.Server.GinMode)

	// 初始化数据库连接。
	// 这里使用 ctx 是为了让连接过程可以被取消。
	db, err := mysql.NewMySQL(ctx, cfg.Database.MySQLDSN)
	if err != nil {
		return nil, fmt.Errorf("connect database: %w", err)
	}

	startedAt := time.Now()

	// 创建 HTTP router。
	// Dependencies 是一种显式依赖注入方式。
	// 这样 handler 不需要自己创建 DB 或 logger，而是从外部传入。
	router := httptransport.NewRouter(httptransport.Dependencies{
		Logger:    logger,
		DB:        db,
		Config:    cfg,
		StartedAt: startedAt,
	})

	// 创建标准库 http.Server。
	// Gin 本质上是 Handler，真正负责监听端口的是 http.Server。
	server := &http.Server{
		Addr:         cfg.Server.Addr,
		Handler:      router,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
		IdleTimeout:  cfg.Server.IdleTimeout,
	}

	var chainClient *ethinfra.Client
	var auctionIndexer *indexer.Indexer

	// 根据配置决定是否在 API 服务进程中同时启动 Indexer。
	// 这个设计很重要：
	// 本地开发时可以 API + Indexer 跑在一个进程里；
	// 生产环境中也可以把 API 和 Indexer 拆成两个独立进程部署。
	if cfg.Indexer.Enabled {
		chainClient, err = ethinfra.NewClient(ctx, cfg.Chain, logger)
		if err != nil {
			// 如果链客户端初始化失败，需要关闭已经打开的 DB。
			_ = db.Close()
			return nil, fmt.Errorf("create chain client: %w", err)
		}

		auctionIndexer, err = indexer.New(db.Gorm, chainClient, cfg, logger)
		if err != nil {
			// 如果 Indexer 初始化失败，需要清理前面已经创建成功的资源。
			chainClient.Close()
			_ = db.Close()
			return nil, fmt.Errorf("create indexer: %w", err)
		}
	}

	return &App{
		cfg:         cfg,
		logger:      logger,
		server:      server,
		db:          db,
		chainClient: chainClient,
		indexer:     auctionIndexer,
		startedAt:   startedAt,
	}, nil
}

// Run 是整个应用的生命周期入口。
// 它负责：
// 1. 启动 HTTP server goroutine。
// 2. 启动 Indexer goroutine。
// 3. 等待系统退出信号或组件异常。
// 4. 触发统一 shutdown。
func (a *App) Run(ctx context.Context) error {
	a.logger.Info(
		"application starting",
		"http_addr", a.server.Addr,
		"indexer_enabled", a.indexer != nil,
	)

	// runCtx 是传给后台组件的运行上下文。
	// 当需要关闭应用时，调用 cancel()，Indexer 会感知到 ctx.Done() 并退出。
	runCtx, cancel := context.WithCancel(ctx)
	defer cancel()

	// errCh 用于接收后台组件的异常错误。
	// 目前有两个组件可能报错：
	// 1. HTTP server
	// 2. Indexer
	errCh := make(chan error, 2)

	// WaitGroup 用于关闭阶段等待 goroutine 退出。
	// 注意：
	// WaitGroup 不是用来确认组件“启动成功”的；
	// WaitGroup 是用来确认组件“已经退出”的。
	var wg sync.WaitGroup

	// 启动 HTTP server。
	wg.Add(1)
	go func() {
		defer wg.Done()

		a.logger.Info("http server starting", "addr", a.server.Addr)

		// ListenAndServe 会阻塞，直到：
		// 1. HTTP server 被 Shutdown；
		// 2. 启动失败；
		// 3. 监听过程中出现异常。
		err := a.server.ListenAndServe()

		// http.ErrServerClosed 是正常关闭时返回的错误，不应该当作异常。
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			sendComponentError(runCtx, errCh, fmt.Errorf("http server: %w", err))
			return
		}

		a.logger.Info("http server stopped")
	}()

	// 如果开启了 Indexer，则启动 Indexer 后台循环。
	if a.indexer != nil {
		wg.Add(1)
		go func() {
			defer wg.Done()

			a.logger.Info("indexer starting")

			// Start 会长期运行：RunOnce → sleep → RunOnce → sleep ...
			// 当 runCtx 被 cancel 后，Start 会退出。
			if err := a.indexer.Start(runCtx); err != nil {
				// 如果是因为应用主动关闭导致的 ctx 取消，不算异常。
				if runCtx.Err() != nil {
					return
				}
				sendComponentError(runCtx, errCh, fmt.Errorf("indexer: %w", err))
				return
			}
		}()
	}

	var runErr error

	// 等待两类事件：
	// 1. ctx.Done()：说明收到 Ctrl+C / SIGTERM。
	// 2. errCh：说明某个组件异常失败。
	select {
	case <-ctx.Done():
		a.logger.Info("shutdown signal received")
	case componentErr := <-errCh:
		if componentErr != nil {
			runErr = componentErr
			a.logger.Error("application component failed", "error", componentErr)
		}
	}

	// 触发后台组件停止。
	// 这会通知 Indexer 的 Start 循环退出。
	cancel()

	// 创建 shutdownCtx，限制优雅关闭的最长等待时间。
	// 避免某个 goroutine 卡住导致进程永远无法退出。
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), a.cfg.Server.ShutdownTimeout)
	defer shutdownCancel()

	// 执行统一关闭：
	// 1. Shutdown HTTP server
	// 2. 等待 goroutine 退出
	// 3. 关闭 chain client
	// 4. 关闭 DB
	if shutdownErr := a.shutdown(shutdownCtx, &wg); shutdownErr != nil {
		if runErr != nil {
			// 如果组件已经异常失败，同时 shutdown 也失败，
			// 优先返回组件运行错误，shutdown 错误只记录日志。
			a.logger.Error("application shutdown failed", "error", shutdownErr)
			return runErr
		}
		return shutdownErr
	}

	return runErr
}

// shutdown 负责按正确顺序关闭应用组件。
//
// 关闭顺序很重要：
// 1. 先 Shutdown HTTP server，让它停止接收新请求。
// 2. 等 HTTP 和 Indexer goroutine 退出。
// 3. 再关闭 chain client。
// 4. 最后关闭 DB。
//
// 不能先关闭 DB。
// 否则 Indexer 可能还在写库，就会出现 sql: database is closed。
func (a *App) shutdown(ctx context.Context, wg *sync.WaitGroup) error {
	a.logger.Info("application shutting down")

	var shutdownErr error

	// 优雅关闭 HTTP server。
	// Shutdown 会停止接收新请求，并等待正在处理的请求完成。
	if a.server != nil {
		if err := a.server.Shutdown(ctx); err != nil {
			shutdownErr = fmt.Errorf("shutdown http server: %w", err)
		}
	}

	// 等待后台 goroutine 退出。
	//
	// 这里没有直接 wg.Wait()，而是放到 goroutine 里，
	// 是因为 wg.Wait() 本身不支持 context 超时。
	// 所以我们用 done channel + select 实现“带超时等待”。
	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		// 所有被 WaitGroup 管理的 goroutine 已经退出。
		// 当前项目中包括：
		// - HTTP server goroutine
		// - Indexer goroutine
	case <-ctx.Done():
		// 等待超时，说明某些组件没有按预期退出。
		if shutdownErr != nil {
			shutdownErr = fmt.Errorf("wait component shutdown: %w", ctx.Err())
		} else {
			a.logger.Error("wait component shutdown failed", "error", ctx.Err())
		}
	}

	// 关闭链客户端。
	// 必须在 Indexer 退出之后关闭，否则 Indexer 可能还在使用 RPC client。
	if a.chainClient != nil {
		a.chainClient.Close()
	}

	// 关闭数据库连接。
	// 必须在所有 goroutine 停止之后关闭，否则后台任务可能仍在使用 DB。
	if a.db != nil {
		if err := a.db.Close(); err != nil {
			a.logger.Error("close database failed", "error", err)

			if shutdownErr == nil {
				shutdownErr = fmt.Errorf("close database: %w", err)
			}
		}
	}

	a.logger.Info("application shut down")

	return shutdownErr
}

// sendComponentError 用于向 errCh 发送组件错误。
//
// 为什么不用简单的 errCh <- err？
// 因为如果应用已经开始关闭，ctx.Done() 可能已经触发。
// 这时如果 errCh 没有人接收，goroutine 可能阻塞住。
//
// 使用 select 可以避免 shutdown 过程中 goroutine 卡死。
func sendComponentError(ctx context.Context, errCh chan<- error, err error) {
	select {
	case errCh <- err:
	case <-ctx.Done():
	}
}
