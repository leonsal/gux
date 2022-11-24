package main

import (
	"runtime"

	"github.com/leonsal/gux/gl"
)

func main() {

	runtime.LockOSThread()
	win, err := gl.CreateWindow("title", 800, 600)
	if err != nil {
		panic(err)
	}

	drawList := gl.NewDrawList()

	for win.StartFrame(0) {

		drawList.AddCmd(gl.Vec4{1, 2, 3, 4}, 5, 6, 7, 8)
		drawList.AddCmd(gl.Vec4{10, 20, 30, 40}, 50, 60, 70, 80)
		win.RenderFrame(drawList)
		drawList.Clear()
	}
	win.Destroy()
}
