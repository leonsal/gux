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

// AddText adds commands to draw text to the specified DrawList.
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
		//fmt.Printf("code:%v %f/%f %f/%f %f \n", code, gi.Bounds.Min.X, gi.Bounds.Min.Y, gi.Bounds.Max.X, gi.Bounds.Max.Y, gi.Advance)
		if prevC >= 0 {
			pos.X += fa.Kern(prevC, code)
		}

		//fmt.Printf("char: %v Info:%+v\n", c, charInfo)
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

// func (w *Window) CreateTextImage(f *Font, text string) (gb.TextureID, float32, float32) {
//
// 	// Create image and draw text on it
// 	img := f.DrawText(text)
// 	b := img.Bounds()
// 	width := b.Dx()
// 	height := b.Dy()
//
// 	// Creates backend texture to store the image and transfer the image
// 	texID := w.CreateTexture(width, height, (*gb.RGBA)(unsafe.Pointer(&img.Pix[0])))
// 	return texID, float32(width), float32(height)
// }
//
// // AddImage adds command to draw specified image to the DrawList.
// func (w *Window) AddImage(dl *gb.DrawList, texID gb.TextureID, width, height float32, pos gb.Vec2) {
//
// 	//
// 	// UV coordinates adjustment
// 	//
// 	//	  0,1    1,1      0,0    1,0
// 	// 0 +------+ 3       +------+
// 	//	 |\     |         |\     |
// 	//	 | \    |         | \    |
// 	//	 |  \   |  --->   |  \   |
// 	//	 |   \  |         |   \  |
// 	//	 |    \ |         |    \ |
// 	//	 |     \|         |     \|
// 	// 1 +------+ 2       +------+
// 	//	 0,0    1,0       0,1    1,1
//
// 	// Creates command
// 	cmd, bufIdx, bufVtx := w.NewDrawCmd(dl, 6, 4)
// 	cmd.TexID = texID
//
// 	// Set vertices
// 	white := gb.MakeColor(255, 255, 255, 255)
// 	bufVtx[0].Pos = pos
// 	bufVtx[0].UV = gb.Vec2{0, 0}
// 	bufVtx[0].Col = white
//
// 	bufVtx[1].Pos = gb.Vec2{pos.X, pos.Y + height - 1}
// 	bufVtx[1].UV = gb.Vec2{0, 1}
// 	bufVtx[1].Col = white
//
// 	bufVtx[2].Pos = gb.Vec2{pos.X + width - 1, pos.Y + height - 1}
// 	bufVtx[2].UV = gb.Vec2{1, 1}
// 	bufVtx[2].Col = white
//
// 	bufVtx[3].Pos = gb.Vec2{pos.X + width - 1, pos.Y}
// 	bufVtx[3].UV = gb.Vec2{1, 0}
// 	bufVtx[3].Col = white
//
// 	// Set indices
// 	bufIdx[0] = 0
// 	bufIdx[1] = 1
// 	bufIdx[2] = 2
// 	bufIdx[3] = 2
// 	bufIdx[4] = 3
// 	bufIdx[5] = 0
// }
