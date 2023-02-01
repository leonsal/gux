package view

import (
	"github.com/leonsal/gux/gb"
	"github.com/leonsal/gux/window"
)

type Label struct {
	View
	text string
	// Style ??
	//ff     app.FontFamilyType
	color  gb.RGBA
	valign window.TextVAlign
}

func NewLabel(text string) *Label {

	l := new(Label)
	l.SetText(text)
	return l
}

func (l *Label) SetText(text string) {

	l.text = text
	//func (w *Window) AddText(dl *gb.DrawList, fa *FontAtlas, pos *gb.Vec2, color gb.RGBA, align TextVAlign, text string) {
}

func (l *Label) Text() string {

	return l.text
}

func (l *Label) Render(w *window.Window) {

	if !l.visible {
		return
	}
	//func (w *Window) AddText(dl *gb.DrawList, fa *FontAtlas, pos *gb.Vec2, color gb.RGBA, align TextVAlign, text string) {

}
