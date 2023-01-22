package main

import (
	"github.com/leonsal/gux"
	"github.com/leonsal/gux/gb"
)

func init() {

	registerTest("bezier", 8, newTestBezier)
}

type testBezier struct {
	p1     gb.Vec2
	p2     gb.Vec2
	p3     gb.Vec2
	deltaX float32
}

func newTestBezier(w *gux.Window) ITest {

	return &testBezier{
		p1:     gb.Vec2{10, 200},
		p2:     gb.Vec2{500, 20},
		p3:     gb.Vec2{1000, 200},
		deltaX: 4.0,
	}
}

func (t *testBezier) draw(w *gux.Window) {

	dl := w.DrawList()
	w.AddCircleFilled(dl, t.p2, 10, gb.MakeColor(255, 0, 0, 255), 12)
	w.AddBezierQuadratic(dl, t.p1, t.p2, t.p3, gb.MakeColor(0, 0, 0, 255), 8, 0)
	t.p2.X += t.deltaX
	if t.p2.X < t.p1.X || t.p2.X > t.p3.X {
		t.deltaX = -t.deltaX
	}
}

func (t *testBezier) destroy(w *gux.Window) {
}
