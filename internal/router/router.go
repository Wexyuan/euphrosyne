package router

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	appconfig "github.com/Wexyuan/euphrosyne/internal/config"
	"github.com/Wexyuan/euphrosyne/internal/middleware"
	"github.com/Wexyuan/euphrosyne/internal/module/auth"
	"github.com/Wexyuan/euphrosyne/internal/module/user"
	"github.com/Wexyuan/euphrosyne/pkg/logger"
	"github.com/Wexyuan/euphrosyne/pkg/response"
)

// Router manages the HTTP routes and middleware.
type Router struct {
	*gin.Engine

	log          *logger.Logger    // application logger
	authRequired gin.HandlerFunc   // auth required middleware
	auth         *auth.AuthHandler // auth endpoints
	user         *user.UserHandler // user endpoints
}

func New(
	log *logger.Logger,
	cfg *appconfig.Config,
	authRequired gin.HandlerFunc,
	authHandler *auth.AuthHandler,
	userHandler *user.UserHandler,
) *Router {
	switch strings.ToLower(cfg.App.Env) {
	case "release":
		gin.SetMode(gin.ReleaseMode)
	case "test":
		gin.SetMode(gin.TestMode)
	default:
		gin.SetMode(gin.DebugMode)
	}

	return &Router{
		Engine:       gin.New(),
		log:          log,
		authRequired: authRequired,
		auth:         authHandler,
		user:         userHandler,
	}
}

// Init sets up the middleware and the routes.
func (r *Router) Init() {
	r.setupGlobalMiddlewares()
	r.setupDefaultRoutes()
	r.setupV1Routes()
}

// setupGlobalMiddlewares sets up the global middleware.
func (r *Router) setupGlobalMiddlewares() {
	r.Use(
		middleware.Recovery(r.log),
		middleware.CORS(),
	)
}

// setupDefaultRoutes sets up the default routes.
func (r *Router) setupDefaultRoutes() {
	r.GET("/healthz", func(c *gin.Context) {
		// Answer the probe without touching the database so it stays cheap.
		c.JSON(http.StatusOK, response.OK(nil))
	})
}

// setupV1Routes sets up the v1 routes.
func (r *Router) setupV1Routes() {
	api := r.Group("/api/v1")

	r.setupAuthRoutes(api)
	r.setupUserRoutes(api)
}
