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
		Size:    132,
		DPI:     72,
		Hinting: font.HintingNone,
	})
	if err != nil {
		log.Fatalf("NewFace: %v", err)
	}
	defer fontFile.Close()

	// Creates font atlas
	runes := []rune{}
	for r := rune(105); r <= 112; r++ {
		runes = append(runes, r)
	}
	t.fa = gux.NewFontAtlas(face, runes)
	t.fa.PrintInfo()

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
	pos := gb.Vec2{100, 200}
	win.AddPolyLineTextured(dl, []gb.Vec2{{0, pos.Y}, {2000, pos.Y}}, gb.MakeColor(0, 0, 0, 255), 0, 2)
	win.AddText(dl, t.fa, pos, gb.MakeColor(0, 0, 0, 255), gux.TextVAlignBase, "ijikiliminioip")
	fmt.Println()
	// win.AddText(dl, t.fa, gb.Vec2{250, 200}, gux.TextVAlignBase, " base")
	// win.AddText(dl, t.fa, gb.Vec2{550, 200}, gux.TextVAlignBottom, " bottom")
	// win.AddImage(dl, t.texID, t.width, t.height, gb.Vec2{50, 400})
}

func (t *testText) destroy(win *gux.Window) {

	t.fa.Destroy(win)
	log.Println("Destroy font atlas and texture")
}
