package httptransport

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func HealthHandler(deps Dependencies) gin.HandlerFunc {
	return func(c *gin.Context) {
		now := time.Now()

		if deps.DB == nil || deps.DB.Gorm == nil {
			if deps.Logger != nil {
				deps.Logger.Error("health check failed", "reason", "database is nil")
			}

			Error(c, http.StatusServiceUnavailable, CodeDatabaseDown, "database is unavailable")
			return
		}

		sqlDB, err := deps.DB.Gorm.DB()
		if err != nil {
			if deps.Logger != nil {
				deps.Logger.Error("get sql db failed", "error", err)
			}

			Error(c, http.StatusServiceUnavailable, CodeDatabaseDown, "database is unavailable")
			return
		}

		if pingErr := sqlDB.PingContext(c.Request.Context()); pingErr != nil {
			if deps.Logger != nil {
				deps.Logger.Error("database ping failed", "error", pingErr)
			}

			Error(c, http.StatusServiceUnavailable, CodeDatabaseDown, "database is unavailable")
			return
		}

		uptimeSeconds := int64(0)
		if !deps.StartedAt.IsZero() {
			uptimeSeconds = int64(now.Sub(deps.StartedAt).Seconds())
		}

		OK(c, gin.H{
			"status":         "ok",
			"database":       "ok",
			"started_at":     deps.StartedAt,
			"now":            now,
			"uptime_seconds": uptimeSeconds,
		})
	}
}
