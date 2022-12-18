#pragma once
#include <stdbool.h>    
#include <stdint.h> // intptr_t

// Graphics backend configuration 
typedef struct gb_config {
    struct {
        bool es;
        int  msaa;
        struct {
            float r; float g; float b; float a;
        } clearColor;
    } opengl;
} gb_config_t;

// Opaque backend window pointer
typedef void* gb_window_t;

typedef struct gb_vec2 {
    float x;
    float y;
} gb_vec2_t;

typedef struct gb_vec3 {
    float x;
    float y;
    float z;
} gb_vec3_t;

typedef struct gb_vec4 {
    float x;
    float y;
    float z;
    float w;
} gb_vec4_t;

// Packed color containg RGBA components each as an unsigned byte
typedef uint32_t gb_rgba_t;

// Texture id
typedef intptr_t gb_texid_t;

// Vertex info
typedef struct gb_vertex {
    gb_vec2_t  pos;                 // Vertex position in screen coordinates
    gb_vec2_t  uv;                  // Texture coordinates
    gb_rgba_t col;                 // Color as an uint32
} gb_vertex_t;

// Single draw command
typedef struct gb_draw_cmd {
    gb_vec4_t       clip_rect;      // Clip rectangle
    gb_texid_t      texid;          // Texture id
    uint32_t        idx_offset;     // Start offset in index buffer
    uint32_t        vtx_offset;     // Start offset in vertex buffer
    uint32_t        elem_count;     // Number of indices
} gb_draw_cmd_t;

// List of draw commands and buffers of vertices indices/positions
typedef struct gb_draw_list {
	gb_draw_cmd_t*  buf_cmd;        // Draw command buffer
    uint32_t        cmd_count;      // Total number of commands
	uint32_t*       buf_idx;        // Indices buffer
    uint32_t        idx_count;      // Total number of indices
	gb_vertex_t*    buf_vtx;        // Vertices info buffer
    uint32_t        vtx_count;      // Total number of vertices
} gb_draw_list_t;

// Single generic event
typedef struct gb_event {
    uint32_t        type;           // Event type
    int32_t         argint[4];      // Integer arguments
    float           argfloat[2];    // Float parameters
} gb_event_t;

// Event types
enum {
    EVENT_KEY,                      // Key input event
    EVENT_CHAR,                     // Character input event
    EVENT_CURSOR_POS,               // Cursor position change event
    EVENT_CURSOR_ENTER,             // Cursor enter/exit event
    EVENT_MOUSE_BUTTON,             // Mouse button event
    EVENT_SCROLL,                   // Scroll event (mouse wheel)
};

gb_window_t gb_create_window(const char* title, int width, int height, gb_config_t* cfg);
void gb_window_destroy(gb_window_t win);
bool gb_window_start_frame(gb_window_t bw, double timeout);
void gb_window_render_frame(gb_window_t win, gb_draw_list_t dl);
gb_texid_t gb_create_texture();
void gb_delete_texture(gb_texid_t texid);
void gb_transfer_texture(gb_texid_t texid, int width, int height, const gb_rgba_t* data);
int gb_get_events(gb_window_t win, gb_event_t* events, int ev_count);

