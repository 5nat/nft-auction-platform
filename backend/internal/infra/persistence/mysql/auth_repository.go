package mysql

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/5nat/nft-auction-platform/backend/internal/infra/persistence/mysql/model"
	"github.com/5nat/nft-auction-platform/backend/internal/modules/auth"
	"gorm.io/gorm"
)

var (
	_ auth.UserRepository = (*AuthRepository)(nil)
	_ auth.NonceStore     = (*AuthRepository)(nil)
)

type AuthRepository struct {
	db *gorm.DB
}

func NewAuthRepository(db *gorm.DB) *AuthRepository {
	return &AuthRepository{
		db: db,
	}
}

func (r *AuthRepository) FindWalletByAddress(ctx context.Context, address string, chainID int64) (*auth.Wallet, error) {
	var wallet model.Wallet

	err := r.db.WithContext(ctx).
		Where("address = ? AND chain_id = ?", address, chainID).
		First(&wallet).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, auth.ErrWalletNotFound
		}
		return nil, fmt.Errorf("find wallet by address: %w", err)
	}

	return toAuthWallet(&wallet), err
}

func (r *AuthRepository) CreateUserWithWallet(ctx context.Context, address string, chainID int64, now time.Time) (*auth.User, *auth.Wallet, error) {
	var createdUser model.User
	var createdWallet model.Wallet

	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		createdUser = model.User{
			Status:    "active",
			CreatedAt: now,
			UpdatedAt: now,
		}

		if err := tx.Create(&createdUser).Error; err != nil {
			return fmt.Errorf("create user: %w", err)
		}

		createdWallet = model.Wallet{
			UserID:    createdUser.ID,
			Address:   address,
			ChainID:   chainID,
			IsPrimary: true,
			CreatedAt: now,
			UpdatedAt: now,
		}
		if err := tx.Create(&createdWallet).Error; err != nil {
			return fmt.Errorf("create wallet: %w", err)
		}

		return nil
	})
	if err != nil {
		return nil, nil, fmt.Errorf("create user with wallet: %w", err)
	}

	return toAuthUser(&createdUser), toAuthWallet(&createdWallet), nil
}

func (r *AuthRepository) FindUserByID(ctx context.Context, userID uint64) (*auth.User, error) {
	var user model.User

	err := r.db.WithContext(ctx).
		Where("id = ?", userID).
		First(&user).
		Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, auth.ErrUserNotFound
		}
		return nil, fmt.Errorf("find user by id: %w", err)
	}

	return toAuthUser(&user), nil
}

func (r *AuthRepository) ListWalletsByUserID(ctx context.Context, userID uint64) ([]*auth.Wallet, error) {
	var wallets []*model.Wallet

	err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Order("is_primary DESC, id ASC").
		Find(&wallets).Error

	if err != nil {
		return nil, fmt.Errorf("list wallets by user id: %w", err)
	}

	result := make([]*auth.Wallet, 0, len(wallets))
	for _, wallet := range wallets {
		result = append(result, toAuthWallet(wallet))
	}

	return result, nil
}

func (r *AuthRepository) Save(ctx context.Context, n *auth.AuthNonce) error {
	if n == nil {
		return fmt.Errorf("save auth nonce: nil nonce")
	}

	nonce := model.AuthNonce{
		Address:     n.Address,
		ChainID:     n.ChainID,
		Nonce:       n.Nonce,
		MessageHash: n.MessageHash,
		Message:     n.Message,
		Domain:      n.Domain,
		URI:         n.URI,
		IssuedAt:    n.IssuedAt,
		ExpiresAt:   n.ExpiresAt,
		UsedAt:      n.UsedAt,
		CreatedAt:   n.CreatedAt,
	}

	if err := r.db.WithContext(ctx).Create(&nonce).Error; err != nil {
		return fmt.Errorf("save auth nonce: %w", err)
	}

	return nil
}

func (r *AuthRepository) FindByMessageHash(ctx context.Context, messageHash string) (*auth.AuthNonce, error) {
	var nonce model.AuthNonce

	err := r.db.WithContext(ctx).
		Where("message_hash = ?", messageHash).
		First(&nonce).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, auth.ErrNonceNotFound
		}
		return nil, fmt.Errorf("find auth nonce by message hash: %w", err)
	}

	return toAuthNonce(&nonce), nil
}

func (r *AuthRepository) Consume(ctx context.Context, id uint64, now time.Time) error {
	result := r.db.WithContext(ctx).
		Model(&model.AuthNonce{}).
		Where("id = ? AND used_at IS NULL AND expires_at > ?", id, now).
		Update("used_at", now)

	if result.Error != nil {
		return fmt.Errorf("consume auth nonce: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return auth.ErrNonceUnavailable
	}

	return nil
}

func toAuthUser(user *model.User) *auth.User {
	if user == nil {
		return nil
	}

	return &auth.User{
		ID:        user.ID,
		Status:    user.Status,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}

func toAuthWallet(wallet *model.Wallet) *auth.Wallet {
	if wallet == nil {
		return nil
	}

	return &auth.Wallet{
		ID:        wallet.ID,
		UserID:    wallet.UserID,
		Address:   wallet.Address,
		ChainID:   wallet.ChainID,
		IsPrimary: wallet.IsPrimary,
		CreatedAt: wallet.CreatedAt,
		UpdatedAt: wallet.UpdatedAt,
	}
}

func toAuthNonce(nonce *model.AuthNonce) *auth.AuthNonce {
	if nonce == nil {
		return nil
	}

	return &auth.AuthNonce{
		ID:          nonce.ID,
		Address:     nonce.Address,
		ChainID:     nonce.ChainID,
		Nonce:       nonce.Nonce,
		MessageHash: nonce.MessageHash,
		Message:     nonce.Message,
		Domain:      nonce.Domain,
		URI:         nonce.URI,
		IssuedAt:    nonce.IssuedAt,
		ExpiresAt:   nonce.ExpiresAt,
		UsedAt:      nonce.UsedAt,
		CreatedAt:   nonce.CreatedAt,
	}
}
