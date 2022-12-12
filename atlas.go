package gux

import (
	"bufio"
	"image"
	"image/png"
	"os"
	"unicode/utf8"

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
	Chars      []CharInfo  // Array of CharInfo indexed by its rune code
	Image      *image.RGBA // Font atlas generated image
	LineHeight int         // Total line height
	Height     int         // Recommended vertical space between two lines of text
	Ascent     int         // Distance from the top of a line to its base line
	Descent    int         // Distance from the bottom of a line to its baseline
}

// NewFontAtlas returns a pointer to a new FontAtlas object
func NewFontAtlas(font *Font, first, last rune) *FontAtlas {

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
	a.Chars = make([]CharInfo, last+1)

	// Get font metrics
	metrics := font.Metrics()
	a.Height = int(metrics.Height >> 6)
	a.Ascent = int(metrics.Ascent >> 6)
	a.Descent = int(metrics.Descent >> 6)
	//fmt.Printf("Font height:%d Font ascent:%d Font descent:%d\n", a.Height, a.Ascent, a.Descent)

	const maxCols = 32
	col := 0
	encoded := make([]byte, 4)
	line := []byte{}
	lines := ""
	maxWidth := 0
	lastX := 0
	lastY := 0
	nlines := 1
	var lineHeight, lineWidth int
	for code := first; code <= last; code++ {

		// Encodes rune into UTF8, appends to current line and measure line width
		count := utf8.EncodeRune(encoded, code)
		line = append(line, encoded[:count]...)
		lineWidth, lineHeight = font.MeasureText(string(line))

		// Sets current code fields
		cinfo := &a.Chars[code]
		cinfo.X = lastX
		cinfo.Y = lastY
		cinfo.Width = lineWidth - lastX - 1
		cinfo.Height = lineHeight
		lastX = lineWidth

		// Updates maximum image width
		if lineWidth > maxWidth {
			maxWidth = lineWidth
		}

		if code == last {
			lines += string(line)
			break
		}

		// Checks end of the current line
		col++
		if col >= maxCols {
			nlines++
			lines += string(line) + "\n"
			line = []byte{}
			col = 0
			lastX = 0
			lastY += lineHeight
		}
	}
	height := nlines * lineHeight
	a.LineHeight = lineHeight

	// Calculate UV coordinates for each char
	imgWidth := float32(maxWidth)
	imgHeight := float32(height)
	for i := 0; i < len(a.Chars); i++ {
		char := &a.Chars[i]
		char.UV[0] = gb.Vec2{float32(char.X) / imgWidth, (float32(char.Y) / imgHeight)}
		char.UV[1] = gb.Vec2{float32(char.X) / imgWidth, (float32(char.Y+char.Height) / imgHeight)}
		char.UV[2] = gb.Vec2{float32(char.X+char.Width) / imgWidth, (float32(char.Y+char.Height) / imgHeight)}
		char.UV[3] = gb.Vec2{float32(char.X+char.Width) / imgWidth, (float32(char.Y) / imgHeight)}
		if i >= 65 {
			//fmt.Printf("i:%d char:%+v\n", i, char)
		}
	}

	// Generates atlas image
	a.Image = font.DrawText(lines)
	return a
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
