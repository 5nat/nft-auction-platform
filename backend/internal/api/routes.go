package api

import (
	"log/slog"
	"time"

	"github.com/5nat/nft-auction-platform/backend/internal/store"
	"github.com/gin-gonic/gin"
)

type Dependencies struct {
	Logger    *slog.Logger
	DB        *store.DB
	StartedAt time.Time
}

func NewRouter(deps Dependencies) *gin.Engine {
	router := gin.New()

	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	router.GET("/health", HealthHandler(deps))

	return router
}
