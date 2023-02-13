package main

import (
	"log"

	"github.com/leonsal/gux/app"
	"github.com/leonsal/gux/color"
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

	group := view.NewGroup()
	l1 := view.NewLabel("This is a label 1")
	l1.SetPos(100, 100)
	l1.SetStyleColor(view.StyleColorText, color.Darkred)

	l2 := view.NewLabel("This is a label 2")
	l2.SetPos(100, 200)
	group.Add(l1)
	group.Add(l2)
	group.SetPos(200, 200)

	a.SetView(w1, group)

	// Second Window
	w2, err := a.NewWindow("AppWin2", 400, 200)
	if err != nil {
		log.Fatal(err)
	}
	w2.SetClearColor(gb.Vec4{0.5, 0.5, 0.5, 1})
	a.SetView(w2, group)

	for !a.Render() {
	}
	a.Close()
}
