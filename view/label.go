package view

import (
	"github.com/leonsal/gux/gb"
	"github.com/leonsal/gux/window"
)

type Label struct {
	View
	text string
	// Style ??
	ff     window.FontFamilyType
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
}

func (l *Label) Text() string {

	return l.text
}

func (l *Label) Render(w *window.Window) {

	if !l.visible {
		return
	}
	dl := w.DrawList()
	w.AddText(dl, w.Font(l.ff, 0), &l.pos, l.color, l.valign, l.text)

}
