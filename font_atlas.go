package gux

import (
	"bufio"
	"fmt"
	"image"
	"image/draw"
	"image/png"
	"io"
	"io/ioutil"
	"os"
	"sort"
	"unicode"
	"unsafe"

	"github.com/leonsal/gux/gb"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
	"golang.org/x/image/math/fixed"
)

type GlyphInfo struct {
	Advance float32    // Amount to add to glyph origin to draw next Glyph
	Bounds  gb.Rect    // Glyph bounds relative to its origin point
	UV      [4]gb.Vec2 // UV coordinates for glyph quad vertices
}

// FontAtlas represents an image containing characters and the information about their location in the image
type FontAtlas struct {
	Face       font.Face          // The font face used to generate the atlas
	Glyphs     map[rune]GlyphInfo // Maps rune code to correspondent Glyph info
	Image      *image.RGBA        // Font atlas generated image
	Ascent     float32            // Distance from the top of a line to its baseline
	Descent    float32            // Distance from the bottom of a line to its baseline
	LineHeight float32            // Total line height
	TexID      gb.TextureID       // Texture ID (valid only after texture was created)
}

// NewFontFaceFromFile creates and returns a Font face with the specified options from
// parsing the specified TrueType or OpenType font file.
func NewFontFaceFromFile(filepath string, opts *opentype.FaceOptions) (font.Face, error) {

	f, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return NewFontFaceFromReader(f, opts)
}

func NewFontFaceFromReader(r io.Reader, opts *opentype.FaceOptions) (font.Face, error) {

	fdata, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}
	return NewFontFace(fdata, opts)
}

func NewFontFace(fontData []byte, opts *opentype.FaceOptions) (font.Face, error) {

	fparsed, err := opentype.Parse(fontData)
	if err != nil {
		return nil, nil
	}
	return opentype.NewFace(fparsed, opts)
}

