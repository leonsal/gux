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
	SetPos(gb.Vec2)
}

type View struct {
	pos       gb.Vec2
	transform gb.Mat3
	//drawList  gb.DrawList
	visible  bool
	children []IView
}

func (v *View) SetPos(pos gb.Vec2) {

	v.pos = pos
}

func DispatchEvents(w *window.Window, v IView) {

}
