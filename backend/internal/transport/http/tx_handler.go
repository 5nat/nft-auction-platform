package httptransport

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/5nat/nft-auction-platform/backend/internal/modules/auth"
	txmodule "github.com/5nat/nft-auction-platform/backend/internal/modules/tx"
	"github.com/5nat/nft-auction-platform/backend/internal/transport/http/middleware"
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
	currentUser, ok := requireCurrentUser(c)
	if !ok {
		return
	}

	var req txmodule.BuildApproveNFTTxRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		Error(c, http.StatusBadRequest, CodeBadRequest, "invalid request body")
		return
	}

	req.Actor = txmodule.ActorFromAuth(currentUser)

	result, err := h.service.BuildApproveNFTTx(c.Request.Context(), req)
	if err != nil {
		writeTxError(c, h.logger, "build approve nft tx failed", err)
		return
	}

	OK(c, result)
}

func (h *TxHandler) BuildCreateAuctionTx(c *gin.Context) {
	currentUser, ok := requireCurrentUser(c)
	if !ok {
		return
	}

	var req txmodule.BuildCreateAuctionTxRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		Error(c, http.StatusBadRequest, CodeBadRequest, "invalid request body")
		return
	}

	req.Actor = txmodule.ActorFromAuth(currentUser)

	result, err := h.service.BuildCreateAuctionTx(c.Request.Context(), req)
	if err != nil {
		writeTxError(c, h.logger, "build create auction tx failed", err)
		return
	}

	OK(c, result)
}

func (h *TxHandler) BuildPlaceBidTx(c *gin.Context) {
	currentUser, ok := requireCurrentUser(c)
	if !ok {
		return
	}

	var req txmodule.BuildPlaceBidTxRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		Error(c, http.StatusBadRequest, CodeBadRequest, "invalid request body")
		return
	}

	req.Actor = txmodule.ActorFromAuth(currentUser)

	result, err := h.service.BuildPlaceBidTx(c.Request.Context(), req)
	if err != nil {
		writeTxError(c, h.logger, "build place bid tx failed", err)
		return
	}

	OK(c, result)
}

func (h *TxHandler) BuildCancelAuctionTx(c *gin.Context) {
	currentUser, ok := requireCurrentUser(c)
	if !ok {
		return
	}

	var req txmodule.BuildCancelAuctionTxRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		Error(c, http.StatusBadRequest, CodeBadRequest, "invalid request body")
		return
	}

	req.Actor = txmodule.ActorFromAuth(currentUser)

	result, err := h.service.BuildCancelAuctionTx(c.Request.Context(), req)
	if err != nil {
		writeTxError(c, h.logger, "build cancel auction tx failed", err)
		return
	}

	OK(c, result)
}

func (h *TxHandler) BuildEndAuctionTx(c *gin.Context) {
	currentUser, ok := requireCurrentUser(c)
	if !ok {
		return
	}

	var req txmodule.BuildEndAuctionTxRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		Error(c, http.StatusBadRequest, CodeBadRequest, "invalid request body")
		return
	}

	req.Actor = txmodule.ActorFromAuth(currentUser)

	result, err := h.service.BuildEndAuctionTx(c.Request.Context(), req)
	if err != nil {
		writeTxError(c, h.logger, "build end auction tx failed", err)
		return
	}

	OK(c, result)
}

func requireCurrentUser(c *gin.Context) (*auth.CurrentUser, bool) {
	currentUser, ok := middleware.CurrentUser(c)
	if !ok || currentUser == nil {
		Error(c, http.StatusUnauthorized, CodeUnauthorized, "unauthorized")
		return nil, false
	}

	return currentUser, true
}
