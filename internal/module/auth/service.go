package auth

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"gorm.io/gorm"

	"github.com/Wexyuan/euphrosyne/internal/ecode"
	"github.com/Wexyuan/euphrosyne/internal/module/base"
	"github.com/Wexyuan/euphrosyne/internal/module/user"
)

// AuthService defines the authentication business operations.
type AuthService interface {
	// Register registers a user.
	Register(ctx context.Context, req *RegisterReq) (*user.UserInfoResp, error)
	// Login authenticates the user.
	Login(ctx context.Context, req *LoginReq) (*TokenInfoResp, error)
	// Refresh refreshes the token pair.
	Refresh(ctx context.Context, req *RefreshReq) (*TokenInfoResp, error)
	// Logout logs out the user.
	Logout(ctx context.Context, userID int64) error
	// ChangeUserPassword changes the user password.
	ChangeUserPassword(ctx context.Context, userID int64, req *ChangeUserPasswordReq) error
}

// authService implements AuthService.
type authService struct {
	*base.Service
	userRepo user.UserRepository // user data operations
	authRepo AuthRepository      // refresh token data operations
	tokens   *TokenManager       // token manager
}

func NewAuthService(
	baseSvc *base.Service,
	userRepo user.UserRepository,
	authRepo AuthRepository,
	tokens *TokenManager,
) AuthService {
	return &authService{
		Service:  baseSvc,
		userRepo: userRepo,
		authRepo: authRepo,
		tokens:   tokens,
	}
}

// Register registers a user.
func (s *authService) Register(ctx context.Context, req *RegisterReq) (*user.UserInfoResp, error) {
	username := strings.TrimSpace(req.Username)

	taken, err := s.userRepo.ExistsUserByUsername(ctx, username)
	if err != nil {
		return nil, fmt.Errorf("[auth] check user exists error: username=%s: %w", username, err)
	}
	if taken {
		return nil, ecode.ErrUserExists
	}

	hash, err := user.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	u := &user.User{
		Username: username,
		Password: hash,
		Nickname: username,
		Status:   user.UserStatusNormal,
	}
	if err := s.userRepo.CreateUser(ctx, u); err != nil {
		// Report the unique index violation as the same error as the check above
		// so a concurrent registration never surfaces as an internal failure.
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return nil, ecode.ErrUserExists
		}
		return nil, fmt.Errorf("[auth] create user error: username=%s: %w", username, err)
	}
	return toUserInfoResp(u), nil
}

// Login authenticates the user.
func (s *authService) Login(ctx context.Context, req *LoginReq) (*TokenInfoResp, error) {
	username := strings.TrimSpace(req.Username)

	u, err := s.userRepo.GetUserByUsername(ctx, username)
	if err != nil {
		return nil, fmt.Errorf("[auth] get user by username error: username=%s: %w", username, err)
	}
	if u == nil {
		return nil, ecode.ErrAuthInvalidCredentials
	}

	if err := user.VerifyPassword(u.Password, req.Password); err != nil {
		return nil, ecode.ErrAuthInvalidCredentials
	}
	// Check the status last so a wrong password never reveals the account state.
	if u.Status != user.UserStatusNormal {
		return nil, ecode.ErrUserDisabled
	}
	return s.issueTokens(ctx, u.ID, u.Username)
}

// Refresh refreshes the token pair.
func (s *authService) Refresh(ctx context.Context, req *RefreshReq) (*TokenInfoResp, error) {
	userID, ok := s.tokens.ParseRefreshToken(req.RefreshToken)
	if !ok {
		return nil, ecode.ErrAuthInvalidRefreshToken
	}

	stored, err := s.authRepo.GetRefreshToken(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("[auth] get refresh token error: user_id=%d: %w", userID, err)
	}
	// An empty digest means the session is gone: logged out, rotated or expired.
	if stored == "" || stored != s.tokens.HashRefreshToken(req.RefreshToken) {
		return nil, ecode.ErrAuthInvalidRefreshToken
	}

	u, err := s.userRepo.GetUserByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("[auth] get user by id error: id=%d: %w", userID, err)
	}
	if u == nil {
		return nil, ecode.ErrAuthInvalidRefreshToken
	}
	return s.issueTokens(ctx, u.ID, u.Username)
}

// Logout logs out the user.
func (s *authService) Logout(ctx context.Context, userID int64) error {
	if err := s.authRepo.DeleteRefreshToken(ctx, userID); err != nil {
		return fmt.Errorf("[auth] delete refresh token error: user_id=%d: %w", userID, err)
	}
	return nil
}

// ChangeUserPassword changes the user password.
func (s *authService) ChangeUserPassword(ctx context.Context, userID int64, req *ChangeUserPasswordReq) error {
	u, err := s.userRepo.GetUserByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("[auth] get user by id error: id=%d: %w", userID, err)
	}
	if u == nil {
		return ecode.ErrUserNotFound
	}

	if err := user.VerifyPassword(u.Password, req.OldPassword); err != nil {
		return ecode.ErrAuthInvalidOldPassword
	}

	hash, err := user.HashPassword(req.NewPassword)
	if err != nil {
		return err
	}

	// Drop the session before the password write so a failure here cannot leave
	// a live session next to a changed password.
	if err := s.authRepo.DeleteRefreshToken(ctx, u.ID); err != nil {
		return fmt.Errorf("[auth] delete refresh token error: user_id=%d: %w", u.ID, err)
	}
	if err := s.userRepo.UpdateUserPassword(ctx, u.ID, hash); err != nil {
		return fmt.Errorf("[auth] update user password error: id=%d: %w", u.ID, err)
	}
	return nil
}

// issueTokens issues a token pair.
func (s *authService) issueTokens(ctx context.Context, userID int64, username string) (*TokenInfoResp, error) {
	access, ttl, err := s.tokens.GenerateAccessToken(userID, username)
	if err != nil {
		return nil, err
	}

	refresh, err := s.tokens.GenerateRefreshToken(userID)
	if err != nil {
		return nil, err
	}
	if err := s.authRepo.SetRefreshToken(ctx, userID, s.tokens.HashRefreshToken(refresh), s.tokens.RefreshTTL()); err != nil {
		return nil, fmt.Errorf("[auth] set refresh token error: user_id=%d: %w", userID, err)
	}
	return &TokenInfoResp{
		AccessToken:  access,
		RefreshToken: refresh,
		ExpiresIn:    int64(ttl.Seconds()),
	}, nil
}

// toUserInfoResp converts the user entity into the user info response.
func toUserInfoResp(u *user.User) *user.UserInfoResp {
	return &user.UserInfoResp{
		ID:        u.ID,
		Username:  u.Username,
		Nickname:  u.Nickname,
		Avatar:    u.Avatar,
		Status:    u.Status.String(),
		CreatedAt: u.CreatedAt.Format(time.RFC3339),
		UpdatedAt: u.UpdatedAt.Format(time.RFC3339),
	}
}
