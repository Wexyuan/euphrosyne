package desktop

import (
	"context"

	"github.com/wailsapp/wails/v3/pkg/application"
	"github.com/wailsapp/wails/v3/pkg/w32"
)

// mainWindowName matches the name of the main window created in cmd/desktop.
const mainWindowName = "main"

// WindowService implements the native window operations.
type WindowService struct{}

func NewWindowService() *WindowService {
	return &WindowService{}
}

// SetTheme repaints the native titlebar to match the web theme.
// Custom colours only apply on Windows 11; older systems just switch dark mode.
func (s *WindowService) SetTheme(ctx context.Context, req *SetThemeReq) error {
	win, ok := application.Get().Window.GetByName(mainWindowName)
	if !ok {
		return nil
	}
	// Repaint on the main thread because the DWM calls must come from it.
	application.InvokeAsync(func() {
		hwnd := uintptr(win.NativeWindow())
		if hwnd == 0 {
			return
		}
		// Colours use the 0x00BBGGRR layout and must match the web layout background.
		bg, text := uint32(0x00F5F5F5), uint32(0x001A1A1A)
		if req.Dark {
			bg, text = uint32(0x000A0A0A), uint32(0x00E5E5E5)
		}
		w32.SetTheme(hwnd, req.Dark)
		w32.SetTitleBarColour(hwnd, bg)
		w32.SetTitleTextColour(hwnd, text)
		w32.SetBorderColour(hwnd, bg)
	})
	return nil
}
