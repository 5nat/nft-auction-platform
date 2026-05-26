package httptransport

import (
	"log/slog"
	"time"

	"github.com/5nat/nft-auction-platform/backend/internal/config"
	mysql "github.com/5nat/nft-auction-platform/backend/internal/infra/persistence/mysql"
	"github.com/5nat/nft-auction-platform/backend/internal/modules/auction"
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

	apiV1 := router.Group("/api/v1")
	auctionRepo := auction.NewGormRepository(deps.DB.Gorm)
	auctionService := auction.NewService(auctionRepo, auction.ServiceConfig{
		DefaultChainID:         deps.Config.Chain.ChainID,
		DefaultContractAddress: deps.Config.Chain.AuctionContract,
	})
	auctionHandler := auction.NewHandler(auctionService, deps.Logger)

	auction.RegisterRoutes(apiV1, auctionHandler)

	return router
}
