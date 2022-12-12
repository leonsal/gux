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
	f.SetFgColor(gb.MakeColor(255, 255, 0, 255))
	f.SetBgColor(gb.MakeColor(0, 0, 0, 100))
	f.SetPointSize(148)
	//for x := rune(0); x <= rune(255); x++ {
	//	fmt.Printf("%v -> %v \n", x, f.Index(x))
	//}

	// Create atlas
	//fa := gux.NewFontAtlas(f, 0, 0xff)
	fa := gux.NewFontAtlas(f, 0x00, 0xFF)
	err = fa.SavePNG("atlas.png")
	if err != nil {
		fmt.Println("SAVE ERROR:", err)
	}

	texId, _, _ := createAtlasTexture(win, fa)

	//	// Create Texture with text
	//	texID, width, height := createText(win, f, `Hello Text: 01234567890
	//abcdefghijklmnopqrstuvwxyz
	//ABCDEFGHIJKLMNOPQRSTUVWXYZ`)

	text := `~!@#$%^&*()_+-={}[]:;'",<.>/?
	1234567890()
	abcdefghijklmnopqrstuvxyz
	1234567890()
	ABCDEFGHIJKLMNjPQRSTUVXYZ
	éú
	`
	//events := make([]gb.Event, 256)
	// Render loop
	var cgoCallsStart int64
	var statsStart runtime.MemStats
	frameCount := 0

	for win.StartFrame(0) {
		//testAtlas(win, fa, texId, "$1AQap")
		testAtlas(win, fa, texId, text)
		//testText(win, texID, width, height)
		//count := win.GetEvents(events)
		//fmt.Println("events:", count)
		//testLines(win)
		//testPolygon(win)
		win.Render()

		// All the allocations are done in the first frame
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

func testAtlas(w *gux.Window, fa *gux.FontAtlas, texId gb.TextureId, text string) {

	white := gb.MakeColor(255, 255, 255, 255)
	posX := float32(0)
	posY := float32(0)

	// For each rune in the text
	for _, c := range text {

		// Process new line
		if c == 0x0A {
			posX = float32(0)
			posY += float32(fa.LineHeight - 1)
			continue
		}

		// Ignore codes with no glyphs
		charInfo, ok := fa.Chars[c]
		if !ok {
			continue
		}

		//fmt.Printf("char: %v Info:%+v\n", c, charInfo)
		dl := w.DrawList()
		cmd, bufIdx, bufVtx := dl.ReserveCmd(6, 4)
		cmd.TexId = texId
		bufVtx[0].Pos = gb.Vec2{posX, posY}
		bufVtx[0].UV = charInfo.UV[0]
		bufVtx[0].Col = white

		bufVtx[1].Pos = gb.Vec2{posX, posY + float32(charInfo.Height-1)}
		bufVtx[1].UV = charInfo.UV[1]
		bufVtx[1].Col = white

		bufVtx[2].Pos = gb.Vec2{posX + float32(charInfo.Width-1), posY + float32(charInfo.Height-1)}
		bufVtx[2].UV = charInfo.UV[2]
		bufVtx[2].Col = white

		bufVtx[3].Pos = gb.Vec2{posX + float32(charInfo.Width-1), posY}
		bufVtx[3].UV = charInfo.UV[3]
		bufVtx[3].Col = white

		bufIdx[0] = 0
		bufIdx[1] = 1
		bufIdx[2] = 2
		bufIdx[3] = 2
		bufIdx[4] = 3
		bufIdx[5] = 0
		dl.AdjustIdx(cmd)
		posX += float32(charInfo.Width - 1)
	}

}
func createAtlasTexture(win *gux.Window, fa *gux.FontAtlas) (gb.TextureId, float32, float32) {

	img := fa.Image
	b := img.Bounds()
	width := b.Dx()
	height := b.Dy()

	// Creates backend texture to store the image and transfer the image
	texID := win.CreateTexture()
	win.TransferTexture(texID, width, height, (*gb.RGBA)(unsafe.Pointer(&img.Pix[0])))
	return texID, float32(width), float32(height)
}

func createText(win *gux.Window, f *gux.Font, text string) (gb.TextureId, float32, float32) {

	// Create image and draw text on it
	img := f.DrawText(text)
	b := img.Bounds()
	width := b.Dx()
	height := b.Dy()

	// Creates backend texture to store the image and transfer the image
	texID := win.CreateTexture()
	win.TransferTexture(texID, width, height, (*gb.RGBA)(unsafe.Pointer(&img.Pix[0])))
	return texID, float32(width), float32(height)
}

func testText(w *gux.Window, texID gb.TextureId, width, height float32) {
	//
	//    OpenGL UV coordinates adjustment
	//
	//	  0,1    1,1      0,0    1,0
	// 0 +------+ 3       +------+
	//	 |\     |         |\     |
	//	 | \    |         | \    |
	//	 |  \   |  --->   |  \   |
	//	 |   \  |         |   \  |
	//	 |    \ |         |    \ |
	//	 |     \|         |     \|
	// 1 +------+ 2       +------+
	//	 0,0    1,0       0,1    1,1

	dl := w.DrawList()
	cmd, bufIdx, bufVtx := dl.ReserveCmd(6, 4)
	cmd.TexId = texID
	white := gb.MakeColor(255, 255, 255, 255)
	bufVtx[0].Pos = gb.Vec2{0, 0}
	bufVtx[0].UV = gb.Vec2{0, 0}
	bufVtx[0].Col = white

	bufVtx[1].Pos = gb.Vec2{0, height - 1}
	bufVtx[1].UV = gb.Vec2{0, 1}
	bufVtx[1].Col = white

	bufVtx[2].Pos = gb.Vec2{width - 1, height - 1}
	bufVtx[2].UV = gb.Vec2{1, 1}
	bufVtx[2].Col = white

	bufVtx[3].Pos = gb.Vec2{width - 1, 0}
	bufVtx[3].UV = gb.Vec2{1, 0}
	bufVtx[3].Col = white

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
