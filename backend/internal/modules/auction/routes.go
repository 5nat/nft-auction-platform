package auction

import "github.com/gin-gonic/gin"

func RegisterRoutes(router gin.IRouter, handler *Handler) {
	router.GET("/auctions", handler.ListAuctions)
	router.GET("/auctions/:auctionId", handler.GetAuction)
	router.GET("/auctions/:auctionId/bids", handler.ListBids)
}
