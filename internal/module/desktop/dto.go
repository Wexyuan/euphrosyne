package desktop

// SetThemeReq is the set theme request.
type SetThemeReq struct {
	Dark bool `json:"dark"` // whether the web UI switched to dark mode
}
