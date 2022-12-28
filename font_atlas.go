package gux

import (
	"bufio"
	"image"
	"image/png"
	"os"
	"strings"
	"unicode/utf8"
	"unsafe"

	"github.com/leonsal/gux/gb"
)

// CharInfo contains the information to locate a character in an FontAtlas
type CharInfo struct {
	X      int        // Position X in pixels in the sheet image from left to right
	Y      int        // Position Y in pixels in the sheet image from top to bottom
	Width  int        // Char width in pixels
	Height int        // Char height in pixels (LINE HEIGHT)
	UV     [4]gb.Vec2 // UV coordinates for char quad vertices
}

// FontAtlas represents an image containing characters and the information about their location in the image
type FontAtlas struct {
	Chars      map[rune]CharInfo // Maps rune code to correspondent CharInfo
	Image      *image.RGBA       // Font atlas generated image
	LineHeight int               // Total line height
	Ascent     int               // Distance from the top of a line to its base line
	Descent    int               // Distance from the bottom of a line to its baseline
	TexID      gb.TextureID      // Texture ID of this atlas
}

// NewFontAtlas creates a font atlas using the specified font and range of character codepoints.
// A Texture is created and sent to the graphics backend.
func (w *Window) NewFontAtlas(font *Font, first, last rune) *FontAtlas {

	//     Vertices indices for for each character quad
	//
	//     0       3
	//     +-------+
	//     |\      |
	//     | \     |
	//     |  \    |
	//     |   \   |
	//     |    \  |
	//     |     \ |
	//     |      \|
	//     +-------+
	//     1       2

	a := new(FontAtlas)
	a.Chars = make(map[rune]CharInfo)

	// Get font metrics
	metrics := font.Metrics()
	a.Ascent = metrics.Ascent.Round()
	a.Descent = metrics.Descent.Round()
	a.LineHeight = a.Ascent + a.Descent

	const maxCols = 32
	col := 0
	encoded := make([]byte, 4)
	line := []byte{}
	lines := strings.Builder{}
	maxWidth := 0
	lastX := 0
	lastY := 0
	nlines := 1
	for code := first; code <= last; code++ {

		// Ignore codes which doesn't have associate glyph in the font
		if font.Index(code) == 0 {
			continue
		}

		// Encodes rune into UTF8, appends to current line and measure line width
		count := utf8.EncodeRune(encoded, code)
		line = append(line, encoded[:count]...)
		lineWidth, _ := font.MeasureText(string(line))

		// Sets current code fields
		var cinfo CharInfo
		cinfo.X = lastX
		cinfo.Y = lastY
		cinfo.Width = lineWidth - lastX - 1
		cinfo.Height = a.LineHeight
		lastX = lineWidth
		a.Chars[code] = cinfo

		// Updates maximum image width
		if lineWidth > maxWidth {
			maxWidth = lineWidth
		}

		if code == last {
			lines.WriteString(string(line))
			break
		}

		// Checks end of the current line
		col++
		if col >= maxCols {
			nlines++
			lines.WriteString(string(line))
			lines.WriteString("\n")
			line = []byte{}
			col = 0
			lastX = 0
			lastY += a.LineHeight
		}
	}
	height := nlines * a.LineHeight

	// Calculate UV coordinates for each char
	imgWidth := float32(maxWidth)
	imgHeight := float32(height)
	for code := first; code <= last; code++ {
		char, ok := a.Chars[code]
		if !ok {
			continue
		}
		char.UV[0] = gb.Vec2{float32(char.X) / imgWidth, (float32(char.Y) / imgHeight)}
		char.UV[1] = gb.Vec2{float32(char.X) / imgWidth, (float32(char.Y+char.Height) / imgHeight)}
		char.UV[2] = gb.Vec2{float32(char.X+char.Width) / imgWidth, (float32(char.Y+char.Height) / imgHeight)}
		char.UV[3] = gb.Vec2{float32(char.X+char.Width) / imgWidth, (float32(char.Y) / imgHeight)}
		a.Chars[code] = char
	}

	// Generates atlas image
	a.Image = font.DrawText(lines.String())

	// Creates backend texture to store the image and transfer the image
	bounds := a.Image.Bounds()
	a.TexID = w.CreateTexture(bounds.Dx(), bounds.Dy(), (*gb.RGBA)(unsafe.Pointer(&a.Image.Pix[0])))
	return a
}

// DestroyFontAtlas destroys the specified FontAtlas
func (w *Window) DestroyFontAtlas(fa *FontAtlas) {

	w.DeleteTexture(fa.TexID)
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
