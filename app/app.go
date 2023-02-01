package app

import (
	"github.com/leonsal/gux/view"
	"github.com/leonsal/gux/window"
)

type App struct {
	windows []*Window
}

type Window struct {
	*window.Window
	fm   *FontManager
	view view.IView
}

// Single App instance
var app *App

// Init initializes the application singleton and returns its reference
func Init() *App {

	if app == nil {
		app = newApp()
	}
	return app
}

func (a *App) Close() {

	if app == nil {
		return
	}

	for i := 0; i < len(a.windows); i++ {
		a.windows[i].Destroy()
	}
	app = nil
}

// Render renders all the application opened windows and
// return true if all windows were closed.
func (a *App) Render() bool {

	if len(a.windows) == 0 {
		return true
	}

	for i := 0; i < len(a.windows); i++ {
		aw := a.windows[i]
		shouldClose := aw.StartFrame()
		if shouldClose {
			a.windows = append(a.windows[:i], a.windows[i+1:]...)
			aw.Destroy()
			continue
		}
		aw.RenderFrame()
	}
	return false
}

// NewWindow creates and returns a new application window
func (a *App) NewWindow(title string, width, height int) (*Window, error) {

	w, err := window.New(title, width, height, nil)
	if err != nil {
		return nil, err
	}
	aw := &Window{Window: w}
	a.windows = append(a.windows, aw)
	return aw, nil
}

// Close closes the specified window
func (aw *Window) Close() {

	for i := 0; i < len(app.windows); i++ {
		if app.windows[i] == aw {
			app.windows = append(app.windows[:i], app.windows[:i+1]...)
			aw.Destroy()
			break
		}
	}
}

// SetView sets the top view of this window
func (aw *Window) SetView(v view.IView) {
	aw.view = v
}

func newApp() *App {

	a := new(App)
	return a
}

func (a *App) defaultFontManager() error {

	//	normalSize := 18
	//	runeSets := [][]rune{}
	//	runeSets = append(runeSets, window.AsciiSet, window.RangeTableSet(unicode.Latin), window.RangeTableSet(unicode.Common))
	//	fm, err := NewFontManager(normalSize, 1, 2, runeSets...)
	//	if err != nil {
	//		return err
	//	}
	//
	return nil
}
