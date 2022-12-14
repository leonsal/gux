//go:build vulkan

package gb

// #cgo CFLAGS: -I../libs/src
// #cgo linux   LDFLAGS: -L../libs -lguxvk -lglfw3 -lm -ldl -lX11
// #cgo windows LDFLAGS: -L../libs -lguxvk -lglfw3 -lgdi32 -limm32
import "C"
