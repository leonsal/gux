package gux

type Widget struct {
}

type IWidget interface {
	Render(*Window)
}
