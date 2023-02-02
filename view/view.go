package view

import (
	"github.com/leonsal/gux/gb"
	"github.com/leonsal/gux/window"
)

type IView interface {
	Render(*window.Window)
	//Event(gb.Event) bool
	//Pos() gb.Vec2
	//Size() gb.Vec2
	SetPos(x, y float32)
}

type View struct {
	visible   bool
	pos       gb.Vec2
	transform gb.Mat3
}

func (v *View) SetPos(x, y float32) {

	v.pos = gb.Vec2{x, y}
}

func DispatchEvents(w *window.Window, v IView) {

}

type Group struct {
	View
	children []IView
}

func NewGroup() *Group {

	g := new(Group)
	g.visible = true
	return g
}

func (g *Group) Add(v IView) {

	g.children = append(g.children, v)
}

func (g *Group) Render(w *window.Window) {

	for _, c := range g.children {
		c.Render(w)
	}
}
