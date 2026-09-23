package router

import "github.com/gin-gonic/gin"

// setupUserRoutes sets up the user routes.
func (r *Router) setupUserRoutes(api *gin.RouterGroup) {
	handler := r.user

	// Register the authenticated endpoints.
	private := api.Group("/users")
	private.Use(r.authRequired)
	private.GET("/me", handler.GetCurrentUser)
	private.PUT("/me/nickname", handler.UpdateUserNickname)
	private.PUT("/me/avatar", handler.UpdateUserAvatar)
}
