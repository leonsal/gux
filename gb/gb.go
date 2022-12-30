package gb

// #include <stdlib.h>
// #include "libgux.h"
import "C"
import (
	"errors"
	"math"
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

// RGBA is a packed color
type RGBA uint32

// Bit masks for RGBA color components
const RGBAMaskR RGBA = 0x00_00_00_FF
const RGBAMaskG RGBA = 0x00_00_FF_00
const RGBAMaskB RGBA = 0x00_FF_00_00
const RGBAMaskA RGBA = 0xFF_00_00_00

// Bit shifts for RGBA color components
const RGBAShiftR = 0
const RGBAShiftG = 8
const RGBAShiftB = 16
const RGBAShiftA = 24

// TextureId is the type for textures identifiers
type TextureID uintptr

// Vertex specifies information about a single vertex
type Vertex struct {
	Pos Vec2 // Position in screen coordinates
	UV  Vec2 // Texture coordinates (only relevant if texture used)
	Col RGBA // RGBA packed color
}

// DrawCmd specifies a single draw command
type DrawCmd struct {
	ClipRect  Vec4      // Clip rectangle
	TexID     TextureID // Texture ID
	idxOffset uint32    // Start offset in index buffer
	vtxOffset uint32    // Start offset in vertex buffer
	elemCount uint32    // Number of indices
}

// DrawList contains lists of commands and buffers for the graphics backend
type DrawList struct {
	bufCmd []DrawCmd // Buffer with draw commands
	bufIdx []uint32  // Buffer with vertices indices
	bufVtx []Vertex  // Buffer with vertices info
	Path   []Vec2    // Temporary list of path points
}

// Event describes an I/O event
type Event struct {
	Type     uint32     // Event type
	ArgInt   [4]int32   // Signed integer arguments
	ArgFloat [2]float32 // Float arguments
}

// Frame parameters contains parameters for start_frame()
type FrameParams struct {
	EvTimeout  float32 // Event timeout in seconds
	ClearColor Vec4    // Window clear color
}

// FrameInfo contains frame information returned by start_frame()
type FrameInfo struct {
	WinClose bool    // Window close request
	WinSize  Vec2    // Window size
	FbSize   Vec2    // Framebuffer size
	FbScale  Vec2    // Framebuffer scale
	Events   []Event // Events array
}

// Graphics backend configuration
type ConfigOpenGL struct {
	ES bool // Use OpenGL ES3.0 instead of OpenGL 3.3
}

type ConfigVulkan struct {
	ValidationLayer bool // Enable Vulkan debug validation layer
}

type Config struct {
	DebugPrintCmds bool
	OpenGL         ConfigOpenGL
	Vulkan         ConfigVulkan
}

// MakeColor makes and returns an RGBA packed color from the specified components
func MakeColor(r, g, b, a byte) RGBA {

	return RGBA(uint32(a)<<24 | uint32(b)<<16 | uint32(g)<<8 | uint32(r))
}

// NewDrawCmd creates and appends a DrawCmd into the DrawList
// reserving space for the specified number of indices and vertices.
// Returns pointer to the command and slices for direct access to the indices and vertices.
func (dl *DrawList) NewDrawCmd(idxCount, vtxCount int) (*DrawCmd, []uint32, []Vertex) {

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

	// Creates and appends new command to the DrawList command buffer
	cmd := DrawCmd{
		idxOffset: uint32(idxOffset),
		vtxOffset: uint32(vtxOffset),
		elemCount: uint32(idxCount),
	}
	dl.bufCmd = append(dl.bufCmd, cmd)
	return &dl.bufCmd[len(dl.bufCmd)-1], dl.bufIdx[idxOffset : idxOffset+idxCount], dl.bufVtx[vtxOffset : vtxOffset+vtxCount]
}

// AddCmd appends a new command to the Draw List
func (dl *DrawList) AddCmd(clipRect Vec4, texId TextureID, indices []uint32, vertices []Vertex) {

	cmd, idx, vtx := dl.NewDrawCmd(len(indices), len(vertices))
	copy(idx, indices)
	copy(vtx, vertices)
	cmd.ClipRect = clipRect
	cmd.TexID = texId
}

// AddList appends the specified DrawList to this one.
// The added DrawList is not modified.
func (dl *DrawList) AddList(src *DrawList) {

	// Append vertices
	vtxOffset := len(dl.bufVtx)
	dl.bufVtx = append(dl.bufVtx, src.bufVtx...)

	// Append indices
	idxOffset := len(dl.bufIdx)
	dl.bufIdx = append(dl.bufIdx, src.bufIdx...)

	// Append commands adjusting offsets
	for i := 0; i < len(src.bufCmd); i++ {
		cmd := src.bufCmd[i]
		cmd.idxOffset += uint32(idxOffset)
		cmd.vtxOffset += uint32(vtxOffset)
		dl.bufCmd = append(dl.bufCmd, cmd)
	}
}

// Clear clears the DrawList commands, indices and vertices buffer without deallocating memory
func (dl *DrawList) Clear() {

	dl.bufCmd = dl.bufCmd[:0]
	dl.bufIdx = dl.bufIdx[:0]
	dl.bufVtx = dl.bufVtx[:0]
	dl.Path = dl.Path[:0]
}

// PathReserve reserves spaces for n points for the DrawList Path slice
func (dl *DrawList) PathReserve(n int) {

	plen := len(dl.Path)
	pcap := cap(dl.Path)
	free := pcap - plen
	if n <= free {
		return
	}
	buf := make([]Vec2, plen, n-free)
	copy(buf, dl.Path)
	dl.Path = buf
}

// PathAppends appends a point to the DrawList Path
func (dl *DrawList) PathAppend(p Vec2) {

	dl.Path = append(dl.Path, p)
}

// PathClear clears the DrawList path without deallocation its memory
func (dl *DrawList) PathClear() {

	dl.Path = dl.Path[:0]
}

// Clone returns a copy of the DrawList
func (dl *DrawList) Clone() DrawList {

	dst := DrawList{}
	dst.bufCmd = make([]DrawCmd, len(dl.bufCmd))
	dst.bufIdx = make([]uint32, len(dl.bufIdx))
	dst.bufVtx = make([]Vertex, len(dl.bufVtx))
	copy(dst.bufCmd, dl.bufCmd)
	copy(dst.bufIdx, dl.bufIdx)
	copy(dst.bufVtx, dl.bufVtx)
	return dst
}

// Translate translates all vertices of the DrawList by the specified delta vector
func (dl *DrawList) Translate(delta Vec2) *DrawList {

	for i := 0; i < len(dl.bufVtx); i++ {
		dl.bufVtx[i].Pos.Add(delta)
	}
	return dl
}

// Scale scales all vertices of the DrawList by the specified scale vector
func (dl *DrawList) Scale(scale Vec2) *DrawList {

	for i := 0; i < len(dl.bufVtx); i++ {
		dl.bufVtx[i].Pos.X *= scale.X
		dl.bufVtx[i].Pos.Y *= scale.Y
	}
	return dl
}

func (dl *DrawList) Rotate(theta float32) *DrawList {

	sin := float32(math.Sin(float64(theta)))
	cos := float32(math.Cos(float64(theta)))
	for i := 0; i < len(dl.bufVtx); i++ {
		v := dl.bufVtx[i].Pos
		dl.bufVtx[i].Pos.X = v.X*cos - v.Y*sin
		dl.bufVtx[i].Pos.Y = v.X*sin + v.Y*cos
	}
	return dl
}

type Window struct {
	c C.gb_window_t
}

func CreateWindow(title string, width, height int, cfg *Config) (*Window, error) {

	var pcfg *C.gb_config_t
	if cfg != nil {
		ccfg := C.gb_config_t{}
		ccfg.debug_print_cmds = C.bool(cfg.DebugPrintCmds)
		ccfg.opengl.es = C.bool(cfg.OpenGL.ES)
		ccfg.vulkan.validation_layer = C.bool(cfg.Vulkan.ValidationLayer)
		pcfg = &ccfg
	}

	ctitle := C.CString(title)
	defer C.free(unsafe.Pointer(ctitle))

	cw := C.gb_create_window(ctitle, C.int(width), C.int(height), pcfg)
	if cw == nil {
		return nil, errors.New("error creating window")
	}
	return &Window{cw}, nil
}

func (w *Window) Destroy() {

	C.gb_window_destroy(w.c)
}

func (w *Window) StartFrame(params *FrameParams) FrameInfo {

	finfo := FrameInfo{}
	cframe := C.gb_window_start_frame(w.c, (*C.gb_frame_params_t)(unsafe.Pointer(params)))
	if cframe.win_close != 0 {
		finfo.WinClose = true
	}
	finfo.WinSize = Vec2{float32(cframe.win_size.x), float32(cframe.win_size.y)}
	finfo.FbSize = Vec2{float32(cframe.fb_size.x), float32(cframe.fb_size.y)}
	finfo.FbScale = Vec2{float32(cframe.fb_scale.x), float32(cframe.fb_scale.y)}
	finfo.Events = unsafe.Slice((*Event)(unsafe.Pointer(cframe.events)), cframe.ev_count)
	return finfo
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

// CreateTexture creates texture with the specified image data and returns the texture id.
func (w *Window) CreateTexture(width, height int, data *RGBA) TextureID {

	return TextureID(C.gb_create_texture(w.c, C.int(width), C.int(height), (*C.gb_rgba_t)(data)))
}

// DeleteTexture deletes the specified texture
func (w *Window) DeleteTexture(texid TextureID) {

	C.gb_delete_texture(w.c, C.gb_texid_t(texid))
}

//func (w *Window) GetEvents(events []Event) int {
//
//	count := C.gb_get_events(w.c, (*C.gb_event_t)(unsafe.Pointer(&events[0])), C.int(len(events)))
//	return int(count)
//}
