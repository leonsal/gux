package main

import (
	"fmt"
	"math"
	"runtime"

	"github.com/leonsal/gux"
	"github.com/leonsal/gux/gb"
)

var colorList []gb.RGBA

func init() {
	colorList = append(colorList,
		gb.MakeColor(255, 0, 0, 255),
		gb.MakeColor(0, 255, 0, 255),
		gb.MakeColor(0, 0, 255, 255),
		gb.MakeColor(0, 0, 0, 255),
		gb.MakeColor(255, 255, 0, 255),
		gb.MakeColor(0, 255, 255, 255),
		gb.MakeColor(255, 255, 255, 255),
		gb.MakeColor(100, 100, 100, 255),
	)
}

func main() {

	runtime.LockOSThread()

	// Create window
	cfg := gb.Config{}
	cfg.DebugPrintCmds = false
	cfg.OpenGL.ES = false
	cfg.Vulkan.ValidationLayer = true
	win, err := gux.NewWindow("title", 2000, 1200, &cfg)
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
	win.SetCursor(gb.CursorIBeam)

	//test := newTestTransform(win)
	for win.StartFrame() {

		testRect(win)
		//testArc(win)
		//test.draw(win)
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

	win.DestroyFontAtlas(fa)
	win.DeleteTexture(texID)
	win.Destroy()
}

type gbtest interface {
	draw(*gux.Window)
}

func testRect(win *gux.Window) {

	dl := win.DrawList()
	startY := float32(50)
	const startX = 50
	const width = 300
	const height = 150
	const deltaX = width + 10
	const deltaY = height + 50
	const rounding = 40.0
	const thickness = 10
	flagList := []gux.DrawFlags{
		gux.DrawFlags_RoundCornersTopLeft,
		gux.DrawFlags_RoundCornersTopRight,
		gux.DrawFlags_RoundCornersBottomRight,
		gux.DrawFlags_RoundCornersBottomLeft,
		gux.DrawFlags_RoundCornersTop,
		gux.DrawFlags_RoundCornersBottom,
		gux.DrawFlags_RoundCornersLeft,
		gux.DrawFlags_RoundCornersRight,
		gux.DrawFlags_RoundCornersAll,
	}

	for idx, flag := range flagList {
		line := float32(idx / 5)
		col := float32(idx % 5)
		min := gb.Vec2{startX + col*deltaX, startY + line*deltaY}
		max := gb.Vec2{width + col*deltaX, startY + height + line*deltaY}
		win.AddRect(dl, min, max, nextColor(idx), rounding, flag, thickness)
	}

	startY += 2 * deltaY
	for idx, flag := range flagList {
		line := float32(idx / 5)
		col := float32(idx % 5)
		min := gb.Vec2{startX + col*deltaX, startY + line*deltaY}
		max := gb.Vec2{width + col*deltaX, startY + height + line*deltaY}
		win.AddRectFilled(dl, min, max, nextColor(idx), rounding, flag)
	}
}

func testArc(win *gux.Window) {

	dl := win.DrawList()
	colors := []gb.RGBA{
		gb.MakeColor(255, 0, 0, 255),
		gb.MakeColor(0, 255, 0, 255),
		gb.MakeColor(0, 0, 255, 255),
		gb.MakeColor(0, 0, 0, 255),
		gb.MakeColor(255, 255, 0, 255),
		gb.MakeColor(0, 255, 255, 255),
		gb.MakeColor(255, 255, 255, 255),
		gb.MakeColor(100, 100, 100, 255),
	}

	nextColor := func(i int) gb.RGBA {
		ci := i % len(colors)
		return colors[ci]
	}

	radius := float32(100)
	center := gb.Vec2{radius + 10, radius + 10}
	segs := 3
	thickness := float32(2)
	deltaY := 2*radius + 50
	deltaX := 2*radius + 40
	for i := 0; i < 12; i++ {
		win.AddCircle(dl, center, radius, nextColor(i), segs, thickness)
		center.X += deltaX
		segs += 3
		thickness += 2
	}

	center = gb.Vec2{radius + 10, center.Y + deltaY}
	segs = 3
	for i := 0; i < 12; i++ {
		win.AddCircleFilled(dl, center, radius, nextColor(i), segs)
		center.X += deltaX
		segs += 3
	}

	center = gb.Vec2{radius + 10, center.Y + deltaY}
	segs = 3
	thickness = float32(2)
	for i := 0; i < 12; i++ {
		amax := 2 * math.Pi * (float64(i) + 1) * 0.08
		win.PathArcTo(dl, center, radius, 0, float32(amax), segs)
		win.PathStroke(dl, nextColor(i), 0, thickness)
		center.X += deltaX
		segs += 3
		thickness += 2
	}

	center = gb.Vec2{radius + 10, center.Y + deltaY}
	segs = 3
	for i := 0; i < 12; i++ {
		amax := 2 * math.Pi * (float64(i) + 1) * 0.08
		win.PathArcTo(dl, center, radius, 0, float32(amax), segs)
		win.PathFillConvex(dl, nextColor(i))
		center.X += deltaX
		segs += 3
		thickness += 2
	}
}

type testTransform struct {
	g1 gb.DrawList
	g2 gb.DrawList
	g3 gb.DrawList
}

func newTestTransform(win *gux.Window) *testTransform {

	t := new(testTransform)

	red := gb.MakeColor(255, 0, 0, 255)
	green := gb.MakeColor(0, 255, 0, 255)
	blue := gb.MakeColor(0, 0, 255, 255)

	// First group
	_, bufIdx, bufVtx := win.NewDrawCmd(&t.g1, 3, 3)
	bufVtx[0] = gb.Vertex{Pos: gb.Vec2{0, -100}, Col: red}
	bufVtx[1] = gb.Vertex{Pos: gb.Vec2{-100, 100}, Col: red}
	bufVtx[2] = gb.Vertex{Pos: gb.Vec2{100, 100}, Col: red}
	bufIdx[0] = 0
	bufIdx[1] = 1
	bufIdx[2] = 2

	// Second group
	_, bufIdx, bufVtx = win.NewDrawCmd(&t.g2, 3, 3)
	bufVtx[0] = gb.Vertex{Pos: gb.Vec2{0, 0}, Col: green}
	bufVtx[1] = gb.Vertex{Pos: gb.Vec2{0, 200}, Col: green}
	bufVtx[2] = gb.Vertex{Pos: gb.Vec2{200, 0}, Col: green}
	bufIdx[0] = 0
	bufIdx[1] = 1
	bufIdx[2] = 2

	// Third group
	_, bufIdx, bufVtx = win.NewDrawCmd(&t.g3, 6, 4)
	bufVtx[0] = gb.Vertex{Pos: gb.Vec2{-100, -100}, Col: blue}
	bufVtx[1] = gb.Vertex{Pos: gb.Vec2{-100, 100}, Col: blue}
	bufVtx[2] = gb.Vertex{Pos: gb.Vec2{100, 100}, Col: blue}
	bufVtx[3] = gb.Vertex{Pos: gb.Vec2{100, -100}, Col: blue}
	bufIdx[0] = 0
	bufIdx[1] = 1
	bufIdx[2] = 2
	bufIdx[3] = 2
	bufIdx[4] = 3
	bufIdx[5] = 0

	return t
}

func (t *testTransform) draw(win *gux.Window) {

	dl := win.DrawList()
	var mat gb.Mat3

	deltaX := float32(210)
	for i := 1; i < 8; i++ {
		sf := 1.0 - float32(i)/10
		mat.SetTranslation(deltaX*float32(i), 100).Rotate(float32(i-1)*float32(math.Pi/16)).Scale(sf, sf)
		dl.AddList2(&t.g1, &mat)
	}

	deltaY := float32(300)
	for i := 1; i < 8; i++ {
		sf := 1.0 - float32(i)/10
		mat.SetTranslation(deltaX*float32(i), deltaY).Rotate(float32(i-1)*float32(math.Pi/16)).Scale(sf, sf)
		dl.AddList2(&t.g2, &mat)
	}

	deltaY += 300
	for i := 1; i < 8; i++ {
		sf := 1.0 - float32(i)/10
		mat.SetTranslation(deltaX*float32(i), deltaY).Rotate(float32(i-1)*float32(math.Pi/16)).Scale(sf, sf)
		dl.AddList2(&t.g3, &mat)
	}
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
		w.AddPolyLineAntiAliased(dl, points1, gb.MakeColor(0, 0, 0, 255), gux.DrawFlags_None, float32(width))
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

func nextColor(i int) gb.RGBA {

	ci := i % len(colorList)
	return colorList[ci]
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
