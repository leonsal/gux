package gux

import "github.com/leonsal/gux/gb"

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

	dlsd := new(DrawListSharedData)

	return dlsd

}

//struct IMGUI_API ImDrawListSharedData
//{
//    ImVec2          TexUvWhitePixel;            // UV of white pixel in the atlas
//    ImFont*         Font;                       // Current/default font (optional, for simplified AddText overload)
//    float           FontSize;                   // Current/default font size (optional, for simplified AddText overload)
//    float           CurveTessellationTol;       // Tessellation tolerance when using PathBezierCurveTo()
//    float           CircleSegmentMaxError;      // Number of circle segments to use per pixel of radius for AddCircle() etc
//    ImVec4          ClipRectFullscreen;         // Value for PushClipRectFullscreen()
//    ImDrawListFlags InitialFlags;               // Initial flags at the beginning of the frame (it is possible to alter flags on a per-drawlist basis afterwards)
//
//    // [Internal] Temp write buffer
//    ImVector<ImVec2> TempBuffer;
//
//    // [Internal] Lookup tables
//    ImVec2          ArcFastVtx[IM_DRAWLIST_ARCFAST_TABLE_SIZE]; // Sample points on the quarter of the circle.
//    float           ArcFastRadiusCutoff;                        // Cutoff radius after which arc drawing will fallback to slower PathArcTo()
//    ImU8            CircleSegmentCounts[64];    // Precomputed segment count for given radius before we calculate it dynamically (to avoid calculation overhead)
//    const ImVec4*   TexUvLines;                 // UV of anti-aliased lines in the atlas
//
//    ImDrawListSharedData();
//    void SetCircleTessellationMaxError(float max_error);
//};

// DrawList return this window DrawList
func (w *Window) DrawList() *gb.DrawList {

	return &w.dl
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
func (w *Window) PathFillConvex(dl *gb.DrawList, col gb.Color) {

	w.AddConvexPolyFilled(dl, dl.Path, col)
	w.PathClear(dl)
}

func (w *Window) PathStroke(dl *gb.DrawList, col gb.Color, flags DrawFlags, thickness float32) {

	w.AddPolyLine(dl, dl.Path, col, flags, thickness)
	w.PathClear(dl)
}

func (w *Window) PathArcTo(dl *gb.DrawList, center gb.Vec2, radius, amin, amax float32, numSegments int) {

}

func (w *Window) PathArcToFast(dl *gb.DrawList, center gb.Vec2, radius, amin12, amax12 float32) {

}

func (w *Window) PathRect(dl *gb.DrawList, rectMin, rectMax gb.Vec2, rounding float32, flags DrawFlags) {

}

func (w *Window) AddLine(dl *gb.DrawList, p1, p2 gb.Vec2, col gb.Color, thickness float32) {

	if (col & gb.ColorMaskA) == 0 {
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

// #define IM_DRAWLIST_CIRCLE_AUTO_SEGMENT_CALC(_RAD,_MAXERROR)    ImClamp(IM_ROUNDUP_TO_EVEN((int)ImCeil(IM_PI / ImAcos(1 - ImMin((_MAXERROR), (_RAD)) / (_RAD)))), IM_DRAWLIST_CIRCLE_AUTO_SEGMENT_MIN, IM_DRAWLIST_CIRCLE_AUTO_SEGMENT_MAX)
func drawListCircleAutoSegmentCalc(rad, maxerror float32) {

}

//#define IM_DRAWLIST_CIRCLE_AUTO_SEGMENT_CALC(_RAD,_MAXERROR)    ImClamp(IM_ROUNDUP_TO_EVEN((int)ImCeil(IM_PI / ImAcos(1 - ImMin((_MAXERROR), (_RAD)) / (_RAD)))), IM_DRAWLIST_CIRCLE_AUTO_SEGMENT_MIN, IM_DRAWLIST_CIRCLE_AUTO_SEGMENT_MAX)

//#define IM_DRAWLIST_CIRCLE_AUTO_SEGMENT_CALC_R(_N,_MAXERROR)    ((_MAXERROR) / (1 - ImCos(IM_PI / ImMax((float)(_N), IM_PI))))
//func drawListCircleAutoSegmentCalcR(n,

// Raw equation from IM_DRAWLIST_CIRCLE_AUTO_SEGMENT_CALC rewritten for 'r' and 'error'.
//#define IM_DRAWLIST_CIRCLE_AUTO_SEGMENT_CALC_ERROR(_N,_RAD)     ((1 - ImCos(IM_PI / ImMax((float)(_N), IM_PI))) / (_RAD))
