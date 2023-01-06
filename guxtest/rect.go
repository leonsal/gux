package main

import (
	"github.com/leonsal/gux"
	"github.com/leonsal/gux/gb"
)

func init() {

	registerTest("rect", 5, newTestRect)
}

type testRect struct{}

func newTestRect(w *gux.Window) ITest {

	return new(testRect)
}

func (t *testRect) draw(win *gux.Window) {

	dl := win.DrawList()
	startY := float32(50)
	const startX = 50
	const width = 300
	const height = 150
	const deltaX = width + 10
	const deltaY = height + 50
	const rounding = 40.0
	const thickness = 10
	flagList := []gux.DrawFlags{
		gux.DrawFlags_RoundCornersTopLeft,
		gux.DrawFlags_RoundCornersTopRight,
		gux.DrawFlags_RoundCornersBottomRight,
		gux.DrawFlags_RoundCornersBottomLeft,
		gux.DrawFlags_RoundCornersTop,
		gux.DrawFlags_RoundCornersBottom,
		gux.DrawFlags_RoundCornersLeft,
		gux.DrawFlags_RoundCornersRight,
		gux.DrawFlags_RoundCornersAll,
	}

	for idx, flag := range flagList {
		line := float32(idx / 5)
		col := float32(idx % 5)
		min := gb.Vec2{startX + col*deltaX, startY + line*deltaY}
		max := gb.Vec2{width + col*deltaX, startY + height + line*deltaY}
		win.AddRect(dl, min, max, nextColor(idx), rounding, flag, thickness)
	}

	startY += 2 * deltaY
	for idx, flag := range flagList {
		line := float32(idx / 5)
		col := float32(idx % 5)
		min := gb.Vec2{startX + col*deltaX, startY + line*deltaY}
		max := gb.Vec2{width + col*deltaX, startY + height + line*deltaY}
		win.AddRectFilled(dl, min, max, nextColor(idx), rounding, flag)
	}
}

func (t *testRect) destroy(w *gux.Window) {

}
