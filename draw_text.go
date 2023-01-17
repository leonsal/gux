package gux

import (
	"unicode"

	"github.com/leonsal/gux/gb"
)

type TextVAlign int

const (
	TextVAlignBase   TextVAlign = 0
	TextVAlignTop    TextVAlign = 1
	TextVAlignBottom TextVAlign = 2
)

// AddText adds command to draw the specifiedraw text.
func (w *Window) AddText(dl *gb.DrawList, fa *FontAtlas, pos gb.Vec2, color gb.RGBA, align TextVAlign, text string) {

	// Each glyph is rendered as Quad
	//
	//  ........................................... ascent
	//                       0          3
	//                       +----------+
	//       0           3   |          |
	//       +-----------+   |          |
	//       |           |   |          |
	//       |           |   |          |
	//       |           |   |          |
	//  ...O.|...........|.O.+----------+.......... baseline
	//       |           |   1          2
	//       |           |
	//       +-----------+
	//       1           2
	//
	//     |---------------|  O=Glyph origin
	//       Advance
	//  .................... ...................... descent

	posX := pos.X
	var posY float32
	switch align {
	case TextVAlignBase:
		posY = pos.Y
	case TextVAlignTop:
		posY = pos.Y - fa.Ascent
	case TextVAlignBottom:
		posY = pos.Y + fa.Descent
	}

	// For each rune in the text
	prevC := rune(-1)
	for _, code := range text {

		// Process new line
		if code == 0x0A {
			posX = pos.X
			posY += float32(fa.LineHeight)
			continue
		}

		// If glyph not found, use replacement char
		gi, ok := fa.Glyphs[code]
		if !ok {
			gi = fa.Glyphs[unicode.ReplacementChar]
		}

		// Adds  horizontal adjustment for the kerning pair (r0, r1) for the FontAtlas face.
		if prevC >= 0 {
			pos.X += fa.Kern(prevC, code)
		}

		// Creates new DrawCmd to draw a Quad for the glyph bounds.
		cmd, bufIdx, bufVtx := w.NewDrawCmd(dl, 6, 4)
		cmd.TexID = fa.TexID
		bufVtx[0].Pos = gb.Vec2{posX + gi.Bounds.Min.X, posY + gi.Bounds.Min.Y}
		bufVtx[0].UV = gi.UV[0]
		bufVtx[0].Col = color

		bufVtx[1].Pos = gb.Vec2{posX + gi.Bounds.Min.X, posY + gi.Bounds.Max.Y}
		bufVtx[1].UV = gi.UV[1]
		bufVtx[1].Col = color

		bufVtx[2].Pos = gb.Vec2{posX + gi.Bounds.Max.X, posY + gi.Bounds.Max.Y}
		bufVtx[2].UV = gi.UV[2]
		bufVtx[2].Col = color

		bufVtx[3].Pos = gb.Vec2{posX + gi.Bounds.Max.X, posY + gi.Bounds.Min.Y}
		bufVtx[3].UV = gi.UV[3]
		bufVtx[3].Col = color

		bufIdx[0] = 0
		bufIdx[1] = 1
		bufIdx[2] = 2
		bufIdx[3] = 2
		bufIdx[4] = 3
		bufIdx[5] = 0
		posX += gi.Advance
		prevC = code
	}
}
