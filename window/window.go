package window

import (
	"github.com/leonsal/gux/gb"
)

const TexLinesWidthMax = 63

type DrawFlags int

const (
	DrawFlags_None   DrawFlags = 0
	DrawFlags_Closed DrawFlags = 1 << iota
	DrawFlags_RoundCornersTopLeft
	DrawFlags_RoundCornersTopRight
	DrawFlags_RoundCornersBottomLeft
	DrawFlags_RoundCornersBottomRight
	DrawFlags_RoundCornersNone
	DrawFlags_RoundCornersTop    = DrawFlags_RoundCornersTopLeft | DrawFlags_RoundCornersTopRight
	DrawFlags_RoundCornersBottom = DrawFlags_RoundCornersBottomLeft | DrawFlags_RoundCornersBottomRight
	DrawFlags_RoundCornersLeft   = DrawFlags_RoundCornersBottomLeft | DrawFlags_RoundCornersTopLeft
	DrawFlags_RoundCornersRight  = DrawFlags_RoundCornersBottomRight | DrawFlags_RoundCornersTopRight
	DrawFlags_RoundCornersAll    = DrawFlags_RoundCornersTopLeft | DrawFlags_RoundCornersTopRight | DrawFlags_RoundCornersBottomLeft | DrawFlags_RoundCornersBottomRight
	DrawFlags_RoundCornersMask_  = DrawFlags_RoundCornersAll | DrawFlags_RoundCornersNone
)

type DrawListFlags int

const (
	DrawListFlags_AntiAliasedFill DrawListFlags = 1 << iota
)
const (
	DrawListCircleSegmentMax = 512
)

// Window corresponds to a native platform Window
type Window struct {
	gbw                  *gb.Window                    // Graphics backend native window reference
	dl                   gb.DrawList                   // Draw list to render
	fm                   *FontManager                  // Current FontManager
	TexWhiteId           gb.TextureID                  // Texture with white opaque pixel
	TexLinesId           gb.TextureID                  // Texture for lines
	TexUvLines           [TexLinesWidthMax + 1]gb.Vec4 // UV coordinates for textured lines
	FringeScale          float32                       // Used for AA
	bufVec2              []gb.Vec2                     // Temporary Vec2 buffer used by drawing functions (to avoid allocations)
	drawFlags            DrawListFlags                 // Flags, you may poke into these to adjust anti-aliasing settings per-primitive.
	frameParams          gb.FrameParams
	frameInfo            gb.FrameInfo
	CurveTessellationTol float32 // IN STYLES ? Tessellation tolerance when using PathBezierCurveTo() without a specific number of segments. Decrease for highly tessellated curves (higher quality, more polygons), increase to reduce quality.
	clipRect             gb.Rect // Current clip rectangle for Draw Commands
}

// New creates and returns a new Window
func New(title string, width, height int, cfg *gb.Config) (*Window, error) {

	// Creates graphics backend native window
	w := new(Window)
	var err error
	w.gbw, err = gb.CreateWindow(title, width, height, cfg)
	if err != nil {
		return nil, err
	}

	// Create textures
	w.buildTexWhite()
	w.buildTexLines()

	w.drawFlags |= DrawListFlags_AntiAliasedFill
	w.FringeScale = 1.0

	//w.frameParams.ClearColor = gb.Vec4{0.5, 0.5, 0.5, 1.0}
	w.frameParams.ClearColor = gb.Vec4{1.0, 1.0, 1.0, 1.0}
	w.frameInfo.WinSize = gb.Vec2{float32(width), float32(height)}
	w.CurveTessellationTol = 1.25
	return w, nil
}

// Size returns the window content area size
func (w *Window) Size() gb.Vec2 {

	return w.frameInfo.WinSize
}

func (w *Window) SetFontManager(fm *FontManager) {

	w.fm = fm
}

func (w *Window) FontManager() *FontManager {

	return w.fm
}

func (w *Window) Font(ff FontFamilyType, relSize int) *FontAtlas {

	return w.fm.Font(ff, relSize)
}

func (w *Window) ClearClipRect() {

	w.clipRect = gb.Rect{gb.Vec2{0, 0}, gb.Vec2{w.frameInfo.WinSize.X, w.frameInfo.WinSize.Y}}
}

func (w *Window) SetClipRect(r gb.Rect) {

	w.clipRect = r
}

func (w *Window) SetClearColor(color gb.Vec4) {

	w.frameParams.ClearColor = color
}

func (w *Window) SetEvTimeout(timeout float32) {

	w.frameParams.EvTimeout = timeout
}

// StartFrame sets the beginning of a new render frame and returns true if
// the window should be closed or false otherwise.
func (w *Window) StartFrame() bool {

	w.dl.Clear()
	w.bufVec2 = w.bufVec2[:0]
	w.frameInfo = w.gbw.StartFrame(&w.frameParams)
	w.ClearClipRect()
	return w.frameInfo.WinClose
}

