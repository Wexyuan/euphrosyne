package user

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/Wexyuan/kairos/internal/ecode"
	"github.com/Wexyuan/kairos/internal/module/base"
)

// UserHandler exposes the user HTTP endpoints.
type UserHandler struct {
	*base.Handler
	userSvc UserService // user business operations
}

func NewUserHandler(baseHandler *base.Handler, userSvc UserService) *UserHandler {
	return &UserHandler{
		Handler: baseHandler,
		userSvc: userSvc,
	}
}

// UpdateUserNickname updates the authenticated user nickname.
func (h *UserHandler) UpdateUserNickname(c *gin.Context) {
	id, ok := h.GetUserID(c)
	if !ok {
		h.Fail(c, http.StatusInternalServerError, ecode.ErrAuthInvalidAccessToken)
		return
	}

	var req UpdateUserNicknameReq
	if !h.BindJSON(c, &req) {
		return
	}

	resp, err := h.userSvc.UpdateUserNickname(c.Request.Context(), id, &req)
	if err != nil {
		h.Fail(c, http.StatusInternalServerError, err)
		return
	}
	h.OK(c, resp)
}

// UpdateUserAvatar updates the authenticated user avatar.
func (h *UserHandler) UpdateUserAvatar(c *gin.Context) {
	id, ok := h.GetUserID(c)
	if !ok {
		h.Fail(c, http.StatusInternalServerError, ecode.ErrAuthInvalidAccessToken)
		return
	}

	var req UpdateUserAvatarReq
	if !h.BindJSON(c, &req) {
		return
	}

	resp, err := h.userSvc.UpdateUserAvatar(c.Request.Context(), id, &req)
	if err != nil {
		h.Fail(c, http.StatusInternalServerError, err)
		return
	}
	h.OK(c, resp)
}

// GetCurrentUser returns the authenticated user info.
func (h *UserHandler) GetCurrentUser(c *gin.Context) {
	id, ok := h.GetUserID(c)
	if !ok {
		h.Fail(c, http.StatusInternalServerError, ecode.ErrAuthInvalidAccessToken)
		return
	}

	resp, err := h.userSvc.GetUserByID(c.Request.Context(), id)
	if err != nil {
		h.Fail(c, http.StatusInternalServerError, err)
		return
	}
	h.OK(c, resp)
}
