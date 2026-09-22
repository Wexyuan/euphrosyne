package auth

import (
	"context"
	"strconv"
	"time"

	"github.com/Wexyuan/euphrosyne/internal/module/base"
)

// refreshTokenKeyPrefix is the cache key prefix of the refresh token.
const refreshTokenKeyPrefix = "auth:refresh:"

// AuthRepository defines the refresh token data operations.
type AuthRepository interface {
	// SetRefreshToken stores the refresh token digest.
	SetRefreshToken(ctx context.Context, userID int64, hash string, ttl time.Duration) error
	// GetRefreshToken loads the refresh token digest.
	GetRefreshToken(ctx context.Context, userID int64) (string, error)
	// DeleteRefreshToken removes the refresh token digest.
	DeleteRefreshToken(ctx context.Context, userID int64) error
}

// authRepository implements AuthRepository.
type authRepository struct {
	*base.Repository
}

func NewAuthRepository(baseRepo *base.Repository) AuthRepository {
	return &authRepository{Repository: baseRepo}
}

// SetRefreshToken stores the refresh token digest.
func (r *authRepository) SetRefreshToken(ctx context.Context, userID int64, hash string, ttl time.Duration) error {
	return r.Cache().Set(ctx, r.tokenKey(userID), hash, int(ttl.Seconds()))
}

// GetRefreshToken loads the refresh token digest.
func (r *authRepository) GetRefreshToken(ctx context.Context, userID int64) (string, error) {
	return r.Cache().Get(ctx, r.tokenKey(userID))
}

// DeleteRefreshToken removes the refresh token digest.
func (r *authRepository) DeleteRefreshToken(ctx context.Context, userID int64) error {
	return r.Cache().Delete(ctx, r.tokenKey(userID))
}

// tokenKey builds the refresh token cache key.
func (r *authRepository) tokenKey(userID int64) string {
	return refreshTokenKeyPrefix + strconv.FormatInt(userID, 10)
}
