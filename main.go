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

		cmd1 := gb.DrawCmd{
			ClipRect: gb.Vec4{1, 2, 3, 4},
			TexId:    0,
			Indices:  []uint32{0, 1, 2},
			Vertices: []gb.Vertex{
				{Pos: gb.Vec2{-0.5, 0.0}, UV: gb.Vec2{0, 0}, Col: 0xFFFFFF},
				{Pos: gb.Vec2{0.5, 0.0}, UV: gb.Vec2{1, 1}, Col: 0xFFFFFF},
				{Pos: gb.Vec2{0.0, -0.5}, UV: gb.Vec2{1, 1}, Col: 0xFFFFFFF},
			},
		}
		drawList.AddCmd(cmd1)

		//cmd2 := gb.DrawCmd{
		//	ClipRect: gb.Vec4{10, 20, 30, 40},
		//	TexId:    1,
		//	Indices:  []uint32{50, 60, 70},
		//	Vertices: []gb.Vertex{
		//		{Pos: gb.Vec2{80, 90}, UV: gb.Vec2{100, 110}, Col: 120},
		//		{Pos: gb.Vec2{130, 140}, UV: gb.Vec2{150, 160}, Col: 170},
		//	},
		//}
		//drawList.AddCmd(cmd2)

		//type DrawCmd struct {
		//	ClipRect Vec4     // Clip rectangle
		//	TexId    int      // Texture ID
		//	Indices  []uint32 // Array of vertices indices
		//	Vertices []Vertex // Array of vertices info
		//}

		//		drawList.AddCmd(gb.DrawCmd{gb.Vec4{5, 6, 7, 8}, 20, []uint32{7, 8, 9, 10}, []float32{30, 40, 50, 60}})
		win.RenderFrame(drawList)
		drawList.Clear()
	}
	win.Destroy()
}
