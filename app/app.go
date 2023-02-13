package app

import (
	"runtime"
	"unicode"

	"github.com/leonsal/gux/gb"
	"github.com/leonsal/gux/util"
	"github.com/leonsal/gux/view"
	"github.com/leonsal/gux/window"
	"golang.org/x/image/font/gofont/gobold"
	"golang.org/x/image/font/gofont/goitalic"
	"golang.org/x/image/font/gofont/goregular"
)

type windowInfo struct {
	w    *window.Window // Native window reference
	view view.IView     // Current top view for the window
}

// App is a singleton with the context of the entire Application
type App struct {
	windows []*windowInfo // List of opened native windows
}

// Single Application instance
var app *App

// Get returns the application singleton
func Get() *App {

	if app == nil {
		runtime.LockOSThread()
		app = new(App)
	}
	return app
}

// Close closes the application and its windows releasing acquired resources.
func (a *App) Close() {

	if app == nil {
		return
	}

	for _, wi := range a.windows {
		wi.w.Destroy()
	}
}

// Render renders all the application opened windows.
// Returns return true if all windows were closed or
// false if at least one window is opened.
func (a *App) Render() bool {

	if len(a.windows) == 0 {
		return true
	}

	for i := 0; i < len(a.windows); i++ {

		// Starts frame and checks if window should close
		wi := a.windows[i]
		shouldClose := wi.w.StartFrame()
		if shouldClose {
			a.windows = append(a.windows[:i], a.windows[i+1:]...)
			wi.w.Destroy()
			continue
		}
		if wi.view != nil {
			// Dispatch events to the top view
			view.DispatchEvents(wi.w, wi.view)
			// Render top view and its children
			wi.view.Render(wi.w)
		}
		wi.w.RenderFrame()
	}
	return false
}

// NewWindow creates and returns a new application window
func (a *App) NewWindow(title string, width, height int) (*window.Window, error) {

	return a.NewWindowEx(title, width, height, nil)
}

func (a *App) NewWindowEx(title string, width, height int, cfg *gb.Config) (*window.Window, error) {

	// Creates native window
	w, err := window.New(title, width, height, cfg)
	if err != nil {
		return nil, err
	}

	// Creates default font manager
	err = a.createFontManager(w)
	if err != nil {
		return nil, err
	}

	wi := &windowInfo{w: w}
	a.windows = append(a.windows, wi)
	return w, nil
}

// Close removes the specified window from the application and then closes it,
// releasing all its resources.
func (a *App) CloseWindow(w *window.Window) {

	for i := 0; i < len(app.windows); i++ {
		wi := app.windows[i]
		if wi.w == w {
			app.windows = append(app.windows[:i], app.windows[:i+1]...)
			w.Destroy()
			break
		}
	}
	panic("CloseWindow(): Invalid window")
}

// SetView sets the top view of the specified window
func (a *App) SetView(w *window.Window, v view.IView) {

	for _, wi := range a.windows {
		if wi.w == w {
			wi.view = v
			return
		}
	}
	panic("SetView(): Invalid window")
}

// createFontManager creates the default FontManager for the window
func (a *App) createFontManager(w *window.Window) error {

	normalSize := 48.0
	runeSets := [][]rune{}
	runeSets = append(runeSets, util.AsciiSet(), util.RangeTableSet(unicode.Latin), util.RangeTableSet(unicode.Common))
	fm, err := window.NewFontManager(normalSize, 1, 2, runeSets...)
	if err != nil {
		return err
	}

	err = fm.AddStyle(window.FontRegular, goregular.TTF)
	if err != nil {
		return err
	}

	err = fm.AddStyle(window.FontBold, gobold.TTF)
	if err != nil {
		return err
	}

	err = fm.AddStyle(window.FontItalic, goitalic.TTF)
	if err != nil {
		return err
	}

	err = fm.BuildFonts(w)
	if err != nil {
		return err
	}
	w.SetFontManager(fm)
	return nil
}
