package main

import (
	"github.com/leonsal/gux/gb"
	"github.com/leonsal/gux/window"
)

func init() {

	registerTest("rect", 5, newTestRect)
}

type testRect struct{}

func newTestRect(w *window.Window) ITest {

	return new(testRect)
}

func (t *testRect) draw(win *window.Window) {

	dl := win.DrawList()
	startY := float32(50)
	const startX = 50
	const width = 300
	const height = 150
	const deltaX = width + 10
	const deltaY = height + 50
	const rounding = 40.0
	const thickness = 10
	flagList := []window.DrawFlags{
		window.DrawFlags_RoundCornersTopLeft,
		window.DrawFlags_RoundCornersTopRight,
		window.DrawFlags_RoundCornersBottomRight,
		window.DrawFlags_RoundCornersBottomLeft,
		window.DrawFlags_RoundCornersTop,
		window.DrawFlags_RoundCornersBottom,
		window.DrawFlags_RoundCornersLeft,
		window.DrawFlags_RoundCornersRight,
		window.DrawFlags_RoundCornersAll,
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

func (t *testRect) destroy(w *window.Window) {

}
