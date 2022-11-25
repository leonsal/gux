#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include "GL/gl3w.h"
#include <GLFW/glfw3.h>
#include "libgux.h"


// Internal state
typedef struct {
    GLFWwindow* w;
    struct {
        float r; float g; float b; float a;
    } clearColor;
    GLuint  shaderHandle;
    GLint   locTex;
    GLint   locProjMtx;
    GLint   locVtxPos;
    GLint   locVtxUV;
    GLint   locVtxColor;
    unsigned int vboHandle;
    unsigned int elementsHandle;
} gb_state_t;


// Forward declarations of internal functions
static void _gc_render(gb_state_t* s, gb_draw_list_t dl);
static bool _gc_init(gb_state_t* s, const char* glsl_version);
static void _gc_set_state(gb_state_t* s);
static bool _gc_create_objects(gb_state_t* s);
static void _gc_destroy_objects(gb_state_t* s);
static bool _gc_check_shader(GLuint handle, const char* desc, const char* src);
static bool _gc_check_program(GLuint handle, const char* desc);

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
                         
    // Creates and initializes backend state
    gb_state_t* s = malloc(sizeof(gb_state_t));
    if (s == NULL) {
        return NULL;
    }
    memset(s, 0, sizeof(gb_state_t));

    // Initialize OpenGL
    bool res = _gc_init(s, NULL);
    if (!res) {
        fprintf(stderr, "OpenGL initialization error");
        return NULL;
    }

    s->w = win;
    s->clearColor.r = 0.5;
    s->clearColor.g = 0.5;
    s->clearColor.b = 0.5;
    s->clearColor.a = 1.0;
    return s;
}

// Destroy the specified window
void gb_window_destroy(gb_window_t bw) {

    gb_state_t* s = (gb_state_t*)(bw);
    glfwDestroyWindow(s->w);
    glfwTerminate();
    free(s);
}

// Starts the frame or returns false if the window should be closed
bool gb_window_start_frame(gb_window_t bw, double timeout) {

    // Checks if user requested window close
    gb_state_t* s = (gb_state_t*)(bw);
    if (glfwWindowShouldClose(s->w)) {
        return false;
    }

    // Poll and handle events, blocking if no events for the specified timeout
    glfwWaitEventsTimeout(timeout);
    return true;
}

// Renders the frame
void gb_window_render_frame(gb_window_t bw, gb_draw_list_t dl) {

    // Sets the OpenGL viewport from the framebuffer size
    gb_state_t* s = (gb_state_t*)(bw);
    int width, height;
    glfwGetFramebufferSize(s->w, &width, &height);
    glViewport(0, 0, width, height);

    // Clears the framebuffer
    glScissor(0, 0, width, height);
    glClearColor(s->clearColor.r, s->clearColor.g, s->clearColor.b, s->clearColor.a);
    glClear(GL_COLOR_BUFFER_BIT);

    // Render commands and swap buffers
    _gc_render(s, dl);
    glfwSwapBuffers(s->w);
}

// Executes draw commands
static void _gc_render(gb_state_t* s, gb_draw_list_t dl)  {

    // Upload vertices and indices buffers
    // glBufferData(GL_ARRAY_BUFFER, vtx_buffer_size, (const GLvoid*)cmd_list->VtxBuffer.Data, GL_STREAM_DRAW));
    //        GL_CALL(glBufferData(GL_ELEMENT_ARRAY_BUFFER, idx_buffer_size, (const GLvoid*)cmd_list->IdxBuffer.Data, GL_STREAM_DRAW));

    //printf("RENDER:idx_count:%d, vtx_count:%d\n", dl.idx_count, dl.vtx_count);
    for (int i = 0; i < dl.cmd_count; i++) {
        gb_draw_cmd_t cmd = dl.bufCmd[i];

        //printf("x:%f, y:%f, z:%f, w:%f, texid:%d, vtx_offset:%d, idx_offset:%d, elem_count:%d\n",
        //    cmd.clip_rect.x, cmd.clip_rect.y, cmd.clip_rect.z, cmd.clip_rect.w,
        //    cmd.texid, cmd.vtx_offset, cmd.idx_offset, cmd.elem_count);
    }

}

