package gux

import (
	"math"

	"github.com/leonsal/gux/gb"
)

const (
	DrawListCircleAutoSegmentMin = 4
	DrawListCircleAutoSegmentMax = 512

	// Lookup table size for adaptive arc drawing, cover full circle.
	DrawListArcFastTableSize = 48

	// Sample index _PathArcToFastEx() for 360 angle.
	DrawListArcFastSampleMax = DrawListArcFastTableSize
)

type DrawListSharedData struct {
	CircleSegmentMaxError float32                           // Number of circle segments to use per pixel of radius for AddCircle()
	ArcFastVtx            [DrawListArcFastTableSize]gb.Vec2 // Sample points on the quarter of the circle
	ArcFastRadiusCutoff   float32                           // Cutoff radius after which arc drawing will fallback to slower PathArcTo()
	CircleSegmentCounts   [64]byte                          // Precomputed segment count for given radius before we calculated it dynamically
}

func NewDrawListSharedData() *DrawListSharedData {

	sd := new(DrawListSharedData)

	for i := 0; i < len(sd.ArcFastVtx); i++ {
		a := (float32(i) * 2 * math.Pi) / float32(len(sd.ArcFastVtx))
		sd.ArcFastVtx[i] = gb.Vec2{Cos(a), Sin(a)}

	}
	sd.ArcFastRadiusCutoff = float32(drawListCircleAutoSegmentCalc(DrawListArcFastSampleMax, sd.CircleSegmentMaxError))

	return sd

}

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

// PathClear clears the draw list path
func (w *Window) PathClear(dl *gb.DrawList) {

	dl.Path = dl.Path[:0]
}

// PathLineTo adds a point to the current draw list path
func (w *Window) PathLineTo(dl *gb.DrawList, pos gb.Vec2) {

	dl.Path = append(dl.Path, pos)
}

// PathFillConvex adds to the drawlist a filled convex polygon with the points in the drawlist current path.
// The path is cleared.
func (w *Window) PathFillConvex(dl *gb.DrawList, col gb.RGBA) {

	w.AddConvexPolyFilled(dl, dl.Path, col)
	w.PathClear(dl)
}

func (w *Window) PathStroke(dl *gb.DrawList, col gb.RGBA, flags DrawFlags, thickness float32) {

	w.AddPolyLine(dl, dl.Path, col, flags, thickness)
	w.PathClear(dl)
}

func (w *Window) PathArcTo(dl *gb.DrawList, center gb.Vec2, radius, amin, amax float32, numSegments int) {

}

func (w *Window) PathArcToFast(dl *gb.DrawList, center gb.Vec2, radius, amin12, amax12 float32) {

}

func (w *Window) PathRect(dl *gb.DrawList, rectMin, rectMax gb.Vec2, rounding float32, flags DrawFlags) {

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

func roundupToEven(v int) int {
	return ((v + 1) / 2 * 2)
}

// #define IM_DRAWLIST_CIRCLE_AUTO_SEGMENT_CALC(_RAD,_MAXERROR)
//
//	ImClamp(IM_ROUNDUP_TO_EVEN((int)ImCeil(IM_PI / ImAcos(1 - ImMin((_MAXERROR), (_RAD)) / (_RAD)))), IM_DRAWLIST_CIRCLE_AUTO_SEGMENT_MIN, IM_DRAWLIST_CIRCLE_AUTO_SEGMENT_MAX)
func drawListCircleAutoSegmentCalc(rad, maxerror float32) int {

	return Clamp(roundupToEven(int(math.Ceil(math.Pi/math.Cos(float64(1-Min(maxerror, rad)/rad))))), DrawListCircleAutoSegmentMin, DrawListCircleAutoSegmentMax)

}

//#define IM_DRAWLIST_CIRCLE_AUTO_SEGMENT_CALC(_RAD,_MAXERROR)    ImClamp(IM_ROUNDUP_TO_EVEN((int)ImCeil(IM_PI / ImAcos(1 - ImMin((_MAXERROR), (_RAD)) / (_RAD)))), IM_DRAWLIST_CIRCLE_AUTO_SEGMENT_MIN, IM_DRAWLIST_CIRCLE_AUTO_SEGMENT_MAX)

//#define IM_DRAWLIST_CIRCLE_AUTO_SEGMENT_CALC_R(_N,_MAXERROR)    ((_MAXERROR) / (1 - ImCos(IM_PI / ImMax((float)(_N), IM_PI))))
//func drawListCircleAutoSegmentCalcR(n,

// Raw equation from IM_DRAWLIST_CIRCLE_AUTO_SEGMENT_CALC rewritten for 'r' and 'error'.
//#define IM_DRAWLIST_CIRCLE_AUTO_SEGMENT_CALC_ERROR(_N,_RAD)     ((1 - ImCos(IM_PI / ImMax((float)(_N), IM_PI))) / (_RAD))
