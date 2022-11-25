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

// DrawCmd specifies a single draw command
type DrawCmd struct {
	ClipRect Vec4      // Clip rectangle
	TexId    int       // Texture ID
	Indices  []uint32  // Array of vertices indices
	Vertices []float32 // Array of vertices positions
}

// DrawList contains a list of commands for the graphics backend
type DrawList struct {
	bufCmd []C.gb_draw_cmd_t // Buffer of draw commands
	bufIdx []uint32          // Buffer of vertices indices
	bufVtx []float32         // Buffer of vertices positions
}

// NewDrawList creates and returns an empty DrawList
func NewDrawList() *DrawList {

	dl := new(DrawList)
	return dl
}

// AddCmd appends a new command to the Draw List
func (dl *DrawList) AddCmd(cmd DrawCmd) {

	cc := C.gb_draw_cmd_t{
		clip_rect:  C.gb_vec4_t{C.float(cmd.ClipRect.X), C.float(cmd.ClipRect.Y), C.float(cmd.ClipRect.Z), C.float(cmd.ClipRect.W)},
		texid:      C.int(cmd.TexId),
		idx_offset: C.int(len(dl.bufIdx)),
		vtx_offset: C.int(len(dl.bufVtx)),
		elem_count: C.int(len(cmd.Indices)),
	}
	dl.bufCmd = append(dl.bufCmd, cc)
	dl.bufIdx = append(dl.bufIdx, cmd.Indices...)
	dl.bufVtx = append(dl.bufVtx, cmd.Vertices...)
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
