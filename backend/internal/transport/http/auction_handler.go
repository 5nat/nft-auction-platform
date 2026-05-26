package httptransport

import (
	"context"
	"log/slog"

	"github.com/5nat/nft-auction-platform/backend/internal/modules/auction"
	"github.com/gin-gonic/gin"
)

type AuctionService interface {
	ListAuctions(ctx context.Context, query auction.ListAuctionsQuery) (auction.PageResult[auction.AuctionDTO], error)
	GetAuction(ctx context.Context, query auction.GetAuctionQuery) (auction.AuctionDTO, error)
	ListBids(ctx context.Context, query auction.ListBidsQuery) (auction.PageResult[auction.BidDTO], error)
}

type AuctionHandler struct {
	service AuctionService
	logger  *slog.Logger
}

func NewAuctionHandler(service AuctionService, logger *slog.Logger) *AuctionHandler {
	return &AuctionHandler{
		service: service,
		logger:  logger,
	}
}

func (h *AuctionHandler) ListAuctions(c *gin.Context) {
	chainID, ok := parseInt64Query(c, "chain_id", 0)
	if !ok {
		return
	}

	page, ok := parseIntQuery(c, "page", 1)
	if !ok {
		return
	}

	pageSize, ok := parseIntQuery(c, "page_size", auction.DefaultPageSize)
	if !ok {
		return
	}

	result, err := h.service.ListAuctions(c.Request.Context(), auction.ListAuctionsQuery{
		ChainID:         chainID,
		ContractAddress: c.Query("contract_address"),
		Status:          c.Query("status"),
		Seller:          c.Query("seller"),
		NFT:             c.Query("nft"),
		Page:            page,
		PageSize:        pageSize,
		Sort:            c.Query("sort"),
	})
	if err != nil {
		writeAuctionError(c, h.logger, "list auctions failed", err)
		return
	}

	OKWithMeta(c, result.Items, result.Meta)
}

func (h *AuctionHandler) GetAuction(c *gin.Context) {
	auctionID, ok := parseUint64Param(c, "auctionId")
	if !ok {
		return
	}

	chainID, ok := parseInt64Query(c, "chain_id", 0)
	if !ok {
		return
	}

	result, err := h.service.GetAuction(c.Request.Context(), auction.GetAuctionQuery{
		ChainID:         chainID,
		ContractAddress: c.Query("contract_address"),
		AuctionID:       auctionID,
	})
	if err != nil {
		writeAuctionError(c, h.logger, "list auctions failed", err)
		return
	}

	OK(c, result)
}

func (h *AuctionHandler) ListBids(c *gin.Context) {
	auctionID, ok := parseUint64Param(c, "auctionId")
	if !ok {
		return
	}

	chainID, ok := parseInt64Query(c, "chain_id", 0)
	if !ok {
		return
	}

	page, ok := parseIntQuery(c, "page", 1)
	if !ok {
		return
	}

	pageSize, ok := parseIntQuery(c, "page_size", auction.DefaultPageSize)
	if !ok {
		return
	}

	result, err := h.service.ListBids(c.Request.Context(), auction.ListBidsQuery{
		ChainID:         chainID,
		ContractAddress: c.Query("contract_address"),
		AuctionID:       auctionID,
		Page:            page,
		PageSize:        pageSize,
	})
	if err != nil {
		writeAuctionError(c, h.logger, "list auctions failed", err)
		return
	}

	OKWithMeta(c, result.Items, result.Meta)
}
