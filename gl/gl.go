package gl

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

//    gl_vec4_t       clip_rect;  // Clip rectangle
//    unsigned int    texid;      // Texture id
//    unsigned int    vtx_offset; // Start offset in vertex buffer
//    unsigned int    idx_offset; // Start offset in index buffer
//    unsigned int    elem_count; // Number of indices

type DrawCmd struct {
	clipRect  Vec4
	textId    uint32
	vtxOffset uint32
	idxOffset uint32
	elemCount uint32
}

type DrawList struct {
	bufCmd   *DrawCmd
	cmdCount uint32
	idxBuf   *uint32
	vtxBuf   *uint32
}

type Window struct {
	c C.gl_window_t
}

func CreateWindow(title string, width, height int) (*Window, error) {

	ctitle := C.CString(title)
	defer C.free(unsafe.Pointer(ctitle))

	cw := C.gl_create_window(ctitle, C.int(width), C.int(height), nil)
	if cw == nil {
		return nil, errors.New("error creating window")
	}
	return &Window{cw}, nil
}

func (w *Window) Destroy() {

	C.gl_window_destroy(w.c)
}

func (w *Window) StartFrame(timeout float64) bool {

	return bool(C.gl_window_start_frame(w.c, C.double(timeout)))
}

func (w *Window) RenderFrame(dl DrawList) {

	C.gl_window_render_frame(w.c)
}
