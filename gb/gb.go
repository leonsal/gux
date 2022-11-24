package gb

// #include <stdlib.h>
// #include "libgux.h"
import "C"
import (
	"errors"
	"unsafe"
)

type Vec4 struct {
	X float32
	Y float32
	Z float32
	W float32
}

type drawCmd struct {
	clipRect  Vec4
	textId    float32
	vtxOffset float32
	idxOffset float32
	elemCount float32
}

type DrawList struct {
	bufCmd []drawCmd // List of draw commands
	bufIdx []uint32  // List of vertices indices
	bufVtx []float32 // List of vertices positions
}

func NewDrawList() *DrawList {

	dl := new(DrawList)
	return dl
}

func (dl *DrawList) AddCmd(clipRect Vec4, texId uint, vtxOffset uint, idxOffset uint, elemCount uint) {

	cmd := drawCmd{
		clipRect, float32(texId), float32(vtxOffset), float32(idxOffset), float32(elemCount),
	}
	dl.bufCmd = append(dl.bufCmd, cmd)
	dl.bufIdx = append(dl.bufIdx, 1)   // just for test
	dl.bufVtx = append(dl.bufVtx, 1.0) // just for test
}

func (dl *DrawList) Clear() {

	dl.bufCmd = dl.bufCmd[:0]
	dl.bufIdx = dl.bufIdx[:0]
	dl.bufVtx = dl.bufVtx[:0]
}

type Window struct {
	c C.gb_window_t
}

func CreateWindow(title string, width, height int) (*Window, error) {

	ctitle := C.CString(title)
	defer C.free(unsafe.Pointer(ctitle))

	cw := C.gb_create_window(ctitle, C.int(width), C.int(height), nil)
	if cw == nil {
		return nil, errors.New("error creating window")
	}
	return &Window{cw}, nil
}

func (w *Window) Destroy() {

	C.gb_window_destroy(w.c)
}

func (w *Window) StartFrame(timeout float64) bool {

	return bool(C.gb_window_start_frame(w.c, C.double(timeout)))
}

func (w *Window) RenderFrame(dl *DrawList) {

	if len(dl.bufCmd) == 0 {
		return
	}

	C.gb_window_render_frame(w.c,
		(*C.gb_draw_cmd_t)(unsafe.Pointer(&dl.bufCmd[0])),
		C.int(len(dl.bufCmd)),
		(*C.int)(unsafe.Pointer(&dl.bufIdx[0])),
		(*C.float)(unsafe.Pointer(&dl.bufVtx[0])),
	)
}
