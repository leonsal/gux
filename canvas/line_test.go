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

	// Scale the supplied array of points
	scale := func(points []gb.Vec2, scale float32) {
		for i := range points {
			(&points[i]).MultScalar(scale)
		}
	}
	// Translate the supplied array of points
	translate := func(points []gb.Vec2, trans gb.Vec2) {
		for i := range points {
			(&points[i]).Add(trans)
		}
	}

	c := New()
	points := []gb.Vec2{{0, 10}, {10, 0}, {20, 10}, {30, 0}, {40, 10}, {50, 0}, {60, 10}}

	points1 := make([]gb.Vec2, len(points))
	copy(points1, points)
	scale(points1, 12)
	translate(points1, gb.Vec2{10, 10})
	for width := 1; width < 60; width += 8 {
		c.AddPolyLineAntiAliased(points1, 0xFF_00_00_00, 0, float32(width))
		translate(points1, gb.Vec2{0, 120})
	}

	points2 := make([]gb.Vec2, len(points))
	copy(points2, points)
	scale(points2, 12)
	translate(points2, gb.Vec2{800, 10})
	for width := 1; width < 60; width += 8 {
		c.AddPolyLineTextured(points2, 0xFF_00_00_00, 0, float32(width))
		translate(points2, gb.Vec2{0, 120})
	}

	for win.StartFrame(0) {
		win.RenderFrame(&c.DrawList)
	}
	win.Destroy()
}
