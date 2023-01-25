package app

import (
	"github.com/leonsal/gux/view"
	"github.com/leonsal/gux/window"
)

type App struct {
	Windows []*Window
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
	aw := &Window{Window: w}
	a.Windows = append(a.Windows, aw)
	return aw, nil
}

func (a *App) Render() bool {

	////toclose := []*Window{}
	//fmt.Printf("windows:%+v\n", a.Windows)
	//aw := a.Windows[0]
	//shouldClose := aw.StartFrame()

	////for i := 0; i < len(a.windows); i++ {
	////	aw := a.windows[i]
	////	shouldClose := aw.StartFrame()
	////	if shouldClose {
	////		toclose = append(toclose, aw)
	////		continue
	////	}
	////	aw.RenderFrame()
	////}
	////for _, aw := range toclose {
	////	aw.Close()
	////}
	return true
}

func (aw *Window) SetView(v view.IView) {
	aw.view = v
}

func (aw *Window) Close() {

	for i := 0; i < len(app.Windows); i++ {
		if app.Windows[i] == aw {
			app.Windows = append(app.Windows[:i], app.Windows[:i+1]...)
			break
		}
	}
	aw.Destroy()
}

func (aw *Window) Render() bool {

	//aw.win.Render()
	return true
}
