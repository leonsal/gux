//go:build !vulkan

package gb

// #cgo CFLAGS: -I../libs/src
// #cgo linux   LDFLAGS: -L../libs -lguxgl -lglfw3 -lm -ldl -lX11
// #cgo windows LDFLAGS: -L../libs -lguxgl -lglfw3 -lgdi32 -limm32
import "C"
