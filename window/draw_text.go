package window

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

// AddGlyph adds a Quad to the DrawList showing Glyph specified by its font atlas and code.
// The origin of the Glyph is given by 'pos' which will be updated to the origin of the next Glyph to add.
// 'prev' should be the code of the previous Glyph of the line or -1.
func (w *Window) AddGlyph(dl *gb.DrawList, fa *FontAtlas, pos *gb.Vec2, color gb.RGBA, align TextVAlign, prev, code rune) {

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

	var posY float32
	switch align {
	case TextVAlignBase:
		posY = pos.Y
	case TextVAlignTop:
		posY = pos.Y + fa.ascent
	case TextVAlignBottom:
		posY = pos.Y - fa.descent
	}

	// If glyph not found, use replacement char
	gi, ok := fa.glyphs[code]
	if !ok {
		gi = fa.glyphs[unicode.ReplacementChar]
	}

	// Adds  horizontal adjustment for the kerning pair (r0, r1) for the FontAtlas face.
	if prev >= 0 {
		pos.X += fa.Kern(prev, code)
	}

	// Creates new DrawCmd to draw a Quad for the glyph bounds.
	cmd, bufIdx, bufVtx := w.NewDrawCmd(dl, 6, 4)
	cmd.TexID = fa.texID
	bufVtx[0].Pos = gb.Vec2{pos.X + gi.Bounds.Min.X, posY + gi.Bounds.Min.Y}
	bufVtx[0].UV = gi.UV[0]
	bufVtx[0].Col = color

	bufVtx[1].Pos = gb.Vec2{pos.X + gi.Bounds.Min.X, posY + gi.Bounds.Max.Y}
	bufVtx[1].UV = gi.UV[1]
	bufVtx[1].Col = color

	bufVtx[2].Pos = gb.Vec2{pos.X + gi.Bounds.Max.X, posY + gi.Bounds.Max.Y}
	bufVtx[2].UV = gi.UV[2]
	bufVtx[2].Col = color

	bufVtx[3].Pos = gb.Vec2{pos.X + gi.Bounds.Max.X, posY + gi.Bounds.Min.Y}
	bufVtx[3].UV = gi.UV[3]
	bufVtx[3].Col = color

	bufIdx[0] = 0
	bufIdx[1] = 1
	bufIdx[2] = 2
	bufIdx[3] = 2
	bufIdx[4] = 3
	bufIdx[5] = 0
	pos.X += gi.Advance
}

// AddText adds command to draw text line
func (w *Window) AddText(dl *gb.DrawList, fa *FontAtlas, pos *gb.Vec2, color gb.RGBA, align TextVAlign, text string) {

	// For each rune in the text
	prev := rune(-1)
	for _, code := range text {

		// Process new line
		if code == 0x0A {
			pos.Y += fa.Height()
			continue
		}
		w.AddGlyph(dl, fa, pos, color, align, prev, code)
		prev = code
	}
}
