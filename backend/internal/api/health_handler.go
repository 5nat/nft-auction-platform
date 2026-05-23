package api

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type HealthResponse struct {
	Status    string `json:"status"`
	Database  string `json:"database"`
	Uptime    string `json:"uptime"`
	Timestamp string `json:"timestamp"`
}

func HealthHandler(deps Dependencies) gin.HandlerFunc {
	return func(c *gin.Context) {
		dbStatus := "ok"
		statusCode := http.StatusOK
		code := CodeOk
		message := "ok"

		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		if err := deps.DB.SQLDB.PingContext(ctx); err != nil {
			dbStatus = "down"
			statusCode = http.StatusServiceUnavailable
			code = CodeDatabaseDown
			message = "database down"
			deps.Logger.Error("database health check failed", "error", err)
		}

		resp := HealthResponse{
			Status:    "ok",
			Database:  dbStatus,
			Uptime:    time.Since(deps.StartedAt).String(),
			Timestamp: time.Now().Format(time.RFC3339),
		}

		if dbStatus != "ok" {
			resp.Status = "degraded"
		}

		JSON(c, statusCode, code, message, resp)
	}
}
