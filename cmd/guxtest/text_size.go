package main

import (
	"fmt"
	"log"

	"github.com/leonsal/gux/gb"
	"github.com/leonsal/gux/window"
	"golang.org/x/image/font"
	"golang.org/x/image/font/gofont/gomedium"
	"golang.org/x/image/font/opentype"
)

func init() {

	registerTest("text_size", 12, newTestText)
}

type testText struct {
	fonts []*window.FontAtlas
}

func newTestText(win *window.Window) ITest {

	t := new(testText)

	// Initial font face options
	opts := opentype.FaceOptions{
		Size:    0,
		DPI:     96,
		Hinting: font.HintingNone,
	}

	// Creates font atlas
	runes := []rune{}
	for r := rune(32); r <= 126; r++ {
		runes = append(runes, r)
	}

	// Creates array of Font Atlases with the specified font sizes
	sizes := []int{12, 18, 22, 28, 32, 40, 48, 64, 72, 144}
	for _, size := range sizes {
		opts.Size = float64(size)
		//fa, err := gux.NewFontAtlasFromFile(win, "/usr/share/fonts/truetype/ubuntu/Ubuntu-R.ttf", &opts, runes)
		fa, err := window.NewFontAtlas(win, gomedium.TTF, &opts, runes)
		if err != nil {
			log.Fatal(err)
		}
		t.fonts = append(t.fonts, fa)
		if false {
			err := fa.SavePNG(fmt.Sprintf("atlas_%d.png", size))
			if err != nil {
				log.Fatal(err)
			}
		}
		fa.ReleaseImage()
	}
	return t
}

func (t *testText) draw(win *window.Window) {

	dl := win.DrawList()
	pos := gb.Vec2{10, 0}
	pos.Y += t.fonts[0].Height()
	text1 := `The quick brown fox jumps over the lazy dog.`

	for _, fa := range t.fonts {
		origin := pos
		win.AddText(dl, fa, &origin, gb.MakeColor(0, 0, 0, 255), window.TextVAlignBase, text1)
		pos.Y += fa.Height() * 2
	}

	//win.AddPolyLineTextured(dl, []gb.Vec2{{0, pos.Y}, {2000, pos.Y}}, gb.MakeColor(0, 0, 0, 100), 0, 2)
	//win.AddText(dl, t.fonts[0], pos, gb.MakeColor(0, 0, 0, 255), gux.TextVAlignBase, text1)
}

func (t *testText) destroy(win *window.Window) {

	for _, f := range t.fonts {
		f.Destroy(win)
	}
	log.Println("Destroy font atlas and texture")
}
