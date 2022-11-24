#pragma once
#include <stdbool.h>

typedef struct gl_config {
    struct {
        bool es;
        int  msaa;
        struct {
            float r; float g; float b; float a;
        } clearColor;
    } opengl;
} gl_config_t;

typedef void* gl_window_t;

typedef struct gl_color {
    float r; float g; float b; float w;
} gl_color_t;

typedef struct gl_vec4 {
    float x; float y; float z; float w;
} gl_vec4_t;

typedef struct gl_draw_cmd {
    gl_vec4_t       clip_rect;  // Clip rectangle
    unsigned int    texid;      // Texture id
    unsigned int    vtx_offset; // Start offset in vertex buffer
    unsigned int    idx_offset; // Start offset in index buffer
    unsigned int    elem_count; // Number of indices
} gl_draw_cmd_t;

typedef struct gl_draw_list {
    gl_draw_cmd_t*  buf_cmd;    // Array of draw commands
    unsigned int    cmd_count;  // Number of draw commands in the array
    int*            buf_idx;    // Index buffer
    float*          buf_vtx;    // Vertices buffer
} gl_draw_list_t;


gl_window_t gl_create_window(const char* title, int width, int height, gl_config_t* cfg);
void gl_window_destroy(gl_window_t win);
bool gl_window_start_frame(gl_window_t bw, double timeout);
void gl_window_render_frame(gl_window_t win, gl_draw_list_t dl);



