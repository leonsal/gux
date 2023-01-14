package gux

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"testing"

	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
)

func TestMain(t *testing.T) {

	f, err := os.Open("guxtest/assets/Roboto-Medium.ttf")
	if err != nil {
		log.Fatalf("Open: %v", err)
	}
	defer f.Close()
	fbytes, err := ioutil.ReadAll(f)
	if err != nil {
		log.Fatalf("ReadAll: %v", err)
	}

	fparsed, err := opentype.Parse(fbytes)
	if err != nil {
		log.Fatalf("Parse: %v", err)
	}

	face, err := opentype.NewFace(fparsed, &opentype.FaceOptions{
		Size:    128,
		DPI:     72,
		Hinting: font.HintingNone,
	})
	if err != nil {
		log.Fatalf("NewFace: %v", err)
	}

	runes := []rune{}
	for r := rune(65); r < 108; r++ {
		runes = append(runes, r)
	}

	fa := NewFontAtlas(face, runes)
	fmt.Println("Ascent:", fa.Ascent, "Descent", fa.Descent, "Lineheight:", fa.LineHeight)
	for code, gi := range fa.Glyphs {
		fmt.Printf("code:%v info:%+v\n", code, gi)
	}

	err = fa.SavePNG("atlas.png")
	if err != nil {
		log.Fatalf("SavePNG: %v", err)
	}

}
