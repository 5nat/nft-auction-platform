package middleware

import (
	"log/slog"
	"net/http"
	"runtime/debug"

	"github.com/gin-gonic/gin"
)

type ErrorResponder func(c *gin.Context, status int, code int, message string)

// Recovery 捕获 Handler 中的 panic，避免整个进程崩溃。
// 捕获 panic 后，会记录结构化错误日志，并返回统一 500 响应。
func Recovery(logger *slog.Logger, respond ErrorResponder, internalCode int) gin.HandlerFunc {
	if logger == nil {
		logger = slog.Default()
	}

	return func(c *gin.Context) {
		defer func() {
			if recovered := recover(); recovered != nil {
				requestID := c.GetString(RequestIDKey)

				logger.Error(
					"panic recovered",
					"request_id", requestID,
					"panic", recovered,
					"method", c.Request.Method,
					"path", c.Request.URL.Path,
					"route", c.FullPath(),
					"client_ip", c.ClientIP(),
					"stack", string(debug.Stack()),
				)

				if !c.Writer.Written() && respond != nil {
					respond(c, http.StatusInternalServerError, internalCode, "internal server error")
				}

				c.Abort()
			}
		}()

		c.Next()
	}
}
