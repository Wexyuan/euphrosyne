package main

import (
	"flag"
	"log"
	"os"

	appconfig "github.com/Wexyuan/euphrosyne/internal/config"
	"github.com/Wexyuan/euphrosyne/internal/module/user"
	"github.com/Wexyuan/euphrosyne/pkg/config"
	"github.com/Wexyuan/euphrosyne/pkg/database"
	"github.com/Wexyuan/euphrosyne/pkg/logger"
)

// defaultConfigPath is the default config file path.
const defaultConfigPath = "configs/config.yaml"

// main migrates the database schema.
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

	// 2. Load the config.
	var cfg appconfig.Config
	if err := config.Load(cfgPath, &cfg); err != nil {
		log.Fatalf("[migrate] load config %s error: %v", cfgPath, err)
	}

	// 3. Create the logger.
	lg, err := logger.New(logger.Options{
		Level:  cfg.Logger.Level,
		Format: cfg.Logger.Format,
		Output: cfg.Logger.Output,
	})
	if err != nil {
		log.Fatalf("[migrate] create logger error: %v", err)
	}
	defer func() { _ = lg.Sync() }()

	// 4. Open the database.
	db, err := database.New(database.Options{
		Driver:          cfg.Database.Driver,
		DSN:             cfg.Database.DSN,
		MaxOpenConns:    cfg.Database.MaxOpenConns,
		MaxIdleConns:    cfg.Database.MaxIdleConns,
		ConnMaxLifetime: cfg.Database.ConnMaxLifetime,
	})
	if err != nil {
		log.Fatalf("[migrate] init database error: driver=%s: %v", cfg.Database.Driver, err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			lg.Errorf("close database error: %v", err)
		}
	}()

	// 5. Migrate the models.
	if err := db.AutoMigrate(&user.User{}); err != nil {
		log.Fatalf("[migrate] migrate models error: %v", err)
	}
	lg.Infof("migrated models successfully on driver %s", cfg.Database.Driver)
}
