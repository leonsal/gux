package gux

import (
	"image"
	"image/color"
	"image/draw"
	"io/ioutil"
	"strings"

	"github.com/golang/freetype/truetype"
	"github.com/leonsal/gux/gb"
	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
)

type Font struct {
	ttf         *truetype.Font // The TrueType font
	face        font.Face      // The font face
	pointSize   float64        // Size of the font in points
	dpi         float64        // Resolution of the font in dots per inch
	lineSpacing float64        // Spacing between lines (relative to the font height)
	hinting     font.Hinting   // Font hinting
	fg          *image.Uniform // Text color cache
	bg          *image.Uniform // Background color cache
	changed     bool           // Whether attributes have changed and the font face needs to be recreated
}

// NewFont creates and returns a new font object using the specified TrueType font file.
func NewFont(ttfFile string) (*Font, error) {

	// Reads font bytes
	fontBytes, err := ioutil.ReadFile(ttfFile)
	if err != nil {
		return nil, err
	}
	return NewFontFromData(fontBytes)
}

// NewFontFromData creates and returns a new font object from the specified TTF data.
func NewFontFromData(fontData []byte) (*Font, error) {

	// Parses the font data
	ttf, err := truetype.Parse(fontData)
	if err != nil {
		return nil, err
	}

	// Create font and initialize with default values
	f := new(Font)
	f.ttf = ttf
	//f.pointSize = 18
	f.pointSize = 60
	f.dpi = 72
	f.lineSpacing = 1.0
	f.hinting = font.HintingNone
	f.fg = image.Black
	f.bg = image.Transparent

	// Create font face
	f.face = truetype.NewFace(f.ttf, &truetype.Options{
		Size:    f.pointSize,
		DPI:     f.dpi,
		Hinting: f.hinting,
	})

	return f, nil
}

// Metrics returns the font metrics.
func (f *Font) Metrics() font.Metrics {

	f.updateFace()
	return f.face.Metrics()
}

// SetFgColor sets the text color.
func (f *Font) SetFgColor(col gb.RGBA) {

	f.fg = image.NewUniform(ConvRGBA(col))
}

// SetBgColor sets the background color.
func (f *Font) SetBgColor(col gb.RGBA) {

	f.bg = image.NewUniform(ConvRGBA(col))
}

// SetPointSize sets the point size of the font.
func (f *Font) SetPointSize(size float64) {

	f.pointSize = size
	f.changed = true
}

// MeasureText returns the minimum width and height in pixels necessary for an image to contain
// the specified text. The supplied text string can contain line break escape sequences (\n).
func (f *Font) MeasureText(text string) (int, int) {

	// Create font drawer
	f.updateFace()
	d := &font.Drawer{Dst: nil, Src: f.fg, Face: f.face}

	// Draw text
	var width, height int
	metrics := f.face.Metrics()
	lineHeight := (metrics.Ascent + metrics.Descent).Ceil()
	lineGap := int((f.lineSpacing - float64(1)) * float64(lineHeight))

	lines := strings.Split(text, "\n")
	for i, s := range lines {
		d.Dot = fixed.P(0, height)
		lineWidth := d.MeasureString(s).Ceil()
		if lineWidth > width {
			width = lineWidth
		}
		height += lineHeight
		if i > 1 {
			height += lineGap
		}
	}
	return width, height
}

// DrawText draws the specified text on a new, tightly fitting image, and returns a pointer to the image.
func (f *Font) DrawText(text string) *image.RGBA {

	// Create image with the dimensions of the text area
	width, height := f.MeasureText(text)
	img := image.NewRGBA(image.Rect(0, 0, width, height))

	// Draws the background color of image
	draw.Draw(img, img.Bounds(), f.bg, image.ZP, draw.Src)

	// Draws the text as the foreground color
	f.DrawTextOnImage(text, 0, 0, img)
	return img
}

// DrawTextOnImage draws the specified text on the specified image at the specified coordinates.
func (f *Font) DrawTextOnImage(text string, x, y int, dst *image.RGBA) {

	f.updateFace()
	d := &font.Drawer{Dst: dst, Src: f.fg, Face: f.face}

	// Draw text
	metrics := f.face.Metrics()
	py := y + metrics.Ascent.Round()
	lineHeight := (metrics.Ascent + metrics.Descent).Ceil()
	lineGap := int((f.lineSpacing - float64(1)) * float64(lineHeight))
	lines := strings.Split(text, "\n")
	for i, s := range lines {
		d.Dot = fixed.P(x, py)
		d.DrawString(s)
		py += lineHeight
		if i > 1 {
			py += lineGap
		}
	}
}

// ConvRGBA converts an gb.RGBA color to image/color.RGBA
func ConvRGBA(col gb.RGBA) color.RGBA {

	red := uint8((col & gb.RGBAMaskR) >> gb.RGBAShiftR)
	green := uint8((col & gb.RGBAMaskG) >> gb.RGBAShiftG)
	blue := uint8((col & gb.RGBAMaskB) >> gb.RGBAShiftB)
	alpha := uint8((col & gb.RGBAMaskA) >> gb.RGBAShiftA)
	return color.RGBA{red, green, blue, alpha}
}

// updateFace updates the font face if parameters have changed.
func (f *Font) updateFace() {

	if f.changed {
		f.face = truetype.NewFace(f.ttf, &truetype.Options{
			Size:    f.pointSize,
			DPI:     f.dpi,
			Hinting: f.hinting,
		})
		f.changed = false
	}
}
