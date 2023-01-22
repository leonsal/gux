package main

import (
	"bufio"
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

	registerTest("text_book", 8, newTestTextBook)
}

type testTextBook struct {
	fa        *gux.FontAtlas
	lines     []string
	firstLine int
	posY      float32
}

func newTestTextBook(win *gux.Window) ITest {

	t := new(testTextBook)

	// Creates FontAtlas
	opts := opentype.FaceOptions{
		Size:    36,
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

	// Reads all lines of the book file
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

	// Initial font dot position
	t.posY = t.fa.Ascent()
	return t
}

func (t *testTextBook) draw(win *gux.Window) {

	dl := win.DrawList()
	color := gb.MakeColor(0, 0, 0, 255)
	origin := gb.Vec2{40, t.posY}

	// Calculates number of lines remaining
	maxLines := int(win.Size().Y/t.fa.Height()) + 2
	nlines := maxLines
	remain := len(t.lines) - t.firstLine
	if remain < nlines {
		nlines = remain
	}

	// Draw lines
	for l := t.firstLine; l < t.firstLine+nlines; l++ {
		pos := origin
		win.AddText(dl, t.fa, &pos, color, gux.TextVAlignBase, t.lines[l])
		origin.Y += t.fa.Height()
	}

	// Scroll origin position up some pixels
	t.posY -= 1
	if t.posY <= -t.fa.Descent() {
		t.posY = t.fa.Ascent()
		t.firstLine++
		if t.firstLine >= len(t.lines) {
			t.firstLine = 0
		}
		fmt.Println("firstLine:", t.firstLine, "maxLines:", maxLines, "idxCount:", dl.IdxCount())
	}
}

func (t *testTextBook) destroy(win *gux.Window) {

	t.fa.Destroy(win)
	log.Println("Destroy font atlas and texture")
}