func NewFontAtlas(face font.Face, runeSets ...[]rune) *FontAtlas {

	// Builds array of unique runes from all the specified rune sets
	seen := make(map[rune]bool)
	runes := []rune{unicode.ReplacementChar}
	//runes := []rune{}
	for _, set := range runeSets {
		for _, r := range set {
			if !seen[r] {
				runes = append(runes, r)
				seen[r] = true
			}
		}
	}

	fixedMapping, fixedBounds := makeSquareMapping(face, runes, fixed.I(2))

	// Creates font atlas image
	atlasImg := image.NewRGBA(image.Rect(
		fixedBounds.Min.X.Floor(),
		fixedBounds.Min.Y.Floor(),
		fixedBounds.Max.X.Ceil(),
		fixedBounds.Max.Y.Ceil(),
	))

	// Draw all glyphs to the font atlas image
	for r, fg := range fixedMapping {
		if dr, mask, maskp, _, ok := face.Glyph(fg.dot, r); ok {
			draw.Draw(atlasImg, dr, mask, maskp, draw.Src)
		}
	}

	// Image bounds
	boundsMinX := i2f(fixedBounds.Min.X)
	boundsMaxX := i2f(fixedBounds.Max.X)
	boundsMinY := i2f(fixedBounds.Min.Y)
	boundsMaxY := i2f(fixedBounds.Max.Y)
	imageWidth := boundsMaxX - boundsMinX
	imageHeight := boundsMaxY - boundsMinY

	fmt.Println("image:", boundsMinX, boundsMaxX, boundsMinY, boundsMaxY, imageWidth, imageHeight)
	fmt.Printf("image rect:%+v\n", atlasImg.Rect)

	glyphs := make(map[rune]GlyphInfo)
	for r, fg := range fixedMapping {
		gi := GlyphInfo{}

		// Get Glyph bounds and advance converting from fixed to float
		bounds, advance, _ := face.GlyphBounds(r)
		gi.Advance = i2f(advance)
		gi.Bounds.Min = gb.Vec2{i2f(bounds.Min.X), i2f(bounds.Min.Y)}
		gi.Bounds.Max = gb.Vec2{i2f(bounds.Max.X), i2f(bounds.Max.Y)}

		// Transform glyphs image coordinates to UV coordinates
		minX := i2f(fg.frame.Min.X)
		minY := -boundsMinY + i2f(fg.frame.Min.Y)
		maxX := i2f(fg.frame.Max.X)
		maxY := -boundsMinY + i2f(fg.frame.Max.Y)
		//fmt.Printf("code:%v minX:%f minY:%f maxX:%f maxY:%f width:%f\n", r, minX, minY, maxX, maxY, maxX-minX)
		fmt.Printf("code:%v minX:%f minY:%f maxX:%f maxY:%f\n", r,
			i2f(fg.frame.Min.X), i2f(fg.frame.Min.Y), i2f(fg.frame.Max.X), i2f(fg.frame.Max.Y))
		gi.UV[0] = gb.Vec2{minX / imageWidth, minY / imageHeight}
		gi.UV[1] = gb.Vec2{minX / imageWidth, maxY / imageHeight}
		gi.UV[2] = gb.Vec2{maxX / imageWidth, maxY / imageHeight}
		gi.UV[3] = gb.Vec2{maxX / imageWidth, minY / imageHeight}
		glyphs[r] = gi
		//glyphs[r] = GlyphInfo{
		//	Dot: gb.Vec2{i2f(fg.dot.X), boundsMaxY - (i2f(fg.dot.Y) - boundsMinY)},
		//	//char.UV[0] = gb.Vec2{float32(char.X) / imgWidth, (float32(char.Y) / imgHeight)}
		//	//char.UV[1] = gb.Vec2{float32(char.X) / imgWidth, (float32(char.Y+char.Height) / imgHeight)}
		//	//char.UV[2] = gb.Vec2{float32(char.X+char.Width) / imgWidth, (float32(char.Y+char.Height) / imgHeight)}
		//	//char.UV[3] = gb.Vec2{float32(char.X+char.Width) / imgWidth, (float32(char.Y) / imgHeight)}
		//	//			Frame: pixel.R(
		//	//				i2f(fg.frame.Min.X),
		//	//				bounds.Max.Y-(i2f(fg.frame.Min.Y)-bounds.Min.Y),
		//	//				i2f(fg.frame.Max.X),
		//	//				bounds.Max.Y-(i2f(fg.frame.Max.Y)-bounds.Min.Y),
		//	//			).Norm(),
		//	Advance: i2f(fg.advance),
		//}
	}

	return &FontAtlas{
		Face:       face,
		Glyphs:     glyphs,
		Image:      atlasImg,
		Ascent:     i2f(face.Metrics().Ascent),
		Descent:    i2f(face.Metrics().Descent),
		LineHeight: i2f(face.Metrics().Height),
	}
}

// SavePNG saves the current atlas image as a PNG image file
func (a *FontAtlas) SavePNG(filename string) error {

	// Save that RGBA image to disk.
	outFile, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer outFile.Close()

	b := bufio.NewWriter(outFile)
	err = png.Encode(b, a.Image)
	if err != nil {
		return err
	}
	err = b.Flush()
	if err != nil {
		return err
	}
	return nil
}

func (a *FontAtlas) CreateTexture(win *Window) error {

	rect := a.Image.Rect
	width := rect.Max.X - rect.Min.X
	height := rect.Max.Y - rect.Min.Y
	fmt.Println("CreateTexture", width, height)
	a.TexID = win.CreateTexture(width, height, (*gb.RGBA)(unsafe.Pointer(&a.Image.Pix[0])))
	return nil
}

func (a *FontAtlas) Destroy(win *Window) error {

	if a.TexID != 0 {
		win.DeleteTexture(a.TexID)
		a.TexID = 0
	}
	return nil
}

func (a *FontAtlas) PrintInfo() {

	fmt.Println("Ascent:", a.Ascent, "Descent", a.Descent, "Lineheight:", a.LineHeight)
	runes := make([]rune, 0, len(a.Glyphs))
	for r := range a.Glyphs {
		runes = append(runes, r)
	}
	sort.Slice(runes, func(i, j int) bool {
		return runes[i] < runes[j]
	})
	for _, r := range runes {
		fmt.Printf("code:%v glyph:[%c] info:%+v\n", r, r, a.Glyphs[r])
	}
}

type fixedGlyph struct {
	dot     fixed.Point26_6
	frame   fixed.Rectangle26_6
	advance fixed.Int26_6
}