// RenderFrame sends this Windows' DrawList to the Graphics Backend for rendering
func (w *Window) RenderFrame() {

	w.gbw.RenderFrame(&w.dl)
}

// Adds specified draw list to this Window's draw list
func (w *Window) AddList(src *gb.DrawList) {

	w.dl.AddList(src)
}

func (w *Window) Destroy() {

	if w.fm != nil {
		w.fm.DestroyFonts(w)
	}
	w.DeleteTexture(w.TexWhiteId)
	w.DeleteTexture(w.TexLinesId)
	w.gbw.Destroy()
}

func (w *Window) SetCursor(cursor gb.Cursor) {

	w.gbw.SetCursor(cursor)
}

func (w *Window) FrameInfo() *gb.FrameInfo {

	return &w.frameInfo
}

func (w *Window) CreateTexture(width, height int, data *gb.RGBA) gb.TextureID {

	return w.gbw.CreateTexture(width, height, data)
}

func (w *Window) DeleteTexture(texid gb.TextureID) {

	w.gbw.DeleteTexture(texid)
}

// ReserveVec2 reserves 'count' gb.Vec2 entries in internal Vec2 buffer
// returning a slice to access these entries
func (w *Window) ReserveVec2(count int) []gb.Vec2 {

	idx := len(w.bufVec2)
	for i := 0; i < count; i++ {
		w.bufVec2 = append(w.bufVec2, gb.Vec2{})
	}
	return w.bufVec2[idx : idx+count]
}

// buildTexWhite generates a 1x1 texture with one white opaque pixel.
// It is used as a default texture for commands which don't use a texture.
func (w *Window) buildTexWhite() {

	// Creates image with one white opaque pixel
	var rect [1]gb.RGBA
	rect[0] = gb.MakeColor(255, 255, 255, 255)

	// Creates and transfer texture
	w.TexWhiteId = w.gbw.CreateTexture(1, 1, &rect[0])
	//fmt.Println("texWhiteId", w.TexWhiteId)
}

// buildTexLines generates a texture with a triangular shape with various line widths
// stacked on top of each other to allow interpolation between them.
func (w *Window) buildTexLines() {

	/*
		Example for TexLinesWidthMax = 9
		T - transparent white
		O - opaque white

		Line Width		Texels
		0				TTTTTTTTTTT
		1               TTTTTOTTTTT
		2               TTTTOOTTTTT
		3               TTTTOOOTTTT
		4               TTTOOOOOTTT
		5               TTTOOOOOTTT
		6               TTOOOOOOTTT
		7               TTOOOOOOOTT
		8               TOOOOOOOOTT
		9               TOOOOOOOOOT
		...
	*/

	width := TexLinesWidthMax + 2
	height := TexLinesWidthMax + 1
	rect := make([]gb.RGBA, width*height)
	uvScale := gb.Vec2{1 / float32(width), 1 / float32(height)}
	for n := 0; n < height; n++ {

		// Each line consists of at least one transparent pixel at each side, with a line of solid pixels in the middle
		lineWidth := n
		padLeft := (width - lineWidth) / 2
		padRight := width - (padLeft + lineWidth)
		pos := n * width

		for i := 0; i < padLeft; i++ {
			rect[pos+i] = gb.MakeColor(255, 255, 255, 0)
		}
		for i := 0; i < lineWidth; i++ {
			rect[pos+padLeft+i] = gb.MakeColor(255, 255, 255, 255)
		}
		for i := 0; i < padRight; i++ {
			rect[pos+padLeft+lineWidth+i] = gb.MakeColor(255, 255, 255, 0)
		}

		// Calculate UVs for this line
		uv0 := gb.Vec2Mult(gb.Vec2{float32(padLeft - 1), float32(n)}, uvScale)
		uv1 := gb.Vec2Mult(gb.Vec2{float32(padLeft + lineWidth + 1), float32(n + 1)}, uvScale)
		halfV := (uv0.Y + uv1.Y) * 0.5 // Calculate a constant V in the middle of the row to avoid sampling artifacts
		w.TexUvLines[n] = gb.Vec4{uv0.X, halfV, uv1.X, halfV}
	}

	// Creates and transfer texture
	w.TexLinesId = w.gbw.CreateTexture(width, height, &rect[0])
	//fmt.Println("texture id", w.TexLinesId)

	// // Print image data
	//
	//	for n := 0; n < height; n++ {
	//		pos := n * width
	//		for c := 0; c < width; c++ {
	//			fmt.Printf("%d ", rect[pos+c])
	//		}
	//		fmt.Println()
	//	}
	//
	// // Print UVs
	//
	//	for n := 0; n < height; n++ {
	//		fmt.Println(w.TexUvLines[n])
	//	}
}
