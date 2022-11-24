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

// drawCmd specifies a single graphics backend draw command
// All fields are 'float32' to facilitate the transfer to C function
type drawCmd struct {
	clipRect  Vec4    // Clip rectangle
	textId    float32 // Texture Id (integer value)
	vtxOffset float32 // Vertex offset in 'bufVtx' (integer value)
	idxOffset float32 // Index offset in 'bufIdx' (integer value)
	elemCount float32 // Number of elements (integer value)
}

// DrawList contains a list of commands for the graphics backend
type DrawList struct {
	bufCmd []drawCmd // Draw commands buffer
	bufIdx []uint32  // Vertices indices buffer
	bufVtx []float32 // Vertices positions buffer
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

// Clear clears the DrawList commands, indices and vertices buffer withou deallocating memory
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

	// Builds C draw list struct and calls backend render
	var cdl C.gb_draw_list_t
	cdl.bufCmd = (*C.gb_draw_cmd_t)(unsafe.Pointer(&dl.bufCmd[0]))
	cdl.cmd_count = C.int(len(dl.bufCmd))
	cdl.bufIdx = (*C.int)(unsafe.Pointer(&dl.bufIdx[0]))
	cdl.idx_count = C.int(len(dl.bufIdx))
	cdl.bufVtx = (*C.float)(unsafe.Pointer(&dl.bufVtx[0]))
	cdl.vtx_count = C.int(len(dl.bufVtx))
	C.gb_window_render_frame(w.c, cdl)
}
