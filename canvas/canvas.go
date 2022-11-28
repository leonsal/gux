package canvas

import "github.com/leonsal/gux/gb"

type Canvas struct {
	DrawList gb.DrawList
}

func New() *Canvas {

	c := new(Canvas)
	return c
}
