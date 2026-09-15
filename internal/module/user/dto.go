package user

import "time"

// UpdateNicknameReq is the request to update the display name.
type UpdateNicknameReq struct {
	// Nickname is the display name shown to other users.
	Nickname *string `json:"nickname" binding:"omitempty,min=1,max=32"`
}

// UpdateAvatarReq is the request to update the avatar url.
type UpdateAvatarReq struct {
	// Avatar is the URL to the user's avatar image.
	Avatar string `json:"avatar" binding:"required,url"`
}

// UserInfoResp is the response of the user info.
type UserInfoResp struct {
	// ID is the primary key.
	ID int64 `json:"id"`
	// Username is the unique username of the user.
	Username string `json:"username"`
	// Nickname is the display name shown to other users.
	Nickname string `json:"nickname"`
	// Avatar is the URL to the user's avatar image.
	Avatar string `json:"avatar"`
	// Status is the user status.
	Status UserStatus `json:"status"`
	// CreatedAt records the creation time.
	CreatedAt time.Time `json:"created_at"`
	// UpdatedAt records the last update time.
	UpdatedAt time.Time `json:"updated_at"`
}
