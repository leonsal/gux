#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <stddef.h>
#include "GL/gl3w.h"
#include <GLFW/glfw3.h>
#include "libgux.h"

// Uncomment next line to enable error checking after each OpenGL call.
#define GUX_GB_DEBUG
#ifdef GUX_GB_DEBUG
#include <stdio.h>
#define GL_CALL(_CALL) do { _CALL; GLenum gl_err = glGetError(); if (gl_err != 0) fprintf(stderr, "GL error 0x%x returned from '%s'.\n", gl_err, #_CALL); } while (0)  // Call with error check
#else
#define GL_CALL(_CALL) _CALL   // Call without error check
#endif


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
static void _gb_render(gb_state_t* s, gb_draw_list_t dl);
static bool _gb_init(gb_state_t* s, const char* glsl_version);
static void _gb_set_state(gb_state_t* s);
static bool _gb_create_objects(gb_state_t* s);
static void _gb_destroy_objects(gb_state_t* s);
static bool _gb_check_shader(GLuint handle, const char* desc, const char* src);
static bool _gb_check_program(GLuint handle, const char* desc);
static void _gb_print_draw_list(gb_draw_list_t dl);

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
    bool res = _gb_init(s, NULL);
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
    GL_CALL(glViewport(0, 0, width, height));

    // Clears the framebuffer
    GL_CALL(glScissor(0, 0, width, height));
    GL_CALL(glClearColor(s->clearColor.r, s->clearColor.g, s->clearColor.b, s->clearColor.a));
    GL_CALL(glClear(GL_COLOR_BUFFER_BIT));

    // Render commands and swap buffers
    _gb_render(s, dl);
    glfwSwapBuffers(s->w);
}

// Executes draw commands
static void _gb_render(gb_state_t* s, gb_draw_list_t dl)  {

    // Upload vertices and indices buffers
    // glBufferData(GL_ARRAY_BUFFER, vtx_buffer_size, (const GLvoid*)cmd_list->VtxBuffer.Data, GL_STREAM_DRAW));
    //        GL_CALL(glBufferData(GL_ELEMENT_ARRAY_BUFFER, idx_buffer_size, (const GLvoid*)cmd_list->IdxBuffer.Data, GL_STREAM_DRAW));
    //
    _gb_print_draw_list(dl);

    for (int i = 0; i < dl.cmd_count; i++) {
        gb_draw_cmd_t cmd = dl.bufCmd[i];
    }

}

// Load OpenGL functions and initialize its state
static bool _gb_init(gb_state_t* s, const char* glsl_version) {

    // Load OpenGL functions
    int res = gl3wInit();
    if (res != 0) {
        fprintf(stderr, "Failed to initialize OpenGL loader!\n");
        return false;
    }

    if (!_gb_create_objects(s)) {
        return false;
    }

    // Sets initial state
    _gb_set_state(s);
    return true;
}

// Sets required OpenGL state
static void _gb_set_state(gb_state_t* s) {

    GL_CALL(glEnable(GL_BLEND));
    GL_CALL(glBlendEquation(GL_FUNC_ADD));
    GL_CALL(glBlendFuncSeparate(GL_SRC_ALPHA, GL_ONE_MINUS_SRC_ALPHA, GL_ONE, GL_ONE_MINUS_SRC_ALPHA));
    GL_CALL(glDisable(GL_CULL_FACE));
    GL_CALL(glDisable(GL_DEPTH_TEST));
    GL_CALL(glDisable(GL_STENCIL_TEST));
    GL_CALL(glEnable(GL_SCISSOR_TEST));

    GL_CALL(glUseProgram(s->shaderHandle));
    GL_CALL(glUniform1i(s->locTex, 0));
    //glUniformMatrix4fv(bd->AttribLocationProjMtx, 1, GL_FALSE, &ortho_projection[0][0]);

    GL_CALL(glBindSampler(0, 0));
    //(void)vertex_array_object;
    //glBindVertexArray(vertex_array_object);

    // Bind vertex/index buffers and setup attributes for ImDrawVert
    GL_CALL(glBindBuffer(GL_ARRAY_BUFFER, s->vboHandle));
    GL_CALL(glBindBuffer(GL_ELEMENT_ARRAY_BUFFER, s->elementsHandle));
    GL_CALL(glEnableVertexAttribArray(s->locVtxPos));
    GL_CALL(glEnableVertexAttribArray(s->locVtxUV));
    GL_CALL(glEnableVertexAttribArray(s->locVtxColor));
    GL_CALL(glVertexAttribPointer(s->locVtxPos,   2, GL_FLOAT, GL_FALSE, sizeof(gb_vertex_t), (GLvoid*)offsetof(gb_vertex_t, pos)));
    GL_CALL(glVertexAttribPointer(s->locVtxUV,    2, GL_FLOAT, GL_FALSE, sizeof(gb_vertex_t), (GLvoid*)offsetof(gb_vertex_t, uv)));
    GL_CALL(glVertexAttribPointer(s->locVtxColor, 4, GL_UNSIGNED_BYTE, GL_TRUE, sizeof(gb_vertex_t), (GLvoid*)offsetof(gb_vertex_t, col)));
}

