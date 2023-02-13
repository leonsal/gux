package view

import (
	"fmt"

	"github.com/leonsal/gux/gb"
	"github.com/leonsal/gux/window"
)

type IView interface {
	Render(*window.Window) // Renders the view at the specified window
	//Event(gb.Event) bool
	//Pos() gb.Vec2
	//Size() gb.Vec2
	SetPos(x, y float32) // Sets the view position relative to its parent
	SetTransform(t *gb.Mat3)
	Styles() StyleMap
}

type View struct {
	iview     IView    // Associated IView
	visible   bool     // Visibility state
	pos       gb.Vec2  // View position relative to its parent
	scale     gb.Vec2  // View scale
	rotation  float32  // Rotation in radians
	transform gb.Mat3  // Current transform matrix used  in AddList2()
	styles    StyleMap // Custom styles map
	parent    IView    // Parent IView (maybe nil)
	children  []IView  // List of child views
}

func (v *View) Init(iv IView) {

	v.iview = iv
	v.visible = true
	v.scale = gb.Vec2{1, 1}
	v.transform.Identity()
}

func (v *View) Styles() StyleMap {
	return v.styles
}

func (v *View) SetVisible(visible bool) {

	v.visible = visible
}

func (v *View) Visible() bool {

	return v.visible
}

func (v *View) SetPos(x, y float32) {

	v.pos = gb.Vec2{x, y}
}

func (v *View) Pos() gb.Vec2 {

	return v.pos
}

func (v *View) SetScale(x, y float32) *View {

	v.scale = gb.Vec2{x, y}
	return v
}

func (v *View) Scale() gb.Vec2 {

	return v.scale
}

func (v *View) SetRotation(r float32) {

	v.rotation = r
}

func (v *View) Rotation() float32 {

	return v.rotation
}

// SetTransform updates the transform matrix of this view and applies
func (v *View) SetTransform(t *gb.Mat3) {

	v.transform.SetTranslationVec(v.pos).Rotate(v.rotation).ScaleVec(v.scale)
	v.transform.Mult(t)
}

func (v *View) Add(iv IView) {

	v.children = append(v.children, iv)
}

func (v *View) RenderChildren(w *window.Window) {

	// Sets parent world transform matrix
	var mat gb.Mat3
	mat.SetTranslationVec(v.pos).Rotate(v.rotation).ScaleVec(v.scale)
	fmt.Printf("mat:%+v\n", mat)
	for _, c := range v.children {
		c.SetTransform(&mat)
		c.Render(w)
	}

}

func DispatchEvents(w *window.Window, v IView) {

}
