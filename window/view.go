package window

// IView is the interface for all Views
type IView interface {
	Render(*Window)
}

// A view can be a control such as a button, check box or a layout container.
type View struct {
	super    IView
	visible  bool
	children []*View
	render   func(*Window)
}

// Init initializes this View
func (v *View) Init(super IView) *View {

	v.super = super
	v.visible = true
	return v
}

// SetRender sets the render function of this view
func (v *View) SetRender(render func(*Window)) *View {

	v.render = render
	return v
}

// SetVisible sets the visibility state of this view
func (v *View) SetVisible(visible bool) *View {

	v.visible = visible
	return v
}

// Render renders this view
func (v *View) Render(w *Window) {

	if !v.visible {
		return
	}
	v.render(w)
}

// Render renders this view's children
func (v *View) RenderChildren(w *Window) {

	for _, cv := range v.children {
		cv.Render(w)
	}
}
