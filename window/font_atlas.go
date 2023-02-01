package window

import (
	"bufio"
	"errors"
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
	"golang.org/x/image/font/sfnt"
	"golang.org/x/image/math/fixed"
)

type GlyphInfo struct {
	Advance float32    // Amount to add to glyph origin to draw next Glyph
	Bounds  gb.Rect    // Glyph bounds relative to its origin point at the baseline
	UV      [4]gb.Vec2 // UV coordinates for glyph quad vertices
}

// FontAtlas represents an image containing characters and the information about their location in the image
type FontAtlas struct {
	face    font.Face          // The font face used to generate the atlas
	glyphs  map[rune]GlyphInfo // Maps rune code to correspondent Glyph info
	image   *image.RGBA        // Font atlas generated image
	ascent  float32            // Distance from the top of a line to its baseline
	descent float32            // Distance from the bottom of a line to its baseline
	height  float32            // Total line height
	texID   gb.TextureID       // Texture ID (valid only after texture was created)
}

func NewFontAtlasFromFile(w *Window, filepath string, opts *opentype.FaceOptions, runeSets ...[]rune) (*FontAtlas, error) {

	f, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return NewFontAtlasFromReader(w, f, opts, runeSets...)
}

func NewFontAtlasFromReader(w *Window, r io.Reader, opts *opentype.FaceOptions, runeSets ...[]rune) (*FontAtlas, error) {

	fdata, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}
	return NewFontAtlas(w, fdata, opts, runeSets...)
}

func NewFontAtlas(w *Window, fontData []byte, opts *opentype.FaceOptions, runeSets ...[]rune) (*FontAtlas, error) {

	// Parses font data and creates the font face
	fnt, err := opentype.Parse(fontData)
	if err != nil {
		return nil, nil
	}
	face, err := opentype.NewFace(fnt, opts)
	if err != nil {
		return nil, nil
	}

	// Builds array of unique runes from all the specified rune sets
	seen := make(map[rune]bool)
	runes := []rune{unicode.ReplacementChar}
	for _, set := range runeSets {
		for _, r := range set {
			// Checks if there is a Glyph for this rune
			var b sfnt.Buffer
			x, err := fnt.GlyphIndex(&b, r)
			if x == 0 || err != nil {
				continue
			}
			// Appends unique rune
			if !seen[r] {
				runes = append(runes, r)
				seen[r] = true
			}
		}
	}
	fixedMapping, fixedBounds := makeSquareMapping(face, runes, fixed.I(2))

	// Creates font atlas image
	img := image.NewRGBA(image.Rect(
		fixedBounds.Min.X.Floor(),
		fixedBounds.Min.Y.Floor(),
		fixedBounds.Max.X.Ceil(),
		fixedBounds.Max.Y.Ceil(),
	))

	// Draw all glyphs to the font atlas image
	for r, fg := range fixedMapping {
		if dr, mask, maskp, _, ok := face.Glyph(fg.dot, r); ok {
			draw.Draw(img, dr, mask, maskp, draw.Src)
		}
	}

	// Creates Font Atlas texture
	width := img.Rect.Max.X - img.Rect.Min.X
	height := img.Rect.Max.Y - img.Rect.Min.Y
	texID := w.CreateTexture(width, height, (*gb.RGBA)(unsafe.Pointer(&img.Pix[0])))

	// Image bounds
	imgMinX := i2f(fixedBounds.Min.X)
	imgMaxX := i2f(fixedBounds.Max.X)
	imgMinY := i2f(fixedBounds.Min.Y)
	imgMaxY := i2f(fixedBounds.Max.Y)
	imgWidth := imgMaxX - imgMinX
	imgHeight := imgMaxY - imgMinY

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
		minY := -imgMinY + i2f(fg.frame.Min.Y)
		maxX := i2f(fg.frame.Max.X)
		maxY := -imgMinY + i2f(fg.frame.Max.Y)
		gi.UV[0] = gb.Vec2{minX / imgWidth, minY / imgHeight}
		gi.UV[1] = gb.Vec2{minX / imgWidth, maxY / imgHeight}
		gi.UV[2] = gb.Vec2{maxX / imgWidth, maxY / imgHeight}
		gi.UV[3] = gb.Vec2{maxX / imgWidth, minY / imgHeight}
		glyphs[r] = gi
	}

	return &FontAtlas{
		face:    face,
		glyphs:  glyphs,
		image:   img,
		ascent:  i2f(face.Metrics().Ascent),
		descent: i2f(face.Metrics().Descent),
		height:  i2f(face.Metrics().Height),
		texID:   texID,
	}, nil
}

// Face returns the font face of the FontAtlas
func (a *FontAtlas) Face() font.Face {

	return a.face
}

// Glyph returns the GlyphInfo for the specified rune from the FontAtlas
func (a *FontAtlas) Glyph(r rune) (GlyphInfo, bool) {

	gi, ok := a.glyphs[r]
	return gi, ok
}

// Ascent returns the Ascent of the font face used in the FontAtlas
func (a *FontAtlas) Ascent() float32 {

	return a.ascent
}

// Descent returns the Descent of the font face used in the FontAtlas
func (a *FontAtlas) Descent() float32 {

	return a.descent
}

// Height returns the recommended amount of vertical space between two lines of text.
func (a *FontAtlas) Height() float32 {

	return a.height
}

// Kern returns the horizontal adjustment for the kerning pair (r0, r1) for the FontAtlas face.
// A positive kern means to move the glyphs further apart.
func (a *FontAtlas) Kern(r0, r1 rune) float32 {

	return i2f(a.face.Kern(r0, r1))
}

// ReleaseImage releases the memory allocated to the image created to build the font atlas texture.
// After the FontAtlas is created, the image is necessary only if the user wants to save the
// image in a PNG file for inspection using SavePNG().
func (a *FontAtlas) ReleaseImage() {

	a.image = nil
}

// SavePNG saves the current atlas image as a PNG image file
func (a *FontAtlas) SavePNG(filename string) error {

	if a.image == nil {
		return errors.New("FontAtlas image was released")
	}

	// Save that RGBA image to disk.
	outFile, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer outFile.Close()

	b := bufio.NewWriter(outFile)
	err = png.Encode(b, a.image)
	if err != nil {
		return err
	}
	err = b.Flush()
	if err != nil {
		return err
	}
	return nil
}

func (a *FontAtlas) Destroy(win *Window) error {

	if a.texID != 0 {
		win.DeleteTexture(a.texID)
		a.texID = 0
	}
	return nil
}

func (a *FontAtlas) PrintInfo() {

	fmt.Println("Ascent:", a.ascent, "Descent", a.descent, "Lineheight:", a.height)
	runes := make([]rune, 0, len(a.glyphs))
	for r := range a.glyphs {
		runes = append(runes, r)
	}
	sort.Slice(runes, func(i, j int) bool {
		return runes[i] < runes[j]
	})
	for _, r := range runes {
		fmt.Printf("code:%v glyph:[%c] info:%+v\n", r, r, a.glyphs[r])
	}
}

// MeasureString returns how far dot would advance by drawing s with f.
func (a *FontAtlas) MeasureString(s string) float32 {

	var advance float32
	prevC := rune(-1)
	for _, c := range s {
		if prevC >= 0 {
			advance += a.Kern(prevC, c)
		}
		gi, ok := a.glyphs[c]
		if !ok {
			continue
		}
		advance += gi.Advance
		prevC = c
	}
	return advance
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
