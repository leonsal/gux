package view

import "github.com/leonsal/gux/window"

type Group struct {
	Parent
}

func NewGroup() *Group {

	g := new(Group)
	g.Init()
	return g
}

func (g *Group) Render(w *window.Window) {

	for _, c := range g.children {
		c.Render(w)
	}
}
