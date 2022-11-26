#pragma once
#include <stdbool.h>

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
    float x; float y;
} gb_vec2_t;

typedef struct gb_vec3 {
    float x; float y; float z;
} gb_vec3_t;

typedef struct gb_vec4 {
    float x; float y; float z; float w;
} gb_vec4_t;

// Packed color
typedef int gb_col32_t;

// Vertex info
typedef struct gb_vertex {
    gb_vec2_t  pos;             // Vertex position in screen coordinates
    gb_vec2_t  uv;              // Texture coordinates
    gb_col32_t col;             // Color as an int32
} gb_vertex_t;

// Single draw command
typedef struct gb_draw_cmd {
    gb_vec4_t   clip_rect;      // Clip rectangle
    int         texid;          // Texture id
    int         idx_offset;     // Start offset in index buffer
    int         vtx_offset;     // Start offset in vertex buffer
    int         elem_count;     // Number of indices
} gb_draw_cmd_t;

// List of draw commands and buffers of vertices indices/positions
typedef struct gb_draw_list {
	gb_draw_cmd_t*  buf_cmd;    // Draw command buffer
    int             cmd_count;  // Total number of commands
	unsigned int*   buf_idx;    // Indices buffer
    int             idx_count;  // Total number of indices
	gb_vertex_t*    buf_vtx;    // Vertices info buffer
    int             vtx_count;  // Total number of vertices
} gb_draw_list_t;

typedef struct gb_draw_data {
    gb_draw_list_t  buf_list;
    int             list_count;
} gb_draw_data_t;

gb_window_t gb_create_window(const char* title, int width, int height, gb_config_t* cfg);
void gb_window_destroy(gb_window_t win);
bool gb_window_start_frame(gb_window_t bw, double timeout);
void gb_window_render_frame(gb_window_t win, gb_draw_list_t dl);



