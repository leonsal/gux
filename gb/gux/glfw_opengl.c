#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <stddef.h>
#include "GL/gl3w.h"
#include <GLFW/glfw3.h>
#include "libgux.h"

// Uncomment next line to enable check for OpenGL calls return errors.
#define GUX_GB_DEBUG
#ifdef GUX_GB_DEBUG
#include <stdio.h>
#define GL_CALL(_CALL) do { \
    _CALL; \
    GLenum gl_err = glGetError(); \
    if (gl_err != 0) { \
        fprintf(stderr, "GL error 0x%x returned from '%s'.\n", gl_err, #_CALL); \
        abort(); \
    }} while (0);
#else
#define GL_CALL(_CALL) _CALL
#endif


// Internal state per Window
typedef struct {
    gb_config_t     cfg;                // User configuration
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
    GLint           vao;                // Single VAO
    gb_frame_info_t frame;              // Frame info returned by gb_window_start_frame()
} gb_state_t;


// Forward declarations of internal functions
static void _gb_render(gb_state_t* s, gb_draw_list_t dl);
static bool _gb_init(gb_state_t* s, const char* glsl_version);
static void _gb_set_state(gb_state_t* s);
static bool _gb_create_objects(gb_state_t* s);
static void _gb_destroy_objects(gb_state_t* s);
static bool _gb_check_shader(GLuint handle, const char* desc, const char* src);
static bool _gb_check_program(GLuint handle, const char* desc);

// Include common internal functions
#include "common.c"

// Creates Graphics Backend window
gb_window_t gb_create_window(const char* title, int width, int height, gb_config_t* pcfg) {

    // Initialize configuration
    gb_config_t cfg = {};
    if (pcfg != NULL) {
        cfg = *pcfg;
    }

    // Setup error callback and initializes GLFW
    glfwSetErrorCallback(_gb_glfw_error_callback);
    if (glfwInit() == 0) {
        fprintf(stderr, "Error initializing GLFW\n");
        return NULL;
    }

    // Determines the api version from configuration
    int clientApi;
    int vmajor;
    int vminor;
    if (cfg.opengl.es) {
        clientApi = GLFW_OPENGL_ES_API;
        vmajor = 3;
        vminor = 1;
    } else {
        clientApi = GLFW_OPENGL_API;
        vmajor = 3;
        vminor = 3;
    }

    // Set GLFW hints
    glfwWindowHint(GLFW_CLIENT_API, clientApi);
    glfwWindowHint(GLFW_CONTEXT_VERSION_MAJOR, vmajor);
    glfwWindowHint(GLFW_CONTEXT_VERSION_MINOR, vminor);
    glfwWindowHint(GLFW_VISIBLE, GLFW_TRUE);
    if (!cfg.opengl.es) {
        glfwWindowHint(GLFW_OPENGL_PROFILE, GLFW_OPENGL_CORE_PROFILE);  // 3.2+ only
    }
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
	if (cfg.unlimited_rate) {
    	glfwSwapInterval(0);
	} else {
    	glfwSwapInterval(1);
	}
         
    // Creates and initializes backend state
    gb_state_t* s = _gb_alloc(sizeof(gb_state_t));
    if (s == NULL) {
        return NULL;
    }
    memset(s, 0, sizeof(gb_state_t));
    s->cfg = cfg;
    s->w = win;
    glfwSetWindowUserPointer(win, s);

    // Initialize OpenGL
    bool res = _gb_init(s, NULL);
    if (!res) {
        fprintf(stderr, "OpenGL initialization error");
        return NULL;
    }

    // Set window event handlers
    _gb_set_ev_handlers(s);

	// Creates cursors only once when the first window is opened
	if (g_window_count == 0) {
    	_gb_create_cursors();
	}
	g_window_count++;
    return s;
}

// Destroy the specified window
void gb_window_destroy(gb_window_t bw) {

    gb_state_t* s = (gb_state_t*)(bw);
    glfwDestroyWindow(s->w);
    _gb_free(s->frame.events);
    s->frame.events = NULL;
    _gb_free(s);

	// If all windows were closed, terminates GLFW
	g_window_count--;
	if (g_window_count <= 0) {
    	_gb_destroy_cursors();
    	glfwTerminate();
	}
}

// Starts the frame returning frame information
gb_frame_info_t* gb_window_start_frame(gb_window_t bw, gb_frame_params_t* params) {

    gb_state_t* s = (gb_state_t*)(bw);
    s->clear_color = params->clear_color;
    _gb_update_frame_info(s, params->ev_timeout);
    return &s->frame;
}

