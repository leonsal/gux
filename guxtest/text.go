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
	fonts  []*gux.FontAtlas
	width  float32
	height float32
}

func newTestText(win *gux.Window) ITest {

	t := new(testText)

	opts := opentype.FaceOptions{
		Size:    0,
		DPI:     96,
		Hinting: font.HintingNone,
	}
	sizes := []int{12, 18, 22, 28, 32, 40, 48, 64}
	for _, size := range sizes {
		opts.Size = float64(size)
		fa := t.createFontAtlas(win, "assets/Roboto-Medium.ttf", &opts)
		t.fonts = append(t.fonts, fa)
	}
	return t
}

func (t *testText) createFontAtlas(win *gux.Window, filePath string, opts *opentype.FaceOptions) *gux.FontAtlas {

	// Opens font file from embedded filesystem
	var err error
	fontFile, err := embedfs.Open(filePath)
	if err != nil {
		log.Fatal(err)
	}

	// Creates font face from file reader
	face, err := gux.NewFontFaceFromReader(fontFile, opts)
	if err != nil {
		log.Fatalf("NewFace: %v", err)
	}
	defer fontFile.Close()

	// Creates font atlas
	runes := []rune{}
	for r := rune(32); r <= 126; r++ {
		runes = append(runes, r)
	}
	fa := gux.NewFontAtlas(face, runes)

	// Optionally save font atlas png for debugging
	if true {
		name := fmt.Sprintf("atlas_%d.png", int(opts.Size))
		err = fa.SavePNG(name)
		if err != nil {
			log.Fatalf("SavePNG: %v", err)
		}
	}

	// Creates font atlas texture
	fa.CreateTexture(win)
	return fa
}

func (t *testText) draw(win *gux.Window) {

	dl := win.DrawList()
	pos := gb.Vec2{10, 0}
	pos.Y += t.fonts[0].LineHeight

	text1 := `We are merely picking up pebbles on the beach 
while the great ocean of truth
lays completely undiscovered before us.`

	for _, fa := range t.fonts {
		win.AddText(dl, fa, pos, gb.MakeColor(0, 0, 0, 255), gux.TextVAlignBase, text1)
		pos.Y += fa.LineHeight * 4
	}

	//win.AddPolyLineTextured(dl, []gb.Vec2{{0, pos.Y}, {2000, pos.Y}}, gb.MakeColor(0, 0, 0, 100), 0, 2)
	//win.AddText(dl, t.fonts[0], pos, gb.MakeColor(0, 0, 0, 255), gux.TextVAlignBase, text1)
}

func (t *testText) destroy(win *gux.Window) {

	for _, f := range t.fonts {
		f.Destroy(win)
	}
	log.Println("Destroy font atlas and texture")
}
