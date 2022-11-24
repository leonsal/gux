#include <stdio.h>
#include <stdlib.h>
#include "GL/gl3w.h"
#include <GLFW/glfw3.h>
#include "libgux.h"


typedef struct {
    GLFWwindow* w;
    struct {
        float r; float g; float b; float a;
    } clearColor;

    GLuint vao;
} window_state_t;


bool _gl_init(const char* glsl_version);

gl_window_t gl_create_window(const char* title, int width, int height, gl_config_t* cfg) {

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

    // Load OpenGL functions
    int res = gl3wInit();
    if (res != 0) {
        fprintf(stderr, "GL3W error");
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
void gl_window_destroy(gl_window_t bw) {

    window_state_t* s = (window_state_t*)(bw);
    glfwDestroyWindow(s->w);
    glfwTerminate();
    free(s);
}

// Starts the frame or returns false if the window should be closed
bool gl_window_start_frame(gl_window_t bw, double timeout) {

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
void gl_window_render_frame(gl_window_t bw, gl_draw_list_t dl) {

    window_state_t* s = (window_state_t*)(bw);
    int display_w, display_h;
    glfwGetFramebufferSize(s->w, &display_w, &display_h);
    glViewport(0, 0, display_w, display_h);
    glClearColor(s->clearColor.r, s->clearColor.g, s->clearColor.b, s->clearColor.a);
    glClear(GL_COLOR_BUFFER_BIT);

    // Renders data here
    // TODO
    glfwSwapBuffers(s->w);
}




bool _gl_init(const char* glsl_version) {

    // Load OpenGL functions
    int res = gl3wInit();
    if (res != 0) {
        fprintf(stderr, "Failed to initialize OpenGL loader!\n");
        return false;
    }



    return true;
}



void _gl_set_state() {

    glEnable(GL_BLEND);
    glBlendEquation(GL_FUNC_ADD);
    glBlendFuncSeparate(GL_SRC_ALPHA, GL_ONE_MINUS_SRC_ALPHA, GL_ONE, GL_ONE_MINUS_SRC_ALPHA);
    glDisable(GL_CULL_FACE);
    glDisable(GL_DEPTH_TEST);
    glDisable(GL_STENCIL_TEST);
    glEnable(GL_SCISSOR_TEST);
}

bool _gl_createDeviceObjects() {






}



//
//// Ends rendering the frame
//void backend_window_end_frame(backend_window_t bw, backend_color_t clear) {
//
//    auto win = reinterpret_cast<GLFWwindow*>(bw);
//    ImGui::Render();
//    int display_w, display_h;
//    glfwGetFramebufferSize(win, &display_w, &display_h);
//    glViewport(0, 0, display_w, display_h);
//    glClearColor(clear.r, clear.g, clear.b, clear.w);
//    glClear(GL_COLOR_BUFFER_BIT);
//    ImGui_ImplOpenGL3_RenderDrawData(ImGui::GetDrawData());
//    glfwSwapBuffers(win);
//}
//
//void backend_update_fonts() {
//
//    ImGui_ImplOpenGL3_DestroyFontsTexture();
//    ImGui_ImplOpenGL3_CreateFontsTexture();
//}
//
//// Creates and returns an OpenGL texture idenfifier
//ImTextureID backend_create_texture() {
//
//    // Create a OpenGL texture identifier
//    GLuint image_texture;
//    glGenTextures(1, &image_texture);
//    glBindTexture(GL_TEXTURE_2D, image_texture);
//
//    // Setup filtering parameters for display
//    glTexParameteri(GL_TEXTURE_2D, GL_TEXTURE_MIN_FILTER, GL_LINEAR);
//    glTexParameteri(GL_TEXTURE_2D, GL_TEXTURE_MAG_FILTER, GL_LINEAR);
//
//    return (void*)(intptr_t)image_texture;
//}
//
//// Deletes previously created texture
//void backend_delete_texture(ImTextureID tid) {
//
//    GLuint tex = reinterpret_cast<intptr_t>(tid);
//    glDeleteTextures(1, &tex); 
//}
//
//// Transfer data for the specified texture
//void backend_transfer_texture(ImTextureID tid, int width, int height, const void* data) {
//
//    GLuint tex = reinterpret_cast<intptr_t>(tid);
//    glBindTexture(GL_TEXTURE_2D, tex);
//    glPixelStorei(GL_UNPACK_ROW_LENGTH, 0);
//    glTexImage2D(GL_TEXTURE_2D, 0, GL_RGBA, width, height, 0, GL_RGBA, GL_UNSIGNED_BYTE, data);
//}
//
//
//
//
//
//
