package server

import (
	"github.com/Wexyuan/kairos/internal/config"
	"github.com/Wexyuan/kairos/internal/router"
	httpserver "github.com/Wexyuan/kairos/pkg/server"
)

// NewHTTPServer creates the HTTP server.
func NewHTTPServer(cfg *config.Config, appRouter *router.Router) *httpserver.HTTPServer {
	appRouter.Init()
	return httpserver.NewHTTP(httpserver.HTTPOptions{
		Addr:            cfg.Server.Addr,
		ReadTimeout:     cfg.Server.ReadTimeout,
		WriteTimeout:    cfg.Server.WriteTimeout,
		ShutdownTimeout: cfg.Server.ShutdownTimeout,
		Handler:         appRouter,
	})
}
