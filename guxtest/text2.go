package main

import (
	"fmt"
	"log"

	"github.com/leonsal/gux"
	"github.com/leonsal/gux/gb"
)

func init() {

	registerTest("text2", 7, newTestText2)
}

type testText2 struct {
	f  *gux.Font
	fa *gux.FontAtlas
}

func newTestText2(win *gux.Window) ITest {

	t := new(testText2)

	// Load fonts from embedded filesystem
	var err error
	textFont, err := embedfs.Open("assets/Roboto-Medium.ttf")
	if err != nil {
		panic(err)
	}
	defer textFont.Close()

	// Create font
	t.f, err = gux.NewFont(textFont)
	if err != nil {
		panic(err)
	}
	t.f.SetFgColor(gb.MakeColor(0, 0, 0, 255))
	//t.f.SetBgColor(gb.MakeColor(0, 0, 0, 100))
	t.f.SetPointSize(80)

	// Create font atlas
	t.fa = win.NewFontAtlas(t.f, 0x00, 0xFF)
	log.Println("Created atlas: LineHeight:", t.fa.LineHeight, "Ascent:", t.fa.Ascent, "Descent:", t.fa.Descent)
	if true {
		err = t.fa.SavePNG("atlas.png")
		if err != nil {
			fmt.Println("SAVE ERROR:", err)
		}
	}
	return t
}

func (t *testText2) draw(win *gux.Window) {

	dl := win.DrawList()
	text := `igijigijlgiglg|g
Small text aligned to the top
`
	win.AddText(dl, t.fa, gb.Vec2{0, 0}, gux.TextVAlignTop, text)
}

func (t *testText2) destroy(win *gux.Window) {

	win.DestroyFontAtlas(t.fa)
	log.Println("Destroy font atlas and texture")
}
