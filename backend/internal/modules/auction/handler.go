package auction

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AuctionService interface {
	ListAuctions(ctx context.Context, query ListAuctionsQuery) (PageResult[AuctionDTO], error)
	GetAuction(ctx context.Context, query GetAuctionQuery) (AuctionDTO, error)
	ListBids(ctx context.Context, query ListBidsQuery) (PageResult[BidDTO], error)
}

type Handler struct {
	service AuctionService
	logger  *slog.Logger
}

func NewHandler(service AuctionService, logger *slog.Logger) *Handler {
	return &Handler{
		service: service,
		logger:  logger,
	}
}

func (h *Handler) ListAuctions(c *gin.Context) {
	chainID, ok := parseInt64Query(c, "chain_id", 0)
	if !ok {
		return
	}

	page, ok := parseIntQuery(c, "page", 1)
	if !ok {
		return
	}

	pageSize, ok := parseIntQuery(c, "page_size", DefaultPageSize)
	if !ok {
		return
	}

	result, err := h.service.ListAuctions(c.Request.Context(), ListAuctionsQuery{
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
		h.writeError(c, "list auctions failed", err)
	}

	c.JSON(http.StatusOK, gin.H{
		"data": result.Items,
		"meta": result.Meta,
	})
}

func (h *Handler) GetAuction(c *gin.Context) {
	auctionID, ok := parseUint64Param(c, "auctionId")
	if !ok {
		return
	}

	chainID, ok := parseInt64Query(c, "chain_id", 0)
	if !ok {
		return
	}

	auction, err := h.service.GetAuction(c.Request.Context(), GetAuctionQuery{
		ChainID:         chainID,
		ContractAddress: c.Query("contract_address"),
		AuctionID:       auctionID,
	})
	if err != nil {
		h.writeError(c, "get auction failed", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": auction,
	})
}

func (h *Handler) ListBids(c *gin.Context) {
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

	pageSize, ok := parseIntQuery(c, "page_size", DefaultPageSize)
	if !ok {
		return
	}

	result, err := h.service.ListBids(c.Request.Context(), ListBidsQuery{
		ChainID:         chainID,
		ContractAddress: c.Query("contract_address"),
		AuctionID:       auctionID,
		Page:            page,
		PageSize:        pageSize,
	})
	if err != nil {
		h.writeError(c, "list bids failed", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": result.Items,
		"meta": result.Meta,
	})
}

func (h *Handler) writeError(c *gin.Context, message string, err error) {
	if validationErr, ok := errors.AsType[*ValidationError](err); ok {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": validationErr.Error(),
		})
		return
	}

	if errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "auction not found",
		})
		return
	}

	h.logger.Error(message, "error", err)

	c.JSON(http.StatusInternalServerError, gin.H{
		"error": message,
	})
}

func parseUint64Param(c *gin.Context, name string) (uint64, bool) {
	raw := c.Param(name)

	value, err := strconv.ParseUint(raw, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid" + name,
		})
		return 0, false
	}

	return value, true
}

func parseIntQuery(c *gin.Context, name string, defaultValue int) (int, bool) {
	raw := c.Query(name)
	if raw == "" {
		return defaultValue, true
	}

	value, err := strconv.Atoi(raw)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid" + name,
		})
		return 0, false
	}

	return value, true
}

func parseInt64Query(c *gin.Context, name string, defaultValue int64) (int64, bool) {
	raw := c.Query(name)
	if raw == "" {
		return defaultValue, true
	}
	value, err := strconv.ParseInt(raw, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid" + name,
		})
		return 0, false
	}

	return value, true
}
