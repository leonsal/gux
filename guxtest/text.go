package main

import (
	"fmt"
	"log"

	"github.com/leonsal/gux"
	"github.com/leonsal/gux/gb"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
)

func init() {

	registerTest("text", 6, newTestText)
}

type testText struct {
	fa     *gux.FontAtlas
	width  float32
	height float32
}

func newTestText(win *gux.Window) ITest {

	t := new(testText)

	// Opens font file from embedded filesystem
	var err error
	fontFile, err := embedfs.Open("assets/Roboto-Medium.ttf")
	if err != nil {
		log.Fatal(err)
	}

	// Creates font face from file reader
	face, err := gux.NewFontFaceFromReader(fontFile, &opentype.FaceOptions{
		Size:    128,
		DPI:     72,
		Hinting: font.HintingNone,
	})
	if err != nil {
		log.Fatalf("NewFace: %v", err)
	}
	defer fontFile.Close()

	// Creates font atlas
	runes := []rune{}
	for r := rune(105); r < 110; r++ {
		runes = append(runes, r)
	}
	t.fa = gux.NewFontAtlas(face, runes)
	fmt.Println("Ascent:", t.fa.Ascent, "Descent", t.fa.Descent, "Lineheight:", t.fa.LineHeight)
	for code, gi := range t.fa.Glyphs {
		fmt.Printf("code:%v info:%+v\n", code, gi)
	}

	// Optionally save font atlas png for debugging
	if true {
		err = t.fa.SavePNG("atlas.png")
		if err != nil {
			log.Fatalf("SavePNG: %v", err)
		}
	}

	// Creates font atlas texture
	t.fa.CreateTexture(win)
	return t
}

func (t *testText) draw(win *gux.Window) {

	dl := win.DrawList()
	win.AddText(dl, t.fa, gb.Vec2{100, 100}, gux.TextVAlignTop, "ij")
	// win.AddText(dl, t.fa, gb.Vec2{250, 200}, gux.TextVAlignBase, " base")
	// win.AddText(dl, t.fa, gb.Vec2{550, 200}, gux.TextVAlignBottom, " bottom")
	// win.AddImage(dl, t.texID, t.width, t.height, gb.Vec2{50, 400})
}

func (t *testText) destroy(win *gux.Window) {

	t.fa.Destroy(win)
	log.Println("Destroy font atlas and texture")
}
