package gb

import (
	"fmt"
	"runtime"
	"testing"
)

func Test1(t *testing.T) {

	runtime.LockOSThread()
	cfg := Config{}
	cfg.DebugPrintCmds = false
	cfg.OpenGL.ES = false
	cfg.Vulkan.ValidationLayer = true
	win, err := CreateWindow("title", 1000, 1000, &cfg)
	if err != nil {
		panic(err)
	}

	// Creates image with one white opaque pixel
	var rect [1]RGBA
	rect[0] = MakeColor(255, 255, 255, 255)

	// Creates and transfer 1 pixel opaque white texture needed for all commands
	texId := win.CreateTexture(1, 1, &rect[0])

	// DrawList 1
	drawList1 := DrawList{}
	{
		cmd, idxBuf, vtxBuf := drawList1.NewDrawCmd(6, 4)
		cmd.ClipRect = Vec4{0, 0, 4000, 4000}
		cmd.TexID = texId
		idxBuf[0] = 0
		idxBuf[1] = 1
		idxBuf[2] = 2
		idxBuf[3] = 2
		idxBuf[4] = 3
		idxBuf[5] = 0
		vtxBuf[0] = Vertex{Pos: Vec2{10, 10}, Col: 0xFF_FF_00_00}
		vtxBuf[1] = Vertex{Pos: Vec2{10, 100}, Col: 0xFF_FF_00_00}
		vtxBuf[2] = Vertex{Pos: Vec2{200, 100}, Col: 0xFF_FF_00_00}
		vtxBuf[3] = Vertex{Pos: Vec2{200, 10}, Col: 0xFF_FF_00_00}
	}
	{
		cmd, idxBuf, vtxBuf := drawList1.NewDrawCmd(6, 4)
		cmd.ClipRect = Vec4{0, 0, 4000, 4000}
		cmd.TexID = texId
		idxBuf[0] = 0
		idxBuf[1] = 1
		idxBuf[2] = 2
		idxBuf[3] = 2
		idxBuf[4] = 3
		idxBuf[5] = 0
		vtxBuf[0] = Vertex{Pos: Vec2{500, 0}, Col: 0xFF_00_00_FF}
		vtxBuf[1] = Vertex{Pos: Vec2{500, 250}, Col: 0xFF_00_00_FF}
		vtxBuf[2] = Vertex{Pos: Vec2{750, 250}, Col: 0xFF_00_00_FF}
		vtxBuf[3] = Vertex{Pos: Vec2{750, 0}, Col: 0xFF_00_00_FF}
	}

	// DrawList 2
	drawList2 := DrawList{}
	{
		cmd, idxBuf, vtxBuf := drawList1.NewDrawCmd(3, 3)
		cmd.ClipRect = Vec4{0, 0, 4000, 4000}
		cmd.TexID = texId
		idxBuf[0] = 0
		idxBuf[1] = 1
		idxBuf[2] = 2
		vtxBuf[0] = Vertex{Pos: Vec2{200, 800}, Col: 0xFF_00_00_FF}
		vtxBuf[1] = Vertex{Pos: Vec2{400, 800}, Col: 0xFF_00_FF_00}
		vtxBuf[2] = Vertex{Pos: Vec2{300, 600}, Col: 0xFF_FF_00_FF}
	}
	{
		cmd, idxBuf, vtxBuf := drawList1.NewDrawCmd(3, 3)
		cmd.ClipRect = Vec4{0, 0, 4000, 4000}
		cmd.TexID = texId
		idxBuf[0] = 0
		idxBuf[1] = 1
		idxBuf[2] = 2
		vtxBuf[0] = Vertex{Pos: Vec2{700, 800}, Col: 0xFF_00_00_FF}
		vtxBuf[1] = Vertex{Pos: Vec2{900, 800}, Col: 0xFF_00_FF_00}
		vtxBuf[2] = Vertex{Pos: Vec2{800, 500}, Col: 0xFF_FF_00_FF}
	}

	// Create new DrawList from concatenation of DrawList1 and 2
	drawList := DrawList{}
	drawList.AddList(&drawList1)
	drawList.AddList(&drawList2)

	// Render loop
	frames := 0
	cgoCallsStart := runtime.NumCgoCall()
	frameParams := FrameParams{}
	frameParams.ClearColor = Vec4{0.5, 0.5, 0.5, 1}
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
	win.DeleteTexture(texId)
	win.Destroy()

}
