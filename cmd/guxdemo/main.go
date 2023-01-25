package main

import (
	"fmt"
	"log"

	"github.com/leonsal/gux/app"
)

func main() {

	a := app.Init()

	w1, err := a.NewWindow("AppWin1", 800, 600)
	fmt.Println("window1 ", w1)
	if err != nil {
		log.Fatal(err)
	}

	w2, err := a.NewWindow("AppWin2", 200, 100)
	fmt.Println("window2 ", w2)
	if err != nil {
		log.Fatal(err)
	}

	//for {
	//	w1.StartFrame()
	//	w1.RenderFrame()
	//	w2.StartFrame()
	//	w2.RenderFrame()
	//}
	for {
		for _, aw := range a.Windows {
			aw.StartFrame()
			aw.RenderFrame()
		}
	}
}
