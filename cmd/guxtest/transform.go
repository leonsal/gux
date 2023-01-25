package main

import (
	"math"

	"github.com/leonsal/gux/gb"
	"github.com/leonsal/gux/window"
)

func init() {

	registerTest("transform", 3, newTestTransform)
}

type testTransform struct {
	g1 gb.DrawList
	g2 gb.DrawList
	g3 gb.DrawList
}

func newTestTransform(win *window.Window) ITest {

	t := new(testTransform)

	red := gb.MakeColor(255, 0, 0, 255)
	green := gb.MakeColor(0, 255, 0, 255)
	blue := gb.MakeColor(0, 0, 255, 255)

	// First group
	_, bufIdx, bufVtx := win.NewDrawCmd(&t.g1, 3, 3)
	bufVtx[0] = gb.Vertex{Pos: gb.Vec2{0, -100}, Col: red}
	bufVtx[1] = gb.Vertex{Pos: gb.Vec2{-100, 100}, Col: red}
	bufVtx[2] = gb.Vertex{Pos: gb.Vec2{100, 100}, Col: red}
	bufIdx[0] = 0
	bufIdx[1] = 1
	bufIdx[2] = 2

	// Second group
	_, bufIdx, bufVtx = win.NewDrawCmd(&t.g2, 3, 3)
	bufVtx[0] = gb.Vertex{Pos: gb.Vec2{0, 0}, Col: green}
	bufVtx[1] = gb.Vertex{Pos: gb.Vec2{0, 200}, Col: green}
	bufVtx[2] = gb.Vertex{Pos: gb.Vec2{200, 0}, Col: green}
	bufIdx[0] = 0
	bufIdx[1] = 1
	bufIdx[2] = 2

	// Third group
	_, bufIdx, bufVtx = win.NewDrawCmd(&t.g3, 6, 4)
	bufVtx[0] = gb.Vertex{Pos: gb.Vec2{-100, -100}, Col: blue}
	bufVtx[1] = gb.Vertex{Pos: gb.Vec2{-100, 100}, Col: blue}
	bufVtx[2] = gb.Vertex{Pos: gb.Vec2{100, 100}, Col: blue}
	bufVtx[3] = gb.Vertex{Pos: gb.Vec2{100, -100}, Col: blue}
	bufIdx[0] = 0
	bufIdx[1] = 1
	bufIdx[2] = 2
	bufIdx[3] = 2
	bufIdx[4] = 3
	bufIdx[5] = 0

	return t
}

func (t *testTransform) draw(win *window.Window) {

	dl := win.DrawList()
	var mat gb.Mat3

	deltaX := float32(210)
	for i := 1; i < 8; i++ {
		sf := 1.0 - float32(i)/10
		mat.SetTranslation(deltaX*float32(i), 100).Rotate(float32(i-1)*float32(math.Pi/16)).Scale(sf, sf)
		dl.AddList2(&t.g1, &mat)
	}

	deltaY := float32(300)
	for i := 1; i < 8; i++ {
		sf := 1.0 - float32(i)/10
		mat.SetTranslation(deltaX*float32(i), deltaY).Rotate(float32(i-1)*float32(math.Pi/16)).Scale(sf, sf)
		dl.AddList2(&t.g2, &mat)
	}

	deltaY += 300
	for i := 1; i < 8; i++ {
		sf := 1.0 - float32(i)/10
		mat.SetTranslation(deltaX*float32(i), deltaY).Rotate(float32(i-1)*float32(math.Pi/16)).Scale(sf, sf)
		dl.AddList2(&t.g3, &mat)
	}
}

func (t *testTransform) destroy(*window.Window) {

}
