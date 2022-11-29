package canvas

import "github.com/leonsal/gux/gb"

type Canvas struct {
	DrawList gb.DrawList
}

type Flags int

const (
	Flag_Closed = (1 << iota)
)

func New() *Canvas {

	c := new(Canvas)
	return c
}