// makeSquareMapping finds an optimal glyph arrangement of the given runes, so that their common
// bounding box is as square as possible.
func makeSquareMapping(face font.Face, runes []rune, padding fixed.Int26_6) (map[rune]fixedGlyph, fixed.Rectangle26_6) {

	width := sort.Search(int(fixed.I(1024*1024)), func(i int) bool {
		width := fixed.Int26_6(i)
		_, bounds := makeMapping(face, runes, padding, width)
		return bounds.Max.X-bounds.Min.X >= bounds.Max.Y-bounds.Min.Y
	})
	return makeMapping(face, runes, padding, fixed.Int26_6(width))
}

// makeMapping arranges glyphs of the given runes into rows in such a way, that no glyph is located
// fully to the right of the specified width. Specifically, it places glyphs in a row one by one and
// once it reaches the specified width, it starts a new row.
func makeMapping(face font.Face, runes []rune, padding, width fixed.Int26_6) (map[rune]fixedGlyph, fixed.Rectangle26_6) {

	mapping := make(map[rune]fixedGlyph)
	bounds := fixed.Rectangle26_6{}

	dot := fixed.P(0, 0)

	for _, r := range runes {
		b, advance, ok := face.GlyphBounds(r)
		if !ok {
			continue
		}

		// this is important for drawing, artifacts arise otherwise
		frame := fixed.Rectangle26_6{
			Min: fixed.P(b.Min.X.Floor(), b.Min.Y.Floor()),
			Max: fixed.P(b.Max.X.Ceil(), b.Max.Y.Ceil()),
		}

		dot.X -= frame.Min.X
		frame = frame.Add(dot)

		mapping[r] = fixedGlyph{
			dot:     dot,
			frame:   frame,
			advance: advance,
		}
		bounds = bounds.Union(frame)

		dot.X = frame.Max.X

		// padding + align to integer
		dot.X += padding
		dot.X = fixed.I(dot.X.Ceil())

		// width exceeded, new row
		if frame.Max.X >= width {
			dot.X = 0
			dot.Y += face.Metrics().Ascent + face.Metrics().Descent

			// padding + align to integer
			dot.Y += padding
			dot.Y = fixed.I(dot.Y.Ceil())
		}
	}

	return mapping, bounds
}

func i2f(i fixed.Int26_6) float32 {
	return float32(i.Floor())
}

