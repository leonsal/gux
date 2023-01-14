package gux

// import "github.com/leonsal/gux/gb"
//
// type TextVAlign int
//
// const (
// 	TextVAlignTop    TextVAlign = 0
// 	TextVAlignBase   TextVAlign = 1
// 	TextVAlignBottom TextVAlign = 2
// )
//
// // AddText adds commands to draw text to the specified DrawList.
// func (w *Window) AddText(dl *gb.DrawList, fa *FontAtlas, pos gb.Vec2, align TextVAlign, text string) {
//
// 	white := gb.MakeColor(255, 255, 255, 255)
//
// 	posX := pos.X
// 	var posY float32
// 	switch align {
// 	case TextVAlignTop:
// 		posY = pos.Y
// 	case TextVAlignBase:
// 		posY = pos.Y - float32(fa.Ascent)
// 	case TextVAlignBottom:
// 		posY = pos.Y - float32(fa.LineHeight)
// 	}
//
// 	// For each rune in the text
// 	for _, code := range text {
//
// 		// Process new line
// 		if code == 0x0A {
// 			posX = pos.X
// 			posY += float32(fa.LineHeight)
// 			continue
// 		}
//
// 		// Ignore codes with no glyphs
// 		charInfo, ok := fa.Glyphs[code]
// 		if !ok {
// 			continue
// 		}
//
// 		//fmt.Printf("char: %v Info:%+v\n", c, charInfo)
// 		cmd, bufIdx, bufVtx := w.NewDrawCmd(dl, 6, 4)
// 		cmd.TexID = fa.TexID
// 		bufVtx[0].Pos = gb.Vec2{posX, posY}
// 		bufVtx[0].UV = charInfo.UV[0]
// 		bufVtx[0].Col = white
//
// 		bufVtx[1].Pos = gb.Vec2{posX, posY + float32(charInfo.Height-1)}
// 		bufVtx[1].UV = charInfo.UV[1]
// 		bufVtx[1].Col = white
//
// 		bufVtx[2].Pos = gb.Vec2{posX + float32(charInfo.Width-1), posY + float32(charInfo.Height-1)}
// 		bufVtx[2].UV = charInfo.UV[2]
// 		bufVtx[2].Col = white
//
// 		bufVtx[3].Pos = gb.Vec2{posX + float32(charInfo.Width-1), posY}
// 		bufVtx[3].UV = charInfo.UV[3]
// 		bufVtx[3].Col = white
//
// 		bufIdx[0] = 0
// 		bufIdx[1] = 1
// 		bufIdx[2] = 2
// 		bufIdx[3] = 2
// 		bufIdx[4] = 3
// 		bufIdx[5] = 0
// 		posX += float32(charInfo.Width - 1)
// 	}
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
