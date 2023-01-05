package gux

import (
	"math"

	"github.com/leonsal/gux/gb"
)

// DrawList return this window DrawList
func (w *Window) DrawList() *gb.DrawList {

	return &w.dl
}

// NewDrawCmd creates and appends a new empty DrawCmd in the specified DrawList
// and returns pointer to the command and slices for setting vertex indices and info.
// The returned command ClipRect is initialized with current window size and the texture id is set to the opaque white dot.
func (w *Window) NewDrawCmd(dl *gb.DrawList, idxCount, vtxCount int) (*gb.DrawCmd, []uint32, []gb.Vertex) {

	cmd, bufIdx, bufVtx := dl.NewDrawCmd(idxCount, vtxCount)
	cmd.ClipRect = gb.Vec4{0, 0, w.frameInfo.WinSize.X, w.frameInfo.WinSize.Y}
	cmd.TexID = w.TexWhiteId
	return cmd, bufIdx, bufVtx
}

//// PathClear clears the draw list path
//func (w *Window) PathClear(dl *gb.DrawList) {
//
//	dl.Path = dl.Path[:0]
//}

// PathLineTo adds a point to the current draw list path
func (w *Window) PathLineTo(dl *gb.DrawList, pos gb.Vec2) {

	dl.PathAppend(pos)
}

// PathFillConvex adds to the drawlist a filled convex polygon with the points in the drawlist current path.
// The path is cleared.
func (w *Window) PathFillConvex(dl *gb.DrawList, col gb.RGBA) {

	w.AddConvexPolyFilled(dl, dl.Path, col)
	dl.PathClear()
}

func (w *Window) PathStroke(dl *gb.DrawList, col gb.RGBA, flags DrawFlags, thickness float32) {

	w.AddPolyLine(dl, dl.Path, col, flags, thickness)
	dl.PathClear()
}

func (w *Window) PathArcTo(dl *gb.DrawList, center gb.Vec2, radius, amin, amax float32, numSegments int) {

	if radius < 0.5 {
		dl.PathAppend(center)
		return
	}

	dl.PathReserve(numSegments + 1)
	// Note that we are adding a point at both a_min and a_max.
	// If you are trying to draw a full closed circle you don't want the overlapping points!
	for i := 0; i <= numSegments; i++ {
		a := amin + (float32(i)/float32(numSegments))*(amax-amin)
		dl.PathAppend(gb.Vec2{center.X + Cos(a)*radius, center.Y + Sin(a)*radius})
	}
}

//void ImDrawList::PathRect(const ImVec2& a, const ImVec2& b, float rounding, ImDrawFlags flags)
//{
//    flags = FixRectCornerFlags(flags);
//    rounding = ImMin(rounding, ImFabs(b.x - a.x) * ( ((flags & ImDrawFlags_RoundCornersTop)  == ImDrawFlags_RoundCornersTop)  || ((flags & ImDrawFlags_RoundCornersBottom) == ImDrawFlags_RoundCornersBottom) ? 0.5f : 1.0f ) - 1.0f);
//    rounding = ImMin(rounding, ImFabs(b.y - a.y) * ( ((flags & ImDrawFlags_RoundCornersLeft) == ImDrawFlags_RoundCornersLeft) || ((flags & ImDrawFlags_RoundCornersRight)  == ImDrawFlags_RoundCornersRight)  ? 0.5f : 1.0f ) - 1.0f);
//
//    if (rounding < 0.5f || (flags & ImDrawFlags_RoundCornersMask_) == ImDrawFlags_RoundCornersNone)
//    {
//        PathLineTo(a);
//        PathLineTo(ImVec2(b.x, a.y));
//        PathLineTo(b);
//        PathLineTo(ImVec2(a.x, b.y));
//    }
//    else
//    {
//        const float rounding_tl = (flags & ImDrawFlags_RoundCornersTopLeft)     ? rounding : 0.0f;
//        const float rounding_tr = (flags & ImDrawFlags_RoundCornersTopRight)    ? rounding : 0.0f;
//        const float rounding_br = (flags & ImDrawFlags_RoundCornersBottomRight) ? rounding : 0.0f;
//        const float rounding_bl = (flags & ImDrawFlags_RoundCornersBottomLeft)  ? rounding : 0.0f;
//        PathArcToFast(ImVec2(a.x + rounding_tl, a.y + rounding_tl), rounding_tl, 6, 9);
//        PathArcToFast(ImVec2(b.x - rounding_tr, a.y + rounding_tr), rounding_tr, 9, 12);
//        PathArcToFast(ImVec2(b.x - rounding_br, b.y - rounding_br), rounding_br, 0, 3);
//        PathArcToFast(ImVec2(a.x + rounding_bl, b.y - rounding_bl), rounding_bl, 3, 6);
//    }
//}

