package view

import (
	"github.com/leonsal/gux/window"
)

type Label struct {
	View
	text string
	// Style ??
	ff     window.FontStyleType
	valign window.TextVAlign
}

func NewLabel(text string) *Label {

	l := new(Label)
	l.Init(l)
	//l.pos = gb.Vec2{100, 100}
	l.ff = window.FontRegular
	l.valign = window.TextVAlignTop
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
	color := l.StyleColor(w, StyleColorText).RGBA()
	pos := l.pos
	w.AddText(w.DrawList(), w.Font(l.ff, 0), &pos, color, l.valign, l.text)
}
