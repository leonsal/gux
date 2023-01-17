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
	Bounds  gb.Rect    // Glyph bounds relative to its origin point at the baseline
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

	// Builds draw information for each Glyph in the atlas
	glyphs := make(map[rune]GlyphInfo)
	for r, fg := range fixedMapping {

		// Get Glyph bounds and advance converting from fixed to float
		bounds, advance, _ := face.GlyphBounds(r)
		gi := GlyphInfo{}
		gi.Advance = i2f(advance)
		gi.Bounds.Min = gb.Vec2{i2f(bounds.Min.X), i2f(bounds.Min.Y)}
		gi.Bounds.Max = gb.Vec2{i2f(bounds.Max.X), i2f(bounds.Max.Y)}

		// Transform glyphs image coordinates to UV coordinates
		minX := i2f(fg.frame.Min.X)
		minY := -boundsMinY + i2f(fg.frame.Min.Y)
		maxX := i2f(fg.frame.Max.X)
		maxY := -boundsMinY + i2f(fg.frame.Max.Y)
		gi.UV[0] = gb.Vec2{minX / imageWidth, minY / imageHeight}
		gi.UV[1] = gb.Vec2{minX / imageWidth, maxY / imageHeight}
		gi.UV[2] = gb.Vec2{maxX / imageWidth, maxY / imageHeight}
		gi.UV[3] = gb.Vec2{maxX / imageWidth, minY / imageHeight}
		glyphs[r] = gi
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

// Kern returns the horizontal adjustment for the kerning pair (r0, r1) for the FontAtlas face.
// A positive kern means to move the glyphs further apart.
func (a *FontAtlas) Kern(r0, r1 rune) float32 {

	return i2f(a.Face.Kern(r0, r1))
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
