package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/Wexyuan/euphrosyne/internal/ecode"
	"github.com/Wexyuan/euphrosyne/internal/module/base"
)

// AuthHandler exposes the authentication HTTP endpoints.
type AuthHandler struct {
	*base.Handler
	authSvc AuthService // auth business operations
}

func NewAuthHandler(baseHandler *base.Handler, authSvc AuthService) *AuthHandler {
	return &AuthHandler{
		Handler: baseHandler,
		authSvc: authSvc,
	}
}

// Register registers a user.
func (h *AuthHandler) Register(c *gin.Context) {
	var req RegisterReq
	if !h.BindJSON(c, &req) {
		return
	}

	resp, err := h.authSvc.Register(c.Request.Context(), &req)
	if err != nil {
		h.Fail(c, http.StatusInternalServerError, err)
		return
	}
	h.OK(c, resp)
}

// Login authenticates the user.
func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginReq
	if !h.BindJSON(c, &req) {
		return
	}

	resp, err := h.authSvc.Login(c.Request.Context(), &req)
	if err != nil {
		h.Fail(c, http.StatusInternalServerError, err)
		return
	}
	h.OK(c, resp)
}

// Refresh refreshes the token pair.
func (h *AuthHandler) Refresh(c *gin.Context) {
	var req RefreshReq
	if !h.BindJSON(c, &req) {
		return
	}

	resp, err := h.authSvc.Refresh(c.Request.Context(), &req)
	if err != nil {
		h.Fail(c, http.StatusInternalServerError, err)
		return
	}
	h.OK(c, resp)
}

// Logout logs out the user.
func (h *AuthHandler) Logout(c *gin.Context) {
	id, ok := h.GetUserID(c)
	if !ok {
		h.Fail(c, http.StatusInternalServerError, ecode.ErrAuthInvalidAccessToken)
		return
	}

	if err := h.authSvc.Logout(c.Request.Context(), id); err != nil {
		h.Fail(c, http.StatusInternalServerError, err)
		return
	}
	h.OK(c, gin.H{"message": "logged out"})
}

// ChangeUserPassword changes the authenticated user password.
func (h *AuthHandler) ChangeUserPassword(c *gin.Context) {
	id, ok := h.GetUserID(c)
	if !ok {
		h.Fail(c, http.StatusInternalServerError, ecode.ErrAuthInvalidAccessToken)
		return
	}

	var req ChangeUserPasswordReq
	if !h.BindJSON(c, &req) {
		return
	}

	if err := h.authSvc.ChangeUserPassword(c.Request.Context(), id, &req); err != nil {
		h.Fail(c, http.StatusInternalServerError, err)
		return
	}
	h.OK(c, gin.H{"message": "password changed"})
}
