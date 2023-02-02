package main

import (
	"log"

	"github.com/leonsal/gux/app"
	"github.com/leonsal/gux/gb"
	"github.com/leonsal/gux/view"
)

func main() {

	a := app.Get()

	// First Window
	cfg := gb.Config{}
	cfg.DebugPrintCmds = false
	cfg.Vulkan.ValidationLayer = true
	w1, err := a.NewWindowEx("AppWin1", 1200, 800, &cfg)
	if err != nil {
		log.Fatal(err)
	}
	w1.SetClearColor(gb.Vec4{0.6, 0.6, 0.6, 1})

	v := view.NewLabel("This is a label jpq")
	a.SetView(w1, v)

	//// Second Window
	//w2, err := a.NewWindow("AppWin2", 400, 200)
	//if err != nil {
	//	log.Fatal(err)
	//}
	//w2.SetClearColor(gb.Vec4{0.5, 0.5, 0.5, 1})

	for !a.Render() {
	}
	a.Close()
}
