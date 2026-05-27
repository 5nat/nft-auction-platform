package httptransport

import (
	"log/slog"
	"time"

	"github.com/5nat/nft-auction-platform/backend/internal/config"
	ethinfra "github.com/5nat/nft-auction-platform/backend/internal/infra/blockchain/ethereum"
	mysql "github.com/5nat/nft-auction-platform/backend/internal/infra/persistence/mysql"
	"github.com/5nat/nft-auction-platform/backend/internal/modules/auction"
	txmodule "github.com/5nat/nft-auction-platform/backend/internal/modules/tx"
	"github.com/gin-gonic/gin"
)

type Dependencies struct {
	Logger *slog.Logger
	DB     *mysql.DB
	Config config.Config

	StartedAt time.Time
}

func NewRouter(deps Dependencies) *gin.Engine {
	router := gin.New()

	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	router.GET("/health", HealthHandler(deps))

	// 拍卖 API
	apiV1 := router.Group("/api/v1")
	auctionRepo := mysql.NewAuctionRepository(deps.DB.Gorm)
	auctionService := auction.NewService(auctionRepo, auction.ServiceConfig{
		DefaultChainID:         deps.Config.Chain.ChainID,
		DefaultContractAddress: deps.Config.Chain.AuctionContract,
	})
	auctionHandler := NewAuctionHandler(auctionService, deps.Logger)
	apiV1.GET("/auctions", auctionHandler.ListAuctions)
	apiV1.GET("/auctions/:auctionId", auctionHandler.GetAuction)
	apiV1.GET("/auctions/:auctionId/bids", auctionHandler.ListBids)

	// 交易请求参数 API
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

	apiV1.POST("/tx/approve-nft", txHandler.BuildApproveNFTTx)
	apiV1.POST("/tx/create-auction", txHandler.BuildCreateAuctionTx)
	apiV1.POST("/tx/place-bid", txHandler.BuildPlaceBidTx)
	apiV1.POST("/tx/cancel-auction", txHandler.BuildCancelAuctionTx)
	apiV1.POST("/tx/end-auction", txHandler.BuildEndAuctionTx)

	return router
}
