package httptransport

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// 统一业务响应风格
/*
	{
	  "code": 0,
	  "message": "success",
	  "data": {}
	}
*/

// REST 风格
/*
	{
	  "data": {}
	}

	{
	  "error": {
		"code": "BAD_REQUEST",
		"message": "invalid chain_id"
	  }
	}
*/

const (
	CodeSuccess         = 0
	CodeBadRequest      = 40001
	CodeAuctionNotFound = 40401
	CodeInternalError   = 50000
	CodeDatabaseDown    = 50301
)

type Response struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data"`
	Meta    any    `json:"meta,omitempty"`
}

func OK(c *gin.Context, data any) {
	c.JSON(http.StatusOK, Response{
		Code:    CodeSuccess,
		Message: "success",
		Data:    data,
	})
}

func OKWithMeta(c *gin.Context, data any, meta any) {
	c.JSON(http.StatusOK, Response{
		Code:    CodeSuccess,
		Message: "success",
		Data:    data,
		Meta:    meta,
	})
}

func Error(c *gin.Context, status int, code int, message string) {
	c.JSON(status, Response{
		Code:    code,
		Message: message,
		Data:    nil,
	})
}
