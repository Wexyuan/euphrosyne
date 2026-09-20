package user

// UpdateUserNicknameReq is the update user nickname request.
type UpdateUserNicknameReq struct {
	Nickname *string `json:"nickname" binding:"omitempty,min=1,max=32"` // new display name
}

// UpdateUserAvatarReq is the update user avatar request.
type UpdateUserAvatarReq struct {
	Avatar string `json:"avatar" binding:"required,url"` // new avatar URL
}

// UserInfoResp is the user info response.
type UserInfoResp struct {
	ID        int64  `json:"id"`         // user ID
	Username  string `json:"username"`   // username
	Nickname  string `json:"nickname"`   // display name
	Avatar    string `json:"avatar"`     // avatar URL
	Status    string `json:"status"`     // user status
	CreatedAt string `json:"created_at"` // creation time
	UpdatedAt string `json:"updated_at"` // last update time
}
