//go:build wireinject
// +build wireinject

package wire

import (
	"github.com/google/wire"

	"github.com/Wexyuan/kairos/internal/config"
	"github.com/Wexyuan/kairos/internal/module/auth"
	"github.com/Wexyuan/kairos/internal/module/base"
	"github.com/Wexyuan/kairos/internal/module/user"
	"github.com/Wexyuan/kairos/internal/router"
	"github.com/Wexyuan/kairos/internal/server"
	"github.com/Wexyuan/kairos/pkg/app"
	"github.com/Wexyuan/kairos/pkg/logger"
	httpserver "github.com/Wexyuan/kairos/pkg/server"
)

// InfraProviderSet defines the infrastructure providers.
var InfraProviderSet = wire.NewSet(
	provideDatabase,
	provideSnowflake,
	provideCache,
	provideTokenManager,
)

// ModuleProviderSet defines the business module providers.
var ModuleProviderSet = wire.NewSet(
	base.NewRepository,
	base.NewService,
	base.NewHandler,
	user.NewUserRepository,
	user.NewUserService,
	user.NewUserHandler,
	auth.NewAuthRepository,
	auth.NewAuthService,
	auth.NewAuthHandler,
)

// MiddlewareProviderSet defines the middleware providers.
var MiddlewareProviderSet = wire.NewSet(
	provideAuthRequired,
)

// RouterProviderSet defines the router providers.
var RouterProviderSet = wire.NewSet(
	router.New,
)

// ServerProviderSet defines the server providers.
var ServerProviderSet = wire.NewSet(
	server.NewHTTPServer,
)

func NewApp(cfg *config.Config, httpServer *httpserver.HTTPServer) *app.App {
	return app.NewApp(
		app.WithName(cfg.App.Name),
		app.WithServers(httpServer),
	)
}

func NewWireApp(cfg *config.Config, log *logger.Logger) (*app.App, func(), error) {
	panic(wire.Build(
		InfraProviderSet,
		ModuleProviderSet,
		MiddlewareProviderSet,
		RouterProviderSet,
		ServerProviderSet,
		NewApp,
	))
}
