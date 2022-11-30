package canvas

import "github.com/leonsal/gux/gb"

type Canvas struct {
	DrawList gb.DrawList // DrawList with all the canvas...
	bufVec2  []gb.Vec2   // Temporary Vec2 buffer used by drawing functions (to avoid allocations)
}

type Flags int

const (
	Flag_Closed = (1 << iota)
)

func New() *Canvas {

	c := new(Canvas)
	return c
}

// ReserveVec2 reserves 'count' gb.Vec2 entries in internal Vec2 buffer
// returning a slice to access these entries
func (c *Canvas) ReserveVec2(count int) []gb.Vec2 {

	idx := len(c.bufVec2)
	for i := 0; i < count; i++ {
		c.bufVec2 = append(c.bufVec2, gb.Vec2{})
	}
	return c.bufVec2[idx : idx+count]
}
