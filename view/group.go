package view

//
import "github.com/leonsal/gux/window"

type Group struct {
	View
}

func NewGroup() *Group {

	g := new(Group)
	g.Init(g)
	return g
}

func (g *Group) Render(w *window.Window) {

	g.RenderChildren(w)
	//for _, c := range g.children {
	//	c.Render(w)
	//}
}
