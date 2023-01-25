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
	view view.IView
}

// Single App instance
var app *App

// Init initializes the application and returns reference to its single instance
func Init() *App {

	if app == nil {
		app = newApp()
	}
	return app
}

func newApp() *App {

	a := new(App)
	return a
}

// NewWindow creates and returns a new application window
func (a *App) NewWindow(title string, width, height int) (*Window, error) {

	w, err := window.New(title, width, height, nil)
	if err != nil {
		return nil, err
	}
	aw := &Window{
		Window: w,
	}
	a.windows = append(a.windows, aw)
	return aw, nil
}

func (a *App) Render() {

	//for _, aw := range a.windows {
	//}
}

func (aw *Window) SetView(v view.IView) {
	aw.view = v
}

func (aw *Window) Render() bool {

	//aw.win.Render()
	return true
}
