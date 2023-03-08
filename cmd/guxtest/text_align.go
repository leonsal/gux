package main

import (
	"fmt"
	"log"
	"strconv"
	"unicode"

	"github.com/leonsal/gux/gb"
	"github.com/leonsal/gux/util"
	"github.com/leonsal/gux/window"
	"golang.org/x/image/font"
	"golang.org/x/image/font/gofont/goregular"
	"golang.org/x/image/font/opentype"
)

func init() {

	registerTest("text_align", 10, newTestTextAlign)
}

type testTextAlign struct {
	fa     *window.FontAtlas
	fcount uint64
	fbuf   []byte
}

func newTestTextAlign(win *window.Window) ITest {

	t := new(testTextAlign)

	// Creates FontAtlas
	opts := opentype.FaceOptions{
		Size:    72,
		DPI:     72,
		Hinting: font.HintingNone,
	}
	fa, err := window.NewFontAtlas(win, goregular.TTF, &opts,
		util.AsciiSet(), util.RangeTableSet(unicode.Latin), util.RangeTableSet(unicode.Common))
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
	t.fbuf = make([]byte, 0, 10)
	return t
}

func (t *testTextAlign) draw(win *window.Window) {

	dl := win.DrawList()
	pos := gb.Vec2{10, 200}
	textColor := gb.MakeColor(0, 0, 0, 255)
	win.AddPolyLineTextured(dl, []gb.Vec2{{0, pos.Y}, {2000, pos.Y}}, gb.MakeColor(255, 0, 0, 255), 0, 2)
	win.AddPolyLineTextured(dl, []gb.Vec2{{0, pos.Y - t.fa.Ascent()}, {2000, pos.Y - t.fa.Ascent()}}, gb.MakeColor(0, 0, 0, 50), 0, 2)
	win.AddPolyLineTextured(dl, []gb.Vec2{{0, pos.Y + t.fa.Descent()}, {2000, pos.Y + t.fa.Descent()}}, gb.MakeColor(0, 0, 0, 50), 0, 2)

	dot := pos
	textBaseline := "TextVAlignBase"
	win.AddText(dl, t.fa, &dot, textColor, window.TextVAlignBase, textBaseline)

	dot = gb.Vec2{dot.X, pos.Y}
	textTop := "TextVAlignTop"
	win.AddText(dl, t.fa, &dot, textColor, window.TextVAlignTop, textTop)

	dot = gb.Vec2{dot.X, pos.Y}
	textBottom := "TextVAlignBottom"
	win.AddText(dl, t.fa, &dot, textColor, window.TextVAlignBottom, textBottom)

	// Formats number of frames with no allocation
	t.fcount++
	t.fbuf = strconv.AppendUint(t.fbuf, t.fcount, 10)
	// Shows number of frames and resets its buffer
	dot = gb.Vec2{0, 0}
	win.AddTextBytes(dl, t.fa, &dot, textColor, window.TextVAlignTop, t.fbuf)
	t.fbuf = t.fbuf[:0]
}

func (t *testTextAlign) destroy(win *window.Window) {

	t.fa.Destroy(win)
	log.Println("Destroy font atlas and texture")
}
