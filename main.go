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

	for win.StartFrame(0) {

		win.RenderFrame()
	}
	win.Destroy()
}
