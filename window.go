package gux

import "github.com/leonsal/gux/gb"

type Window struct {
	gbw *gb.Window  // Graphics backend native window reference
	dl  gb.DrawList // Draw list to render
}

func NewWindow(title string, width, height int) (*Window, error) {

	// Creates graphics backend native window
	w := new(Window)
	var err error
	w.gbw, err = gb.CreateWindow(title, width, height)
	if err != nil {
		return nil, err
	}

	// Create line texture and transfer to backend

	return w, nil
}

func (w *Window) StartFrame(timeout float64) bool {

	return w.gbw.StartFrame(timeout)
}

func (w *Window) RenderFrame(widget IWidget) {

	widget.Render(w)
	w.gbw.RenderFrame(&w.dl)
}
