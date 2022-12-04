package canvas

import (
	"runtime"
	"testing"

	"github.com/leonsal/gux"
	"github.com/leonsal/gux/gb"
)

func TestLine(t *testing.T) {

	runtime.LockOSThread()

	// Create window
	win, err := gux.NewWindow("title", 2000, 1200)
	if err != nil {
		panic(err)
	}

	// Function to scale the supplied array of points
	scale := func(points []gb.Vec2, scale float32) {
		for i := range points {
			(&points[i]).MultScalar(scale)
		}
	}
	// Function to translate the supplied array of points
	translate := func(points []gb.Vec2, trans gb.Vec2) {
		for i := range points {
			(&points[i]).Add(trans)
		}
	}

	// Create canvas
	canvas := New(win)

	// Draw lines into the canvas
	points := []gb.Vec2{{0, 10}, {10, 0}, {20, 10}, {30, 0}, {40, 10}, {50, 0}, {60, 10}}
	points1 := make([]gb.Vec2, len(points))
	points2 := make([]gb.Vec2, len(points))

	drawLines := func(c *Canvas) {

		copy(points1, points)
		scale(points1, 12)
		translate(points1, gb.Vec2{10, 10})
		for width := 1; width < 60; width += 8 {
			c.AddPolyLineAntiAliased(points1, gb.MakeColor(0, 0, 0, 255), 0, float32(width))
			translate(points1, gb.Vec2{0, 120})
		}

		copy(points2, points)
		scale(points2, 12)
		translate(points2, gb.Vec2{800, 10})
		for width := 1; width < 60; width += 8 {
			c.AddPolyLineTextured(points2, gb.MakeColor(0, 0, 0, 255), 0, float32(width))
			translate(points2, gb.Vec2{0, 120})
		}
	}

	for win.StartFrame(0) {
		drawLines(canvas)
		win.RenderFrame(canvas)
	}
	win.Destroy()
}
