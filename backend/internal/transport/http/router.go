package httptransport

import (
	"log/slog"
	"time"

	"github.com/5nat/nft-auction-platform/backend/internal/config"
	ethinfra "github.com/5nat/nft-auction-platform/backend/internal/infra/blockchain/ethereum"
	mysql "github.com/5nat/nft-auction-platform/backend/internal/infra/persistence/mysql"
	"github.com/5nat/nft-auction-platform/backend/internal/modules/auction"
	"github.com/5nat/nft-auction-platform/backend/internal/modules/auth"
	txmodule "github.com/5nat/nft-auction-platform/backend/internal/modules/tx"
	"github.com/5nat/nft-auction-platform/backend/internal/transport/http/middleware"
	"github.com/gin-gonic/gin"
)

type Dependencies struct {
	Logger *slog.Logger
	DB     *mysql.DB
	Config config.Config

	StartedAt time.Time
}

type Handlers struct {
	Auction *AuctionHandler
	Tx      *TxHandler
	Auth    *AuthHandler
}

// BuildHandlers 负责组装 HTTP 层需要的 Handler。
// 这里是当前 HTTP 层的组合入口：创建 Repository、Service 和 Handler。
// NewRouter 只负责注册路由，不再承担依赖创建职责。
func BuildHandlers(deps Dependencies) (Handlers, error) {
	// Auction
	auctionRepo := mysql.NewAuctionRepository(deps.DB.Gorm)

	auctionService := auction.NewService(auctionRepo, auction.ServiceConfig{
		DefaultChainID:         deps.Config.Chain.ChainID,
		DefaultContractAddress: deps.Config.Chain.AuctionContract,
	})

	auctionHandler := NewAuctionHandler(auctionService, deps.Logger)

	// Tx
	txBuilder, err := ethinfra.NewTxCalldataBuilder()
	if err != nil {
		panic("initializing tx calldata builder: " + err.Error())
	}

	txService := txmodule.NewService(
		txBuilder,
		txmodule.ServiceConfig{
			DefaultChainID:         deps.Config.Chain.ChainID,
			DefaultContractAddress: deps.Config.Chain.AuctionContract,
		},
		auctionRepo,
	)
	txHandler := NewTxHandler(txService, deps.Logger)

	// Auth
	authRepo := mysql.NewAuthRepository(deps.DB.Gorm)

	authService := auth.NewService(authRepo,
		authRepo,
		auth.ServiceConfig{
			Domain:         deps.Config.Auth.Domain,
			URI:            deps.Config.Auth.URI,
			Statement:      "Sign in to NFT Auction Platform.",
			NonceTTL:       10 * time.Minute,
			JWTSecret:      deps.Config.Auth.JWTSecret,
			AccessTokenTTL: deps.Config.Auth.AccessTokenTTL,
		})

	authHandler := NewAuthHandler(authService, deps.Logger)

	return Handlers{
		Auction: auctionHandler,
		Tx:      txHandler,
		Auth:    authHandler,
	}, nil
}

func NewRouter(deps Dependencies, handlers Handlers) *gin.Engine {
	router := gin.New()

	router.Use(middleware.Logger(deps.Logger))
	router.Use(middleware.Recovery(deps.Logger, Error, CodeInternalError))

	router.GET("/health", HealthHandler(deps))

	if gin.Mode() != gin.ReleaseMode {
		router.StaticFile("/dev/auth-test.html", "./dev/auth-test.html")
	}

	apiV1 := router.Group("/api/v1")

	registerAuctionRoutes(apiV1, handlers.Auction)
	registerTxRoutes(apiV1, handlers.Tx, handlers.Auth)
	registerAuthRoutes(apiV1, handlers.Auth)

	return router
}

func registerAuctionRoutes(apiV1 *gin.RouterGroup, handler *AuctionHandler) {
	apiV1.GET("/auctions", handler.ListAuctions)
	apiV1.GET("/auctions/:auctionId", handler.GetAuction)
	apiV1.GET("/auctions/:auctionId/bids", handler.ListBids)
}

func registerTxRoutes(apiV1 *gin.RouterGroup, txHandler *TxHandler, authHandler *AuthHandler) {
	authRequired := middleware.Auth(authHandler.service, authHandler.logger, writeAuthError)

	txGroup := apiV1.Group("/tx")
	txGroup.Use(authRequired)
	{
		txGroup.POST("/approve-nft", txHandler.BuildApproveNFTTx)
		txGroup.POST("/create-auction", txHandler.BuildCreateAuctionTx)
		txGroup.POST("/place-bid", txHandler.BuildPlaceBidTx)
		txGroup.POST("/cancel-auction", txHandler.BuildCancelAuctionTx)
		txGroup.POST("/end-auction", txHandler.BuildEndAuctionTx)
	}

}

func registerAuthRoutes(apiV1 *gin.RouterGroup, handler *AuthHandler) {
	authGroup := apiV1.Group("/auth")
	{
		authGroup.POST("/nonce", handler.CreateNonce)
		authGroup.POST("/verify", handler.VerifyNonce)
	}

	authRequired := middleware.Auth(handler.service, handler.logger, writeAuthError)

	apiV1.GET("/me", authRequired, handler.Me)
}