static bool _gb_create_objects(gb_state_t* s) {

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
    GL_CALL(glShaderSource(vert_handle, 1, &vertex_shader, NULL));
    GL_CALL(glCompileShader(vert_handle));
    if (!_gb_check_shader(vert_handle, "vertex_shader", vertex_shader)) {
        return false;
    }

    // Create fragment shader
    GLuint frag_handle = glCreateShader(GL_FRAGMENT_SHADER);
    GL_CALL(glShaderSource(frag_handle, 1, &fragment_shader, NULL));
    GL_CALL(glCompileShader(frag_handle));
    if (!_gb_check_shader(frag_handle, "fragment shader", fragment_shader)) {
        return false;
    }

    // Create program
    s->shaderHandle = glCreateProgram();
    GL_CALL(glAttachShader(s->shaderHandle, vert_handle));
    GL_CALL(glAttachShader(s->shaderHandle, frag_handle));
    GL_CALL(glLinkProgram(s->shaderHandle));
    if (!_gb_check_program(s->shaderHandle, "shader program")) {
        return false;
    }

    // Discard shaders
    GL_CALL(glDetachShader(s->shaderHandle, vert_handle));
    GL_CALL(glDetachShader(s->shaderHandle, frag_handle));
    GL_CALL(glDeleteShader(vert_handle));
    GL_CALL(glDeleteShader(frag_handle));

    // Get uniform locations from shader progrm
    s->locTex = glGetUniformLocation(s->shaderHandle, "Texture");
    s->locProjMtx = glGetUniformLocation(s->shaderHandle, "ProjMtx");
    s->locVtxPos = (GLuint)glGetAttribLocation(s->shaderHandle, "Position");
    s->locVtxUV = (GLuint)glGetAttribLocation(s->shaderHandle, "UV");
    s->locVtxColor = (GLuint)glGetAttribLocation(s->shaderHandle, "Color");
    printf("LOCS: %d/%d/%d/%d/%d\n", s->locTex, s->locProjMtx, s->locVtxPos, s->locVtxUV, s->locVtxColor);

    // Create buffers
    GL_CALL(glGenBuffers(1, &s->vboHandle));
    GL_CALL(glGenBuffers(1, &s->elementsHandle));
    return true;
}


static void _gb_destroy_objects(gb_state_t* s) {

    if (s->vboHandle) {
        GL_CALL(glDeleteBuffers(1, &s->vboHandle));
        s->vboHandle = 0;
    }
    if (s->elementsHandle) {
        GL_CALL(glDeleteBuffers(1, &s->elementsHandle));
        s->elementsHandle = 0;
    }
    if (s->shaderHandle) {
        GL_CALL(glDeleteProgram(s->shaderHandle));
        s->shaderHandle = 0;
    }
}

static bool _gb_check_shader(GLuint handle, const char* desc, const char* src) {

    GLint status = 0, log_length = 0;
    GL_CALL(glGetShaderiv(handle, GL_COMPILE_STATUS, &status));
    GL_CALL(glGetShaderiv(handle, GL_INFO_LOG_LENGTH, &log_length));
    if (status == GL_FALSE || log_length > 0) {
        fprintf(stderr, "ERROR: gb/_gb_check_shader: error compiling %s\n", desc);
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

static bool _gb_check_program(GLuint handle, const char* desc) {

    GLint status = 0, log_length = 0;
    GL_CALL(glGetProgramiv(handle, GL_LINK_STATUS, &status));
    GL_CALL(glGetProgramiv(handle, GL_INFO_LOG_LENGTH, &log_length));
    if (status == GL_FALSE || log_length > 0) {
        fprintf(stderr, "ERROR: gb/_gb_check_program: error linking %s\n", desc);
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

static void _gb_print_draw_list(gb_draw_list_t dl) {

    printf("DrawList: idx_count:%d, vtx_count:%d\n", dl.idx_count, dl.vtx_count);
    printf("Indices :");
    for (int i = 0; i < dl.idx_count; i++) {
        printf("%d,", dl.bufIdx[i]);
    }
    printf("\n");

    printf("Vertices:\n");
    for (int i = 0; i < dl.vtx_count; i++) {
        gb_vertex_t v = dl.bufVtx[i];
        printf("\tx:%f, y:%f u:%f v:%f col:%06X\n", v.pos.x, v.pos.y, v.uv.x, v.uv.y, v.col);
    }

    printf("Commands:\n");
    for (int i = 0; i < dl.cmd_count; i++) {
        gb_draw_cmd_t cmd = dl.bufCmd[i];
        printf("\tx:%f, y:%f, z:%f, w:%f, texid:%d, vtx_offset:%d, idx_offset:%d, elem_count:%d\n",
            cmd.clip_rect.x, cmd.clip_rect.y, cmd.clip_rect.z, cmd.clip_rect.w,
            cmd.texid, cmd.vtx_offset, cmd.idx_offset, cmd.elem_count);
    }
    printf("\n");
}


