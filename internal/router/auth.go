package router

import "github.com/gin-gonic/gin"

// setupAuthRoutes sets up the auth routes.
func (r *Router) setupAuthRoutes(api *gin.RouterGroup) {
	handler := r.auth

	// Register the public endpoints.
	public := api.Group("/auth")
	public.POST("/register", handler.Register)
	public.POST("/login", handler.Login)
	public.POST("/refresh", handler.Refresh)

	// Register the authenticated endpoints.
	private := api.Group("/auth")
	private.Use(r.authRequired)
	private.POST("/logout", handler.Logout)
	private.PUT("/password", handler.ChangeUserPassword)
}
