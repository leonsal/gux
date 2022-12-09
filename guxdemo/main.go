package main

import (
	"fmt"
	"runtime"
	"unsafe"

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
	fmt.Printf("font:%+v\n", f)
	f.SetFgColor(gb.MakeColor(255, 255, 0, 255))
	f.SetPointSize(180)

	// Create image with font
	img := f.DrawText("Hello gjpq\nLine2\nLine3")
	b := img.Bounds()
	width := b.Dx()
	height := b.Dy()

	// Creates backend texture to store the image and transfer the image
	texID := win.CreateTexture()
	win.TransferTexture(texID, width, height, (*gb.RGBA)(unsafe.Pointer(&img.Pix[0])))
	fmt.Println("width: ", width, "height:", height, "texID", texID)

	//events := make([]gb.Event, 256)
	// Render loop
	for win.StartFrame(0) {
		testText(win, texID, width, height)
		//count := win.GetEvents(events)
		//fmt.Println("events:", count)
		//testLines(win)
		//testPolygon(win)
		win.Render()
	}
	win.Destroy()
}

func testText(w *gux.Window, texID gb.TextureId, width, height int) {

	dl := w.DrawList()
	cmd, bufIdx, bufVtx := dl.ReserveCmd(6, 4)
	cmd.TexId = texID
	bufVtx[0].Pos = gb.Vec2{0, 0}
	bufVtx[0].UV = gb.Vec2{0, 0}
	bufVtx[0].Col = gb.MakeColor(255, 255, 255, 255)

	bufVtx[1].Pos = gb.Vec2{0, float32(height)}
	bufVtx[1].UV = gb.Vec2{0, 1}
	bufVtx[1].Col = gb.MakeColor(255, 255, 255, 255)

	bufVtx[2].Pos = gb.Vec2{float32(width), float32(height)}
	bufVtx[2].UV = gb.Vec2{1, 1}
	bufVtx[2].Col = gb.MakeColor(255, 255, 255, 255)

	bufVtx[3].Pos = gb.Vec2{float32(width), 0}
	bufVtx[3].UV = gb.Vec2{1, 0}
	bufVtx[3].Col = gb.MakeColor(255, 255, 255, 255)

	bufIdx[0] = 0
	bufIdx[1] = 1
	bufIdx[2] = 2
	bufIdx[3] = 2
	bufIdx[4] = 3
	bufIdx[5] = 0
	dl.AdjustIdx(cmd)
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
