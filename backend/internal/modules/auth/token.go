package auth

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

// TokenClaims 是 Auth Service 内部签发 access_token 所需的业务声明。
// 不直接把 jwt.Claims 暴露给 Service，是为了避免业务逻辑被 JWT 库细节污染。
type TokenClaims struct {
	UserID        uint64
	WalletID      uint64
	WalletAddress string
	ChainID       int64
	IssuedAt      time.Time
	ExpiresAt     time.Time
}

type accessTokenClaims struct {
	UserID        uint64 `json:"user_id"`
	WalletID      uint64 `json:"wallet_id"`
	WalletAddress string `json:"wallet_address"`
	ChainID       int64  `json:"chain_id"`

	jwt.RegisteredClaims
}

// TokenIssuer 定义 access_token 签发能力。
// 后续如果要替换 JWT，或在测试中 mock token 签发，只需要替换这个接口实现。
type TokenIssuer interface {
	IssueAccessToken(ctx context.Context, claims TokenClaims) (string, error)
}

// TokenParser 定义 access_token 解析能力。
// AuthMiddleware 会通过它把 Authorization Header 中的 token 解析为 CurrentUser。
type TokenParser interface {
	ParseAccessToken(ctx context.Context, rawToken string) (*CurrentUser, error)
}

// TokenManager 同时具备签发和解析 token 的能力。
// 当前 JWTIssuer 会实现这个接口。
type TokenManager interface {
	TokenIssuer
	TokenParser
}

// JWTIssuer 使用 HS256 签发 JWT。
type JWTIssuer struct {
	secret []byte
}

func NewJWTIssuer(secret string) *JWTIssuer {
	secret = strings.TrimSpace(secret)
	if secret == "" {
		secret = "dev-secret-change-me"
	}

	return &JWTIssuer{
		secret: []byte(secret),
	}
}

// IssueAccessToken 签发 token
func (i *JWTIssuer) IssueAccessToken(ctx context.Context, claims TokenClaims) (string, error) {
	if i == nil || len(i.secret) == 0 {
		return "", fmt.Errorf("jwt secret is not configured")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, accessTokenClaims{
		UserID:        claims.UserID,
		WalletID:      claims.WalletID,
		WalletAddress: claims.WalletAddress,
		ChainID:       claims.ChainID,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(claims.IssuedAt.UTC()),
			ExpiresAt: jwt.NewNumericDate(claims.ExpiresAt.UTC()),
		},
	})

	signed, err := token.SignedString(i.secret)
	if err != nil {
		return "", fmt.Errorf("sign jwt access token: %w", err)
	}

	return signed, nil
}

// ParseAccessToken 解析并校验 access_token。
// 如果 token 过期、签名错误、格式错误或缺少关键字段，统一返回 ErrInvalidToken。
func (i *JWTIssuer) ParseAccessToken(ctx context.Context, rawToken string) (*CurrentUser, error) {
	if i == nil || len(i.secret) == 0 {
		return nil, fmt.Errorf("jwt secret is not configured")
	}

	rowToken := strings.TrimSpace(rawToken)
	if rowToken == "" {
		return nil, ErrInvalidToken
	}

	var claims accessTokenClaims

	token, err := jwt.ParseWithClaims(rowToken, &claims, func(token *jwt.Token) (any, error) {
		if token.Method != jwt.SigningMethodHS256 {
			return nil, ErrInvalidToken
		}

		return i.secret, nil
	})
	if err != nil {
		return nil, ErrInvalidToken
	}

	if token == nil || !token.Valid {
		return nil, ErrInvalidToken
	}

	if claims.UserID == 0 || claims.WalletID == 0 || claims.WalletAddress == "" || claims.ChainID == 0 {
		return nil, ErrInvalidToken
	}

	return &CurrentUser{
		UserID:        claims.UserID,
		WalletID:      claims.WalletID,
		WalletAddress: claims.WalletAddress,
		ChainID:       claims.ChainID,
	}, nil
}
