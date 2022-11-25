package main

import (
	"runtime"

	"github.com/leonsal/gux/gb"
)

func main() {

	runtime.LockOSThread()
	win, err := gb.CreateWindow("title", 800, 600)
	if err != nil {
		panic(err)
	}

	drawList := gb.NewDrawList()
	for win.StartFrame(0) {

		drawList.AddCmd(gb.Vec4{1, 2, 3, 4}, 10, []uint32{4, 5, 6}, []float32{10, 20, 30})
		drawList.AddCmd(gb.Vec4{5, 6, 7, 8}, 20, []uint32{7, 8, 9, 10}, []float32{30, 40, 50, 60})
		win.RenderFrame(drawList)
		drawList.Clear()
	}
	win.Destroy()
}
