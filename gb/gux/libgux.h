#pragma once
#include <stdbool.h>    
#include <stdint.h> // intptr_t

// Opaque backend window pointer
typedef void* gb_window_t;

// Vector with 2 components
typedef struct gb_vec2 {
    float x;
    float y;
} gb_vec2_t;

// Vector with 3 components
typedef struct gb_vec3 {
    float x;
    float y;
    float z;
} gb_vec3_t;

// Vector with 4 components
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
    gb_vec2_t  uv;                  // Vertex texture coordinates
    gb_rgba_t col;                  // Vertex color
} gb_vertex_t;

// Type for draw buffer index
typedef uint32_t gb_index_t;

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
	gb_index_t*     buf_idx;        // Indices buffer
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

// Frame parameters for gb_window_start_frame()
typedef struct gb_frame_params {
    float           ev_timeout;     // Event timeout in seconds
    gb_vec4_t       clear_color;    // Window clear color
} gb_frame_params_t;

// Frame information returned by gb_window_start_frame()
typedef struct gb_frame_info {
    uint32_t        win_close;      // Window close request
    gb_vec2_t       win_size;       // Window size
    gb_vec2_t       fb_size;        // Framebuffer size
    gb_vec2_t       fb_scale;       // Framebuffer scale
    uint32_t        ev_cap;         // Event buffer current capacity
    uint32_t        ev_count;       // Number of frame events in the following array
    gb_event_t*     events;         // Pointer to array of events
} gb_frame_info_t;

// Graphics backend configuration 
typedef struct gb_config_open_t {
    bool es;
} gb_config_opengl_t;

typedef struct gb_config_vulkan {
    bool validation_layer;
} gb_config_vulkan_t;

typedef struct gb_config {
    bool		debug_print_cmds;   // Print draw commands for debugging
	bool		unlimited_rate;		// Unlimited frame rate if true
    gb_config_opengl_t  opengl;     // OpenGL configuration
    gb_config_vulkan_t  vulkan;     // Vulkan configuration
} gb_config_t;

// Cursor types
enum {
    CURSOR_DEFAULT,
    CURSOR_ARROW,
    CURSOR_IBEAM,
    CURSOR_CROSSHAIR,
    CURSOR_HAND,
    CURSOR_HRESIZE,
    CURSOR_VRESIZE,
    _CURSOR_COUNT,
};

// Event types
enum {
    EVENT_KEY,                      // Key input event
    EVENT_CHAR,                     // Character input event
    EVENT_CURSOR_POS,               // Cursor position change event
    EVENT_CURSOR_ENTER,             // Cursor enter/exit event
    EVENT_MOUSE_BUTTON,             // Mouse button event
    EVENT_SCROLL,                   // Scroll event (mouse wheel)
};

// Public API
gb_window_t gb_create_window(const char* title, int width, int height, gb_config_t* pcfg);
void gb_window_destroy(gb_window_t win);
gb_frame_info_t* gb_window_start_frame(gb_window_t bw, gb_frame_params_t* params);
void gb_window_render_frame(gb_window_t win, gb_draw_list_t dl);
void gb_set_cursor(gb_window_t win, int cursor);
gb_texid_t gb_create_texture(gb_window_t win, int width, int height, const gb_rgba_t* data);
void gb_delete_texture(gb_window_t win, gb_texid_t texid);