func (w *Window) PathRect(dl *gb.DrawList, min, max gb.Vec2, rounding float32, flags DrawFlags) {

	if rounding < 0.5 || (flags&DrawFlags_RoundCornersMask_ == 0) {
		w.PathLineTo(dl, min)
		w.PathLineTo(dl, gb.Vec2{max.X, min.Y})
		w.PathLineTo(dl, max)
		w.PathLineTo(dl, gb.Vec2{min.X, max.Y})
		return
	}
	rtl := float32(0)
	if flags&DrawFlags_RoundCornersTopLeft != 0 {
		rtl = rounding
	}
	rtr := float32(0)
	if flags&DrawFlags_RoundCornersTopRight != 0 {
		rtr = rounding
	}
	rbr := float32(0)
	if flags&DrawFlags_RoundCornersBottomRight != 0 {
		rbr = rounding
	}
	rbl := float32(0)
	if flags&DrawFlags_RoundCornersBottomLeft != 0 {
		rbl = rounding
	}
	w.PathArcTo(dl, gb.Vec2{min.X + rtl, min.Y + rtl}, rtl, -2*math.Pi/2, -2*math.Pi/4, 32)
	w.PathArcTo(dl, gb.Vec2{max.X - rtr, min.Y + rtr}, rtr, -2*math.Pi/4, 0, 16)
	w.PathArcTo(dl, gb.Vec2{max.X - rbr, max.Y - rbr}, rbr, 0, 2*math.Pi/4, 16)
	w.PathArcTo(dl, gb.Vec2{min.X + rbl, max.Y - rbl}, rbl, 2*math.Pi/4, math.Pi, 16)

}

func (w *Window) AddLine(dl *gb.DrawList, p1, p2 gb.Vec2, col gb.RGBA, thickness float32) {

	if (col & gb.RGBAMaskA) == 0 {
		return
	}
	w.PathLineTo(dl, gb.Vec2Add(p1, gb.Vec2{0.5, 0.5}))
	w.PathLineTo(dl, gb.Vec2Add(p2, gb.Vec2{0.5, 0.5}))
	w.PathStroke(dl, col, 0, thickness)
}

func (w *Window) AddRect(dl *gb.DrawList, pmin, max gb.Vec2, rounding float32, flags DrawFlags, thickness float32) {

}

func (w *Window) AddCircle(dl *gb.DrawList, center gb.Vec2, radius float32, col gb.RGBA, numSegments int, thickness float32) {

	if (col&gb.RGBAMaskA) == 0 || radius < 0.5 {
		return
	}
	numSegments = Clamp(numSegments, 3, DrawListCircleSegmentMax)

	// Because we are filling a closed shape we remove 1 from the count of segments/points
	amax := (2 * math.Pi) * (float64(numSegments) - 1.0) / float64(numSegments)
	w.PathArcTo(dl, center, radius-0.5, 0.0, float32(amax), numSegments-1)
	w.PathStroke(dl, col, DrawFlags_Closed, thickness)
}

func (w *Window) AddCircleFilled(dl *gb.DrawList, center gb.Vec2, radius float32, col gb.RGBA, numSegments int) {

	if (col&gb.RGBAMaskA) == 0 || radius < 0.5 {
		return
	}
	numSegments = Clamp(numSegments, 3, DrawListCircleSegmentMax)

	// Because we are filling a closed shape we remove 1 from the count of segments/points
	amax := (2 * math.Pi) * (float64(numSegments) - 1.0) / float64(numSegments)
	w.PathArcTo(dl, center, radius-0.5, 0.0, float32(amax), numSegments-1)
	w.PathFillConvex(dl, col)
}

func roundupToEven(v int) int {
	return ((v + 1) / 2 * 2)
}