// ---------------------------------------------------------------------------------------------
// package text
//
// import (
// 	"image"
// 	"image/draw"
// 	"sort"
// 	"unicode"
//
// 	"github.com/faiface/pixel"
// 	"golang.org/x/image/font"
// 	"golang.org/x/image/math/fixed"
// )
//
// // Atlas7x13 is an Atlas using basicfont.Face7x13 with the ASCII rune set
// var Atlas7x13 *Atlas
//
// // Glyph describes one glyph in an Atlas.
// type Glyph struct {
// 	Dot     pixel.Vec
// 	Frame   pixel.Rect
// 	Advance float64
// }
//
// // Atlas is a set of pre-drawn glyphs of a fixed set of runes. This allows for efficient text drawing.
// type Atlas struct {
// 	face       font.Face
// 	pic        pixel.Picture
// 	mapping    map[rune]Glyph
// 	ascent     float64
// 	descent    float64
// 	lineHeight float64
// }
//
// // NewAtlas creates a new Atlas containing glyphs of the union of the given sets of runes (plus
// // unicode.ReplacementChar) from the given font face.
// //
// // Creating an Atlas is rather expensive, do not create a new Atlas each frame.
// //
// // Do not destroy or close the font.Face after creating the Atlas. Atlas still uses it.
// func NewAtlas(face font.Face, runeSets ...[]rune) *Atlas {
// 	seen := make(map[rune]bool)
// 	runes := []rune{unicode.ReplacementChar}
// 	for _, set := range runeSets {
// 		for _, r := range set {
// 			if !seen[r] {
// 				runes = append(runes, r)
// 				seen[r] = true
// 			}
// 		}
// 	}
//
// 	fixedMapping, fixedBounds := makeSquareMapping(face, runes, fixed.I(2))
//
// 	atlasImg := image.NewRGBA(image.Rect(
// 		fixedBounds.Min.X.Floor(),
// 		fixedBounds.Min.Y.Floor(),
// 		fixedBounds.Max.X.Ceil(),
// 		fixedBounds.Max.Y.Ceil(),
// 	))
//
// 	for r, fg := range fixedMapping {
// 		if dr, mask, maskp, _, ok := face.Glyph(fg.dot, r); ok {
// 			draw.Draw(atlasImg, dr, mask, maskp, draw.Src)
// 		}
// 	}
//
// 	bounds := pixel.R(
// 		i2f(fixedBounds.Min.X),
// 		i2f(fixedBounds.Min.Y),
// 		i2f(fixedBounds.Max.X),
// 		i2f(fixedBounds.Max.Y),
// 	)
//
// 	mapping := make(map[rune]Glyph)
// 	for r, fg := range fixedMapping {
// 		mapping[r] = Glyph{
// 			Dot: pixel.V(
// 				i2f(fg.dot.X),
// 				bounds.Max.Y-(i2f(fg.dot.Y)-bounds.Min.Y),
// 			),
// 			Frame: pixel.R(
// 				i2f(fg.frame.Min.X),
// 				bounds.Max.Y-(i2f(fg.frame.Min.Y)-bounds.Min.Y),
// 				i2f(fg.frame.Max.X),
// 				bounds.Max.Y-(i2f(fg.frame.Max.Y)-bounds.Min.Y),
// 			).Norm(),
// 			Advance: i2f(fg.advance),
// 		}
// 	}
//
// 	return &Atlas{
// 		face:       face,
// 		pic:        pixel.PictureDataFromImage(atlasImg),
// 		mapping:    mapping,
// 		ascent:     i2f(face.Metrics().Ascent),
// 		descent:    i2f(face.Metrics().Descent),
// 		lineHeight: i2f(face.Metrics().Height),
// 	}
// }
//
// // Picture returns the underlying Picture containing an arrangement of all the glyphs contained
// // within the Atlas.
// func (a *Atlas) Picture() pixel.Picture {
// 	return a.pic
// }
//
// // Contains reports wheter r in contained within the Atlas.
// func (a *Atlas) Contains(r rune) bool {
// 	_, ok := a.mapping[r]
// 	return ok
// }
//
// // Glyph returns the description of r within the Atlas.
// func (a *Atlas) Glyph(r rune) Glyph {
// 	return a.mapping[r]
// }
//
// // Kern returns the kerning distance between runes r0 and r1. Positive distance means that the
// // glyphs should be further apart.
// func (a *Atlas) Kern(r0, r1 rune) float64 {
// 	return i2f(a.face.Kern(r0, r1))
// }
//
// // Ascent returns the distance from the top of the line to the baseline.
// func (a *Atlas) Ascent() float64 {
// 	return a.ascent
// }
//
// // Descent returns the distance from the baseline to the bottom of the line.
// func (a *Atlas) Descent() float64 {
// 	return a.descent
// }
//
// // LineHeight returns the recommended vertical distance between two lines of text.
// func (a *Atlas) LineHeight() float64 {
// 	return a.lineHeight
// }
//
// // DrawRune returns parameters necessary for drawing a rune glyph.
// //
// // Rect is a rectangle where the glyph should be positioned. Frame is the glyph frame inside the
// // Atlas's Picture. NewDot is the new position of the dot.
// func (a *Atlas) DrawRune(prevR, r rune, dot pixel.Vec) (rect, frame, bounds pixel.Rect, newDot pixel.Vec) {
// 	if !a.Contains(r) {
// 		r = unicode.ReplacementChar
// 	}
// 	if !a.Contains(unicode.ReplacementChar) {
// 		return pixel.Rect{}, pixel.Rect{}, pixel.Rect{}, dot
// 	}
// 	if !a.Contains(prevR) {
// 		prevR = unicode.ReplacementChar
// 	}
//
// 	if prevR >= 0 {
// 		dot.X += a.Kern(prevR, r)
// 	}
//
// 	glyph := a.Glyph(r)
//
// 	rect = glyph.Frame.Moved(dot.Sub(glyph.Dot))
// 	bounds = rect
//
// 	if bounds.W()*bounds.H() != 0 {
// 		bounds = pixel.R(
// 			bounds.Min.X,
// 			dot.Y-a.Descent(),
// 			bounds.Max.X,
// 			dot.Y+a.Ascent(),
// 		)
// 	}
//
// 	dot.X += glyph.Advance
//
// 	return rect, glyph.Frame, bounds, dot
// }
//
// type fixedGlyph struct {
// 	dot     fixed.Point26_6
// 	frame   fixed.Rectangle26_6
// 	advance fixed.Int26_6
// }
//
// // makeSquareMapping finds an optimal glyph arrangement of the given runes, so that their common
// // bounding box is as square as possible.
// func makeSquareMapping(face font.Face, runes []rune, padding fixed.Int26_6) (map[rune]fixedGlyph, fixed.Rectangle26_6) {
// 	width := sort.Search(int(fixed.I(1024*1024)), func(i int) bool {
// 		width := fixed.Int26_6(i)
// 		_, bounds := makeMapping(face, runes, padding, width)
// 		return bounds.Max.X-bounds.Min.X >= bounds.Max.Y-bounds.Min.Y
// 	})
// 	return makeMapping(face, runes, padding, fixed.Int26_6(width))
// }
//
// // makeMapping arranges glyphs of the given runes into rows in such a way, that no glyph is located
// // fully to the right of the specified width. Specifically, it places glyphs in a row one by one and
// // once it reaches the specified width, it starts a new row.
// func makeMapping(face font.Face, runes []rune, padding, width fixed.Int26_6) (map[rune]fixedGlyph, fixed.Rectangle26_6) {
// 	mapping := make(map[rune]fixedGlyph)
// 	bounds := fixed.Rectangle26_6{}
//
// 	dot := fixed.P(0, 0)
//
// 	for _, r := range runes {
// 		b, advance, ok := face.GlyphBounds(r)
// 		if !ok {
// 			continue
// 		}
//
// 		// this is important for drawing, artifacts arise otherwise
// 		frame := fixed.Rectangle26_6{
// 			Min: fixed.P(b.Min.X.Floor(), b.Min.Y.Floor()),
// 			Max: fixed.P(b.Max.X.Ceil(), b.Max.Y.Ceil()),
// 		}
//
// 		dot.X -= frame.Min.X
// 		frame = frame.Add(dot)
//
// 		mapping[r] = fixedGlyph{
// 			dot:     dot,
// 			frame:   frame,
// 			advance: advance,
// 		}
// 		bounds = bounds.Union(frame)
//
// 		dot.X = frame.Max.X
//
// 		// padding + align to integer
// 		dot.X += padding
// 		dot.X = fixed.I(dot.X.Ceil())
//
// 		// width exceeded, new row
// 		if frame.Max.X >= width {
// 			dot.X = 0
// 			dot.Y += face.Metrics().Ascent + face.Metrics().Descent
//
// 			// padding + align to integer
// 			dot.Y += padding
// 			dot.Y = fixed.I(dot.Y.Ceil())
// 		}
// 	}
//
// 	return mapping, bounds
// }
//
// func i2f(i fixed.Int26_6) float64 {
// 	return float64(i) / (1 << 6)
// }
//
//
//
//

