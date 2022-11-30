package gb

import (
	"fmt"
	"runtime"
	"testing"
)

func Test1(t *testing.T) {

	runtime.LockOSThread()
	win, err := CreateWindow("title", 1000, 1000)
	if err != nil {
		panic(err)
	}

	// DrawList 1
	drawList1 := DrawList{}
	drawList1.AddCmd(Vec4{}, 0,
		[]uint32{0, 1, 2, 2, 3, 0},
		[]Vertex{
			Vertex{Pos: Vec2{10, 10}, Col: 0xFF_FF_00_00},
			Vertex{Pos: Vec2{10, 100}, Col: 0xFF_FF_00_00},
			Vertex{Pos: Vec2{200, 100}, Col: 0xFF_FF_00_00},
			Vertex{Pos: Vec2{200, 10}, Col: 0xFF_FF_00_00},
		},
	)
	drawList1.AddCmd(Vec4{}, 0,
		[]uint32{0, 1, 2, 2, 3, 0},
		[]Vertex{
			Vertex{Pos: Vec2{500, 0}, Col: 0xFF_00_00_FF},
			Vertex{Pos: Vec2{500, 250}, Col: 0xFF_00_00_FF},
			Vertex{Pos: Vec2{750, 250}, Col: 0xFF_00_00_FF},
			Vertex{Pos: Vec2{750, 0}, Col: 0xFF_00_00_FF},
		},
	)

	// DrawList 2
	drawList2 := DrawList{}
	drawList2.AddCmd(Vec4{}, 0,
		[]uint32{0, 1, 2},
		[]Vertex{
			Vertex{Pos: Vec2{200, 800}, Col: 0xFF_00_00_FF},
			Vertex{Pos: Vec2{400, 800}, Col: 0xFF_00_FF_00},
			Vertex{Pos: Vec2{300, 600}, Col: 0xFF_FF_00_00},
		},
	)
	drawList2.AddCmd(Vec4{}, 0,
		[]uint32{0, 1, 2},
		[]Vertex{
			Vertex{Pos: Vec2{700, 800}, Col: 0xFF_00_00_FF},
			Vertex{Pos: Vec2{900, 800}, Col: 0xFF_00_FF_00},
			Vertex{Pos: Vec2{800, 500}, Col: 0xFF_FF_00_00},
		},
	)

	// Create new DrawList from concatenation of DrawList1 and 2
	drawList := DrawList{}
	drawList.AddList(drawList1)
	drawList.AddList(drawList2)

	// Render loop
	frames := 0
	cgoCallsStart := runtime.NumCgoCall()
	for win.StartFrame(0) {
		win.RenderFrame(&drawList)
		frames++
		if frames > 500 {
			break
		}
	}
	cgoCalls := runtime.NumCgoCall() - cgoCallsStart
	fmt.Println("cgo calls/frame", cgoCalls/int64(frames))
	win.Destroy()

}
