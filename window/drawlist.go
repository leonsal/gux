package window

import "github.com/leonsal/gux/gb"

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