// import (
// 	"bufio"
// 	"fmt"
// 	"image"
// 	"image/png"
// 	"os"
// 	"strings"
// 	"unicode/utf8"
// 	"unsafe"
//
// 	"github.com/leonsal/gux/gb"
// )
//
// // CharInfo contains the information to locate a character in an FontAtlas
// type CharInfo struct {
// 	X      int        // Position X in pixels in the sheet image from left to right
// 	Y      int        // Position Y in pixels in the sheet image from top to bottom
// 	Width  int        // Char width in pixels
// 	Height int        // Char height in pixels (LINE HEIGHT)
// 	UV     [4]gb.Vec2 // UV coordinates for char quad vertices
// }
//
// // FontAtlas represents an image containing characters and the information about their location in the image
// type FontAtlas struct {
// 	Chars      map[rune]CharInfo // Maps rune code to correspondent CharInfo
// 	Image      *image.RGBA       // Font atlas generated image
// 	LineHeight int               // Total line height
// 	Ascent     int               // Distance from the top of a line to its base line
// 	Descent    int               // Distance from the bottom of a line to its baseline
// 	TexID      gb.TextureID      // Texture ID of this atlas
// }
//
// // NewFontAtlas creates a font atlas using the specified font and range of character codepoints.
// // A Texture is created and sent to the graphics backend.
// func (w *Window) NewFontAtlas(font *Font, first, last rune) *FontAtlas {
//
// 	//     Vertices indices for for each character quad
// 	//
// 	//     0       3
// 	//     +-------+
// 	//     |\      |
// 	//     | \     |
// 	//     |  \    |
// 	//     |   \   |
// 	//     |    \  |
// 	//     |     \ |
// 	//     |      \|
// 	//     +-------+
// 	//     1       2
//
// 	a := new(FontAtlas)
// 	a.Chars = make(map[rune]CharInfo)
//
// 	// Get font metrics
// 	metrics := font.Metrics()
// 	a.Ascent = metrics.Ascent.Round()
// 	a.Descent = metrics.Descent.Round()
// 	a.LineHeight = a.Ascent + a.Descent
//
// 	const maxCols = 32
// 	col := 0
// 	encoded := make([]byte, 4)
// 	line := []byte{}
// 	lines := strings.Builder{}
// 	maxWidth := 0
// 	lastX := 0
// 	lastY := 0
// 	nlines := 1
// 	for code := first; code <= last; code++ {
//
// 		// Ignore codes which doesn't have associate glyph in the font
// 		if font.Index(code) == 0 {
// 			continue
// 		}
//
// 		// Encodes rune into UTF8, appends to current line and measure line width
// 		count := utf8.EncodeRune(encoded, code)
// 		line = append(line, encoded[:count]...)
// 		lineWidth, _ := font.MeasureText(string(line))
//
// 		// Sets current code fields
// 		var cinfo CharInfo
// 		cinfo.X = lastX
// 		cinfo.Y = lastY
// 		cinfo.Width = lineWidth - lastX
// 		cinfo.Height = a.LineHeight
// 		lastX = lineWidth
// 		a.Chars[code] = cinfo
//
// 		// Updates maximum image width
// 		if lineWidth > maxWidth {
// 			maxWidth = lineWidth
// 		}
//
// 		if code == last {
// 			lines.WriteString(string(line))
// 			break
// 		}
//
// 		// Checks end of the current line
// 		col++
// 		if col >= maxCols {
// 			nlines++
// 			lines.WriteString(string(line))
// 			lines.WriteString("\n")
// 			line = []byte{}
// 			col = 0
// 			lastX = 0
// 			lastY += a.LineHeight
// 		}
// 	}
// 	height := nlines * a.LineHeight
//
// 	// Calculate UV coordinates for each char
// 	imgWidth := float32(maxWidth)
// 	imgHeight := float32(height)
// 	for code := first; code <= last; code++ {
// 		char, ok := a.Chars[code]
// 		if !ok {
// 			continue
// 		}
// 		char.UV[0] = gb.Vec2{float32(char.X) / imgWidth, (float32(char.Y) / imgHeight)}
// 		char.UV[1] = gb.Vec2{float32(char.X) / imgWidth, (float32(char.Y+char.Height) / imgHeight)}
// 		char.UV[2] = gb.Vec2{float32(char.X+char.Width) / imgWidth, (float32(char.Y+char.Height) / imgHeight)}
// 		char.UV[3] = gb.Vec2{float32(char.X+char.Width) / imgWidth, (float32(char.Y) / imgHeight)}
// 		a.Chars[code] = char
// 		fmt.Printf("code: %v cinfo: %+v\n", code, char)
// 	}
//
// 	// Generates atlas image
// 	a.Image = font.DrawText(lines.String())
//
// 	// Creates backend texture to store the image and transfer the image
// 	bounds := a.Image.Bounds()
// 	a.TexID = w.CreateTexture(bounds.Dx(), bounds.Dy(), (*gb.RGBA)(unsafe.Pointer(&a.Image.Pix[0])))
// 	return a
// }
//
// // DestroyFontAtlas destroys the specified FontAtlas
// func (w *Window) DestroyFontAtlas(fa *FontAtlas) {
//
// 	w.DeleteTexture(fa.TexID)
// }
//
// // SavePNG saves the current atlas image as a PNG image file
// func (a *FontAtlas) SavePNG(filename string) error {
//
// 	// Save that RGBA image to disk.
// 	outFile, err := os.Create(filename)
// 	if err != nil {
// 		return err
// 	}
// 	defer outFile.Close()
//
// 	b := bufio.NewWriter(outFile)
// 	err = png.Encode(b, a.Image)
// 	if err != nil {
// 		return err
// 	}
// 	err = b.Flush()
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }
