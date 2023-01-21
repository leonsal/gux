package main

import (
	"bufio"
	"fmt"
	"log"

	"github.com/leonsal/gux"
	"github.com/leonsal/gux/gb"
	"golang.org/x/image/font"
	"golang.org/x/image/font/gofont/gomedium"
	"golang.org/x/image/font/opentype"
)

func init() {

	registerTest("text_book", 6, newTestTextBook)
}

type testTextBook struct {
	fa        *gux.FontAtlas
	lines     []string
	firstLine int
	frames    int
}

func newTestTextBook(win *gux.Window) ITest {

	t := new(testTextBook)

	// Initial font face options
	opts := opentype.FaceOptions{
		Size:    28,
		DPI:     72,
		Hinting: font.HintingNone,
	}

	// Creates font atlas
	runes := []rune{}
	for r := rune(32); r <= 126; r++ {
		runes = append(runes, r)
	}
	fa, err := gux.NewFontAtlas(win, gomedium.TTF, &opts, runes)
	if err != nil {
		log.Fatal(err)
	}
	if false {
		err := fa.SavePNG(fmt.Sprintf("atlas_%d.png", int(opts.Size)))
		if err != nil {
			log.Fatal(err)
		}
	}
	fa.ReleaseImage()
	t.fa = fa

	// Reads book file
	file, err := embedfs.Open("assets/book.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		t.lines = append(t.lines, line)
	}
	err = scanner.Err()
	if err != nil {
		log.Fatal(err)
	}
	return t
}

func (t *testTextBook) draw(win *gux.Window) {

	dl := win.DrawList()
	color := gb.MakeColor(0, 0, 0, 255)
	origin := gb.Vec2{40, t.fa.Ascent()}

	nlines := 80
	remain := len(t.lines) - t.firstLine
	if remain < nlines {
		nlines = remain
	}
	for l := t.firstLine; l < t.firstLine+nlines; l++ {
		pos := origin
		win.AddText(dl, t.fa, &pos, color, gux.TextVAlignBase, t.lines[l])
		origin.Y += t.fa.Height()
	}
	t.frames++
	if t.frames >= 50 {
		t.frames = 0
		t.firstLine++
	}
}

func (t *testTextBook) destroy(win *gux.Window) {

	t.fa.Destroy(win)
	log.Println("Destroy font atlas and texture")
}
