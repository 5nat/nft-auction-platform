package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

const (
	CodeOk           = 0
	CodeBadRequest   = 40000
	CodeUnauthorized = 40100
	CodeNotFound     = 40400
	CodeDatabaseDown = 50001
	CodeInternal     = 50000
)

type Response struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

func JSON(c *gin.Context, httpStatus int, code int, message string, data any) {
	c.JSON(httpStatus, Response{
		Code:    code,
		Message: message,
		Data:    data,
	})
}

func OK(c *gin.Context, data any) {
	JSON(c, http.StatusOK, CodeOk, "ok", data)
}

func Error(c *gin.Context, httpStatus int, code int, message string) {
	JSON(c, httpStatus, code, message, nil)
}