// Renders the frame
void gb_window_render_frame(gb_window_t bw, gb_draw_list_t dl) {

    // Sets OpenGL context from window
    gb_state_t* s = (gb_state_t*)(bw);
    glfwMakeContextCurrent(s->w);

    // Sets the OpenGL viewport from the framebuffer size
    int width, height;
    glfwGetFramebufferSize(s->w, &width, &height);
    GL_CALL(glViewport(0, 0, width, height));

    // Clears the framebuffer
    GL_CALL(glScissor(0, 0, width, height));
    GL_CALL(glClearColor(s->clear_color.x, s->clear_color.y, s->clear_color.z, s->clear_color.w));
    GL_CALL(glClear(GL_COLOR_BUFFER_BIT));

    // Render commands and swap buffers
    _gb_render(s, dl);
    glfwSwapBuffers(s->w);
}

// Sets window cursor
void gb_set_cursor(gb_window_t win, int cursor) {

    gb_state_t* s = (gb_state_t*)(win);
    _gb_set_cursor(s, cursor);
}

// Creates and returns an OpenGL texture identifier
gb_texid_t gb_create_texture(gb_window_t w, int width, int height, const gb_rgba_t* data) {

    // Sets OpenGL context from window
    gb_state_t* s = (gb_state_t*)(w);
    glfwMakeContextCurrent(s->w);

    // Create a OpenGL texture identifier
    GLuint image_texture;
    glGenTextures(1, &image_texture);
    glBindTexture(GL_TEXTURE_2D, image_texture);

    // Setup filtering parameters for display
    glTexParameteri(GL_TEXTURE_2D, GL_TEXTURE_MIN_FILTER, GL_LINEAR);
    glTexParameteri(GL_TEXTURE_2D, GL_TEXTURE_MAG_FILTER, GL_LINEAR);

    // Transfer data
    glPixelStorei(GL_UNPACK_ROW_LENGTH, 0);
    glTexImage2D(GL_TEXTURE_2D, 0, GL_RGBA, width, height, 0, GL_RGBA, GL_UNSIGNED_BYTE, data);
    return (intptr_t)image_texture;
}

// Deletes previously created texture
void gb_delete_texture(gb_window_t w, gb_texid_t texid) {

    // Sets OpenGL context from window
    gb_state_t* s = (gb_state_t*)(w);

    glfwMakeContextCurrent(s->w);
    GLuint tex = (GLuint)texid;
    glDeleteTextures(1, &tex); 
}


//-----------------------------------------------------------------------------
// Internal functions
//-----------------------------------------------------------------------------


