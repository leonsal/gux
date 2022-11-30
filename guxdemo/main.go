package main

import (
	"runtime"

	"github.com/leonsal/gux/gb"
)

func main() {

	runtime.LockOSThread()
	win, err := gb.CreateWindow("title", 1000, 1000)
	if err != nil {
		panic(err)
	}

	drawList1 := gb.DrawList{}

	cmd1 := gb.DrawCmd{}
	cmd1.AddIndices(0, 1, 2, 2, 3, 0)
	cmd1.AddVertices(
		gb.Vertex{Pos: gb.Vec2{10, 10}, Col: 0xFF_FF_00_00},
		gb.Vertex{Pos: gb.Vec2{10, 100}, Col: 0xFF_FF_00_00},
		gb.Vertex{Pos: gb.Vec2{200, 100}, Col: 0xFF_FF_00_00},
		gb.Vertex{Pos: gb.Vec2{200, 10}, Col: 0xFF_FF_00_00},
	)
	drawList1.AddCmd(cmd1)

	cmd2 := gb.DrawCmd{}
	cmd2.AddIndices(0, 1, 2, 2, 3, 0)
	cmd2.AddVertices(
		gb.Vertex{Pos: gb.Vec2{500, 0}, Col: 0xFF_00_00_FF},
		gb.Vertex{Pos: gb.Vec2{500, 250}, Col: 0xFF_00_00_FF},
		gb.Vertex{Pos: gb.Vec2{750, 250}, Col: 0xFF_00_00_FF},
		gb.Vertex{Pos: gb.Vec2{750, 0}, Col: 0xFF_00_00_FF},
	)
	drawList1.AddCmd(cmd2)

	drawList2 := gb.DrawList{}
	cmd3 := gb.DrawCmd{}
	cmd3.AddIndices(0, 1, 2)
	cmd3.AddVertices(
		gb.Vertex{Pos: gb.Vec2{200, 800}, Col: 0xFF_00_00_FF},
		gb.Vertex{Pos: gb.Vec2{400, 800}, Col: 0xFF_00_FF_00},
		gb.Vertex{Pos: gb.Vec2{300, 600}, Col: 0xFF_FF_00_00},
	)
	drawList2.AddCmd(cmd3)

	cmd4 := gb.DrawCmd{}
	cmd4.AddIndices(0, 1, 2)
	cmd4.AddVertices(
		gb.Vertex{Pos: gb.Vec2{700, 800}, Col: 0xFF_00_00_FF},
		gb.Vertex{Pos: gb.Vec2{900, 800}, Col: 0xFF_00_FF_00},
		gb.Vertex{Pos: gb.Vec2{800, 500}, Col: 0xFF_FF_00_00},
	)
	drawList2.AddCmd(cmd4)

	drawList := gb.DrawList{}
	drawList.AddList(drawList1)
	drawList.AddList(drawList2)
	for win.StartFrame(0) {
		win.RenderFrame(&drawList)
	}
	win.Destroy()
}