package gux

import "github.com/leonsal/gux/gb"

type TextVAlign int

const (
	TextVAlignTop    TextVAlign = 0
	TextVAlignBase   TextVAlign = 1
	TextVAlignBottom TextVAlign = 2
)

func (w *Window) AddText(dl *gb.DrawList, fa *FontAtlas, pos gb.Vec2, align TextVAlign, text string) {

}
