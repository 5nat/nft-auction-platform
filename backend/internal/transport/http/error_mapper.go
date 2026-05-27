package httptransport

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/5nat/nft-auction-platform/backend/internal/modules/auction"
	txmodule "github.com/5nat/nft-auction-platform/backend/internal/modules/tx"
	"github.com/gin-gonic/gin"
)

func writeAuctionError(c *gin.Context, logger *slog.Logger, logMessage string, err error) {
	var validationErr *auction.ValidationError
	if errors.As(err, &validationErr) {
		Error(c, http.StatusBadRequest, CodeBadRequest, validationErr.Message)
		return
	}

	if errors.Is(err, auction.ErrAuctionNotFound) {
		Error(c, http.StatusNotFound, CodeAuctionNotFound, "auction not found")
		return
	}

	if logger != nil {
		logger.Error(logMessage, "error", err)
	}

	Error(c, http.StatusInternalServerError, CodeInternalError, "internal server error")
}

func writeTxError(c *gin.Context, logger *slog.Logger, logMessage string, err error) {
	if isTxValidationError(err) {
		Error(c, http.StatusBadRequest, CodeBadRequest, err.Error())
		return
	}

	if errors.Is(err, auction.ErrAuctionNotFound) {
		Error(c, http.StatusNotFound, CodeAuctionNotFound, "auction not found")
		return
	}

	if errors.Is(err, auction.ErrAuctionNotActive) ||
		errors.Is(err, auction.ErrAuctionExpired) ||
		errors.Is(err, auction.ErrAuctionNotExpired) {
		Error(c, http.StatusConflict, CodeTxPreconditionFailed, err.Error())
		return
	}

	if logger != nil {
		logger.Error(logMessage, "error", err)
	}

	Error(c, http.StatusInternalServerError, CodeInternalError, "internal server error")
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
