package gux

import (
	"fmt"
	"log"
	"testing"

	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
)

func TestMain(t *testing.T) {

	face, err := NewFontFaceFromFile("guxtest/assets/Roboto-Medium.ttf", &opentype.FaceOptions{
		Size:    128,
		DPI:     72,
		Hinting: font.HintingNone,
	})
	if err != nil {
		log.Fatalf("NewFace: %v", err)
	}

	runes := []rune{}
	for r := rune(105); r < 108; r++ {
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
