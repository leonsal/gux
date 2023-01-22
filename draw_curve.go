package gux

import "github.com/leonsal/gux/gb"

func (w *Window) AddBezierQuadratic(dl *gb.DrawList, p1, p2, p3 gb.Vec2, col gb.RGBA, thickness float32, numSegments int) {

	if (col & gb.RGBAMaskA) == 0 {
		return
	}
	w.PathLineTo(dl, p1)
	w.PathBezierQuadraticCurveTo(dl, p2, p3, numSegments)
	w.PathStroke(dl, col, 0, thickness)
}

func (w *Window) PathBezierQuadraticCurveTo(dl *gb.DrawList, p2, p3 gb.Vec2, numSegments int) {

}

//void ImDrawList::PathBezierQuadraticCurveTo(const ImVec2& p2, const ImVec2& p3, int num_segments)
//{
//    ImVec2 p1 = _Path.back();
//    if (num_segments == 0)
//    {
//        IM_ASSERT(_Data->CurveTessellationTol > 0.0f);
//        PathBezierQuadraticCurveToCasteljau(&_Path, p1.x, p1.y, p2.x, p2.y, p3.x, p3.y, _Data->CurveTessellationTol, 0);// Auto-tessellated
//    }
//    else
//    {
//        float t_step = 1.0f / (float)num_segments;
//        for (int i_step = 1; i_step <= num_segments; i_step++)
//            _Path.push_back(ImBezierQuadraticCalc(p1, p2, p3, t_step * i_step));
//    }
//}

func bezierQuadraticCalc(p1, p2, p3 gb.Vec2, t float32) gb.Vec2 {

	u := float32(1.0 - t)
	w1 := u * u
	w2 := 2 * u * t
	w3 := t * t
	return gb.Vec2{w1*p1.X + w2*p2.X + w3*p3.X, w1*p1.Y + w2*p2.Y + w3*p3.Y}
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

//
