package user

import (
	"context"
	"fmt"
	"time"

	"github.com/Wexyuan/kairos/internal/ecode"
	"github.com/Wexyuan/kairos/internal/module/base"
)

// UserService defines the user business operations.
type UserService interface {
	// UpdateUserNickname updates the user nickname.
	UpdateUserNickname(ctx context.Context, id int64, req *UpdateUserNicknameReq) (*UserInfoResp, error)
	// UpdateUserAvatar updates the user avatar.
	UpdateUserAvatar(ctx context.Context, id int64, req *UpdateUserAvatarReq) (*UserInfoResp, error)
	// GetUserByID retrieves the user info by ID.
	GetUserByID(ctx context.Context, id int64) (*UserInfoResp, error)
}

// userService implements UserService.
type userService struct {
	*base.Service
	userRepo UserRepository // user data operations
}

func NewUserService(baseSvc *base.Service, userRepo UserRepository) UserService {
	return &userService{
		Service:  baseSvc,
		userRepo: userRepo,
	}
}

// UpdateUserNickname updates the user nickname.
func (s *userService) UpdateUserNickname(ctx context.Context, id int64, req *UpdateUserNicknameReq) (*UserInfoResp, error) {
	if req.Nickname != nil {
		if err := s.userRepo.UpdateUserNickname(ctx, id, *req.Nickname); err != nil {
			return nil, fmt.Errorf("[user] update user nickname error: id=%d: %w", id, err)
		}
	}
	return s.GetUserByID(ctx, id)
}

// UpdateUserAvatar updates the user avatar.
func (s *userService) UpdateUserAvatar(ctx context.Context, id int64, req *UpdateUserAvatarReq) (*UserInfoResp, error) {
	if err := s.userRepo.UpdateUserAvatar(ctx, id, req.Avatar); err != nil {
		return nil, fmt.Errorf("[user] update user avatar error: id=%d: %w", id, err)
	}
	return s.GetUserByID(ctx, id)
}

// GetUserByID retrieves the user info by ID.
func (s *userService) GetUserByID(ctx context.Context, id int64) (*UserInfoResp, error) {
	u, err := s.getUserByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return s.toUserInfoResp(u), nil
}

// getUserByID retrieves the user by ID.
func (s *userService) getUserByID(ctx context.Context, id int64) (*User, error) {
	u, err := s.userRepo.GetUserByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("[user] get user by id error: id=%d: %w", id, err)
	}
	// Translate a missing row into the business error the handlers expect.
	if u == nil {
		return nil, ecode.ErrUserNotFound
	}
	return u, nil
}

// toUserInfoResp converts the user entity into the user info response.
func (s *userService) toUserInfoResp(u *User) *UserInfoResp {
	return &UserInfoResp{
		ID:        u.ID,
		Username:  u.Username,
		Nickname:  u.Nickname,
		Avatar:    u.Avatar,
		Status:    u.Status.String(),
		CreatedAt: u.CreatedAt.Format(time.RFC3339),
		UpdatedAt: u.UpdatedAt.Format(time.RFC3339),
	}
}
