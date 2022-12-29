package gb

import (
	"fmt"
	"runtime"
	"testing"
)

func Test1(t *testing.T) {

	runtime.LockOSThread()
	win, err := CreateWindow("title", 1000, 1000, nil)
	if err != nil {
		panic(err)
	}

	// Creates image with one white opaque pixel
	var rect [1]RGBA
	rect[0] = MakeColor(255, 255, 255, 255)

	// Creates and transfer texture
	texId := win.CreateTexture(1, 1, &rect[0])

	// DrawList 1
	drawList1 := DrawList{}
	cmd1, idxBuf1, vtxBuf1 := drawList1.NewDrawCmd(6, 4)
	cmd1.ClipRect = Vec4{0, 0, 4000, 4000}
	cmd1.TexID = texId
	idxBuf1[0] = 0
	idxBuf1[1] = 1
	idxBuf1[2] = 2
	idxBuf1[3] = 2
	idxBuf1[4] = 3
	idxBuf1[5] = 0
	vtxBuf1[0] = Vertex{Pos: Vec2{10, 10}, Col: 0xFF_FF_00_00}
	vtxBuf1[1] = Vertex{Pos: Vec2{10, 100}, Col: 0xFF_FF_00_00}
	vtxBuf1[2] = Vertex{Pos: Vec2{200, 100}, Col: 0xFF_FF_00_00}
	vtxBuf1[2] = Vertex{Pos: Vec2{200, 10}, Col: 0xFF_FF_00_00}

	cmd2, idxBuf2, vtxBuf2 := drawList1.NewDrawCmd(6, 4)
	cmd2.ClipRect = Vec4{0, 0, 4000, 4000}
	cmd2.TexID = texId
	idxBuf2[0] = 0
	idxBuf2[1] = 1
	idxBuf2[2] = 2
	idxBuf2[3] = 2
	idxBuf2[4] = 3
	idxBuf2[5] = 0
	vtxBuf2[0] = Vertex{Pos: Vec2{500, 0}, Col: 0xFF_00_00_FF}
	vtxBuf2[1] = Vertex{Pos: Vec2{500, 250}, Col: 0xFF_00_00_FF}
	vtxBuf2[2] = Vertex{Pos: Vec2{750, 250}, Col: 0xFF_00_00_FF}
	vtxBuf2[3] = Vertex{Pos: Vec2{750, 0}, Col: 0xFF_00_00_FF}

	//	drawList1.AddCmd(Vec4{}, texId,
	//		[]uint32{0, 1, 2, 2, 3, 0},
	//		[]Vertex{
	//			Vertex{Pos: Vec2{10, 10}, Col: 0xFF_FF_00_00},
	//			Vertex{Pos: Vec2{10, 100}, Col: 0xFF_FF_00_00},
	//			Vertex{Pos: Vec2{200, 100}, Col: 0xFF_FF_00_00},
	//			Vertex{Pos: Vec2{200, 10}, Col: 0xFF_FF_00_00},
	//		},
	//	)
	//	drawList1.AddCmd(Vec4{}, texId,
	//		[]uint32{0, 1, 2, 2, 3, 0},
	//		[]Vertex{
	//			Vertex{Pos: Vec2{500, 0}, Col: 0xFF_00_00_FF},
	//			Vertex{Pos: Vec2{500, 250}, Col: 0xFF_00_00_FF},
	//			Vertex{Pos: Vec2{750, 250}, Col: 0xFF_00_00_FF},
	//			Vertex{Pos: Vec2{750, 0}, Col: 0xFF_00_00_FF},
	//		},
	//	)
	//
	//	// DrawList 2
	//	drawList2 := DrawList{}
	//	drawList2.AddCmd(Vec4{}, texId,
	//		[]uint32{0, 1, 2},
	//		[]Vertex{
	//			Vertex{Pos: Vec2{200, 800}, Col: 0xFF_00_00_FF},
	//			Vertex{Pos: Vec2{400, 800}, Col: 0xFF_00_FF_00},
	//			Vertex{Pos: Vec2{300, 600}, Col: 0xFF_FF_00_00},
	//		},
	//	)
	//	drawList2.AddCmd(Vec4{}, texId,
	//		[]uint32{0, 1, 2},
	//		[]Vertex{
	//			Vertex{Pos: Vec2{700, 800}, Col: 0xFF_00_00_FF},
	//			Vertex{Pos: Vec2{900, 800}, Col: 0xFF_00_FF_00},
	//			Vertex{Pos: Vec2{800, 500}, Col: 0xFF_FF_00_00},
	//		},
	//	)
	//
	//	// Create new DrawList from concatenation of DrawList1 and 2
	//	drawList := DrawList{}
	//	drawList.AddList(drawList1)
	//	drawList.AddList(drawList2)

	// Render loop
	frames := 0
	cgoCallsStart := runtime.NumCgoCall()
	frameParams := FrameParams{}
	for {
		finfo := win.StartFrame(&frameParams)
		if finfo.WinClose {
			break
		}
		win.RenderFrame(&drawList1)
		frames++
		if frames > 500 {
			break
		}
	}
	cgoCalls := runtime.NumCgoCall() - cgoCallsStart
	fmt.Println("cgo calls/frame", cgoCalls/int64(frames))
	win.Destroy()

}
