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
    GLFWwindow*     w;                  // GLFW window pointer
    gb_vec4_t       clear_color;        // Current color to clear color buffer before rendering
    GLuint          handle_shader;      // Handle of compiled shader program
    GLint           uni_tex;            // Location of texture id uniform in the shader
    GLint           uni_projmtx;        // Location of projection matrix uniform the shader
    GLint           attrib_vtx_pos;     // Location of vertex position attribute in the shader
    GLint           attrib_vtx_uv;      // Location of vertex uv attribute in the shader
    GLint           attrib_vtx_color;   // Location of vertex color attribute in the shader
    unsigned int    handle_vbo;         // Handle of vertex buffer object
    unsigned int    handle_elems;       // Handle of vertex elements object
    GLint           vao;
    gb_event_t*     events;             // Pointer to events array
    int             ev_count;           // Current number of valid events in the events array
    int             ev_cap;             // Current capacity of events array
} gb_state_t;


// Forward declarations of internal functions
static void _gb_render(gb_state_t* s, gb_vec2_t disp_pos, gb_vec2_t disp_size,  gb_draw_list_t dl);
static bool _gb_init(gb_state_t* s, const char* glsl_version);
static void _gb_set_state(gb_state_t* s);
static bool _gb_create_objects(gb_state_t* s);
static void _gb_destroy_objects(gb_state_t* s);
static bool _gb_check_shader(GLuint handle, const char* desc, const char* src);
static bool _gb_check_program(GLuint handle, const char* desc);
static void _gb_print_draw_list(gb_draw_list_t dl);
static void _gb_set_ev_handlers(gb_state_t* s);
static gb_event_t* _gb_ev_reserve(gb_state_t* s);
static void _gb_key_callback(GLFWwindow* win, int key, int scancode, int action, int mods);
static void _gb_char_callback(GLFWwindow* win, unsigned int codepoint);
static void _gb_cursor_pos_callback(GLFWwindow* win, double xpos, double ypos);
static void _gb_cursor_enter_callback(GLFWwindow* win, int entered);
static void _gb_mouse_button_callback(GLFWwindow* win, int button, int action, int mods);
static void _gb_scroll_callback(GLFWwindow* win, double xoffset, double yoffset);

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
    s->w = win;
    glfwSetWindowUserPointer(win, s);

    // Initialize OpenGL
    bool res = _gb_init(s, NULL);
    if (!res) {
        fprintf(stderr, "OpenGL initialization error");
        return NULL;
    }

    // Sets default clear color
    s->clear_color.x = 0.5;
    s->clear_color.y = 0.5;
    s->clear_color.z = 0.5;
    s->clear_color.w = 1.0;

    // Allocates initial events array
    s->ev_count = 0;
    s->ev_cap = 1024;
    s->events = malloc(sizeof(gb_event_t) * s->ev_cap);
    if (s->events == NULL) {
        fprintf(stderr, "No memory for events array");
        return NULL;
    }
    _gb_set_ev_handlers(s);
    return s;
}

// Destroy the specified window
void gb_window_destroy(gb_window_t bw) {

    gb_state_t* s = (gb_state_t*)(bw);
    glfwDestroyWindow(s->w);
    free(s->events);
    s->events = NULL;

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
    GL_CALL(glClearColor(s->clear_color.x, s->clear_color.y, s->clear_color.z, s->clear_color.w));
    GL_CALL(glClear(GL_COLOR_BUFFER_BIT));

    // Render commands and swap buffers
    gb_vec2_t disp_pos = {0,0};
    gb_vec2_t disp_size = {width, height};
    _gb_render(s, disp_pos, disp_size, dl);
    glfwSwapBuffers(s->w);
}

// Creates and returns an OpenGL texture idenfifier
gb_texid_t gb_create_texture() {

    // Create a OpenGL texture identifier
    GLuint image_texture;
    glGenTextures(1, &image_texture);
    glBindTexture(GL_TEXTURE_2D, image_texture);

    // Setup filtering parameters for display
    glTexParameteri(GL_TEXTURE_2D, GL_TEXTURE_MIN_FILTER, GL_LINEAR);
    glTexParameteri(GL_TEXTURE_2D, GL_TEXTURE_MAG_FILTER, GL_LINEAR);

    return (intptr_t)image_texture;
}

// Deletes previously created texture
void gb_delete_texture(gb_texid_t texid) {

    GLuint tex = (GLuint)texid;
    glDeleteTextures(1, &tex); 
}