// Load OpenGL functions and initialize its state
static bool _gc_init(gb_state_t* s, const char* glsl_version) {

    // Load OpenGL functions
    int res = gl3wInit();
    if (res != 0) {
        fprintf(stderr, "Failed to initialize OpenGL loader!\n");
        return false;
    }

    if (!_gc_create_objects(s)) {
        return false;
    }

    // Sets initial state and checks fo error
    _gc_set_state(s);
    int err = glGetError();
    if (err != GL_NO_ERROR) {
        fprintf(stderr, "OpenGL returned error:%d", err);
        return false;
    }
    return true;
}

// Sets required OpenGL state
static void _gc_set_state(gb_state_t* s) {

    glEnable(GL_BLEND);
    glBlendEquation(GL_FUNC_ADD);
    glBlendFuncSeparate(GL_SRC_ALPHA, GL_ONE_MINUS_SRC_ALPHA, GL_ONE, GL_ONE_MINUS_SRC_ALPHA);
    glDisable(GL_CULL_FACE);
    glDisable(GL_DEPTH_TEST);
    glDisable(GL_STENCIL_TEST);
    glEnable(GL_SCISSOR_TEST);
}

static bool _gc_create_objects(gb_state_t* s) {

//    const GLchar* vertex_shader_glsl_300_es =
//        "precision highp float;\n"
//        "layout (location = 0) in vec2 Position;\n"
//        "layout (location = 1) in vec2 UV;\n"
//        "layout (location = 2) in vec4 Color;\n"
//        "uniform mat4 ProjMtx;\n"
//        "out vec2 Frag_UV;\n"
//        "out vec4 Frag_Color;\n"
//        "void main()\n"
//        "{\n"
//        "    Frag_UV = UV;\n"
//        "    Frag_Color = Color;\n"
//        "    gl_Position = ProjMtx * vec4(Position.xy,0,1);\n"
//        "}\n";
//
//    const GLchar* fragment_shader_glsl_300_es =
//        "precision mediump float;\n"
//        "uniform sampler2D Texture;\n"
//        "in vec2 Frag_UV;\n"
//        "in vec4 Frag_Color;\n"
//        "layout (location = 0) out vec4 Out_Color;\n"
//        "void main()\n"
//        "{\n"
//        "    Out_Color = Frag_Color * texture(Texture, Frag_UV.st);\n"
//        "}\n";
//
    const GLchar* vertex_shader_glsl_330_core =
        "#version 330 core\n"
        "layout (location = 0) in vec2 Position;\n"
        "layout (location = 1) in vec2 UV;\n"
        "layout (location = 2) in vec4 Color;\n"
        "uniform mat4 ProjMtx;\n"
        "out vec2 Frag_UV;\n"
        "out vec4 Frag_Color;\n"
        "void main()\n"
        "{\n"
        "    Frag_UV = UV;\n"
        "    Frag_Color = Color;\n"
        "    gl_Position = ProjMtx * vec4(Position.xy,0,1);\n"
        "}\n";

    const GLchar* fragment_shader_glsl_330_core =
        "#version 330 core\n"
        "in vec2 Frag_UV;\n"
        "in vec4 Frag_Color;\n"
        "uniform sampler2D Texture;\n"
        "layout (location = 0) out vec4 Out_Color;\n"
        "void main()\n"
        "{\n"
        "    Out_Color = Frag_Color * texture(Texture, Frag_UV.st);\n"
        "}\n";

    const GLchar* vertex_shader = vertex_shader_glsl_330_core;
    const GLchar* fragment_shader = fragment_shader_glsl_330_core;

    // Create vertex shader
    GLuint vert_handle = glCreateShader(GL_VERTEX_SHADER);
    glShaderSource(vert_handle, 1, &vertex_shader, NULL);
    glCompileShader(vert_handle);
    if (!_gc_check_shader(vert_handle, "vertex_shader", vertex_shader)) {
        return false;
    }

    // Create fragment shader
    GLuint frag_handle = glCreateShader(GL_FRAGMENT_SHADER);
    glShaderSource(frag_handle, 1, &fragment_shader, NULL);
    glCompileShader(frag_handle);
    if (!_gc_check_shader(frag_handle, "fragment shader", fragment_shader)) {
        return false;
    }

    // Create program
    s->shaderHandle = glCreateProgram();
    glAttachShader(s->shaderHandle, vert_handle);
    glAttachShader(s->shaderHandle, frag_handle);
    glLinkProgram(s->shaderHandle);
    if (!_gc_check_program(s->shaderHandle, "shader program")) {
        return false;
    }

    // Discard shaders
    glDetachShader(s->shaderHandle, vert_handle);
    glDetachShader(s->shaderHandle, frag_handle);
    glDeleteShader(vert_handle);
    glDeleteShader(frag_handle);

    // Get uniform locations from shader progrm
    s->locTex = glGetUniformLocation(s->shaderHandle, "Texture");
    s->locProjMtx = glGetUniformLocation(s->shaderHandle, "ProjMtx");
    s->locVtxPos = (GLuint)glGetAttribLocation(s->shaderHandle, "Position");
    s->locVtxUV = (GLuint)glGetAttribLocation(s->shaderHandle, "UV");
    s->locVtxColor = (GLuint)glGetAttribLocation(s->shaderHandle, "Color");
    printf("LOCS: %d/%d/%d/%d/%d\n", s->locTex, s->locProjMtx, s->locVtxPos, s->locVtxUV, s->locVtxColor);

    // Create buffers
    glGenBuffers(1, &s->vboHandle);
    glGenBuffers(1, &s->elementsHandle);
    return true;
}


