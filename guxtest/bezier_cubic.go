package main

import (
	"github.com/leonsal/gux"
	"github.com/leonsal/gux/gb"
)

func init() {

	registerTest("bezier_cubic", 7, newTestBezierCubic)
}

type testBezierCubic struct {
	p1      gb.Vec2
	p2      gb.Vec2
	p3      gb.Vec2
	p4      gb.Vec2
	deltaX2 float32
	deltaX3 float32
}

func newTestBezierCubic(w *gux.Window) ITest {

	t := new(testBezierCubic)
	t.p1 = gb.Vec2{10, w.Size().Y / 2}
	t.p2 = gb.Vec2{10, 10}
	t.p3 = gb.Vec2{10, w.Size().Y - 10}
	t.p4 = gb.Vec2{w.Size().X - 10, t.p1.Y}

	t.deltaX2 = 6.0
	t.deltaX3 = 4.0
	return t
}

func (t *testBezierCubic) draw(w *gux.Window) {

	dl := w.DrawList()
	t.p1.Y = w.Size().Y / 2
	t.p3.Y = w.Size().Y - 10
	t.p4 = gb.Vec2{w.Size().X - 10, t.p1.Y}

	colorPoint := gb.MakeColor(255, 0, 0, 255)
	w.AddCircleFilled(dl, t.p1, 10, colorPoint, 12)
	w.AddCircleFilled(dl, t.p2, 10, colorPoint, 12)
	w.AddCircleFilled(dl, t.p3, 10, colorPoint, 12)
	w.AddCircleFilled(dl, t.p4, 10, colorPoint, 12)
	w.AddBezierCubic(dl, t.p1, t.p2, t.p3, t.p4, gb.MakeColor(0, 0, 0, 255), 4, 0)
	t.p2.X += t.deltaX2
	if t.p2.X >= t.p4.X && t.deltaX2 > 0 {
		t.deltaX2 = -t.deltaX2
	}
	if t.p2.X <= t.p1.X && t.deltaX2 < 0 {
		t.deltaX2 = -t.deltaX2
	}
	t.p3.X += t.deltaX3
	if t.p3.X >= t.p4.X && t.deltaX3 > 0 {
		t.deltaX3 = -t.deltaX3
	}
	if t.p3.X <= t.p1.X && t.deltaX3 < 0 {
		t.deltaX3 = -t.deltaX3
	}
}

func (t *testBezierCubic) destroy(w *gux.Window) {
}
