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

	drawList := gb.DrawList{}

	cmd1 := gb.DrawCmd{}
	cmd1.AddIndices(0, 1, 2)
	cmd1.AddVertices(
		gb.Vertex{Pos: gb.Vec2{400, 500}, Col: 0xFF_00_00_FF},
		gb.Vertex{Pos: gb.Vec2{600, 500}, Col: 0xFF_00_FF_00},
		gb.Vertex{Pos: gb.Vec2{500, 250}, Col: 0xFF_FF_00_00},
	)
	drawList.AddCmd(cmd1)

	cmd2 := gb.DrawCmd{}
	cmd2.AddIndices(0, 1, 2)
	cmd2.AddVertices(
		gb.Vertex{Pos: gb.Vec2{0, 0}, Col: 0xFF_00_00_FF},
		gb.Vertex{Pos: gb.Vec2{100, 0}, Col: 0xFF_00_00_FF},
		gb.Vertex{Pos: gb.Vec2{100, 200}, Col: 0xFF_00_00_FF},
	)
	drawList.AddCmd(cmd2)

	for win.StartFrame(0) {
		win.RenderFrame(&drawList)
	}
	win.Destroy()
}
