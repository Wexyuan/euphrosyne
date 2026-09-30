package kairos

import "embed"

// Assets holds the embedded frontend build output.
//
//go:embed all:frontend/dist
var Assets embed.FS

// TrayIcon holds the embedded tray icon from the build assets.
//
//go:embed build/windows/icon.ico
var TrayIcon []byte
