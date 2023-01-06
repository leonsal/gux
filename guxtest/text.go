package main

import (
	"fmt"

	"github.com/leonsal/gux"
	"github.com/leonsal/gux/gb"
)

func init() {

	registerTest("text", 6, newTestText)
}

type testText struct {
	f      *gux.Font
	fa     *gux.FontAtlas
	texID  gb.TextureID
	width  float32
	height float32
}

func newTestText(win *gux.Window) ITest {

	t := new(testText)

	// Create font
	var err error
	t.f, err = gux.NewFont("assets/Roboto-Medium.ttf")
	if err != nil {
		panic(err)
	}
	t.f.SetFgColor(gb.MakeColor(255, 255, 0, 255))
	t.f.SetBgColor(gb.MakeColor(0, 0, 0, 100))
	t.f.SetPointSize(148)

	// Create font atlas
	t.fa = win.NewFontAtlas(t.f, 0x00, 0xFF)
	fmt.Println("ATLAS: LineHeight:", t.fa.LineHeight, "Ascent:", t.fa.Ascent, "Descent:", t.fa.Descent)
	if false {
		err = t.fa.SavePNG("atlas.png")
		if err != nil {
			fmt.Println("SAVE ERROR:", err)
		}
	}

	// Creates text image
	t.texID, t.width, t.height = win.CreateTextImage(t.f, "text image")
	fmt.Println("TextImage:", t.texID, t.width, t.height)
	return t
}

func (t *testText) draw(win *gux.Window) {

	dl := win.DrawList()
	win.AddText(dl, t.fa, gb.Vec2{50, 200}, gux.TextVAlignTop, "top ")
	win.AddText(dl, t.fa, gb.Vec2{250, 200}, gux.TextVAlignBase, " base")
	win.AddText(dl, t.fa, gb.Vec2{550, 200}, gux.TextVAlignBottom, " bottom")
	win.AddImage(dl, t.texID, t.width, t.height, gb.Vec2{50, 400})
}

func (t *testText) destroy(win *gux.Window) {

	win.DestroyFontAtlas(t.fa)
	win.DeleteTexture(t.texID)
	fmt.Println("destroy")
}
