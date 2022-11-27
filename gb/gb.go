package gb

// #include <stdlib.h>
// #include "libgux.h"
import "C"
import (
	"errors"
	"unsafe"
)

type Vec2 struct {
	X float32
	Y float32
}

type Vec3 struct {
	X float32
	Y float32
	Z float32
}
type Vec4 struct {
	X float32
	Y float32
	Z float32
	W float32
}

// Packed RGBA color from LSB to MSB
type RGBA32 uint32

// Vertex specifies information about a single vertex
type Vertex struct {
	Pos Vec2   // Position in screen coordinates
	UV  Vec2   // Texture coordinates (only relevant if texture used)
	Col RGBA32 // RGBA packed color
}

// DrawCmd specifies a single draw command
type DrawCmd struct {
	ClipRect Vec4     // Clip rectangle
	TexId    int      // Texture ID
	Indices  []uint32 // Array of vertices indices
	Vertices []Vertex // Array of vertices info
}

// AddIndices adds the specified indices elements to the draw command
func (cmd *DrawCmd) AddIndices(indices ...uint32) {

	for _, idx := range indices {
		cmd.Indices = append(cmd.Indices, idx)
	}
}

// AddVertices adds the specified vertices to the draw command
func (cmd *DrawCmd) AddVertices(vertices ...Vertex) {

	for _, vtx := range vertices {
		cmd.Vertices = append(cmd.Vertices, vtx)
	}
}

// DrawList contains a list of commands for the graphics backend
type DrawList struct {
	bufCmd []C.gb_draw_cmd_t // Buffer of draw commands
	bufIdx []uint32          // Buffer of vertices indices
	bufVtx []C.gb_vertex_t   // Buffer of vertices info
}

// AddCmd appends a new command to the Draw List
func (dl *DrawList) AddCmd(cmd DrawCmd) {

	// Convert command to C struct and appends to commands buffer
	cc := C.gb_draw_cmd_t{
		clip_rect:  C.gb_vec4_t{C.float(cmd.ClipRect.X), C.float(cmd.ClipRect.Y), C.float(cmd.ClipRect.Z), C.float(cmd.ClipRect.W)},
		texid:      C.int(cmd.TexId),
		idx_offset: C.int(len(dl.bufIdx)),
		vtx_offset: C.int(len(dl.bufVtx)),
		elem_count: C.int(len(cmd.Indices)),
	}
	dl.bufCmd = append(dl.bufCmd, cc)

	// Appends command indices to indices buffer
	idxOffset := uint32(len(dl.bufIdx))
	for i := range cmd.Indices {
		idx := cmd.Indices[i]
		dl.bufIdx = append(dl.bufIdx, idxOffset+idx)
	}

	// Convert vertex info to C struct and appends to vertices buffer
	for i := range cmd.Vertices {
		v := &cmd.Vertices[i]
		dl.bufVtx = append(dl.bufVtx, C.gb_vertex_t{
			C.gb_vec2_t{C.float(v.Pos.X), C.float(v.Pos.Y)},
			C.gb_vec2_t{C.float(v.UV.X), C.float(v.UV.Y)},
			C.int(v.Col),
		})
	}
}

// AddList appends the specified DrawList to this one
func (dl *DrawList) AddList(src DrawList) {
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

	// Builds C draw list struct and calls backend render
	var cdl C.gb_draw_list_t
	if len(dl.bufCmd) > 0 {
		cdl.buf_cmd = (*C.gb_draw_cmd_t)(unsafe.Pointer(&dl.bufCmd[0]))
		cdl.cmd_count = C.int(len(dl.bufCmd))
		cdl.buf_idx = (*C.uint)(unsafe.Pointer(&dl.bufIdx[0]))
		cdl.idx_count = C.int(len(dl.bufIdx))
		cdl.buf_vtx = (*C.gb_vertex_t)(unsafe.Pointer(&dl.bufVtx[0]))
		cdl.vtx_count = C.int(len(dl.bufVtx))
	}
	C.gb_window_render_frame(w.c, cdl)
}
