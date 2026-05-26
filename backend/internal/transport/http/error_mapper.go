package httptransport

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/5nat/nft-auction-platform/backend/internal/modules/auction"
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
