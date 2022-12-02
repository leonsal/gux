package canvas

import (
	"runtime"
	"testing"

	"github.com/leonsal/gux/gb"
)

func TestLine(t *testing.T) {

	runtime.LockOSThread()
	win, err := gb.CreateWindow("title", 1000, 1000)
	if err != nil {
		panic(err)
	}

	c := New()
	//points := []gb.Vec2{{100, 500}, {500, 500}, {500, 900}}
	//points := []gb.Vec2{{100, 500}, {500, 500}}
	points := []gb.Vec2{
		{100, 500}, {200, 100},
		{300, 500}, {400, 100},
		{500, 500}, {600, 100},
		{700, 500}, {800, 100},
	}
	c.AddPolyLineTextured(points, 0xFF_00_00_00, 0, 10)

	points2 := []gb.Vec2{
		{100, 900}, {200, 500},
		{300, 900}, {400, 500},
		{500, 900}, {600, 500},
		{700, 900}, {800, 500},
	}
	c.AddPolyLineTextured(points2, 0xFF_00_00_00, 0, 10)
	for win.StartFrame(0) {
		win.RenderFrame(&c.DrawList)
	}
	win.Destroy()
}
