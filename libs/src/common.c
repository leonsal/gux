//
// This source file should be INCLUDED in glfw_opengl.c and glfw_vulkan.c
// It DOES NOT contain any graphics api calls.
//

static void _gb_print_draw_list(gb_draw_list_t dl);
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

