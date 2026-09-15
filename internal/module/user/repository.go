package user

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"github.com/Wexyuan/euphrosyne/internal/module/base"
)

// UserRepository defines the user data operations.
type UserRepository interface {
	// CreateUser generates ID and creates user in database.
	CreateUser(ctx context.Context, user *User) error
	// UpdateUserPassword updates user password in database.
	UpdateUserPassword(ctx context.Context, id int64, hash string) error
	// UpdateUserNickname updates user nickname in database.
	UpdateUserNickname(ctx context.Context, id int64, nickname string) error
	// UpdateUserAvatar updates user avatar in database.
	UpdateUserAvatar(ctx context.Context, id int64, avatar string) error
	// GetUserByID retrieves user by ID, returns nil if not found.
	GetUserByID(ctx context.Context, id int64) (*User, error)
	// GetUserByUsername retrieves user by username, returns nil if not found.
	GetUserByUsername(ctx context.Context, username string) (*User, error)
	// ExistsUserByUsername checks whether the username exists in database.
	ExistsUserByUsername(ctx context.Context, username string) (bool, error)
}

// userRepository implements UserRepository interface.
type userRepository struct {
	*base.Repository
}

// NewUserRepository creates the user repository.
func NewUserRepository(repo *base.Repository) UserRepository {
	return &userRepository{Repository: repo}
}

// CreateUser generates ID and creates user in database.
func (r *userRepository) CreateUser(ctx context.Context, user *User) error {
	user.ID = r.GenerateID()
	return r.DB(ctx).Create(user).Error
}

// UpdateUserPassword updates user password in database.
func (r *userRepository) UpdateUserPassword(ctx context.Context, id int64, hash string) error {
	return r.DB(ctx).Model(&User{}).Where("id = ?", id).Update("password", hash).Error
}

// UpdateUserNickname updates user nickname in database.
func (r *userRepository) UpdateUserNickname(ctx context.Context, id int64, nickname string) error {
	return r.DB(ctx).Model(&User{}).Where("id = ?", id).Update("nickname", nickname).Error
}

// UpdateUserAvatar updates user avatar in database.
func (r *userRepository) UpdateUserAvatar(ctx context.Context, id int64, avatar string) error {
	return r.DB(ctx).Model(&User{}).Where("id = ?", id).Update("avatar", avatar).Error
}

// GetUserByID retrieves user by ID, returns nil if not found.
func (r *userRepository) GetUserByID(ctx context.Context, id int64) (*User, error) {
	var user User
	err := r.DB(ctx).Where("id = ?", id).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

// GetUserByUsername retrieves user by username, returns nil if not found.
func (r *userRepository) GetUserByUsername(ctx context.Context, username string) (*User, error) {
	var user User
	err := r.DB(ctx).Where("username = ?", username).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

// ExistsUserByUsername checks whether the username exists in database.
func (r *userRepository) ExistsUserByUsername(ctx context.Context, username string) (bool, error) {
	var count int64
	if err := r.DB(ctx).Model(&User{}).Where("username = ?", username).Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}
