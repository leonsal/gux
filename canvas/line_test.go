package canvas

import (
	"runtime"
	"testing"

	"github.com/leonsal/gux/gb"
)

func TestLine(t *testing.T) {

	runtime.LockOSThread()
	win, err := gb.CreateWindow("title", 2000, 1200)
	if err != nil {
		panic(err)
	}

	// Translate and scale the supplied array of points
	transform := func(points []gb.Vec2, trans gb.Vec2, scale float32) {
		for i := range points {
			(&points[i]).Add(trans).MultScalar(scale)
		}
	}

	c := New()
	points := []gb.Vec2{{0, 10}, {10, 0}, {20, 10}, {30, 0}, {40, 10}, {50, 0}, {60, 10}}
	transform(points, gb.Vec2{10, 0}, 20)

	for width := 1; width < 60; width += 8 {
		c.AddPolyLineAntiAliased(points, 0xFF_00_00_00, 0, float32(width))
		transform(points, gb.Vec2{0, 120}, 1)
	}

	for win.StartFrame(0) {
		win.RenderFrame(&c.DrawList)
	}
	win.Destroy()
}
