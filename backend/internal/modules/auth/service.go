package auth

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/common"
)

// Clock 抽象当前时间来源。
// Service 不直接调用 time.Now()，是为了让同一次业务操作中的时间一致，
// 同时方便后续单元测试注入固定时间。
type Clock interface {
	Now() time.Time
}

// RealClock 是生产环境使用的真实时钟。
type RealClock struct{}

func (RealClock) Now() time.Time {
	return time.Now().UTC()
}

type Service struct {
	users  UserRepository
	nonces NonceStore
	clock  Clock

	tokenIssuer    TokenIssuer
	tokenParser    TokenParser
	accessTokenTTL time.Duration

	domain    string
	uri       string
	statement string
	nonceTTL  time.Duration
}

type ServiceConfig struct {
	Domain    string
	URI       string
	Statement string
	NonceTTL  time.Duration

	Clock Clock

	JWTSecret string

	AccessTokenTTL time.Duration
	TokenIssuer    TokenIssuer
	TokenParser    TokenParser
}

func NewService(users UserRepository, nonces NonceStore, cfg ServiceConfig) *Service {
	clock := cfg.Clock
	if clock == nil {
		clock = RealClock{}
	}

	ttl := cfg.NonceTTL
	if ttl <= 0 {
		ttl = 10 * time.Minute
	}

	accessTokenTTL := cfg.AccessTokenTTL
	if accessTokenTTL <= 0 {
		accessTokenTTL = 24 * time.Hour
	}

	statement := cfg.Statement
	if statement == "" {
		statement = "Sign in to NFT Auction Platform"
	}

	tokenIssuer := cfg.TokenIssuer
	tokenParser := cfg.TokenParser
	if tokenIssuer == nil || tokenParser == nil {
		jwtIssuer := NewJWTIssuer(cfg.JWTSecret)

		if tokenIssuer == nil {
			tokenIssuer = jwtIssuer
		}

		if tokenParser == nil {
			tokenParser = jwtIssuer
		}
	}

	return &Service{
		users:  users,
		nonces: nonces,
		clock:  clock,

		tokenIssuer:    tokenIssuer,
		tokenParser:    tokenParser,
		accessTokenTTL: accessTokenTTL,

		domain:    strings.TrimSpace(cfg.Domain),
		uri:       strings.TrimSpace(cfg.URI),
		statement: statement,
		nonceTTL:  ttl,
	}
}

func (s *Service) CreateNonce(ctx context.Context, cmd CreateNonceCommand) (*NonceChallenge, error) {
	if s == nil {
		return nil, fmt.Errorf("auth service is nil")
	}

	if s.clock == nil {
		return nil, fmt.Errorf("auth service clock is nil")
	}

	if s.nonces == nil {
		return nil, fmt.Errorf("auth nonce store is nil")
	}

	address, err := NormalizeAddress(cmd.Address)
	if err != nil {
		return nil, err
	}

	if cmd.ChainID <= 0 {
		return nil, ErrInvalidChainID
	}

	if s.domain == "" || s.uri == "" {
		return nil, fmt.Errorf("auth service domain or URI is empty")
	}

	now := s.clock.Now().UTC()
	expiresAt := now.Add(s.nonceTTL)

	nonce, err := GenerateNonce()
	if err != nil {
		return nil, fmt.Errorf("generate auth nonce: %w", err)
	}

	message := BuildSIWEMessage(
		s.domain,
		address,
		s.statement,
		s.uri,
		cmd.ChainID,
		nonce,
		now,
		expiresAt,
	)

	authNonce := &AuthNonce{
		Address:     address,
		ChainID:     cmd.ChainID,
		Nonce:       nonce,
		MessageHash: HashMessage(message),
		Message:     message,
		Domain:      s.domain,
		URI:         s.uri,
		IssuedAt:    now,
		ExpiresAt:   expiresAt,
		CreatedAt:   now,
	}

	if saveAuthErr := s.nonces.Save(ctx, authNonce); saveAuthErr != nil {
		return nil, fmt.Errorf("save auth nonce: %w", saveAuthErr)
	}

	return &NonceChallenge{
		Nonce:     nonce,
		Message:   message,
		ExpiresAt: expiresAt,
	}, nil
}

// NormalizeAddress 校验并规范化钱包地址。
func NormalizeAddress(address string) (string, error) {
	address = strings.TrimSpace(address)
	if !common.IsHexAddress(address) {
		return "", ErrInvalidAddress
	}

	return common.HexToAddress(address).Hex(), nil
}

// GenerateNonce 生成一次性随机 nonce。
func GenerateNonce() (string, error) {
	var buf [16]byte
	if _, err := rand.Read(buf[:]); err != nil {
		return "", err
	}

	return hex.EncodeToString(buf[:]), nil
}

