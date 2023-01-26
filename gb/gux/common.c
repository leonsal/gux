//
// This source file should be INCLUDED in glfw_opengl.c and glfw_vulkan.c
// It DOES NOT contain any graphics specific api calls.
//

static void _gb_create_cursors();
static void _gb_set_cursor(gb_state_t* s, int cursor);
static void _gb_destroy_cursors();
static void _gb_update_frame_info(gb_state_t* s, double timeout);
static void _gb_print_draw_list(gb_draw_list_t dl);
static void _gb_glfw_error_callback(int error, const char* description);
static void _gb_set_ev_handlers(gb_state_t* s);
static gb_event_t* _gb_ev_reserve(gb_state_t* s);
static void _gb_key_callback(GLFWwindow* win, int key, int scancode, int action, int mods);
static void _gb_char_callback(GLFWwindow* win, unsigned int codepoint);
static void _gb_cursor_pos_callback(GLFWwindow* win, double xpos, double ypos);
static void _gb_cursor_enter_callback(GLFWwindow* win, int entered);
static void _gb_mouse_button_callback(GLFWwindow* win, int button, int action, int mods);
static void _gb_scroll_callback(GLFWwindow* win, double xoffset, double yoffset);
static void* _gb_alloc(size_t count);
static void _gb_free(void* p);

// Global static data
static int g_window_count 		= 0;		// Current number of opened windows
static GLFWcursor** g_cursors	= NULL;		// GLFW cursors

static void _gb_create_cursors() {

    g_cursors = (GLFWcursor**)_gb_alloc(sizeof(GLFWcursor*) * _CURSOR_COUNT);
}

static void _gb_set_cursor(gb_state_t* s, int cursor) {

    if (cursor == 0) {
        glfwSetCursor(s->w, NULL);
        return;
    }
    if (g_cursors[cursor] != NULL) {
        glfwSetCursor(s->w, g_cursors[cursor]);
        return;
    }
    switch (cursor) {
        case CURSOR_ARROW:
            g_cursors[cursor] = glfwCreateStandardCursor(GLFW_ARROW_CURSOR);
            break;
        case CURSOR_IBEAM:
            g_cursors[cursor] = glfwCreateStandardCursor(GLFW_IBEAM_CURSOR);
            break;
        case CURSOR_CROSSHAIR:
            g_cursors[cursor] = glfwCreateStandardCursor(GLFW_CROSSHAIR_CURSOR);
            break;
        case CURSOR_HAND:
            g_cursors[cursor] = glfwCreateStandardCursor(GLFW_HAND_CURSOR);
            break;
        case CURSOR_HRESIZE:
            g_cursors[cursor] = glfwCreateStandardCursor(GLFW_HRESIZE_CURSOR);
            break;
        case CURSOR_VRESIZE:
            g_cursors[cursor] = glfwCreateStandardCursor(GLFW_VRESIZE_CURSOR);
            break;
        default:
            fprintf(stderr, "Invalid cursor:%d\n", cursor);
            abort();
    }
    glfwSetCursor(s->w, g_cursors[cursor]);
}

static void _gb_destroy_cursors() {

    for (int cursor = CURSOR_ARROW; cursor < _CURSOR_COUNT; cursor++) {
        GLFWcursor* c = g_cursors[cursor];
        if (c == NULL) {
            continue;
        }
        glfwDestroyCursor(c);
        g_cursors[cursor] = NULL;
    }
	_gb_free(g_cursors);
	g_cursors = NULL;
}


// Updates frame information at the start of the frame
static void _gb_update_frame_info(gb_state_t* s, double timeout) {

    // Checks if user requested window close
    s->frame.win_close = 0;
    if (glfwWindowShouldClose(s->w)) {
        s->frame.win_close = 1;
    }

    // Get window and framebuffer sizes and calculates framebuffer scale
    int width, height;
    glfwGetWindowSize(s->w, &width, &height);
    s->frame.win_size.x = (float)width;
    s->frame.win_size.y = (float)height;
    glfwGetFramebufferSize(s->w, &width, &height);
    s->frame.fb_size.x = (float)width;
    s->frame.fb_size.y = (float)height;
    if (s->frame.win_size.x > 0 && s->frame.win_size.y > 0) {
        s->frame.fb_scale.x = s->frame.fb_size.x / s->frame.win_size.x;
        s->frame.fb_scale.y = s->frame.fb_size.y / s->frame.win_size.y;
    }

    // Poll and handle events, blocking if no events for the specified timeout
    s->frame.ev_count = 0;
    glfwWaitEventsTimeout(timeout);
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

// GLFW error callback function
static void _gb_glfw_error_callback(int error, const char* description) {

    fprintf(stderr, "GLFW Error %d: %s\n", error, description);
    abort();
}

// Setup GLFW event handlers
static void _gb_set_ev_handlers(gb_state_t* s) {

    // Allocates initial events array
    s->frame.ev_count = 0;
    s->frame.ev_cap = 128;
    s->frame.events = _gb_alloc(sizeof(gb_event_t) * s->frame.ev_cap);

    // Install event callbacks
    glfwSetKeyCallback(s->w, _gb_key_callback);
    glfwSetCharCallback(s->w, _gb_char_callback);
    glfwSetCursorPosCallback(s->w, _gb_cursor_pos_callback);
    glfwSetCursorEnterCallback(s->w, _gb_cursor_enter_callback);
    glfwSetMouseButtonCallback(s->w, _gb_mouse_button_callback);
    glfwSetScrollCallback(s->w, _gb_scroll_callback);
}

// Reserve event at the end of events array allocating memory if necessary and returns its pointer
static gb_event_t* _gb_ev_reserve(gb_state_t* s) {

    if (s->frame.ev_count >= s->frame.ev_cap) {
        int new_cap = s->frame.ev_cap + 128;
        s->frame.events = realloc(s->frame.events, sizeof(gb_event_t) * new_cap);
        if (s->frame.events == NULL) {
            fprintf(stderr, "No memory for events array");
            return NULL;
        }
        s->frame.ev_cap = new_cap;
    }
    //printf("events:%d\n", s->ev_count + 1);
    return &s->frame.events[s->frame.ev_count++];
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

// Allocates and clears memory 
static void* _gb_alloc(size_t count) {

    void *p = malloc(count);
    if (p == NULL) {
        fprintf(stderr, "GB: NO MEMORY\n");
        abort();
    }
    memset(p, 0, count);
    return p;
}

// Free previous allocated memory
static void _gb_free(void* p) {

    free(p);
}


