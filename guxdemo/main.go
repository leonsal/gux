package main

import (
	"runtime"

	"github.com/leonsal/gux"
	"github.com/leonsal/gux/canvas"
	"github.com/leonsal/gux/gb"
)

func main() {

	runtime.LockOSThread()

	// Create window
	win, err := gux.NewWindow("title", 2000, 1200)
	if err != nil {
		panic(err)
	}

	// Creates test view
	test := NewTestLines(win)

	// Render loop
	for win.StartFrame(0) {
		win.RenderFrame(test)
	}
	win.Destroy()
}

// testLines is a View with a canvas to which lines are drawn
type testLines struct {
	gux.View
	c *canvas.Canvas
}

func NewTestLines(w *gux.Window) *testLines {

	// Create test
	t := new(testLines)
	t.Init(t)
	t.c = canvas.New(w)

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

	// Line points
	points := []gb.Vec2{{0, 10}, {10, 0}, {20, 10}, {30, 0}, {40, 10}, {50, 0}, {60, 10}}
	points1 := make([]gb.Vec2, len(points))
	points2 := make([]gb.Vec2, len(points))

	// Sets the render function of this View
	t.SetRender(func(w *gux.Window) {

		// Add poly lines anti aliased
		copy(points1, points)
		scale(points1, 12)
		translate(points1, gb.Vec2{10, 10})
		for width := 1; width < 60; width += 8 {
			t.c.AddPolyLineAntiAliased(points1, gb.MakeColor(0, 0, 0, 255), 0, float32(width))
			translate(points1, gb.Vec2{0, 120})
		}

		// Add poly lines textured
		copy(points2, points)
		scale(points2, 12)
		translate(points2, gb.Vec2{800, 10})
		for width := 1; width < 60; width += 8 {
			t.c.AddPolyLineTextured(points2, gb.MakeColor(0, 0, 0, 255), 0, float32(width))
			translate(points2, gb.Vec2{0, 120})
		}

		t.c.Render(w)
	})

	return t
}
