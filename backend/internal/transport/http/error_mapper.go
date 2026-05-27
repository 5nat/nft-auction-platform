package httptransport

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/5nat/nft-auction-platform/backend/internal/modules/auction"
	"github.com/5nat/nft-auction-platform/backend/internal/modules/auth"
	txmodule "github.com/5nat/nft-auction-platform/backend/internal/modules/tx"
	"github.com/gin-gonic/gin"
)

// mappedError 表示领域错误映射到 HTTP 层后的结果。
// status 是 HTTP 状态码，code 是项目内部业务错误码。
// log 用于控制是否需要记录 error 级别日志。
type mappedError struct {
	status  int
	code    int
	message string
	log     bool
}

func writeAuctionError(c *gin.Context, logger *slog.Logger, logMessage string, err error) {
	writeAppError(c, logger, logMessage, err)
}

func writeTxError(c *gin.Context, logger *slog.Logger, logMessage string, err error) {
	writeAppError(c, logger, logMessage, err)
}

func writeAuthError(c *gin.Context, logger *slog.Logger, logMessage string, err error) {
	writeAppError(c, logger, logMessage, err)
}

// writeAppError 是 HTTP 层统一错误出口。
// 可预期业务错误不打 error 日志；未知错误才记录日志并返回 INTERNAL_ERROR。
func writeAppError(c *gin.Context, logger *slog.Logger, logMessage string, err error) {
	mapped := mapDomainError(err)

	if mapped.log && logger != nil {
		logger.Error(logMessage, "error", err)
	}

	Error(c, mapped.status, mapped.code, mapped.message)
}

func mapDomainError(err error) mappedError {
	if err == nil {
		return internalMappedError()
	}

	var validationErr *auction.ValidationError
	if errors.As(err, &validationErr) {
		return mappedError{
			status:  http.StatusBadRequest,
			code:    CodeBadRequest,
			message: validationErr.Message,
			log:     false,
		}
	}

	if isTxValidationError(err) {
		return mappedError{
			status:  http.StatusBadRequest,
			code:    CodeBadRequest,
			message: err.Error(),
			log:     false,
		}
	}

	if errors.Is(err, auth.ErrInvalidAddress) {
		return mappedError{
			status:  http.StatusBadRequest,
			code:    CodeBadRequest,
			message: "invalid wallet address",
			log:     false,
		}
	}

	if errors.Is(err, auth.ErrInvalidChainID) {
		return mappedError{
			status:  http.StatusBadRequest,
			code:    CodeBadRequest,
			message: "invalid chain id",
			log:     false,
		}
	}

	if errors.Is(err, auth.ErrInvalidMessage) {
		return mappedError{
			status:  http.StatusBadRequest,
			code:    CodeBadRequest,
			message: "invalid auth message",
			log:     false,
		}
	}

	if errors.Is(err, auth.ErrNonceNotFound) {
		return mappedError{
			status:  http.StatusBadRequest,
			code:    CodeAuthNonceNotFound,
			message: "auth nonce not found",
			log:     false,
		}
	}

	if errors.Is(err, auth.ErrNonceExpired) {
		return mappedError{
			status:  http.StatusBadRequest,
			code:    CodeAuthNonceExpired,
			message: "auth nonce expired",
			log:     false,
		}
	}

	if errors.Is(err, auth.ErrNonceUsed) || errors.Is(err, auth.ErrNonceUnavailable) {
		return mappedError{
			status:  http.StatusBadRequest,
			code:    CodeAuthNonceUsed,
			message: "auth nonce unavailable",
			log:     false,
		}
	}

	if errors.Is(err, auth.ErrInvalidSignature) {
		return mappedError{
			status:  http.StatusUnauthorized,
			code:    CodeInvalidSignature,
			message: "invalid signature",
			log:     false,
		}
	}

	if errors.Is(err, txmodule.ErrActorChainMismatch) {
		return mappedError{
			status:  http.StatusForbidden,
			code:    CodeForbidden,
			message: err.Error(),
			log:     false,
		}
	}

	if errors.Is(err, txmodule.ErrUnauthorized) {
		return mappedError{
			status:  http.StatusUnauthorized,
			code:    CodeUnauthorized,
			message: "unauthorized",
			log:     false,
		}
	}

	if errors.Is(err, auth.ErrUnauthorized) {
		return mappedError{
			status:  http.StatusUnauthorized,
			code:    CodeUnauthorized,
			message: "unauthorized",
			log:     false,
		}
	}

	if errors.Is(err, auth.ErrInvalidToken) {
		return mappedError{
			status:  http.StatusUnauthorized,
			code:    CodeInvalidToken,
			message: "invalid token",
			log:     false,
		}
	}

	if errors.Is(err, auth.ErrUserNotFound) {
		return mappedError{
			status:  http.StatusNotFound,
			code:    CodeUserNotFound,
			message: "user not found",
			log:     false,
		}
	}

	if errors.Is(err, auth.ErrWalletNotFound) {
		return mappedError{
			status:  http.StatusNotFound,
			code:    CodeWalletNotFound,
			message: "wallet not found",
			log:     false,
		}
	}

	if errors.Is(err, auction.ErrAuctionNotFound) {
		return mappedError{
			status:  http.StatusNotFound,
			code:    CodeAuctionNotFound,
			message: "auction not found",
			log:     false,
		}
	}

	if errors.Is(err, auction.ErrAuctionNotActive) ||
		errors.Is(err, auction.ErrAuctionExpired) ||
		errors.Is(err, auction.ErrAuctionNotExpired) {
		return mappedError{
			status:  http.StatusConflict,
			code:    CodeTxPreconditionFailed,
			message: err.Error(),
			log:     false,
		}
	}

	return internalMappedError()
}

func internalMappedError() mappedError {
	return mappedError{
		status:  http.StatusInternalServerError,
		code:    CodeInternalError,
		message: "internal server error",
		log:     true,
	}
}

func isTxValidationError(err error) bool {
	return errors.Is(err, txmodule.ErrInvalidChainID) ||
		errors.Is(err, txmodule.ErrInvalidContractAddress) ||
		errors.Is(err, txmodule.ErrInvalidNFTContract) ||
		errors.Is(err, txmodule.ErrInvalidOperator) ||
		errors.Is(err, txmodule.ErrInvalidTokenID) ||
		errors.Is(err, txmodule.ErrInvalidAuctionID) ||
		errors.Is(err, txmodule.ErrInvalidBidToken) ||
		errors.Is(err, txmodule.ErrInvalidAmount) ||
		errors.Is(err, txmodule.ErrInvalidDuration)
}
