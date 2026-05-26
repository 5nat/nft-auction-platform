package httptransport

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func parseUint64Param(c *gin.Context, name string) (uint64, bool) {
	raw := c.Param(name)

	value, err := strconv.ParseUint(raw, 10, 64)
	if err != nil {
		Error(c, http.StatusBadRequest, CodeBadRequest, "invalid "+name)
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
		Error(c, http.StatusBadRequest, CodeBadRequest, "invalid "+name)
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
		Error(c, http.StatusBadRequest, CodeBadRequest, "invalid "+name)
		return 0, false
	}

	return value, true
}
