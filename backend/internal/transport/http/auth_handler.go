package httptransport

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/5nat/nft-auction-platform/backend/internal/modules/auth"
	"github.com/5nat/nft-auction-platform/backend/internal/transport/http/middleware"
	"github.com/gin-gonic/gin"
)

type createAuthNonceRequest struct {
	Address string `json:"address" binding:"required"`
	ChainID int64  `json:"chain_id" binding:"required"`
}

type createAuthNonceResponse struct {
	Nonce     string `json:"nonce"`
	Message   string `json:"message"`
	ExpiresAt string `json:"expires_at"`
}

type verifyAuthRequest struct {
	Address   string `json:"address" binding:"required"`
	ChainID   int64  `json:"chain_id" binding:"required"`
	Message   string `json:"message" binding:"required"`
	Signature string `json:"signature" binding:"required"`
}

type verifyAuthResponse struct {
	AccessToken string         `json:"access_token"`
	TokenType   string         `json:"token_type"`
	ExpiresIn   int64          `json:"expires_in"`
	User        verifyAuthUser `json:"user"`
}

type verifyAuthUser struct {
	ID            uint64 `json:"id"`
	WalletID      uint64 `json:"wallet_id"`
	WalletAddress string `json:"wallet_address"`
	ChainID       int64  `json:"chain_id"`
}

type meResponse struct {
	ID      uint64     `json:"id"`
	Status  string     `json:"status"`
	Wallets []meWallet `json:"wallets"`
}

type meWallet struct {
	ID        uint64 `json:"id"`
	Address   string `json:"address"`
	ChainID   int64  `json:"chain_id"`
	IsPrimary bool   `json:"is_primary"`
}

type AuthHandler struct {
	service *auth.Service
	logger  *slog.Logger
}

func NewAuthHandler(service *auth.Service, logger *slog.Logger) *AuthHandler {
	return &AuthHandler{
		service: service,
		logger:  logger,
	}
}

// CreateNonce 生成钱包登录挑战。
// 这个接口不完成登录，只返回服务端生成的 SIWE 风格 message。
// 前端拿到 message 后，需要调用钱包签名，再把 message + signature 提交给 verify 接口。
func (h *AuthHandler) CreateNonce(c *gin.Context) {
	var req createAuthNonceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		Error(c, http.StatusBadRequest, CodeBadRequest, "invalid request body")
		return
	}

	challenge, err := h.service.CreateNonce(c.Request.Context(), auth.CreateNonceCommand{
		Address: req.Address,
		ChainID: req.ChainID,
	})

	if err != nil {
		writeAuthError(c, h.logger, "create auth nonce failed", err)
		return
	}

	OK(c, createAuthNonceResponse{
		Nonce:     challenge.Nonce,
		Message:   challenge.Message,
		ExpiresAt: challenge.ExpiresAt.Format(time.RFC3339),
	})
}

// VerifyNonce 校验钱包签名并完成登录。
// 这个接口会校验 message + signature，成功后返回 JWT access_token。
// 后续访问需要登录的接口时，前端应带上：
// Authorization: Bearer <access_token>
func (h *AuthHandler) VerifyNonce(c *gin.Context) {
	var req verifyAuthRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		Error(c, http.StatusBadRequest, CodeBadRequest, "invalid request body")
		return
	}

	result, err := h.service.Verify(c.Request.Context(), auth.VerifyCommand{
		Address:   req.Address,
		ChainID:   req.ChainID,
		Message:   req.Message,
		Signature: req.Signature,
	})
	if err != nil {
		writeAuthError(c, h.logger, "verify auth signature failed", err)
		return
	}

	OK(c, verifyAuthResponse{
		AccessToken: result.AccessToken,
		TokenType:   result.TokenType,
		ExpiresIn:   result.ExpiresIn,
		User: verifyAuthUser{
			ID:            result.User.ID,
			WalletID:      result.User.WalletID,
			WalletAddress: result.User.WalletAddress,
			ChainID:       result.User.ChainID,
		},
	})
}

// Me 返回当前登录用户信息。
// 该接口必须经过 AuthMiddleware，因此可以从 gin.Context 中读取 CurrentUser。
func (h *AuthHandler) Me(c *gin.Context) {
	currentUser, ok := middleware.CurrentUser(c)
	if !ok {
		writeAuthError(c, h.logger, "current user missing", auth.ErrUnauthorized)
		return
	}

	profile, err := h.service.GetProfile(c.Request.Context(), currentUser)
	if err != nil {
		writeAuthError(c, h.logger, "get current user profile failed", err)
		return
	}

	wallets := make([]meWallet, 0, len(profile.Wallets))
	for _, wallet := range profile.Wallets {
		wallets = append(wallets, meWallet{
			ID:        wallet.ID,
			Address:   wallet.Address,
			ChainID:   wallet.ChainID,
			IsPrimary: wallet.IsPrimary,
		})
	}

	OK(c, meResponse{
		ID:      profile.ID,
		Status:  profile.Status,
		Wallets: wallets,
	})
}
