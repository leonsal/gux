package main

import (
	"fmt"
	"log"
	"unicode"

	"github.com/leonsal/gux"
	"github.com/leonsal/gux/gb"
	"golang.org/x/image/font"
	"golang.org/x/image/font/gofont/goregular"
	"golang.org/x/image/font/opentype"
)

func init() {

	registerTest("text_align", 6, newTestTextAlign)
}

type testTextAlign struct {
	fa *gux.FontAtlas
}

func newTestTextAlign(win *gux.Window) ITest {

	t := new(testTextAlign)

	// Creates FontAtlas
	opts := opentype.FaceOptions{
		Size:    72,
		DPI:     72,
		Hinting: font.HintingNone,
	}
	fa, err := gux.NewFontAtlas(win, goregular.TTF, &opts,
		gux.AsciiSet(), gux.RangeTableSet(unicode.Latin), gux.RangeTableSet(unicode.Common))
	if err != nil {
		log.Fatal(err)
	}

	// Optionally save PNG
	if true {
		err := fa.SavePNG(fmt.Sprintf("atlas_%d.png", int(opts.Size)))
		if err != nil {
			log.Fatal(err)
		}
	}
	fa.ReleaseImage()
	t.fa = fa
	return t
}

func (t *testTextAlign) draw(win *gux.Window) {

	dl := win.DrawList()
	pos := gb.Vec2{10, 200}
	color := gb.MakeColor(0, 0, 0, 255)
	lineColor := gb.MakeColor(0, 0, 0, 50)
	win.AddPolyLineTextured(dl, []gb.Vec2{{0, pos.Y}, {2000, pos.Y}}, lineColor, 0, 2)
	win.AddPolyLineTextured(dl, []gb.Vec2{{0, pos.Y - t.fa.Ascent()}, {2000, pos.Y - t.fa.Ascent()}}, lineColor, 0, 2)
	win.AddPolyLineTextured(dl, []gb.Vec2{{0, pos.Y + t.fa.Descent()}, {2000, pos.Y + t.fa.Descent()}}, lineColor, 0, 2)

	dot := pos
	textBaseline := "TextVAlignBase"
	win.AddText(dl, t.fa, &dot, color, gux.TextVAlignBase, textBaseline)

	dot = gb.Vec2{dot.X, pos.Y}
	textTop := "TextVAlignTop"
	win.AddText(dl, t.fa, &dot, color, gux.TextVAlignTop, textTop)

	dot = gb.Vec2{dot.X, pos.Y}
	textBottom := "TextVAlignBottom"
	win.AddText(dl, t.fa, &dot, color, gux.TextVAlignBottom, textBottom)
}

func (t *testTextAlign) destroy(win *gux.Window) {

	t.fa.Destroy(win)
	log.Println("Destroy font atlas and texture")
}
