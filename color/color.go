package color

import "github.com/leonsal/gux/gb"

type Color struct {
	R float32
	G float32
	B float32
	A float32
}

func (c Color) RGBA() gb.RGBA {

	return gb.MakeColor(byte(c.R*255.0), byte(c.G*255.0), byte(c.B*255.0), byte(c.A*255.0))
}

func ColorFromRGBA(rgba gb.RGBA) Color {

	return Color{
		R: float32(((rgba >> gb.RGBAShiftR) & 0xFF) * 1.0 / 255.0),
		G: float32(((rgba >> gb.RGBAShiftG) & 0xFF) * 1.0 / 255.0),
		B: float32(((rgba >> gb.RGBAShiftB) & 0xFF) * 1.0 / 255.0),
		A: float32(((rgba >> gb.RGBAShiftA) & 0xFF) * 1.0 / 255.0),
	}
}