// Verify 是 Web3 登录闭环的核心
// 1. 根据 message_hash 找回服务端生成的登录挑战
// 2. 检查 nonce 是否过期、是否已经使用
// 3. 从 signature 中恢复钱包地址，证明用户确实控制该钱包
// 4. 查找或创建平台用户或钱包绑定
// 5. 原子消费 nonce，防止同一签名重复登录
// 6. 签发 JWT access_token
func (s *Service) Verify(ctx context.Context, cmd VerifyCommand) (*LoginResult, error) {
	if s == nil {
		return nil, fmt.Errorf("auth service is nil")
	}

	if s.clock == nil {
		return nil, fmt.Errorf("auth service clock is nil")
	}

	if s.users == nil {
		return nil, fmt.Errorf("auth service users is nil")
	}

	if s.nonces == nil {
		return nil, fmt.Errorf("auth nonce store is nil")
	}

	if s.tokenIssuer == nil {
		return nil, fmt.Errorf("auth service token issuer is nil")
	}

	address, err := NormalizeAddress(cmd.Address)
	if err != nil {
		return nil, err
	}

	if cmd.ChainID <= 0 {
		return nil, ErrInvalidChainID
	}

	if strings.TrimSpace(cmd.Message) == "" {
		return nil, ErrInvalidMessage
	}

	if strings.TrimSpace(cmd.Signature) == "" {
		return nil, ErrInvalidSignature
	}

	now := s.clock.Now().UTC()
	messageHash := HashMessage(cmd.Message)

	authNonce, err := s.nonces.FindByMessageHash(ctx, messageHash)
	if err != nil {
		return nil, err
	}

	if authNonce.Message != cmd.Message {
		return nil, ErrInvalidMessage
	}

	if !strings.EqualFold(authNonce.Address, address) || authNonce.ChainID != cmd.ChainID {
		return nil, ErrInvalidMessage
	}

	if authNonce.UsedAt != nil {
		return nil, ErrNonceUsed
	}

	if !now.Before(authNonce.ExpiresAt) {
		return nil, ErrNonceExpired
	}

	recoveredAddress, err := RecoverAddressFromPersonalSignature(cmd.Message, cmd.Signature)
	if err != nil {
		return nil, err
	}

	if !strings.EqualFold(recoveredAddress, address) {
		return nil, ErrInvalidSignature
	}

	user, wallet, err := s.findOrCreateUserWithWallet(ctx, address, cmd.ChainID, now)
	if err != nil {
		return nil, err
	}

	if user.Status != "active" {
		return nil, ErrUnauthorized
	}

	if err := s.nonces.Consume(ctx, authNonce.ID, now); err != nil {
		return nil, err
	}

	expiresAt := now.Add(s.accessTokenTTL)

	accessToken, err := s.tokenIssuer.IssueAccessToken(ctx, TokenClaims{
		UserID:        user.ID,
		WalletID:      wallet.ID,
		WalletAddress: wallet.Address,
		ChainID:       cmd.ChainID,
		IssuedAt:      now,
		ExpiresAt:     expiresAt,
	})

	if err != nil {
		return nil, fmt.Errorf("issue access token: %w", err)
	}

	return &LoginResult{
		AccessToken: accessToken,
		TokenType:   "Bearer",
		ExpiresIn:   int64(s.accessTokenTTL.Seconds()),
		User: &LoginUser{
			ID:            user.ID,
			WalletID:      wallet.ID,
			WalletAddress: wallet.Address,
			ChainID:       cmd.ChainID,
		},
	}, nil
}

func (s *Service) findOrCreateUserWithWallet(ctx context.Context, address string, chainID int64, now time.Time) (*User, *Wallet, error) {
	wallet, err := s.users.FindWalletByAddress(ctx, address, chainID)
	if err == nil {
		user, err := s.users.FindUserByID(ctx, wallet.UserID)
		if err != nil {
			return nil, nil, err
		}
		return user, wallet, nil
	}

	if !errors.Is(err, ErrWalletNotFound) {
		return nil, nil, err
	}

	user, wallet, err := s.users.CreateUserWithWallet(ctx, address, chainID, now)
	if err != nil {
		return nil, nil, err
	}

	return user, wallet, nil
}

// AuthenticateToken 校验 access_token，并返回当前请求用户。
// 这里只信任服务端签发的 JWT，不信任前端直接传来的 user_id 或 wallet_address。
func (s *Service) AuthenticateToken(ctx context.Context, rawToken string) (*CurrentUser, error) {
	if s == nil {
		return nil, fmt.Errorf("auth service is nil")
	}

	if s.tokenParser == nil {
		return nil, fmt.Errorf("auth token parser is nil")
	}

	if s.users == nil {
		return nil, fmt.Errorf("auth user repository is nil")
	}

	currentUser, err := s.tokenParser.ParseAccessToken(ctx, rawToken)
	if err != nil {
		return nil, err
	}

	user, err := s.users.FindUserByID(ctx, currentUser.UserID)
	if err != nil {
		return nil, err
	}

	if user.Status != "active" {
		return nil, ErrUnauthorized
	}

	return currentUser, nil
}

// GetProfile 返回当前用户资料。
// 这个方法由 /me 接口调用，用于确认 JWT 登录态是否可用。
func (s *Service) GetProfile(ctx context.Context, currentUser *CurrentUser) (*UserProfile, error) {
	if s == nil {
		return nil, fmt.Errorf("auth service is nil")
	}

	if s.users == nil {
		return nil, fmt.Errorf("auth user repository is nil")
	}

	if currentUser == nil || currentUser.UserID == 0 {
		return nil, ErrUnauthorized
	}

	user, err := s.users.FindUserByID(ctx, currentUser.UserID)
	if err != nil {
		return nil, err
	}

	if user.Status != "active" {
		return nil, ErrUnauthorized
	}

	wallets, err := s.users.ListWalletsByUserID(ctx, currentUser.UserID)
	if err != nil {
		return nil, err
	}

	return &UserProfile{
		ID:      user.ID,
		Status:  user.Status,
		Wallets: wallets,
	}, nil
}
