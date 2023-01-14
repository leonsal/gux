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
		Size:    32,
		DPI:     72,
		Hinting: font.HintingNone,
	})
	if err != nil {
		log.Fatalf("NewFace: %v", err)
	}
	fmt.Printf("%v\n", face)

	//fa := NewFont
	//fa := NewFontAtlas(face, rune runeSets ...[]rune) *FontAtlas {

}
