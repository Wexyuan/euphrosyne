package user

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"github.com/Wexyuan/kairos/internal/module/base"
)

// UserRepository defines the user data operations.
type UserRepository interface {
	// CreateUser creates the user.
	CreateUser(ctx context.Context, user *User) error
	// UpdateUserPassword updates the password hash.
	UpdateUserPassword(ctx context.Context, id int64, hash string) error
	// UpdateUserNickname updates the nickname.
	UpdateUserNickname(ctx context.Context, id int64, nickname string) error
	// UpdateUserAvatar updates the avatar url.
	UpdateUserAvatar(ctx context.Context, id int64, avatar string) error
	// GetUserByID retrieves the user by ID.
	GetUserByID(ctx context.Context, id int64) (*User, error)
	// GetUserByUsername retrieves the user by username.
	GetUserByUsername(ctx context.Context, username string) (*User, error)
	// ExistsUserByUsername reports whether the user exists by username.
	ExistsUserByUsername(ctx context.Context, username string) (bool, error)
}

// userRepository implements UserRepository.
type userRepository struct {
	*base.Repository
}

func NewUserRepository(baseRepo *base.Repository) UserRepository {
	return &userRepository{Repository: baseRepo}
}

// CreateUser creates the user.
func (r *userRepository) CreateUser(ctx context.Context, user *User) error {
	user.ID = r.GenerateID()
	return r.DB(ctx).Create(user).Error
}

// UpdateUserPassword updates the password hash.
func (r *userRepository) UpdateUserPassword(ctx context.Context, id int64, hash string) error {
	return r.DB(ctx).Model(&User{}).Where("id = ?", id).Update("password", hash).Error
}

// UpdateUserNickname updates the nickname.
func (r *userRepository) UpdateUserNickname(ctx context.Context, id int64, nickname string) error {
	return r.DB(ctx).Model(&User{}).Where("id = ?", id).Update("nickname", nickname).Error
}

// UpdateUserAvatar updates the avatar url.
func (r *userRepository) UpdateUserAvatar(ctx context.Context, id int64, avatar string) error {
	return r.DB(ctx).Model(&User{}).Where("id = ?", id).Update("avatar", avatar).Error
}

// GetUserByID retrieves the user by ID.
func (r *userRepository) GetUserByID(ctx context.Context, id int64) (*User, error) {
	var user User
	err := r.DB(ctx).Where("id = ?", id).First(&user).Error
	if err != nil {
		// Return nil instead of an error so callers treat a missing row as a state.
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

// GetUserByUsername retrieves the user by username.
func (r *userRepository) GetUserByUsername(ctx context.Context, username string) (*User, error) {
	var user User
	err := r.DB(ctx).Where("username = ?", username).First(&user).Error
	if err != nil {
		// Return nil instead of an error so callers treat a missing row as a state.
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

// ExistsUserByUsername reports whether the user exists by username.
func (r *userRepository) ExistsUserByUsername(ctx context.Context, username string) (bool, error) {
	var count int64
	if err := r.DB(ctx).Model(&User{}).Where("username = ?", username).Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}
