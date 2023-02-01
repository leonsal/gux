package app

import (
	"log"
	"runtime"
	"unicode"

	"github.com/leonsal/gux/util"
	"github.com/leonsal/gux/view"
	"github.com/leonsal/gux/window"
	"golang.org/x/image/font/gofont/gobold"
	"golang.org/x/image/font/gofont/goitalic"
	"golang.org/x/image/font/gofont/goregular"
)

// App is a singleton with the context of the entire Application
type App struct {
	windows []*Window // List of opened native windows
}

type Window struct {
	*window.Window              // Embedded window reference
	fm             *FontManager // Window FontManager
	view           view.IView   // Current top view
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

	for i := 0; i < len(a.windows); i++ {
		a.windows[i].Close()
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
		aw := a.windows[i]
		shouldClose := aw.StartFrame()
		if shouldClose {
			a.windows = append(a.windows[:i], a.windows[i+1:]...)
			aw.close()
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

	log.Println("close window 1", aw, app)
	for i := 0; i < len(app.windows); i++ {
		if app.windows[i] == aw {
			log.Println("close window 2", aw)
			app.windows = append(app.windows[:i], app.windows[:i+1]...)
			aw.close()
			break
		}
	}
}

// CreateFontManager creates the window default font manager
func (aw *Window) CreateFontManager() error {

	normalSize := 18.0
	runeSets := [][]rune{}
	runeSets = append(runeSets, util.AsciiSet(), util.RangeTableSet(unicode.Latin), util.RangeTableSet(unicode.Common))
	fm, err := NewFontManager(normalSize, 1, 2, runeSets...)
	if err != nil {
		return err
	}

	err = fm.AddFamily(FontRegular, goregular.TTF)
	if err != nil {
		return err
	}

	err = fm.AddFamily(FontBold, gobold.TTF)
	if err != nil {
		return err
	}

	err = fm.AddFamily(FontItalic, goitalic.TTF)
	if err != nil {
		return err
	}

	err = fm.BuildFonts(aw)
	if err != nil {
		return err
	}
	aw.fm = fm
	return nil
}

// SetView sets the top view of this window
func (aw *Window) SetView(v view.IView) {
	aw.view = v
}

func (aw *Window) close() {

	if aw.fm != nil {
		aw.fm.DestroyFonts(aw)
		aw.fm = nil
	}
	aw.Destroy()
}
