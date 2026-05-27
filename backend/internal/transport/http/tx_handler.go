package httptransport

import (
	"context"
	"log/slog"
	"net/http"

	txmodule "github.com/5nat/nft-auction-platform/backend/internal/modules/tx"
	"github.com/gin-gonic/gin"
)

type TxService interface {
	BuildApproveNFTTx(ctx context.Context, req txmodule.BuildApproveNFTTxRequest) (txmodule.TransactionRequestDTO, error)
	BuildCreateAuctionTx(ctx context.Context, req txmodule.BuildCreateAuctionTxRequest) (txmodule.TransactionRequestDTO, error)
	BuildPlaceBidTx(ctx context.Context, req txmodule.BuildPlaceBidTxRequest) (txmodule.TransactionRequestDTO, error)
	BuildCancelAuctionTx(ctx context.Context, req txmodule.BuildCancelAuctionTxRequest) (txmodule.TransactionRequestDTO, error)
	BuildEndAuctionTx(ctx context.Context, req txmodule.BuildEndAuctionTxRequest) (txmodule.TransactionRequestDTO, error)
}

type TxHandler struct {
	service TxService
	logger  *slog.Logger
}

func NewTxHandler(service TxService, logger *slog.Logger) *TxHandler {
	return &TxHandler{
		service: service,
		logger:  logger,
	}
}

func (h *TxHandler) BuildApproveNFTTx(c *gin.Context) {
	var req txmodule.BuildApproveNFTTxRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		Error(c, http.StatusBadRequest, CodeBadRequest, "invalid request body")
		return
	}

	result, err := h.service.BuildApproveNFTTx(c.Request.Context(), req)
	if err != nil {
		writeTxError(c, h.logger, "build approve nft tx failed", err)
		return
	}

	OK(c, result)
}

func (h *TxHandler) BuildCreateAuctionTx(c *gin.Context) {
	var req txmodule.BuildCreateAuctionTxRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		Error(c, http.StatusBadRequest, CodeBadRequest, "invalid request body")
		return
	}

	result, err := h.service.BuildCreateAuctionTx(c.Request.Context(), req)
	if err != nil {
		writeTxError(c, h.logger, "build create auction tx failed", err)
		return
	}

	OK(c, result)
}

func (h *TxHandler) BuildPlaceBidTx(c *gin.Context) {
	var req txmodule.BuildPlaceBidTxRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		Error(c, http.StatusBadRequest, CodeBadRequest, "invalid request body")
		return
	}

	result, err := h.service.BuildPlaceBidTx(c.Request.Context(), req)
	if err != nil {
		writeTxError(c, h.logger, "build place bid tx failed", err)
		return
	}

	OK(c, result)
}

func (h *TxHandler) BuildCancelAuctionTx(c *gin.Context) {
	var req txmodule.BuildCancelAuctionTxRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		Error(c, http.StatusBadRequest, CodeBadRequest, "invalid request body")
		return
	}

	result, err := h.service.BuildCancelAuctionTx(c.Request.Context(), req)
	if err != nil {
		writeTxError(c, h.logger, "build cancel auction tx failed", err)
		return
	}

	OK(c, result)
}

func (h *TxHandler) BuildEndAuctionTx(c *gin.Context) {
	var req txmodule.BuildEndAuctionTxRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		Error(c, http.StatusBadRequest, CodeBadRequest, "invalid request body")
		return
	}

	result, err := h.service.BuildEndAuctionTx(c.Request.Context(), req)
	if err != nil {
		writeTxError(c, h.logger, "build end auction tx failed", err)
		return
	}

	OK(c, result)
}
