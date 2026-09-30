package desktop

import (
	"github.com/wailsapp/wails/v3/pkg/application"
	"github.com/wailsapp/wails/v3/pkg/events"
)

// Tray implements the system tray for the main window.
type Tray struct {
	win  application.Window // main window
	icon []byte             // tray icon bytes
}

func NewTray(win application.Window, icon []byte) *Tray {
	return &Tray{
		win:  win,
		icon: icon,
	}
}

// Run runs the system tray.
func (t *Tray) Run() {
	t.hideOnClose()

	tray := application.Get().SystemTray.New()
	tray.SetIcon(t.icon)
	tray.SetTooltip("Kairos")
	tray.SetMenu(t.buildMenu())
	tray.OnClick(t.showWindow)
	tray.Run()
}

// hideOnClose hides the window when it is closed.
func (t *Tray) hideOnClose() {
	t.win.OnWindowEvent(events.Windows.WindowClosing, func(*application.WindowEvent) {
		t.win.Hide()
	})
}

// buildMenu builds the system tray menu.
func (t *Tray) buildMenu() *application.Menu {
	menu := application.NewMenu()
	menu.Add("显示主窗口").OnClick(func(*application.Context) { t.showWindow() })
	menu.Add("退出").OnClick(func(*application.Context) { application.Get().Quit() })
	return menu
}

// showWindow shows the main window.
func (t *Tray) showWindow() {
	t.win.Restore()
	t.win.Show()
	t.win.Focus()
}
