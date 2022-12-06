package main

import (
	"runtime"

	"github.com/leonsal/gux/gb"
	"github.com/leonsal/gux/window"
)

func main() {

	runtime.LockOSThread()

	// Create window
	win, err := window.New("title", 2000, 1200)
	if err != nil {
		panic(err)
	}

	// Render loop
	for win.StartFrame(0) {
		//testLines(win)
		testPolygon(win)
		win.Render()
	}
	win.Destroy()
}

func testLines(w *window.Window) {

	// Line points
	points := []gb.Vec2{{0, 10}, {10, 0}, {20, 10}, {30, 0}, {40, 10}, {50, 0}, {60, 10}}
	points1 := w.ReserveVec2(len(points))
	points2 := w.ReserveVec2(len(points))

	// Add poly lines anti aliased
	copy(points1, points)
	scalePoints(points1, 12)
	translatePoints(points1, gb.Vec2{10, 10})
	dl := w.DrawList()
	for width := 1; width < 60; width += 8 {
		w.AddPolyLineAntiAliased(dl, points1, gb.MakeColor(0, 0, 0, 255), 0, float32(width))
		translatePoints(points1, gb.Vec2{0, 120})
	}

	// Add poly lines textured
	copy(points2, points)
	scalePoints(points2, 12)
	translatePoints(points2, gb.Vec2{800, 10})
	for width := 1; width < 60; width += 8 {
		w.AddPolyLineTextured(dl, points2, gb.MakeColor(0, 0, 0, 255), 0, float32(width))
		translatePoints(points2, gb.Vec2{0, 120})
	}
}

func testPolygon(w *window.Window) {

	dl := w.DrawList()

	triangle := []gb.Vec2{{0, 1000}, {1000, 1000}, {500, 0}}
	//scalePoints(triangle, 4)
	//translatePoints(triangle, gb.Vec2{500, 500})
	w.AddConvexPolyFilled(dl, triangle, gb.MakeColor(0, 0, 0, 255), window.DrawFlag_AntiAliasedFill)

	//rect := []gb.Vec2{{0, 100}, {200, 100}, {200, 0}, {0, 0}}
	//translatePoints(rect, gb.Vec2{120, 10})
	//w.AddConvexPolyFilled(dl, rect, gb.MakeColor(255, 0, 0, 255), 0)
}

// scale the supplied array of points
func scalePoints(points []gb.Vec2, scale float32) {
	for i := range points {
		(&points[i]).MultScalar(scale)
	}
}

// translate the supplied array of points
func translatePoints(points []gb.Vec2, trans gb.Vec2) {

	for i := range points {
		(&points[i]).Add(trans)
	}
}
