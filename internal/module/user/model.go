package user

import (
	"github.com/Wexyuan/kairos/internal/module/base"
)

// User is the user entity stored in the database.
type User struct {
	base.BaseModel
	Username string     `gorm:"type:varchar(32);not null;uniqueIndex:idx_username;comment:unique username" json:"username"`
	Password string     `gorm:"type:varchar(100);not null;comment:bcrypt password hash" json:"-"`
	Nickname string     `gorm:"type:varchar(32);comment:display name" json:"nickname"`
	Avatar   string     `gorm:"type:varchar(255);comment:avatar url" json:"avatar"`
	Status   UserStatus `gorm:"type:tinyint unsigned;not null;default:1;index:idx_status;comment:user status" json:"status"`
}

// UserStatus defines the user status type.
type UserStatus uint

const (
	// UserStatusNormal marks a normal user that may sign in.
	UserStatusNormal UserStatus = 1
	// UserStatusDisabled marks a disabled user that may not sign in.
	UserStatusDisabled UserStatus = 2
)

// String returns the user status string.
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