static void _gc_destroy_objects(gb_state_t* s) {

    if (s->vboHandle) {
        glDeleteBuffers(1, &s->vboHandle);
        s->vboHandle = 0;
    }
    if (s->elementsHandle) {
        glDeleteBuffers(1, &s->elementsHandle);
        s->elementsHandle = 0;
    }
    if (s->shaderHandle) {
        glDeleteProgram(s->shaderHandle);
        s->shaderHandle = 0;
    }
}

static bool _gc_check_shader(GLuint handle, const char* desc, const char* src) {

    GLint status = 0, log_length = 0;
    glGetShaderiv(handle, GL_COMPILE_STATUS, &status);
    glGetShaderiv(handle, GL_INFO_LOG_LENGTH, &log_length);
    if (status == GL_FALSE || log_length > 0) {
        fprintf(stderr, "ERROR: gb/_gc_check_shader: error compiling %s\n", desc);
        fprintf(stderr, "%s\n", src);
        if (log_length > 0) {
            GLchar* buf = malloc(log_length + 1);
            glGetShaderInfoLog(handle, log_length, NULL, buf);
            fprintf(stderr, "%s\n", buf);
            free(buf);
        }
        return false;
    }
    return true;
}

static bool _gc_check_program(GLuint handle, const char* desc) {

    GLint status = 0, log_length = 0;
    glGetProgramiv(handle, GL_LINK_STATUS, &status);
    glGetProgramiv(handle, GL_INFO_LOG_LENGTH, &log_length);
    if (status == GL_FALSE || log_length > 0) {
        fprintf(stderr, "ERROR: gb/_gc_check_program: error linking %s\n", desc);
        if (log_length > 0) {
            GLchar* buf = malloc(log_length + 1);
            glGetProgramInfoLog(handle, log_length, NULL, buf);
            fprintf(stderr, "%s\n", buf);
            free(buf);
        }
        return false;
    }
    return true;
}