// Transfer data for the specified texture
void gb_transfer_texture(gb_texid_t texid, int width, int height, const gb_color_t* data) {

    GLuint tex = (GLuint)(texid);
    glBindTexture(GL_TEXTURE_2D, tex);
    glPixelStorei(GL_UNPACK_ROW_LENGTH, 0);
    glTexImage2D(GL_TEXTURE_2D, 0, GL_RGBA, width, height, 0, GL_RGBA, GL_UNSIGNED_BYTE, data);
}

int gb_get_events(gb_window_t win, gb_event_t* events, int ev_count) {

    // Transfer specified number of events
    gb_state_t* s = (gb_state_t*)(win);
    if (s->ev_count == 0) {
        return 0;
    }
    if (ev_count > s->ev_count) {
        ev_count = s->ev_count;
    }
    memcpy(events, s->events, sizeof(gb_event_t) * ev_count);

    // Shift remaining events to the start of the buffer
    int remain = s->ev_count - ev_count;
    if (remain > 0) {
        memmove(s->events, s->events + (sizeof(gb_event_t) * ev_count), remain);
    }
    s->ev_count = remain;
    return ev_count;
}


//-----------------------------------------------------------------------------
// Internal functions
//-----------------------------------------------------------------------------


// Executes draw commands
static void _gb_render(gb_state_t* s, gb_vec2_t disp_pos, gb_vec2_t disp_size,  gb_draw_list_t dl)  {

    //printf("render-> cmd_count:%d idx_count:%d vtx_count:%d\n", dl.cmd_count, dl.idx_count, dl.vtx_count);

    // Do not render when minimized
    int fb_width = (int)disp_size.x;
    int fb_height = (int)disp_size.y;
    if (fb_width <= 0 || fb_height <= 0) {
        return;
    }

    // Upload vertices and indices buffers
    const GLsizeiptr vtx_buffer_size = (GLsizeiptr)dl.vtx_count * (int)sizeof(gb_vertex_t);
    const GLsizeiptr idx_buffer_size = (GLsizeiptr)dl.idx_count * (int)sizeof(int);
    //printf("buffer sizes:%ld/%ld\n",  vtx_buffer_size, idx_buffer_size);
    GL_CALL(glBufferData(GL_ARRAY_BUFFER, vtx_buffer_size, (const GLvoid*)dl.buf_vtx, GL_STREAM_DRAW));
    GL_CALL(glBufferData(GL_ELEMENT_ARRAY_BUFFER, idx_buffer_size, (const GLvoid*)dl.buf_idx, GL_STREAM_DRAW));

    // Sets orthogonal projection
    float L = disp_pos.x;
    float R = disp_pos.x + disp_size.x;
    float T = disp_pos.y;
    float B = disp_pos.y + disp_size.y;
    const float ortho_projection[4][4] = {
        { 2.0f/(R-L),   0.0f,         0.0f,   0.0f },
        { 0.0f,         2.0f/(T-B),   0.0f,   0.0f },
        { 0.0f,         0.0f,        -1.0f,   0.0f },
        { (R+L)/(L-R),  (T+B)/(B-T),  0.0f,   1.0f },
    };
    glUniformMatrix4fv(s->uni_projmtx, 1, GL_FALSE, &ortho_projection[0][0]);

    //_gb_print_draw_list(dl);

    for (int i = 0; i < dl.cmd_count; i++) {
        gb_draw_cmd_t cmd = dl.buf_cmd[i];
        // Apply scissor/clipping rectangle (Y is inverted in OpenGL)
        //GL_CALL(glScissor((int)clip_min.x, (int)((float)fb_height - clip_max.y), (int)(clip_max.x - clip_min.x), (int)(clip_max.y - clip_min.y)));

        GL_CALL(glBindTexture(GL_TEXTURE_2D, (GLuint)(cmd.texid)));
        GL_CALL(glDrawElements(GL_TRIANGLES, (GLsizei)cmd.elem_count, GL_UNSIGNED_INT, (void*)(intptr_t)(cmd.idx_offset * sizeof(cmd.idx_offset))));
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

    GL_CALL(glUseProgram(s->handle_shader));
    GL_CALL(glUniform1i(s->uni_tex, 0));
    
    GL_CALL(glBindSampler(0, 0));
    GL_CALL(glBindVertexArray(s->vao));

    // Bind vertex/index buffers and setup attributes for ImDrawVert
    GL_CALL(glBindBuffer(GL_ARRAY_BUFFER, s->handle_vbo));
    GL_CALL(glBindBuffer(GL_ELEMENT_ARRAY_BUFFER, s->handle_elems));
    GL_CALL(glEnableVertexAttribArray(s->attrib_vtx_pos));
    GL_CALL(glEnableVertexAttribArray(s->attrib_vtx_uv));
    GL_CALL(glEnableVertexAttribArray(s->attrib_vtx_color));

    printf("loc:%d sizeof:%ld, offset:%ld\n", s->attrib_vtx_pos, sizeof(gb_vertex_t), offsetof(gb_vertex_t, pos));
    printf("loc:%d sizeof:%ld, offset:%ld\n", s->attrib_vtx_uv, sizeof(gb_vertex_t), offsetof(gb_vertex_t, uv));
    printf("loc:%d sizeof:%ld, offset:%ld\n", s->attrib_vtx_color, sizeof(gb_vertex_t), offsetof(gb_vertex_t, col));
    GL_CALL(glVertexAttribPointer(s->attrib_vtx_pos,   2, GL_FLOAT, GL_FALSE, sizeof(gb_vertex_t), (GLvoid*)offsetof(gb_vertex_t, pos)));
    GL_CALL(glVertexAttribPointer(s->attrib_vtx_uv,    2, GL_FLOAT, GL_FALSE, sizeof(gb_vertex_t), (GLvoid*)offsetof(gb_vertex_t, uv)));
    GL_CALL(glVertexAttribPointer(s->attrib_vtx_color, 4, GL_UNSIGNED_BYTE, GL_TRUE, sizeof(gb_vertex_t), (GLvoid*)offsetof(gb_vertex_t, col)));
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
        "    //gl_Position = vec4(Position.xy,0,1);\n"
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
    s->handle_shader = glCreateProgram();
    GL_CALL(glAttachShader(s->handle_shader, vert_handle));
    GL_CALL(glAttachShader(s->handle_shader, frag_handle));
    GL_CALL(glLinkProgram(s->handle_shader));
    if (!_gb_check_program(s->handle_shader, "shader program")) {
        return false;
    }

    // Discard shaders
    GL_CALL(glDetachShader(s->handle_shader, vert_handle));
    GL_CALL(glDetachShader(s->handle_shader, frag_handle));
    GL_CALL(glDeleteShader(vert_handle));
    GL_CALL(glDeleteShader(frag_handle));

    // Get uniform locations from shader progrm
    s->uni_tex = glGetUniformLocation(s->handle_shader, "Texture");
    s->uni_projmtx = glGetUniformLocation(s->handle_shader, "ProjMtx");
    s->attrib_vtx_pos = (GLuint)glGetAttribLocation(s->handle_shader, "Position");
    s->attrib_vtx_uv = (GLuint)glGetAttribLocation(s->handle_shader, "UV");
    s->attrib_vtx_color = (GLuint)glGetAttribLocation(s->handle_shader, "Color");
    printf("LOCS: %d/%d/%d/%d/%d\n", s->uni_tex, s->uni_projmtx, s->attrib_vtx_pos, s->attrib_vtx_uv, s->attrib_vtx_color);

    // Create buffers
    GL_CALL(glGenBuffers(1, &s->handle_vbo));
    GL_CALL(glGenBuffers(1, &s->handle_elems));

    // Create VAO
    GL_CALL(glGenVertexArrays(1, &s->vao));

    return true;
}


static void _gb_destroy_objects(gb_state_t* s) {

    if (s->handle_vbo) {
        GL_CALL(glDeleteBuffers(1, &s->handle_vbo));
        s->handle_vbo = 0;
    }
    if (s->handle_elems) {
        GL_CALL(glDeleteBuffers(1, &s->handle_elems));
        s->handle_elems = 0;
    }
    if (s->handle_shader) {
        GL_CALL(glDeleteProgram(s->handle_shader));
        s->handle_shader = 0;
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

// Prints the specifid draw list for debugging
static void _gb_print_draw_list(gb_draw_list_t dl) {

    printf("DrawList: idx_count:%d, vtx_count:%d\n", dl.idx_count, dl.vtx_count);
    printf("Indices :");
    for (int i = 0; i < dl.idx_count; i++) {
        printf("%d,", dl.buf_idx[i]);
    }
    printf("\n");

    printf("Vertices:\n");
    for (int i = 0; i < dl.vtx_count; i++) {
        gb_vertex_t v = dl.buf_vtx[i];
        printf("\tx:%f, y:%f u:%f v:%f col:%06X\n", v.pos.x, v.pos.y, v.uv.x, v.uv.y, v.col);
    }

    printf("Commands:\n");
    for (int i = 0; i < dl.cmd_count; i++) {
        gb_draw_cmd_t cmd = dl.buf_cmd[i];
        printf("\tx:%f, y:%f, z:%f, w:%f, texid:%lu, vtx_offset:%d, idx_offset:%d, elem_count:%d\n",
            cmd.clip_rect.x, cmd.clip_rect.y, cmd.clip_rect.z, cmd.clip_rect.w,
            cmd.texid, cmd.vtx_offset, cmd.idx_offset, cmd.elem_count);
    }
    printf("\n");
}

// Setup GLFW event handlers
static void _gb_set_ev_handlers(gb_state_t* s) {

    glfwSetKeyCallback(s->w, _gb_key_callback);
    glfwSetCharCallback(s->w, _gb_char_callback);
    glfwSetCursorPosCallback(s->w, _gb_cursor_pos_callback);
    glfwSetCursorEnterCallback(s->w, _gb_cursor_enter_callback);
    glfwSetMouseButtonCallback(s->w, _gb_mouse_button_callback);
    glfwSetScrollCallback(s->w, _gb_scroll_callback);
}

// Reserve event at the end of events array allocating memory if necessary and returns its pointer
static gb_event_t* _gb_ev_reserve(gb_state_t* s) {

    if (s->ev_count >= s->ev_cap) {
        int new_cap = s->ev_cap + 128;
        s->events = realloc(s->events, sizeof(gb_event_t) * new_cap);
        if (s->events == NULL) {
            fprintf(stderr, "No memory for events array");
            return NULL;
        }
        s->ev_cap = new_cap;
    }
    //printf("events:%d\n", s->ev_count + 1);
    return &s->events[s->ev_count++];
}

// Appends GLFW key event to events array
static void _gb_key_callback(GLFWwindow* win, int key, int scancode, int action, int mods) {

    gb_state_t* s = (gb_state_t*)(glfwGetWindowUserPointer(win));
    gb_event_t* ev = _gb_ev_reserve(s);
    if (ev == NULL) {
        return;
    }
    ev->type = EVENT_KEY;
    ev->argint[0] = key;
    ev->argint[1] = scancode;
    ev->argint[2] = action;
    ev->argint[3] = mods;
}

// Appends GLFW characted event to events array
static void _gb_char_callback(GLFWwindow* win, unsigned int codepoint) {

    gb_state_t* s = (gb_state_t*)(glfwGetWindowUserPointer(win));
    gb_event_t* ev = _gb_ev_reserve(s);
    if (ev == NULL) {
        return;
    }
    ev->type = EVENT_CHAR;
    ev->argint[0] = codepoint;
}

// Appends GLFW cursor position event to events array
static void _gb_cursor_pos_callback(GLFWwindow* win, double xpos, double ypos) {

    gb_state_t* s = (gb_state_t*)(glfwGetWindowUserPointer(win));
    gb_event_t* ev = _gb_ev_reserve(s);
    if (ev == NULL) {
        return;
    }
    ev->type = EVENT_CURSOR_POS;
    ev->argfloat[0] = xpos;
    ev->argfloat[1] = ypos;
}

// Appends GLFW cursor enter event to events array
static void _gb_cursor_enter_callback(GLFWwindow* win, int entered) {

    gb_state_t* s = (gb_state_t*)(glfwGetWindowUserPointer(win));
    gb_event_t* ev = _gb_ev_reserve(s);
    if (ev == NULL) {
        return;
    }
    ev->type = EVENT_CURSOR_ENTER;
    ev->argint[0] = entered;
}

// Appends GLFW mouse button event to events array
static void _gb_mouse_button_callback(GLFWwindow* win, int button, int action, int mods) {

    gb_state_t* s = (gb_state_t*)(glfwGetWindowUserPointer(win));
    gb_event_t* ev = _gb_ev_reserve(s);
    if (ev == NULL) {
        return;
    }
    ev->type = EVENT_MOUSE_BUTTON;
    ev->argint[0] = button;
    ev->argint[1] = action;
    ev->argint[2] = mods;
}

// Appends GLFW scroll event to events array
static void _gb_scroll_callback(GLFWwindow* win, double xoffset, double yoffset) {

    gb_state_t* s = (gb_state_t*)(glfwGetWindowUserPointer(win));
    gb_event_t* ev = _gb_ev_reserve(s);
    if (ev == NULL) {
        return;
    }
    ev->type = EVENT_SCROLL;
    ev->argfloat[0] = xoffset;
    ev->argfloat[1] = yoffset;
}


