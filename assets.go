package kairos

import "embed"

// Assets holds the embedded frontend build output.
//
//go:embed all:frontend/dist
var Assets embed.FS
