package main

import (
	"math"

	"github.com/leonsal/gux"
	"github.com/leonsal/gux/gb"
)

func init() {

	registerTest("arc", 4, newTestArc)
}

type testArc struct{}

func newTestArc(w *gux.Window) ITest {

	return new(testArc)
}

func (t *testArc) draw(win *gux.Window) {

	dl := win.DrawList()
	radius := float32(100)
	center := gb.Vec2{radius + 10, radius + 10}
	segs := 3
	thickness := float32(2)
	deltaY := 2*radius + 50
	deltaX := 2*radius + 40
	for i := 0; i < 12; i++ {
		win.AddCircle(dl, center, radius, nextColor(i), segs, thickness)
		center.X += deltaX
		segs += 3
		thickness += 2
	}

	center = gb.Vec2{radius + 10, center.Y + deltaY}
	segs = 3
	for i := 0; i < 12; i++ {
		win.AddCircleFilled(dl, center, radius, nextColor(i), segs)
		center.X += deltaX
		segs += 3
	}

	center = gb.Vec2{radius + 10, center.Y + deltaY}
	segs = 3
	thickness = float32(2)
	for i := 0; i < 12; i++ {
		amax := 2 * math.Pi * (float64(i) + 1) * 0.08
		win.PathArcTo(dl, center, radius, 0, float32(amax), segs)
		win.PathStroke(dl, nextColor(i), 0, thickness)
		center.X += deltaX
		segs += 3
		thickness += 2
	}

	center = gb.Vec2{radius + 10, center.Y + deltaY}
	segs = 3
	for i := 0; i < 12; i++ {
		amax := 2 * math.Pi * (float64(i) + 1) * 0.08
		win.PathArcTo(dl, center, radius, 0, float32(amax), segs)
		win.PathFillConvex(dl, nextColor(i))
		center.X += deltaX
		segs += 3
		thickness += 2
	}
}

func (t *testArc) destroy(w *gux.Window) {

}
