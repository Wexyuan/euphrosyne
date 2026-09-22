package auth

// RegisterReq is the register request.
type RegisterReq struct {
	Username string `json:"username" binding:"required,min=4,max=32"` // unique username
	Password string `json:"password" binding:"required,min=8,max=20"` // plain password
}

// LoginReq is the login request.
type LoginReq struct {
	Username string `json:"username" binding:"required"` // username
	Password string `json:"password" binding:"required"` // plain password
}

// RefreshReq is the refresh token request.
type RefreshReq struct {
	RefreshToken string `json:"refresh_token" binding:"required"` // plain refresh token
}

// ChangeUserPasswordReq is the change user password request.
type ChangeUserPasswordReq struct {
	OldPassword string `json:"old_password" binding:"required"`              // old password
	NewPassword string `json:"new_password" binding:"required,min=8,max=20"` // new password
}

// TokenInfoResp is the token pair response.
type TokenInfoResp struct {
	AccessToken  string `json:"access_token"`  // access token
	RefreshToken string `json:"refresh_token"` // refresh token
	ExpiresIn    int64  `json:"expires_in"`    // access token lifetime in seconds
}
