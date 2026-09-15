package user

import (
	"github.com/Wexyuan/euphrosyne/internal/module/base"
)

// User is the user entity stored in the database.
type User struct {
	base.BaseModel
	// Username is the unique username of the user.
	Username string `gorm:"type:varchar(32);not null;uniqueIndex:idx_username;comment:unique username" json:"username"`
	// Password is the hashed password of the user.
	Password string `gorm:"type:varchar(100);not null;comment:bcrypt password hash" json:"-"`
	// Nickname is the display name shown to other users.
	Nickname string `gorm:"type:varchar(32);comment:display name" json:"nickname"`
	// Avatar is the URL to the user's avatar image.
	Avatar string `gorm:"type:varchar(255);comment:avatar url" json:"avatar"`
	// Status is the user status.
	Status UserStatus `gorm:"type:tinyint unsigned;not null;default:1;index:idx_status;comment:user status" json:"status"`
}

// UserStatus identifies the user lifecycle state.
type UserStatus uint

const (
	// UserStatusNormal marks a normal user that may sign in.
	UserStatusNormal UserStatus = 1
	// UserStatusDisabled marks a disabled user that may not sign in.
	UserStatusDisabled UserStatus = 2
)

// String returns the string representation of the user status.
func (s UserStatus) String() string {
	switch s {
	case UserStatusNormal:
		return "normal"
	case UserStatusDisabled:
		return "disabled"
	default:
		return "unknown"
	}
}
