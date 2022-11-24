#include <stdio.h>
#include <stdlib.h>
#include "GL/gl3w.h"
#include <GLFW/glfw3.h>
#include "libgux.h"


// Internal state
typedef struct {
    GLFWwindow* w;
    struct {
        float r; float g; float b; float a;
    } clearColor;

    GLuint vao;
} window_state_t;


// Forward declarations of internal functions
static bool _gb_init(const char* glsl_version);
static void _gb_set_state();
static void _gb_render(window_state_t* s, gb_draw_list_t dl);

// Creates Graphics Backend window
gb_window_t gb_create_window(const char* title, int width, int height, gb_config_t* cfg) {

	if (glfwInit() == 0) {
        return NULL;
    }

    // Determines the api version from configuration
    int clientApi;
    const char* glVersion;
    int vmajor;
    int vminor;
//    if (cfg->opengl.es) {
//        clientApi = GLFW_OPENGL_ES_API;
//        glVersion ="#version 300 es";
//        vmajor = 3;
//        vminor = 1;
//    } else {
        clientApi = GLFW_OPENGL_API;
        glVersion ="#version 330";
        vmajor = 3;
        vminor = 3;
//    }

    // Set GLFW hints
    glfwWindowHint(GLFW_CLIENT_API, clientApi);
    glfwWindowHint(GLFW_CONTEXT_VERSION_MAJOR, vmajor);
    glfwWindowHint(GLFW_CONTEXT_VERSION_MINOR, vminor);
    glfwWindowHint(GLFW_VISIBLE, GLFW_TRUE);
//    if (!cfg->opengl.es) {
        glfwWindowHint(GLFW_OPENGL_PROFILE, GLFW_OPENGL_CORE_PROFILE);  // 3.2+ only
//    }
//    if (cfg->opengl.msaa > 0 && cfg->opengl.msaa <= 16) {
//        glfwWindowHint(GLFW_SAMPLES, 8);
//    }

    // Create window
    GLFWwindow* win = glfwCreateWindow(width, height, title, NULL, NULL);
    if (win == NULL) {
        fprintf(stderr, "Error creating GLFW Window\n");
        return NULL;
    }

    glfwMakeContextCurrent(win);
    glfwSwapInterval(1); // Enable vsync

    // Initialize OpenGL
    bool res = _gb_init(NULL);
    if (!res) {
        fprintf(stderr, "OpenGL initialization error");
        return NULL;
    }
    
    window_state_t* s = malloc(sizeof(window_state_t));
    if (s == NULL) {
        return NULL;
    }
    s->w = win;
    s->clearColor.r = 0.5;
    s->clearColor.g = 0.5;
    s->clearColor.b = 0.5;
    s->clearColor.a = 1.0;
    printf("%f/%f/%f\n", s->clearColor.r, s->clearColor.g, s->clearColor.b);
    return s;
}

// Destroy the specified window
void gb_window_destroy(gb_window_t bw) {

    window_state_t* s = (window_state_t*)(bw);
    glfwDestroyWindow(s->w);
    glfwTerminate();
    free(s);
}

// Starts the frame or returns false if the window should be closed
bool gb_window_start_frame(gb_window_t bw, double timeout) {

    window_state_t* s = (window_state_t*)(bw);
    // Checks if user requested window close
    if (glfwWindowShouldClose(s->w)) {
        return false;
    }

    // Poll and handle events (inputs, window resize, etc.)
    // Blocks if no events for the specified timeout
    glfwWaitEventsTimeout(timeout);

    // Starts the gl frame...
    // TODO
    return true;
}


// Renders the frame
//void gb_window_render_frame(gb_window_t bw, gb_draw_cmd_t* cmds, int cmd_count, int* buf_idx, float* buf_vtx) {
void gb_window_render_frame(gb_window_t bw, gb_draw_list_t dl) {

    // Sets the OpenGL viewport
    window_state_t* s = (window_state_t*)(bw);
    int display_w, display_h;
    glfwGetFramebufferSize(s->w, &display_w, &display_h);
    glViewport(0, 0, display_w, display_h);

    // Clears the framebuffer
    glClearColor(s->clearColor.r, s->clearColor.g, s->clearColor.b, s->clearColor.a);
    glClear(GL_COLOR_BUFFER_BIT);

    // Render commands and swap buffers
    _gb_render(s, dl);
    glfwSwapBuffers(s->w);
}

// Load OpenGL functions and initialize its state
static bool _gb_init(const char* glsl_version) {

    // Load OpenGL functions
    int res = gl3wInit();
    if (res != 0) {
        fprintf(stderr, "Failed to initialize OpenGL loader!\n");
        return false;
    }

    // Sets initial state and checks fo error
    _gb_set_state();
    int err = glGetError();
    if (err != GL_NO_ERROR) {
        fprintf(stderr, "OpenGL returned error:%d", err);
        return false;
    }
    return true;
}

// Sets desired OpenGL state
static void _gb_set_state() {

    glEnable(GL_BLEND);
    glBlendEquation(GL_FUNC_ADD);
    glBlendFuncSeparate(GL_SRC_ALPHA, GL_ONE_MINUS_SRC_ALPHA, GL_ONE, GL_ONE_MINUS_SRC_ALPHA);
    glDisable(GL_CULL_FACE);
    glDisable(GL_DEPTH_TEST);
    glDisable(GL_STENCIL_TEST);
    // glEnable(GL_SCISSOR_TEST); // DOES NOT CLEAR COMPLETE FRAMEBUFFER
}

static bool _gb_createDeviceObjects() {


}

// Render commands
static void _gb_render(window_state_t* s, gb_draw_list_t dl)  {

    //printf("idx_count:%d, vtx_count:%d\n", dl.idx_count, dl.vtx_count);
    for (int i = 0; i < dl.cmd_count; i++) {

        gb_draw_cmd_t cmd = dl.bufCmd[i];
        int texid = cmd.texid;
        int idx_offset = cmd.idx_offset;
        int vtx_offset = cmd.vtx_offset;
        int elem_count = cmd.elem_count;

        //printf("x:%f, y:%f, z:%f, w:%f, texid:%d, vtx_offset:%d, idx_offset:%d, elem_count:%d\n",
        //    cmd.clip_rect.x, cmd.clip_rect.y, cmd.clip_rect.z, cmd.clip_rect.w,
        //    texid, vtx_offset, idx_offset, elem_count);
    }
}

