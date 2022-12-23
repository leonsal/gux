package main

import (
	"fmt"
	"runtime"

	"github.com/leonsal/gux"
	"github.com/leonsal/gux/gb"
)

func main() {

	runtime.LockOSThread()

	// Create window
	win, err := gux.NewWindow("title", 2000, 1200)
	if err != nil {
		panic(err)
	}

	// Create font
	f, err := gux.NewFont("assets/Roboto-Medium.ttf")
	if err != nil {
		panic(err)
	}
	f.SetFgColor(gb.MakeColor(255, 255, 0, 255))
	f.SetBgColor(gb.MakeColor(0, 0, 0, 100))
	f.SetPointSize(148)

	// Create font atlas
	fa := win.NewFontAtlas(f, 0x00, 0xFF)
	fmt.Println("ATLAS: LineHeight:", fa.LineHeight, "Ascent:", fa.Ascent, "Descent:", fa.Descent)
	if false {
		err = fa.SavePNG("atlas.png")
		if err != nil {
			fmt.Println("SAVE ERROR:", err)
		}
	}
	texID, width, height := win.CreateTextImage(f, "text image")
	fmt.Println("TextImage:", texID, width, height)

	//text := `~!@#$%^&*()_+-={}[]:;'",<.>/?
	//1234567890()
	//abcdefghijklmnopqrstuvxyz
	//1234567890()
	//ABCDEFGHIJKLMNjPQRSTUVXYZ
	//éú
	//`
	//events := make([]gb.Event, 256)
	// Render loop
	var cgoCallsStart int64
	var statsStart runtime.MemStats
	frameCount := 0

	for win.StartFrame() {

		testBasic(win)
		//testText(win, fa, texID, width, height)
		//testLines(win)
		//testPolygon(win)
		win.Render()

		// All the allocations should be done in the first frame
		frameCount++
		if frameCount == 1 {
			cgoCallsStart = runtime.NumCgoCall()
			runtime.ReadMemStats(&statsStart)
		}
	}

	// Calculates and shows allocations and cgo calls per frame
	cgoCalls := runtime.NumCgoCall() - cgoCallsStart
	cgoPerFrame := float64(cgoCalls) / float64(frameCount)
	var stats runtime.MemStats
	runtime.ReadMemStats(&stats)
	allocsPerFrame := float64(stats.Alloc-statsStart.Alloc) / float64(frameCount)
	fmt.Println("Frames:", frameCount, "Allocs per frame:", allocsPerFrame, "CGO calls per frame:", cgoPerFrame)

	win.Destroy()
}

func testBasic(win *gux.Window) {

	dl := win.DrawList()
	cmd, bufIdx, bufVtx := win.NewDrawCmd(dl, 3, 3)
	bufVtx[0] = gb.Vertex{Pos: gb.Vec2{200, 800}, Col: 0xFF_00_00_FF}
	bufVtx[1] = gb.Vertex{Pos: gb.Vec2{400, 800}, Col: 0xFF_00_FF_00}
	bufVtx[2] = gb.Vertex{Pos: gb.Vec2{300, 600}, Col: 0xFF_FF_00_00}
	bufIdx[0] = 0
	bufIdx[1] = 1
	bufIdx[2] = 2
	dl.AdjustIdx(cmd)

	//cmd, bufIdx, bufVtx = win.NewDrawCmd(dl, 3, 3)
	//bufVtx[0] = gb.Vertex{Pos: gb.Vec2{0, 0}, Col: 0xFF_00_00_FF}
	//bufVtx[1] = gb.Vertex{Pos: gb.Vec2{0, 300}, Col: 0xFF_00_FF_00}
	//bufVtx[2] = gb.Vertex{Pos: gb.Vec2{300, 0}, Col: 0xFF_FF_00_00}
	//bufIdx[0] = 0
	//bufIdx[1] = 1
	//bufIdx[2] = 2
	//dl.AdjustIdx(cmd)
}

func testText(win *gux.Window, fa *gux.FontAtlas, texID gb.TextureID, width, height float32) {

	win.AddText(win.DrawList(), fa, gb.Vec2{50, 200}, gux.TextVAlignTop, "top ")
	win.AddText(win.DrawList(), fa, gb.Vec2{250, 200}, gux.TextVAlignBase, " base")
	win.AddText(win.DrawList(), fa, gb.Vec2{550, 200}, gux.TextVAlignBottom, " bottom")
	win.AddImage(win.DrawList(), texID, width, height, gb.Vec2{50, 400})

}

func testLines(w *gux.Window) {

	// Line points
	points := []gb.Vec2{{0, 10}, {10, 0}, {20, 10}, {30, 0}, {40, 10}, {50, 0}, {60, 10}}
	points1 := w.ReserveVec2(len(points))
	points2 := w.ReserveVec2(len(points))

	// Add poly lines anti aliased
	copy(points1, points)
	scalePoints(points1, 12)
	translatePoints(points1, gb.Vec2{10, 10})
	dl := w.DrawList()
	for width := 1; width < 60; width += 8 {
		w.AddPolyLineAntiAliased(dl, points1, gb.MakeColor(0, 0, 0, 255), 0, float32(width))
		translatePoints(points1, gb.Vec2{0, 120})
	}

	// Add poly lines textured
	copy(points2, points)
	scalePoints(points2, 12)
	translatePoints(points2, gb.Vec2{800, 10})
	for width := 1; width < 60; width += 8 {
		w.AddPolyLineTextured(dl, points2, gb.MakeColor(0, 0, 0, 255), 0, float32(width))
		translatePoints(points2, gb.Vec2{0, 120})
	}
}

func testPolygon(w *gux.Window) {

	dl := w.DrawList()

	triangle := []gb.Vec2{{0, 100}, {100, 100}, {50, 0}}
	points := w.ReserveVec2(len(triangle))
	copy(points, triangle)
	scalePoints(points, 2)
	w.AddConvexPolyFilled(dl, points, gb.MakeColor(0, 0, 0, 255))

	rect := []gb.Vec2{{0, 100}, {200, 100}, {200, 0}, {0, 0}}
	points = w.ReserveVec2(len(rect))
	copy(points, rect)
	translatePoints(points, gb.Vec2{220, 10})
	w.AddConvexPolyFilled(dl, points, gb.MakeColor(255, 0, 0, 255))

	points = w.ReserveVec2(len(triangle))
	copy(points, triangle)
	scalePoints(points, 2)
	translatePoints(points, gb.Vec2{0, 300})
	w.AddConvexPolyFilled(dl, points, gb.MakeColor(0, 255, 0, 255))

	points = w.ReserveVec2(len(rect))
	copy(points, rect)
	translatePoints(points, gb.Vec2{300, 300})
	w.AddConvexPolyFilled(dl, points, gb.MakeColor(0, 255, 255, 255))
}

// scale the supplied array of points
func scalePoints(points []gb.Vec2, scale float32) {
	for i := range points {
		(&points[i]).MultScalar(scale)
	}
}

// translate the supplied array of points
func translatePoints(points []gb.Vec2, trans gb.Vec2) {

	for i := range points {
		(&points[i]).Add(trans)
	}
}
