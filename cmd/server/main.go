package main

import (
	"context"
	"flag"
	"log"
	"os"

	"github.com/Wexyuan/euphrosyne/cmd/server/wire"
	appconfig "github.com/Wexyuan/euphrosyne/internal/config"
	"github.com/Wexyuan/euphrosyne/pkg/config"
	"github.com/Wexyuan/euphrosyne/pkg/logger"
)

// defaultConfigPath is the default config file path.
const defaultConfigPath = "configs/config.yaml"

// main starts the application.
func main() {
	// 1. Resolve the config path.
	var cfgPath string
	flag.StringVar(&cfgPath, "config", "", "config file path")
	flag.Parse()
	if cfgPath == "" {
		// Fall back to the env so a container can override the path without flags.
		cfgPath = os.Getenv("CONFIG_PATH")
	}
	if cfgPath == "" {
		cfgPath = defaultConfigPath
	}

	// 2. Load the config file.
	var cfg appconfig.Config
	if err := config.Load(cfgPath, &cfg); err != nil {
		log.Fatalf("[main] load config %s error: %v", cfgPath, err)
	}

	// 3. Create the logger.
	lg, err := logger.New(logger.Options{
		Level:  cfg.Logger.Level,
		Format: cfg.Logger.Format,
		Output: cfg.Logger.Output,
	})
	if err != nil {
		log.Fatalf("[main] create logger error: %v", err)
	}
	defer func() { _ = lg.Sync() }()

	// 4. Build the application.
	application, cleanup, err := wire.NewWireApp(&cfg, lg)
	if err != nil {
		log.Fatalf("[main] build application error: %v", err)
	}
	defer cleanup()

	// 5. Log the startup settings.
	lg.Infof("%s starting on %s (env=%s, db=%s, cache=%s)",
		cfg.App.Name, cfg.Server.Addr, cfg.App.Env, cfg.Database.Driver, cfg.Cache.Driver)

	// 6. Run until a termination signal arrives.
	if err := application.Run(context.Background()); err != nil {
		log.Fatalf("[main] serve on %s error: %v", cfg.Server.Addr, err)
	}
}
