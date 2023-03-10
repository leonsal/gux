package window

import (
	"github.com/leonsal/gux/gb"
	"github.com/leonsal/gux/util"
)

func (w *Window) AddBezierQuadratic(dl *gb.DrawList, p1, p2, p3 gb.Vec2, col gb.RGBA, thickness float32, numSegments int) {

	if (col & gb.RGBAMaskA) == 0 {
		return
	}
	w.PathLineTo(dl, p1)
	w.PathBezierQuadraticCurveTo(dl, p2, p3, numSegments)
	w.PathStroke(dl, col, 0, thickness)
}

func (w *Window) PathBezierQuadraticCurveTo(dl *gb.DrawList, p2, p3 gb.Vec2, numSegments int) {

	p1 := dl.PathBack()
	if numSegments == 0 {
		util.Assert(w.CurveTessellationTol > 0, "")
		pathBezierQuadraticCurveToCasteljau(&dl.Path, p1.X, p1.Y, p2.X, p2.Y, p3.X, p3.Y, w.CurveTessellationTol, 0)
	} else {
		tstep := 1.0 / float32(numSegments)
		for istep := 1; istep <= numSegments; istep++ {
			dl.PathAppend(bezierQuadraticCalc(p1, p2, p3, tstep*float32(istep)))
		}
	}
}

func (w *Window) AddBezierCubic(dl *gb.DrawList, p1, p2, p3, p4 gb.Vec2, col gb.RGBA, thickness float32, numSegments int) {

	if (col & gb.RGBAMaskA) == 0 {
		return
	}
	w.PathLineTo(dl, p1)
	w.PathBezierCubicCurveTo(dl, p2, p3, p4, numSegments)
	w.PathStroke(dl, col, 0, thickness)
}

func (w *Window) PathBezierCubicCurveTo(dl *gb.DrawList, p2, p3, p4 gb.Vec2, numSegments int) {

	p1 := dl.PathBack()
	if numSegments == 0 {
		util.Assert(w.CurveTessellationTol > 0, "")
		pathBezierCubicCurveToCasteljau(&dl.Path, p1.X, p1.Y, p2.X, p2.Y, p3.X, p3.Y, p4.X, p4.Y, w.CurveTessellationTol, 0)
	} else {
		tstep := 1.0 / float32(numSegments)
		for istep := 1; istep <= numSegments; istep++ {
			dl.PathAppend(bezierCubicCalc(p1, p2, p3, p4, tstep*float32(istep)))
		}
	}
}

func bezierCubicCalc(p1, p2, p3, p4 gb.Vec2, t float32) gb.Vec2 {

	u := 1.0 - t
	w1 := u * u * u
	w2 := 3 * u * u * t
	w3 := 3 * u * t * t
	w4 := t * t * t
	return gb.Vec2{w1*p1.X + w2*p2.X + w3*p3.X + w4*p4.X, w1*p1.Y + w2*p2.Y + w3*p3.Y + w4*p4.Y}
}

func bezierQuadraticCalc(p1, p2, p3 gb.Vec2, t float32) gb.Vec2 {

	u := float32(1.0 - t)
	w1 := u * u
	w2 := 2 * u * t
	w3 := t * t
	return gb.Vec2{w1*p1.X + w2*p2.X + w3*p3.X, w1*p1.Y + w2*p2.Y + w3*p3.Y}
}

func pathBezierCubicCurveToCasteljau(path *[]gb.Vec2, x1, y1, x2, y2, x3, y3, x4, y4, tessTol float32, level int) {

	dx := x4 - x1
	dy := y4 - y1
	d2 := (x2-x4)*dy - (y2-y4)*dx
	d3 := (x3-x4)*dy - (y3-y4)*dx
	if d2 < 0 {
		d2 = -d2
	}
	if d3 < 0 {
		d3 = -d3
	}
	if (d2+d3)*(d2+d3) < tessTol*(dx*dx+dy*dy) {
		*path = append(*path, gb.Vec2{x4, y4})
	} else if level < 10 {
		x12 := (x1 + x2) * 0.5
		y12 := (y1 + y2) * 0.5
		x23 := (x2 + x3) * 0.5
		y23 := (y2 + y3) * 0.5
		x34 := (x3 + x4) * 0.5
		y34 := (y3 + y4) * 0.5
		x123 := (x12 + x23) * 0.5
		y123 := (y12 + y23) * 0.5
		x234 := (x23 + x34) * 0.5
		y234 := (y23 + y34) * 0.5
		x1234 := (x123 + x234) * 0.5
		y1234 := (y123 + y234) * 0.5
		pathBezierCubicCurveToCasteljau(path, x1, y1, x12, y12, x123, y123, x1234, y1234, tessTol, level+1)
		pathBezierCubicCurveToCasteljau(path, x1234, y1234, x234, y234, x34, y34, x4, y4, tessTol, level+1)
	}
}

func pathBezierQuadraticCurveToCasteljau(path *[]gb.Vec2, x1, y1, x2, y2, x3, y3, tessTol float32, level int) {

	dx := x3 - x1
	dy := y3 - y1
	det := (x2-x3)*dy - (y2-y3)*dx
	if det*det*4.0 < tessTol*(dx*dx+dy*dy) {
		*path = append(*path, gb.Vec2{x3, y3})
	} else if level < 10 {
		x12 := (x1 + x2) * 0.5
		y12 := (y1 + y2) * 0.5
		x23 := (x2 + x3) * 0.5
		y23 := (y2 + y3) * 0.5
		x123 := (x12 + x23) * 0.5
		y123 := (y12 + y23) * 0.5
		pathBezierQuadraticCurveToCasteljau(path, x1, y1, x12, y12, x123, y123, tessTol, level+1)
		pathBezierQuadraticCurveToCasteljau(path, x123, y123, x23, y23, x3, y3, tessTol, level+1)
	}
}
