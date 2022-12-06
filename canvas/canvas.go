package canvas

//
//import (
//	"github.com/leonsal/gux"
//	"github.com/leonsal/gux/gb"
//)
//
//// Canvas is a View to which draw commands can be be added
//type Canvas struct {
//	gux.View             // Embedded base View
//	w        *gux.Window // Native window
//	dl       gb.DrawList // DrawList with all the canvas...
//	bufVec2  []gb.Vec2   // Temporary Vec2 buffer used by drawing functions (to avoid allocations)
//}
//
//type Flags int
//
//const (
//	Flag_Closed Flags = (1 << iota)
//	Flag_AntiAliasedFill
//)
//
//func New(w *gux.Window) *Canvas {
//
//	c := new(Canvas)
//	c.Init(w)
//	return c
//}
//
//func (c *Canvas) Init(w *gux.Window) {
//
//	c.View.Init(c)
//	c.w = w
//	c.SetRender(func(w *gux.Window) {
//
//		// Add canvas draw list to window draw list
//		w.AddList(c.dl)
//
//		// Clear this canvas draw list and aux buffer (without deallocating memory)
//		c.dl.Clear()
//		c.bufVec2 = c.bufVec2[:0]
//	})
//}
//
//// ReserveVec2 reserves 'count' gb.Vec2 entries in internal Vec2 buffer
//// returning a slice to access these entries
//func (c *Canvas) ReserveVec2(count int) []gb.Vec2 {
//
//	idx := len(c.bufVec2)
//	for i := 0; i < count; i++ {
//		c.bufVec2 = append(c.bufVec2, gb.Vec2{})
//	}
//	return c.bufVec2[idx : idx+count]
//}
