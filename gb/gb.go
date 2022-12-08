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
type Color uint32

type TextureId uintptr

// Mask for Color alpha
const ColorMaskA Color = 0xFF_00_00_00

// Vertex specifies information about a single vertex
type Vertex struct {
	Pos Vec2  // Position in screen coordinates
	UV  Vec2  // Texture coordinates (only relevant if texture used)
	Col Color // RGBA packed color
}

// DrawCmd specifies a single draw command
type DrawCmd struct {
	ClipRect  Vec4      // Clip rectangle
	TexId     TextureId // Texture ID
	idxOffset uint32    // Start offset in index buffer
	vtxOffset uint32    // Start offset in vertex buffer
	elemCount uint32    // Number of indices
}

// DrawList contains lists of commands and buffers for the graphics backend
type DrawList struct {
	bufCmd []DrawCmd // Buffer with draw commands
	bufIdx []uint32  // Buffer with vertices indices
	bufVtx []Vertex  // Buffer with vertices info
	Path   []Vec2    // Path being built
}

type Event struct {
	Type     uint32
	ArgInt   [4]int32
	ArgFloat [2]float32
}

func MakeColor(r, g, b, a byte) Color {

	return Color(uint32(a)<<24 | uint32(b)<<16 | uint32(g)<<8 | uint32(r))
}

// ReserveCmd creates and appends a DrawCmd into the DrawList
// reserving space for the specified number of indices and vertices.
// Returns pointer to the command and slices for direct access to the indices and vertices.
// After the indices were set (starting from 0) 'AdjustIdx()' must be called to adjust
// the indices considering the command idx offset.
func (dl *DrawList) ReserveCmd(idxCount, vtxCount int) (*DrawCmd, []uint32, []Vertex) {

	// Reserve space for indices
	idxOffset := len(dl.bufIdx)
	for i := 0; i < idxCount; i++ {
		dl.bufIdx = append(dl.bufIdx, 0)
	}

	// Reserve space for vertices
	vtxOffset := len(dl.bufVtx)
	for i := 0; i < vtxCount; i++ {
		dl.bufVtx = append(dl.bufVtx, Vertex{})
	}

	// Reserve command
	cmd := DrawCmd{
		ClipRect:  Vec4{},
		TexId:     1, // First texture allocated: white pixel
		idxOffset: uint32(idxOffset),
		vtxOffset: uint32(vtxOffset),
		elemCount: uint32(idxCount),
	}
	dl.bufCmd = append(dl.bufCmd, cmd)
	return &dl.bufCmd[len(dl.bufCmd)-1], dl.bufIdx[idxOffset : idxOffset+idxCount], dl.bufVtx[vtxOffset : vtxOffset+vtxCount]
}

// AdjustIdx must be called with the DrawCmd pointer returned by ReserveCmd() to adjust the indices buffers
func (dl *DrawList) AdjustIdx(cmd *DrawCmd) {

	for i := 0; i < int(cmd.elemCount); i++ {
		dl.bufIdx[i+int(cmd.idxOffset)] += cmd.vtxOffset
	}
}

// AddCmd appends a new command to the Draw List
func (dl *DrawList) AddCmd(clipRect Vec4, texId TextureId, indices []uint32, vertices []Vertex) {

	cmd, idx, vtx := dl.ReserveCmd(len(indices), len(vertices))
	copy(idx, indices)
	copy(vtx, vertices)
	cmd.ClipRect = clipRect
	cmd.TexId = texId
	dl.AdjustIdx(cmd)
}

// AddList appends the specified DrawList to this one
func (dl *DrawList) AddList(src DrawList) {

	// Append vertices
	vtxOffset := len(dl.bufVtx)
	dl.bufVtx = append(dl.bufVtx, src.bufVtx...)

	// Append indices adjusting offset
	idxOffset := len(dl.bufIdx)
	for _, idx := range src.bufIdx {
		dl.bufIdx = append(dl.bufIdx, idx+uint32(vtxOffset))
	}

	// Append commands adjusting offsets
	for i := 0; i < len(src.bufCmd); i++ {
		cmd := &src.bufCmd[i]
		cmd.idxOffset += uint32(idxOffset)
		cmd.vtxOffset += uint32(vtxOffset)
		dl.bufCmd = append(dl.bufCmd, *cmd)
	}
}

// Clear clears the DrawList commands, indices and vertices buffer without deallocating memory
func (dl *DrawList) Clear() {

	dl.bufCmd = dl.bufCmd[:0]
	dl.bufIdx = dl.bufIdx[:0]
	dl.bufVtx = dl.bufVtx[:0]
	dl.Path = dl.Path[:0]
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
		cdl.cmd_count = C.uint(len(dl.bufCmd))
		cdl.buf_idx = (*C.uint)(unsafe.Pointer(&dl.bufIdx[0]))
		cdl.idx_count = C.uint(len(dl.bufIdx))
		cdl.buf_vtx = (*C.gb_vertex_t)(unsafe.Pointer(&dl.bufVtx[0]))
		cdl.vtx_count = C.uint(len(dl.bufVtx))
	}
	C.gb_window_render_frame(w.c, cdl)
}

// CreateTexture creates an empty texture and returns its ID
func (w *Window) CreateTexture() TextureId {

	return TextureId(C.gb_create_texture())
}

// DeleteTexture deletes the specified texture
func (w *Window) DeleteTexture(texid TextureId) {

	C.gb_delete_texture(C.gb_texid_t(texid))
}

// TransferTexture transfers data to the texture
func (w *Window) TransferTexture(texid TextureId, width, height int, data *Color) {

	C.gb_transfer_texture(C.gb_texid_t(texid), C.int(width), C.int(height), (*C.gb_color_t)(data))
}

func (w *Window) GetEvents(events []Event) int {

	count := C.gb_get_events(w.c, (*C.gb_event_t)(unsafe.Pointer(&events[0])), C.int(len(events)))
	return int(count)
}
