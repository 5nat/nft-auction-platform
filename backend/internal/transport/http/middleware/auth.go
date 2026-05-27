package middleware

import (
	"log/slog"
	"strings"

	"github.com/5nat/nft-auction-platform/backend/internal/modules/auth"
	"github.com/gin-gonic/gin"
)

const CurrentUserKey = "current_user"

type ErrorWriter func(c *gin.Context, logger *slog.Logger, logMessage string, err error)

// Auth 是 JWT 鉴权中间件。
// 它从 Authorization: Bearer <token> 中解析 token，
// 校验成功后把 auth.CurrentUser 放入 gin.Context，供后续 Handler 使用。
func Auth(service *auth.Service, logger *slog.Logger, writeError ErrorWriter) gin.HandlerFunc {
	return func(c *gin.Context) {
		if service == nil {
			writeError(c, logger, "auth service is nil", auth.ErrUnauthorized)
			c.Abort()
			return
		}

		rawToken, ok := extractBearerToken(c.GetHeader("Authorization"))
		if !ok {
			writeError(c, logger, "missing authorization token", auth.ErrUnauthorized)
			c.Abort()
			return
		}

		currentUser, err := service.AuthenticateToken(c.Request.Context(), rawToken)
		if err != nil {
			writeError(c, logger, "authenticate token failed", err)
			c.Abort()
			return
		}

		c.Set(CurrentUserKey, currentUser)
		c.Next()
	}
}

func CurrentUser(c *gin.Context) (*auth.CurrentUser, bool) {
	value, exists := c.Get(CurrentUserKey)
	if !exists {
		return nil, false
	}

	currentUser, ok := value.(*auth.CurrentUser)
	return currentUser, ok
}

func extractBearerToken(header string) (string, bool) {
	header = strings.TrimSpace(header)
	if header == "" {
		return "", false
	}

	const prefix = "Bearer "
	if !strings.HasPrefix(header, prefix) {
		return "", false
	}

	token := strings.TrimSpace(strings.TrimPrefix(header, prefix))
	if token == "" {
		return "", false
	}

	return token, true
}
