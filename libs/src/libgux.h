#pragma once
#include <stdbool.h>

typedef struct gb_config {
    struct {
        bool es;
        int  msaa;
        struct {
            float r; float g; float b; float a;
        } clearColor;
    } opengl;
} gb_config_t;

typedef void* gb_window_t;

typedef struct gb_color {
    float r; float g; float b; float w;
} gb_color_t;

typedef struct gb_vec4 {
    float x; float y; float z; float w;
} gb_vec4_t;

typedef struct gb_draw_cmd {
    gb_vec4_t   clip_rect;  // Clip rectangle
    float       texid;      // Texture id
    float       vtx_offset; // Start offset in vertex buffer
    float       idx_offset; // Start offset in index buffer
    float       elem_count; // Number of indices
} gb_draw_cmd_t;

typedef struct gb_draw_list {



} gb_draw_list_t;


gb_window_t gb_create_window(const char* title, int width, int height, gb_config_t* cfg);
void gb_window_destroy(gb_window_t win);
bool gb_window_start_frame(gb_window_t bw, double timeout);
void gb_window_render_frame(gb_window_t win, gb_draw_cmd_t* cmds, int cmd_count, int* buf_idx, float* buf_vtx);



