package main

import (
	"flag"
	"log"
	"os"

	"github.com/wailsapp/wails/v3/pkg/application"

	"github.com/Wexyuan/kairos/cmd/desktop/wire"
	appconfig "github.com/Wexyuan/kairos/internal/config"
	"github.com/Wexyuan/kairos/pkg/config"
	"github.com/Wexyuan/kairos/pkg/logger"
)

// defaultConfigPath is the default config file path.
const defaultConfigPath = "configs/config.yaml"

// main starts the desktop application.
func main() {
	// 1. Resolve the config path.
	var cfgPath string
	flag.StringVar(&cfgPath, "config", "", "config file path")
	flag.Parse()
	if cfgPath == "" {
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
	app, cleanup, err := wire.NewWireApp(&cfg, lg)
	if err != nil {
		log.Fatalf("[main] build application error: %v", err)
	}
	defer cleanup()

	// 5. Create the main window.
	app.Window.NewWithOptions(application.WebviewWindowOptions{
		Name:             "main",
		Title:            "Kairos",
		Width:            1280,
		Height:           800,
		BackgroundColour: application.NewRGB(255, 255, 255),
		URL:              "/",
	})

	// 6. Run until the window is closed.
	if err := app.Run(); err != nil {
		log.Fatalf("[main] run application error: %v", err)
	}
}