// Executes draw commands
static void _gb_render(gb_state_t* s, gb_draw_list_t dl)  {

    // Do not render when minimized
    if (s->frame.fb_size.x <= 0 || s->frame.fb_size.y <= 0) {
        return;
    }

    // For OpenGL ES where glDrawElementsBaseVertex() is not available,
    // adjusts indices buffer for all commands before uploading it.
    // THIS MODIFIES THE USER DRAWLIST. SHOULD A TEMP BUFFER BE USED ???
    if (s->cfg.opengl.es) {
        for (int cmd_i = 0; cmd_i < dl.cmd_count; cmd_i++) {
            gb_draw_cmd_t* pcmd = &dl.buf_cmd[cmd_i];
            for (int idx = 0; idx < pcmd->elem_count; idx++) {
                dl.buf_idx[pcmd->idx_offset + idx] += pcmd->vtx_offset;
            }
        }
    }

    // Upload vertices and indices buffers
    const GLsizeiptr vtx_buffer_size = (GLsizeiptr)dl.vtx_count * (int)sizeof(gb_vertex_t);
    const GLsizeiptr idx_buffer_size = (GLsizeiptr)dl.idx_count * (int)sizeof(int);
    GL_CALL(glBufferData(GL_ARRAY_BUFFER, vtx_buffer_size, (const GLvoid*)dl.buf_vtx, GL_STREAM_DRAW));
    GL_CALL(glBufferData(GL_ELEMENT_ARRAY_BUFFER, idx_buffer_size, (const GLvoid*)dl.buf_idx, GL_STREAM_DRAW));

    // Sets orthogonal projection
    float L = 0;
    float R = 0 + s->frame.fb_size.x;
    float T = 0;
    float B = 0 + s->frame.fb_size.y;
    const float ortho_projection[4][4] = {
        { 2.0f/(R-L),   0.0f,         0.0f,   0.0f },
        { 0.0f,         2.0f/(T-B),   0.0f,   0.0f },
        { 0.0f,         0.0f,        -1.0f,   0.0f },
        { (R+L)/(L-R),  (T+B)/(B-T),  0.0f,   1.0f },
    };
    glUniformMatrix4fv(s->uni_projmtx, 1, GL_FALSE, &ortho_projection[0][0]);

    gb_vec2_t clip_off = {0,0};
    gb_vec2_t clip_scale = s->frame.fb_scale;

    // DrawCmd loop
    for (int i = 0; i < dl.cmd_count; i++) {
        gb_draw_cmd_t* pcmd = &dl.buf_cmd[i];

        // Project scissor/clipping rectangles into framebuffer space
        gb_vec2_t clip_min = {(pcmd->clip_rect.x - clip_off.x) * clip_scale.x, (pcmd->clip_rect.y - clip_off.y) * clip_scale.y};
        gb_vec2_t clip_max = {(pcmd->clip_rect.z - clip_off.x) * clip_scale.x, (pcmd->clip_rect.w - clip_off.y) * clip_scale.y};
        if (clip_max.x <= clip_min.x || clip_max.y <= clip_min.y) {
            continue;
        }
        // Apply scissor/clipping rectangle (Y is inverted in OpenGL)
        GL_CALL(glScissor((int)clip_min.x, (int)(s->frame.fb_size.y - clip_max.y), (int)(clip_max.x - clip_min.x), (int)(clip_max.y - clip_min.y)));

        // Set texture and draw 
        GL_CALL(glBindTexture(GL_TEXTURE_2D, (GLuint)(pcmd->texid)));
        if (s->cfg.opengl.es) {
            GL_CALL(glDrawElements(GL_TRIANGLES, (GLsizei)pcmd->elem_count, GL_UNSIGNED_INT, (void*)(intptr_t)(pcmd->idx_offset * sizeof(gb_index_t))));
        } else {
            GL_CALL(glDrawElementsBaseVertex(GL_TRIANGLES, (GLsizei)pcmd->elem_count, GL_UNSIGNED_INT,
                (void*)(intptr_t)(pcmd->idx_offset * sizeof(gb_index_t)), (GLint)pcmd->vtx_offset));
        }
    }

    if (s->cfg.debug_print_cmds) {
        _gb_print_draw_list(dl);
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

    //GL_CALL(glPolygonMode(GL_FRONT_AND_BACK, GL_LINE));

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

    //printf("loc:%d sizeof:%ld, offset:%ld\n", s->attrib_vtx_pos, sizeof(gb_vertex_t), offsetof(gb_vertex_t, pos));
    //printf("loc:%d sizeof:%ld, offset:%ld\n", s->attrib_vtx_uv, sizeof(gb_vertex_t), offsetof(gb_vertex_t, uv));
    //printf("loc:%d sizeof:%ld, offset:%ld\n", s->attrib_vtx_color, sizeof(gb_vertex_t), offsetof(gb_vertex_t, col));
    GL_CALL(glVertexAttribPointer(s->attrib_vtx_pos,   2, GL_FLOAT, GL_FALSE, sizeof(gb_vertex_t), (GLvoid*)offsetof(gb_vertex_t, pos)));
    GL_CALL(glVertexAttribPointer(s->attrib_vtx_uv,    2, GL_FLOAT, GL_FALSE, sizeof(gb_vertex_t), (GLvoid*)offsetof(gb_vertex_t, uv)));
    GL_CALL(glVertexAttribPointer(s->attrib_vtx_color, 4, GL_UNSIGNED_BYTE, GL_TRUE, sizeof(gb_vertex_t), (GLvoid*)offsetof(gb_vertex_t, col)));
}

static bool _gb_create_objects(gb_state_t* s) {

    const GLchar* vertex_shader_glsl_300_es =
        "#version 300 es\n"
        "precision highp float;\n"
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

    const GLchar* fragment_shader_glsl_300_es =
        "#version 300 es\n"
        "precision mediump float;\n"
        "uniform sampler2D Texture;\n"
        "in vec2 Frag_UV;\n"
        "in vec4 Frag_Color;\n"
        "layout (location = 0) out vec4 Out_Color;\n"
        "void main()\n"
        "{\n"
        "    Out_Color = Frag_Color * texture(Texture, Frag_UV.st);\n"
        "}\n";

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

    // Select shaders from the configuration
    const GLchar* vertex_shader;
    const GLchar* fragment_shader;
    if (s->cfg.opengl.es) {
        vertex_shader = vertex_shader_glsl_300_es;
        fragment_shader = fragment_shader_glsl_300_es;
    } else {
        vertex_shader = vertex_shader_glsl_330_core;
        fragment_shader = fragment_shader_glsl_330_core;
    }

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
    //printf("LOCS: %d/%d/%d/%d/%d\n", s->uni_tex, s->uni_projmtx, s->attrib_vtx_pos, s->attrib_vtx_uv, s->attrib_vtx_color);

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
            GLchar* buf = _gb_alloc(log_length + 1);
            glGetShaderInfoLog(handle, log_length, NULL, buf);
            fprintf(stderr, "%s\n", buf);
            _gb_free(buf);
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
            GLchar* buf = _gb_alloc(log_length + 1);
            glGetProgramInfoLog(handle, log_length, NULL, buf);
            fprintf(stderr, "%s\n", buf);
            _gb_free(buf);
        }
        return false;
    }
    return true;
}

