package main

import (
	"github.com/leonsal/gux"
	"github.com/leonsal/gux/gb"
)

func init() {

	registerTest("bezier_quad", 6, newTestBezier)
}

type testBezier struct {
	p1     gb.Vec2
	p2     gb.Vec2
	p3     gb.Vec2
	deltaX float32
}

func newTestBezier(w *gux.Window) ITest {

	t := new(testBezier)
	t.p1 = gb.Vec2{10, 300}
	t.p2 = gb.Vec2{10, 20}
	t.p3 = gb.Vec2{10, t.p1.Y}
	t.deltaX = 6.0
	return t
}

func (t *testBezier) draw(w *gux.Window) {

	dl := w.DrawList()
	t.p1.Y = w.Size().Y / 2
	t.p3.X = w.Size().X - t.p1.X
	t.p3.Y = t.p1.Y

	colorPoint := gb.MakeColor(255, 0, 0, 255)
	w.AddCircleFilled(dl, t.p1, 10, colorPoint, 12)
	w.AddCircleFilled(dl, t.p2, 10, colorPoint, 12)
	w.AddCircleFilled(dl, t.p3, 10, colorPoint, 12)
	w.AddBezierQuadratic(dl, t.p1, t.p2, t.p3, gb.MakeColor(0, 0, 0, 255), 4, 0)
	t.p2.X += t.deltaX
	if t.p2.X >= t.p3.X && t.deltaX > 0 {
		t.deltaX = -t.deltaX
	}
	if t.p2.X <= t.p1.X && t.deltaX < 0 {
		t.deltaX = -t.deltaX
	}
}

func (t *testBezier) destroy(w *gux.Window) {
}
