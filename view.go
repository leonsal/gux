package gux

// A view can be a control such as a button, check box or a layout container.
type View struct {
	visible  bool
	children []*View
}

type IView interface {
	Render(*Window)
}

func (v *View) SetVisible(visible bool) *View {

	v.visible = visible
	return v
}

func (v *View) Render(w *Window) {

	if !v.visible {
		return
	}
	for _, cv := range v.children {
		cv.Render(w)
	}
}
