package middleware

import (
	"crypto/rand"
	"encoding/hex"
	"log/slog"
	"time"

	"github.com/gin-gonic/gin"
)

const (
	RequestIDKey    = "request_id"
	RequestIDHeader = "X-Request-ID"
)

// Logger 记录每一次 HTTP 请求。
// 它会生成或透传 X-Request-ID，并把 request_id 写入 gin.Context 和响应头。
// 后续排查问题时，可以用 request_id 串联一次请求的所有日志。
func Logger(logger *slog.Logger) gin.HandlerFunc {
	if logger == nil {
		logger = slog.Default()
	}

	return func(c *gin.Context) {
		start := time.Now()

		requestID := c.GetHeader(RequestIDHeader)
		if requestID == "" {
			requestID = generateRequestID()
		}

		c.Set(RequestIDKey, requestID)
		c.Writer.Header().Set(RequestIDHeader, requestID)

		c.Next()

		latency := time.Since(start)

		logger.Info(
			"http request",
			"request_id", requestID,
			"method", c.Request.Method,
			"path", c.Request.URL.Path,
			"route", c.FullPath(),
			"status", c.Writer.Status(),
			"latency", latency.String(),
			"client_ip", c.ClientIP(),
			"user_agent", c.Request.UserAgent(),
			"errors", c.Errors.String(),
		)
	}
}

func generateRequestID() string {
	var buf [16]byte
	if _, err := rand.Read(buf[:]); err != nil {
		return hex.EncodeToString([]byte(time.Now().Format("20060102150405.000000000")))
	}

	return hex.EncodeToString(buf[:])
}
