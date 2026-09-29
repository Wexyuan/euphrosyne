//go:build wireinject
// +build wireinject

package wire

import (
	"github.com/google/wire"
	"github.com/wailsapp/wails/v3/pkg/application"

	"github.com/Wexyuan/kairos"
	"github.com/Wexyuan/kairos/internal/config"
	"github.com/Wexyuan/kairos/internal/module/auth"
	"github.com/Wexyuan/kairos/internal/module/base"
	"github.com/Wexyuan/kairos/internal/module/user"
	"github.com/Wexyuan/kairos/pkg/logger"
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
	user.NewUserRepository,
	user.NewUserService,
	auth.NewAuthRepository,
	auth.NewAuthService,
)

func NewApp(
	userSvc *user.UserService,
	authSvc *auth.AuthService,
) *application.App {
	return application.New(
		application.Options{
			Name: "kairos",
			Services: []application.Service{
				application.NewService(userSvc),
				application.NewService(authSvc),
			},
			Assets: application.AssetOptions{
				Handler: application.AssetFileServerFS(kairos.Assets),
			},
		},
	)
}

func NewWireApp(cfg *config.Config, log *logger.Logger) (*application.App, func(), error) {
	panic(wire.Build(
		InfraProviderSet,
		ModuleProviderSet,
		NewApp,
	))
}
