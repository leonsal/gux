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
	points := []gb.Vec2{{100, 500}, {500, 500}, {500, 900}}
	//points := []gb.Vec2{{100, 500}, {500, 500}}
	c.polyLineTextured(points, 0xFF_00_00_00, 0, 10)
	for win.StartFrame(0) {
		win.RenderFrame(&c.DrawList)
	}
	win.Destroy()
}
