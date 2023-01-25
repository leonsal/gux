package main

import (
	"fmt"
	"log"

	"github.com/leonsal/gux/app"
)

func main() {

	a := app.Init()
	w1, err := a.NewWindow("AppWin1", 800, 600)
	fmt.Println("window1", w1)
	if err != nil {
		log.Fatal(err)
	}
	for a.Render() {

	}

}
