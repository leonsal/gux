//go:build !vulkan

package gl

// #cgo CFLAGS: -I../libs/src
// #cgo linux   LDFLAGS: -L../libs -lguxgl -lGL -lglfw3 -lm -ldl -lX11
// #cgo windows LDFLAGS: -L../libs -lguxgl -lglfw3 -lopengl32 -lgdi32 -limm32
import "C"
