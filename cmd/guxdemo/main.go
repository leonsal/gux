package main

import (
	"log"

	"github.com/leonsal/gux/app"
	"github.com/leonsal/gux/gb"
)

func main() {

	a := app.Init()

	w1, err := a.NewWindow("AppWin1", 800, 600)
	if err != nil {
		log.Fatal(err)
	}
	w1.SetClearColor(gb.Vec4{0, 1, 0, 1})

	w2, err := a.NewWindow("AppWin2", 200, 100)
	if err != nil {
		log.Fatal(err)
	}
	w2.SetClearColor(gb.Vec4{0.5, 0.5, 0.5, 1})

	//for {
	//	w1.StartFrame()
	//	w1.RenderFrame()
	//	w2.StartFrame()
	//	w2.RenderFrame()
	//}
	for len(a.Windows) > 0 {
		for i := 0; i < len(a.Windows); i++ {
			aw := a.Windows[i]
			//fmt.Println("i:", i, "aw:", aw)
			shouldClose := aw.StartFrame()
			if !shouldClose {
				a.Windows = append(a.Windows[:i], a.Windows[i+1:]...)
				aw.Close()
				continue
			}
			aw.RenderFrame()
		}
	}
}
