package httptransport

import (
	"net/http"

	"github.com/5nat/nft-auction-platform/backend/internal/transport/http/middleware"
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
	CodeSuccess = 0

	CodeForbidden = 40301

	// CodeBadRequest 通用请求错误：参数缺失、格式错误、非法 chain_id、非法 address 等。
	CodeBadRequest = 40001

	// CodeUnauthorized Auth 相关错误。
	// 这里单独拆分 Auth 错误码，是因为钱包登录涉及 nonce、签名、JWT 等安全链路，
	// 前端和后端排查问题时需要比普通 BAD_REQUEST 更明确的错误类型。
	CodeUnauthorized     = 40101
	CodeInvalidToken     = 40102
	CodeInvalidSignature = 40103

	CodeAuthNonceNotFound = 40011
	CodeAuthNonceExpired  = 40012
	CodeAuthNonceUsed     = 40013

	CodeUserNotFound   = 40411
	CodeWalletNotFound = 40412

	// CodeAuctionNotFound Auction 相关错误。
	CodeAuctionNotFound = 40401

	// CodeTxPreconditionFailed Tx / Policy 前置条件失败。
	// 例如拍卖已过期、拍卖未激活、未到结束时间等。
	CodeTxPreconditionFailed = 40901

	// CodeInternalError 通用内部错误。
	CodeInternalError = 50000

	// CodeDatabaseDown 依赖服务不可用。
	CodeDatabaseDown = 50301
)

type Response struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data"`
	Meta    any    `json:"meta,omitempty"`
}

type ResponseMeta struct {
	RequestID string `json:"request_id,omitempty"`
}

func OK(c *gin.Context, data any) {
	c.JSON(http.StatusOK, Response{
		Code:    CodeSuccess,
		Message: "success",
		Data:    data,
		Meta:    requestMeta(c),
	})
}

func OKWithMeta(c *gin.Context, data any, meta any) {
	c.JSON(http.StatusOK, Response{
		Code:    CodeSuccess,
		Message: "success",
		Data:    data,
		Meta:    mergeMeta(c, meta),
	})
}

func Error(c *gin.Context, status int, code int, message string) {
	c.JSON(status, Response{
		Code:    code,
		Message: message,
		Data:    nil,
		Meta:    requestMeta(c),
	})
}

func requestMeta(c *gin.Context) ResponseMeta {
	if c == nil {
		return ResponseMeta{}
	}

	return ResponseMeta{
		RequestID: c.GetString(middleware.RequestIDKey),
	}
}

func mergeMeta(c *gin.Context, meta any) any {
	requestID := ""
	if c != nil {
		requestID = c.GetString(middleware.RequestIDKey)
	}

	if requestID == "" {
		return meta
	}

	switch m := meta.(type) {
	case nil:
		return ResponseMeta{
			RequestID: requestID,
		}

	case map[string]any:
		m["request_id"] = requestID
		return m

	case gin.H:
		m["request_id"] = requestID
		return m

	default:
		return gin.H{
			"request_id": requestID,
			"extra":      meta,
		}
	}
}
