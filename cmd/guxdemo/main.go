package main

import (
	"log"

	"github.com/leonsal/gux/app"
	"github.com/leonsal/gux/gb"
)

func main() {

	a := app.Get()

	// First Window
	w1, err := a.NewWindow("AppWin1", 800, 600)
	if err != nil {
		log.Fatal(err)
	}
	w1.SetClearColor(gb.Vec4{0, 1, 0, 1})
	err = w1.CreateFontManager()
	if err != nil {
		log.Fatal(err)
	}

	// Secondt Window
	w2, err := a.NewWindow("AppWin2", 400, 200)
	if err != nil {
		log.Fatal(err)
	}
	w2.SetClearColor(gb.Vec4{0.5, 0.5, 0.5, 1})
	err = w2.CreateFontManager()
	if err != nil {
		log.Fatal(err)
	}

	for !a.Render() {
	}
	//a.Close()
}
