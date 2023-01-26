#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <stddef.h>
#include <assert.h>

#define VK_NO_PROTOTYPES
#include "volk.h"

#define GLFW_INCLUDE_NONE
#include <GLFW/glfw3.h>

#include "libgux.h"

// Size of a static C-style array.
#define GB_ARRAYSIZE(_ARR)  ((int)(sizeof(_ARR) / sizeof(*(_ARR))))

// Assert macro
#define GB_ASSERT(_EXPR)    assert(_EXPR)

// Check Vulkan return error code
#define GB_VK_CHECK(_ERR)   _gb_check_vk_result(_ERR, __LINE__);

// Vulkan texture information
struct vulkan_texinfo {
    VkImage                 vk_image;
    VkImageView             vk_image_view;
    VkDeviceMemory          vk_memory;
    VkDescriptorSet         vk_descriptor_set;
};

// Vulkan frame info for each swapchain image count
struct vulkan_frame {
    VkCommandPool           vk_command_pool;
    VkCommandBuffer         vk_command_buffer;
    VkFence                 vk_fence;
    VkImage                 vk_backbuffer;
    VkImageView             vk_backbuffer_view;
    VkFramebuffer           vk_framebuffer;
    VkSemaphore             vk_image_acquired_sema;
    VkSemaphore             vk_render_complete_sema;
};

// Vulkan frame buffers for each swapchain image count
struct vulkan_frame_buffers {
    VkDeviceMemory          vk_vertex_buffer_memory;
    VkDeviceMemory          vk_index_buffer_memory;
    VkDeviceSize            vk_vertex_buffer_size;
    VkDeviceSize            vk_index_buffer_size;
    VkBuffer                vk_vertex_buffer;
    VkBuffer                vk_index_buffer;
};

// Backend window state
typedef struct {
    gb_config_t                     cfg;
    GLFWwindow*                     w;
    gb_vec4_t                       clear_color;
    uint32_t                        min_image_count;
    uint32_t                        queue_family;
    VkSampleCountFlagBits           vk_msaa_samples;
    VkDeviceSize                    vk_buffer_memory_alignment;
    // Initialization fields
    VkInstance                      vk_instance;
    VkDebugReportCallbackEXT        vk_debug_report;
    VkPhysicalDevice                vk_physical_device;
    VkDevice                        vk_device;
    VkAllocationCallbacks*          vk_allocator;
    VkQueue                         vk_queue;
    VkDescriptorPool                vk_descriptor_pool;
    VkShaderModule                  vk_shader_module_vert;
    VkShaderModule                  vk_shader_module_frag;
    VkDescriptorSetLayout           vk_descriptor_set_layout;
    VkSampler                       vk_sampler;
    VkPipelineLayout                vk_pipeline_layout;
    VkPipelineCreateFlags           vk_pipeline_create_flags;
    VkPipelineCache                 vk_pipeline_cache;
    // Window related
    int                             width;
    int                             height;
    VkSurfaceKHR                    vk_surface;
    VkSurfaceFormatKHR              vk_surface_format;
    VkPresentModeKHR                vk_present_mode;
    VkSwapchainKHR                  vk_swapchain;
    VkRenderPass                    vk_render_pass;
    VkPipeline                      vk_pipeline;
    VkClearValue                    vk_clear_value;
    uint32_t                        subpass;
    uint32_t                        image_count;
    struct vulkan_frame*            vk_frames;
    struct vulkan_frame_buffers*    vk_buffers;
    uint32_t                        frame_index;
    uint32_t                        sema_index;
    uint32_t                        buffers_index;
    bool                            swapchain_rebuild;

    gb_frame_info_t                 frame;          // Frame info returned by gb_window_start_frame()
} gb_state_t;


// Forward declarations of internal functions
static void _gb_render(gb_state_t* s, gb_draw_list_t dl);
static void _gb_vulkan_render_draw_data(gb_state_t* s, gb_draw_list_t dl, VkCommandBuffer command_buffer);
static void _gb_frame_present(gb_state_t* s);
static void _gb_vulkan_setup_render_state(gb_state_t* s, gb_draw_list_t dl);
static void _gb_create_or_resize_buffer(gb_state_t* s, VkBuffer* buffer, VkDeviceMemory* buffer_memory,
    VkDeviceSize* p_buffer_size, size_t new_size, VkBufferUsageFlagBits usage);
static uint32_t _gb_vulkan_memory_type(gb_state_t* s, VkMemoryPropertyFlags properties, uint32_t type_bits);
static void _gb_setup_vulkan(gb_state_t* s, const char** extensions, uint32_t extensions_count);
static void _gb_setup_vulkan_window(gb_state_t* s, int width, int height);
static VkSurfaceFormatKHR _gb_select_surface_format(gb_state_t* s, const VkFormat* request_formats, int request_formats_count, VkColorSpaceKHR request_color_space);
static VkPresentModeKHR _gb_select_present_mode(gb_state_t* s, const VkPresentModeKHR* request_modes, int request_modes_count);
static void _gb_create_or_resize_window(gb_state_t* s, int width, int height);
static void _gb_create_window_swap_chain(gb_state_t* s, int w, int h);
static void _gb_create_window_command_buffers(gb_state_t* s);
static int  _gb_get_min_image_count_from_present_mode(VkPresentModeKHR present_mode);
static void _gb_set_min_image_count(gb_state_t* s, uint32_t min_image_count);
static void _gb_create_pipeline(gb_state_t* s);
static gb_texid_t _gb_create_texture(gb_state_t* s, int width, int height, const gb_rgba_t* pixels);
static void _gb_destroy_texture(gb_state_t* s, struct vulkan_texinfo* tex);
VkDescriptorSet _gb_create_tex_descriptor_set(gb_state_t* s, VkSampler sampler, VkImageView image_view, VkImageLayout image_layout);
static void _gb_create_shader_modules(gb_state_t* s);
static void _gb_destroy_window(gb_state_t* s);
static void _gb_destroy_frame(gb_state_t* s, struct vulkan_frame* fd);
static void _gb_destroy_frame_buffers(gb_state_t* s, struct vulkan_frame_buffers* buffers);
static void _gb_destroy_window_frame_buffers(gb_state_t* s);
static void _gb_destroy_vulkan(gb_state_t* s);
static void _gb_check_vk_result(VkResult err, int line);
static VKAPI_ATTR VkBool32 VKAPI_CALL _gb_debug_report(VkDebugReportFlagsEXT flags, VkDebugReportObjectTypeEXT objectType,
    uint64_t object, size_t location, int32_t messageCode, const char* pLayerPrefix, const char* pMessage, void* pUserData);
       
// Include common internal functions
#include "common.c"

// Creates Graphics Backend window
gb_window_t gb_create_window(const char* title, int width, int height, gb_config_t* cfg) {

    // Setup error callback and initializes GLFW
    glfwSetErrorCallback(_gb_glfw_error_callback);
    if (glfwInit() == 0) {
        fprintf(stderr, "Error initializing GLFW\n");
        return NULL;
    }

    // Creates GLFW window
    glfwWindowHint(GLFW_CLIENT_API, GLFW_NO_API);
    GLFWwindow* win = glfwCreateWindow(width, height, title, NULL, NULL);
    if (win == NULL) {
        return NULL;
    }

    // Setup Vulkan
    if (!glfwVulkanSupported()) {
        printf("GLFW: Vulkan Not Supported\n");
        return NULL;
    }

    // Load Vulkan functions
    VkResult res = volkInitialize();
    if (res != VK_SUCCESS) {
        printf("VOLK: Initialization error\n");
        return NULL;
    }

    // Creates and initializes backend state
    gb_state_t* s = _gb_alloc(sizeof(gb_state_t));
    s->w = win;
    s->min_image_count = 2;
    s->queue_family = (uint32_t)-1;
    s->vk_buffer_memory_alignment = 256;
    glfwSetWindowUserPointer(win, s);
    if (cfg != NULL) {
        s->cfg = *cfg;
    }

    // Get required Vulkan extensions from GLFW (WSI) and initializes Vulkan
    uint32_t extensions_count = 0;
    const char** extensions = glfwGetRequiredInstanceExtensions(&extensions_count);
    _gb_setup_vulkan(s, extensions, extensions_count);

    // Create Window Surface and initializes Vulkan window
    VkResult err = glfwCreateWindowSurface(s->vk_instance, win, s->vk_allocator, &s->vk_surface);
    GB_VK_CHECK(err);
    int w, h;
    glfwGetFramebufferSize(win, &w, &h);
    _gb_setup_vulkan_window(s, w, h);

    // Set window event handlers
    _gb_set_ev_handlers(s);

	// Creates cursors only once when the first window is opened
	if (g_window_count == 0) {
    	_gb_create_cursors();
	}
	g_window_count++;
    return s;
}

// Destroy graphics backend window
void gb_window_destroy(gb_window_t win) {

    gb_state_t* s = (gb_state_t*)(win);
    glfwDestroyWindow(s->w);
    _gb_destroy_window(s);
    _gb_destroy_vulkan(s);
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

    // Checks if user requested window close
    gb_state_t* s = (gb_state_t*)(bw);
    s->clear_color = params->clear_color;
    _gb_update_frame_info(s, params->ev_timeout);

    // Resize swap chain?
    if (s->swapchain_rebuild) {
        int width, height;
        glfwGetFramebufferSize(s->w, &width, &height);
        if (width > 0 && height > 0) {
            _gb_set_min_image_count(s, s->min_image_count);
            _gb_create_or_resize_window(s, width, height);
            s->frame_index  = 0;
            s->swapchain_rebuild = false;
        }
    }
    return &s->frame;
}

// Renders the frame draw list
void gb_window_render_frame(gb_window_t win, gb_draw_list_t dl) {

    gb_state_t* s = (gb_state_t*)(win);
    if (s->frame.win_size.x <= 0 || s->frame.win_size.y <= 0) {
        return;
    }
    s->vk_clear_value.color.float32[0] = s->clear_color.x * s->clear_color.w;
    s->vk_clear_value.color.float32[1] = s->clear_color.y * s->clear_color.w;
    s->vk_clear_value.color.float32[2] = s->clear_color.z * s->clear_color.w;
    s->vk_clear_value.color.float32[3] = s->clear_color.w;
    _gb_render(s, dl);
    _gb_frame_present(s);
}

// Sets window cursor
void gb_set_cursor(gb_window_t win, int cursor) {

    gb_state_t* s = (gb_state_t*)(win);
    _gb_set_cursor(s, cursor);
}

// Creates and returns texture
gb_texid_t gb_create_texture(gb_window_t win, int width, int height, const gb_rgba_t* data) {

    gb_state_t* s = (gb_state_t*)(win);
    return _gb_create_texture(s, width, height, data);
}

// Deletes the specified texture
void gb_delete_texture(gb_window_t win, gb_texid_t texid) {

    gb_state_t* s = (gb_state_t*)(win);
    _gb_destroy_texture(s, (struct vulkan_texinfo*)(texid));
}


//-----------------------------------------------------------------------------
// Internal functions
//-----------------------------------------------------------------------------

// Executes draw commands
static void _gb_render(gb_state_t* s, gb_draw_list_t dl) {

    VkResult err;

    VkSemaphore image_acquired_sema  = s->vk_frames[s->sema_index].vk_image_acquired_sema;
    VkSemaphore render_complete_sema = s->vk_frames[s->sema_index].vk_render_complete_sema;
    err = vkAcquireNextImageKHR(s->vk_device, s->vk_swapchain, UINT64_MAX, image_acquired_sema, VK_NULL_HANDLE, &s->frame_index);
    if (err == VK_ERROR_OUT_OF_DATE_KHR || err == VK_SUBOPTIMAL_KHR) {
        s->swapchain_rebuild = true;
        return;
    }
    GB_VK_CHECK(err);
    struct vulkan_frame* fd = &s->vk_frames[s->frame_index];

    {
        err = vkWaitForFences(s->vk_device, 1, &fd->vk_fence, VK_TRUE, UINT64_MAX);    // wait indefinitely instead of periodically checking
        GB_VK_CHECK(err);

        err = vkResetFences(s->vk_device, 1, &fd->vk_fence);
        GB_VK_CHECK(err);
    }
    {
        err = vkResetCommandPool(s->vk_device, fd->vk_command_pool, 0);
        GB_VK_CHECK(err);
        VkCommandBufferBeginInfo info = {};
        info.sType = VK_STRUCTURE_TYPE_COMMAND_BUFFER_BEGIN_INFO;
        info.flags |= VK_COMMAND_BUFFER_USAGE_ONE_TIME_SUBMIT_BIT;
        err = vkBeginCommandBuffer(fd->vk_command_buffer, &info);
        GB_VK_CHECK(err);
    }
    {
        VkRenderPassBeginInfo info = {};
        info.sType = VK_STRUCTURE_TYPE_RENDER_PASS_BEGIN_INFO;
        info.renderPass = s->vk_render_pass;
        info.framebuffer = fd->vk_framebuffer;
        info.renderArea.extent.width = s->width;
        info.renderArea.extent.height = s->height;
        info.clearValueCount = 1;
        info.pClearValues = &s->vk_clear_value;
        vkCmdBeginRenderPass(fd->vk_command_buffer, &info, VK_SUBPASS_CONTENTS_INLINE);
    }

    // Record primitives into command buffer
    _gb_vulkan_render_draw_data(s, dl, fd->vk_command_buffer);

    // Submit command buffer
    vkCmdEndRenderPass(fd->vk_command_buffer);
    {
        VkPipelineStageFlags wait_stage = VK_PIPELINE_STAGE_COLOR_ATTACHMENT_OUTPUT_BIT;
        VkSubmitInfo info = {};
        info.sType = VK_STRUCTURE_TYPE_SUBMIT_INFO;
        info.waitSemaphoreCount = 1;
        info.pWaitSemaphores = &image_acquired_sema;
        info.pWaitDstStageMask = &wait_stage;
        info.commandBufferCount = 1;
        info.pCommandBuffers = &fd->vk_command_buffer;
        info.signalSemaphoreCount = 1;
        info.pSignalSemaphores = &render_complete_sema;

        err = vkEndCommandBuffer(fd->vk_command_buffer);
        GB_VK_CHECK(err);
        err = vkQueueSubmit(s->vk_queue, 1, &info, fd->vk_fence);
        GB_VK_CHECK(err);
    }

    if (s->cfg.debug_print_cmds) {
        _gb_print_draw_list(dl);
    }
}

static void _gb_vulkan_render_draw_data(gb_state_t* s, gb_draw_list_t dl, VkCommandBuffer command_buffer) {

    // Do not render when minimized
    if (s->frame.fb_size.x <= 0 || s->frame.fb_size.y <= 0) {
        return;
    }

    // Allocate render buffers if necessary
    if (s->vk_buffers == NULL)  {
        s->vk_buffers = (struct vulkan_frame_buffers*)_gb_alloc(sizeof(struct vulkan_frame_buffers) * s->image_count);
    }
    s->buffers_index = (s->buffers_index + 1) % s->image_count;

    struct vulkan_frame* fd = &s->vk_frames[s->frame_index];
    struct vulkan_frame_buffers* buf = &s->vk_buffers[s->buffers_index];

    if (dl.vtx_count > 0) {
        // Create or resize the vertex/index buffers
        size_t vertex_size = dl.vtx_count * sizeof(gb_vertex_t);
        size_t index_size = dl.idx_count * sizeof(gb_index_t);
        if (buf->vk_vertex_buffer == VK_NULL_HANDLE || buf->vk_vertex_buffer_size < vertex_size) {
            _gb_create_or_resize_buffer(s, &buf->vk_vertex_buffer, &buf->vk_vertex_buffer_memory, &buf->vk_vertex_buffer_size, vertex_size, VK_BUFFER_USAGE_VERTEX_BUFFER_BIT);
        }
        if (buf->vk_index_buffer == VK_NULL_HANDLE || buf->vk_index_buffer_size < index_size) {
            _gb_create_or_resize_buffer(s, &buf->vk_index_buffer, &buf->vk_index_buffer_memory, &buf->vk_index_buffer_size, index_size, VK_BUFFER_USAGE_INDEX_BUFFER_BIT);
        }

        // Upload vertex/index data into a single contiguous GPU buffer
        gb_vertex_t* vtx_dst = NULL;
        gb_index_t*  idx_dst = NULL;
        VkResult err = vkMapMemory(s->vk_device, buf->vk_vertex_buffer_memory, 0, buf->vk_vertex_buffer_size, 0, (void**)(&vtx_dst));
        GB_VK_CHECK(err);
        err = vkMapMemory(s->vk_device, buf->vk_index_buffer_memory, 0, buf->vk_index_buffer_size, 0, (void**)(&idx_dst));
        GB_VK_CHECK(err);

        memcpy(vtx_dst, dl.buf_vtx, vertex_size);
        memcpy(idx_dst, dl.buf_idx, index_size);

        VkMappedMemoryRange range[2] = {};
        range[0].sType = VK_STRUCTURE_TYPE_MAPPED_MEMORY_RANGE;
        range[0].memory = buf->vk_vertex_buffer_memory;
        range[0].size = VK_WHOLE_SIZE;
        range[1].sType = VK_STRUCTURE_TYPE_MAPPED_MEMORY_RANGE;
        range[1].memory = buf->vk_index_buffer_memory;
        range[1].size = VK_WHOLE_SIZE;
        err = vkFlushMappedMemoryRanges(s->vk_device, 2, range);
        GB_VK_CHECK(err);
        vkUnmapMemory(s->vk_device, buf->vk_vertex_buffer_memory);
        vkUnmapMemory(s->vk_device, buf->vk_index_buffer_memory);
    }

    // Setup desired Vulkan state
    _gb_vulkan_setup_render_state(s, dl);

    // Will project scissor/clipping rectangles into framebuffer space
    gb_vec2_t clip_off = {0,0};
    gb_vec2_t clip_scale = s->frame.fb_scale;

    // Apply draw commands
    for (int cmd_i = 0; cmd_i < dl.cmd_count; cmd_i++) {
        gb_draw_cmd_t* pcmd = &dl.buf_cmd[cmd_i];
        // Project scissor/clipping rectangles into framebuffer space
        gb_vec2_t clip_min = {(pcmd->clip_rect.x - clip_off.x) * clip_scale.x, (pcmd->clip_rect.y - clip_off.y) * clip_scale.y};
        gb_vec2_t clip_max = {(pcmd->clip_rect.z - clip_off.x) * clip_scale.x, (pcmd->clip_rect.w - clip_off.y) * clip_scale.y};

        // Clamp to viewport as vkCmdSetScissor() won't accept values that are off bounds
        if (clip_min.x < 0.0f) { clip_min.x = 0.0f; }
        if (clip_min.y < 0.0f) { clip_min.y = 0.0f; }
        if (clip_max.x > s->frame.fb_size.x) { clip_max.x = s->frame.fb_size.x; }
        if (clip_max.y > s->frame.fb_size.y) { clip_max.y = s->frame.fb_size.y; }
        if (clip_max.x <= clip_min.x || clip_max.y <= clip_min.y) {
            continue;
        }

        // Apply scissor/clipping rectangle
        VkRect2D scissor;
        scissor.offset.x = (int32_t)(clip_min.x);
        scissor.offset.y = (int32_t)(clip_min.y);
        scissor.extent.width = (uint32_t)(clip_max.x - clip_min.x);
        scissor.extent.height = (uint32_t)(clip_max.y - clip_min.y);
        vkCmdSetScissor(command_buffer, 0, 1, &scissor);

        // Bind DescriptorSet with font or user texture
        struct vulkan_texinfo* texinfo = (struct vulkan_texinfo*)(pcmd->texid);
        VkDescriptorSet desc_set[1] = { texinfo->vk_descriptor_set };
        vkCmdBindDescriptorSets(command_buffer, VK_PIPELINE_BIND_POINT_GRAPHICS, s->vk_pipeline_layout, 0, 1, desc_set, 0, NULL);
        // Draw
        vkCmdDrawIndexed(command_buffer, pcmd->elem_count, 1, pcmd->idx_offset, pcmd->vtx_offset, 0);
    }

    // Note: at this point both vkCmdSetViewport() and vkCmdSetScissor() have been called.
    // Our last values will leak into user/application rendering IF:
    // - Your app uses a pipeline with VK_DYNAMIC_STATE_VIEWPORT or VK_DYNAMIC_STATE_SCISSOR dynamic state
    // - And you forgot to call vkCmdSetViewport() and vkCmdSetScissor() yourself to explicitly set that state.
    // If you use VK_DYNAMIC_STATE_VIEWPORT or VK_DYNAMIC_STATE_SCISSOR you are responsible for setting the values before rendering.
    // In theory we should aim to backup/restore those values but I am not sure this is possible.
    // We perform a call to vkCmdSetScissor() to set back a full viewport which is likely to fix things for 99% users but technically this is not perfect. (See github #4644)
    VkRect2D scissor = { { 0, 0 }, { (uint32_t)s->frame.fb_size.x, (uint32_t)s->frame.fb_size.y } };
    vkCmdSetScissor(command_buffer, 0, 1, &scissor);
}

static void _gb_frame_present(gb_state_t* s) {

    if (s->swapchain_rebuild) {
        return;
    }
    VkSemaphore render_complete_sema = s->vk_frames[s->sema_index].vk_render_complete_sema;
    VkPresentInfoKHR info = {};
    info.sType = VK_STRUCTURE_TYPE_PRESENT_INFO_KHR;
    info.waitSemaphoreCount = 1;
    info.pWaitSemaphores = &render_complete_sema;
    info.swapchainCount = 1;
    info.pSwapchains = &s->vk_swapchain;
    info.pImageIndices = &s->frame_index;
    VkResult err = vkQueuePresentKHR(s->vk_queue, &info);
    if (err == VK_ERROR_OUT_OF_DATE_KHR || err == VK_SUBOPTIMAL_KHR) {
        s->swapchain_rebuild = true;
        return;
    }
    GB_VK_CHECK(err);
    s->sema_index = (s->sema_index + 1) % s->image_count; // Now we can use the next set of semaphores
}

static void _gb_vulkan_setup_render_state(gb_state_t* s, gb_draw_list_t dl) {

    struct vulkan_frame* fd = &s->vk_frames[s->frame_index];
    struct vulkan_frame_buffers* buf = &s->vk_buffers[s->buffers_index];

    // Bind pipeline:
    {
        vkCmdBindPipeline(fd->vk_command_buffer, VK_PIPELINE_BIND_POINT_GRAPHICS, s->vk_pipeline);
    }

    // Bind Vertex And Index Buffer:
    if (dl.vtx_count > 0) {
        VkBuffer vertex_buffers[1] = { buf->vk_vertex_buffer };
        VkDeviceSize vertex_offset[1] = { 0 };
        vkCmdBindVertexBuffers(fd->vk_command_buffer, 0, 1, vertex_buffers, vertex_offset);
        vkCmdBindIndexBuffer(fd->vk_command_buffer, buf->vk_index_buffer, 0, sizeof(gb_index_t) == 2 ? VK_INDEX_TYPE_UINT16 : VK_INDEX_TYPE_UINT32);
    }

    // Setup viewport:
    {
        VkViewport viewport;
        viewport.x = 0;
        viewport.y = 0;
        viewport.width = s->frame.fb_size.x;
        viewport.height = s->frame.fb_size.y;
        viewport.minDepth = 0.0f;
        viewport.maxDepth = 1.0f;
        vkCmdSetViewport(fd->vk_command_buffer, 0, 1, &viewport);
    }

    // Setup scale and translation:
    // Our visible imgui space lies from draw_data->DisplayPps (top left) to draw_data->DisplayPos+data_data->DisplaySize (bottom right). DisplayPos is (0,0) for single viewport apps.
    {
        float scale[2];
        scale[0] = 2.0f / s->frame.fb_size.x;
        scale[1] = 2.0f / s->frame.fb_size.y;
        float translate[2];
        translate[0] = -1.0f - 0 * scale[0];
        translate[1] = -1.0f - 0 * scale[1];
        vkCmdPushConstants(fd->vk_command_buffer, s->vk_pipeline_layout, VK_SHADER_STAGE_VERTEX_BIT, sizeof(float) * 0, sizeof(float) * 2, scale);
        vkCmdPushConstants(fd->vk_command_buffer, s->vk_pipeline_layout, VK_SHADER_STAGE_VERTEX_BIT, sizeof(float) * 2, sizeof(float) * 2, translate);
    }
}

static void _gb_create_or_resize_buffer(gb_state_t* s, VkBuffer* buffer, VkDeviceMemory* buffer_memory,
    VkDeviceSize* p_buffer_size, size_t new_size, VkBufferUsageFlagBits usage) {

    VkResult err;
    if (*buffer != VK_NULL_HANDLE) {
        vkDestroyBuffer(s->vk_device, *buffer, s->vk_allocator);
    }
    if (*buffer_memory != VK_NULL_HANDLE) {
        vkFreeMemory(s->vk_device, *buffer_memory, s->vk_allocator);
    }

    VkDeviceSize vertex_buffer_size_aligned = ((new_size - 1) / s->vk_buffer_memory_alignment + 1) * s->vk_buffer_memory_alignment;
    VkBufferCreateInfo buffer_info = {};
    buffer_info.sType = VK_STRUCTURE_TYPE_BUFFER_CREATE_INFO;
    buffer_info.size = vertex_buffer_size_aligned;
    buffer_info.usage = usage;
    buffer_info.sharingMode = VK_SHARING_MODE_EXCLUSIVE;
    err = vkCreateBuffer(s->vk_device, &buffer_info, s->vk_allocator, buffer);
    GB_VK_CHECK(err);

    VkMemoryRequirements req;
    vkGetBufferMemoryRequirements(s->vk_device, *buffer, &req);
    s->vk_buffer_memory_alignment = (s->vk_buffer_memory_alignment > req.alignment) ? s->vk_buffer_memory_alignment : req.alignment;
    VkMemoryAllocateInfo alloc_info = {};
    alloc_info.sType = VK_STRUCTURE_TYPE_MEMORY_ALLOCATE_INFO;
    alloc_info.allocationSize = req.size;
    alloc_info.memoryTypeIndex = _gb_vulkan_memory_type(s, VK_MEMORY_PROPERTY_HOST_VISIBLE_BIT, req.memoryTypeBits);
    err = vkAllocateMemory(s->vk_device, &alloc_info, s->vk_allocator, buffer_memory);
    GB_VK_CHECK(err);

    err = vkBindBufferMemory(s->vk_device, *buffer, *buffer_memory, 0);
    GB_VK_CHECK(err);
    *p_buffer_size = req.size;
}

static uint32_t _gb_vulkan_memory_type(gb_state_t* s, VkMemoryPropertyFlags properties, uint32_t type_bits) {

    VkPhysicalDeviceMemoryProperties prop;
    vkGetPhysicalDeviceMemoryProperties(s->vk_physical_device, &prop);
    for (uint32_t i = 0; i < prop.memoryTypeCount; i++) {
        if ((prop.memoryTypes[i].propertyFlags & properties) == properties && type_bits & (1 << i)) {
            return i;
        }
    }
    GB_ASSERT(0);
    return 0xFFFFFFFF; // Unable to find memoryType
}

static void _gb_setup_vulkan(gb_state_t* s, const char** extensions, uint32_t extensions_count) {
    VkResult err;

    // Create Vulkan Instance with optional validation layer
    VkInstanceCreateInfo create_info = {
        .sType = VK_STRUCTURE_TYPE_INSTANCE_CREATE_INFO,
        .enabledExtensionCount = extensions_count,
        .ppEnabledExtensionNames = extensions,
    };
    if (s->cfg.vulkan.validation_layer) {
        // Enabling validation layers
        const char* layers[] = { "VK_LAYER_KHRONOS_validation" };
        create_info.enabledLayerCount = 1;
        create_info.ppEnabledLayerNames = layers;

        // Enable debug report extension (we need additional storage, so we duplicate the user array to add our new extension to it)
        const char** extensions_ext = (const char**)_gb_alloc(sizeof(const char*) * (extensions_count + 1));
        memcpy(extensions_ext, extensions, extensions_count * sizeof(const char*));
        extensions_ext[extensions_count] = "VK_EXT_debug_report";
        create_info.enabledExtensionCount = extensions_count + 1;
        create_info.ppEnabledExtensionNames = extensions_ext;

        // Create Vulkan Instance
        err = vkCreateInstance(&create_info, s->vk_allocator, &s->vk_instance);
        GB_VK_CHECK(err);
        _gb_free(extensions_ext);

        // Load Vulkan functions for the instance
        volkLoadInstance(s->vk_instance);

        // Get the function pointer (required for any extensions)
        PFN_vkCreateDebugReportCallbackEXT vkCreateDebugReportCallbackEXT = (PFN_vkCreateDebugReportCallbackEXT)vkGetInstanceProcAddr(s->vk_instance, "vkCreateDebugReportCallbackEXT");
        GB_ASSERT(vkCreateDebugReportCallbackEXT != NULL);

        // Setup the debug report callback
        VkDebugReportCallbackCreateInfoEXT debug_report_ci = {};
        debug_report_ci.sType = VK_STRUCTURE_TYPE_DEBUG_REPORT_CALLBACK_CREATE_INFO_EXT;
        debug_report_ci.flags = VK_DEBUG_REPORT_ERROR_BIT_EXT | VK_DEBUG_REPORT_WARNING_BIT_EXT | VK_DEBUG_REPORT_PERFORMANCE_WARNING_BIT_EXT;
        debug_report_ci.pfnCallback = _gb_debug_report;
        debug_report_ci.pUserData = NULL;
        err = vkCreateDebugReportCallbackEXT(s->vk_instance, &debug_report_ci, s->vk_allocator, &s->vk_debug_report);
        GB_VK_CHECK(err);
        printf("VULKAN validation layer installed\n");
    } else {
        // Create Vulkan Instance without any debug feature
        err = vkCreateInstance(&create_info, s->vk_allocator, &s->vk_instance);
        GB_VK_CHECK(err);
        
        // Load Vulkan functions for the instance
        volkLoadInstance(s->vk_instance);
    }

    // Select GPU
    uint32_t gpu_count;
    err = vkEnumeratePhysicalDevices(s->vk_instance, &gpu_count, NULL);
    GB_VK_CHECK(err);
    GB_ASSERT(gpu_count > 0);

    VkPhysicalDevice* gpus = (VkPhysicalDevice*)_gb_alloc(sizeof(VkPhysicalDevice) * gpu_count);
    err = vkEnumeratePhysicalDevices(s->vk_instance, &gpu_count, gpus);
    GB_VK_CHECK(err);

    // If a number >1 of GPUs got reported, find discrete GPU if present, or use first one available. This covers
    // most common cases (multi-gpu/integrated+dedicated graphics). Handling more complicated setups (multiple
    // dedicated GPUs) is out of scope of this sample.
    int use_gpu = 0;
    for (int i = 0; i < (int)gpu_count; i++) {
        VkPhysicalDeviceProperties properties;
        vkGetPhysicalDeviceProperties(gpus[i], &properties);
        if (properties.deviceType == VK_PHYSICAL_DEVICE_TYPE_DISCRETE_GPU) {
            use_gpu = i;
            break;
        }
    }
    s->vk_physical_device = gpus[use_gpu];
    _gb_free(gpus);

    // Select graphics queue family
    uint32_t count;
    vkGetPhysicalDeviceQueueFamilyProperties(s->vk_physical_device, &count, NULL);
    VkQueueFamilyProperties* queues = (VkQueueFamilyProperties*)_gb_alloc(sizeof(VkQueueFamilyProperties) * count);
    vkGetPhysicalDeviceQueueFamilyProperties(s->vk_physical_device, &count, queues);
    for (uint32_t i = 0; i < count; i++) {
        if (queues[i].queueFlags & VK_QUEUE_GRAPHICS_BIT) {
            s->queue_family = i;
            break;
        }
    }
    _gb_free(queues);
    GB_ASSERT(s->queue_family != (uint32_t)-1);

    // Create Logical Device (with 1 queue)
    int device_extension_count = 1;
    const char* device_extensions[] = { "VK_KHR_swapchain" };
    const float queue_priority[] = { 1.0f };
    VkDeviceQueueCreateInfo queue_info[1] = {};
    queue_info[0].sType = VK_STRUCTURE_TYPE_DEVICE_QUEUE_CREATE_INFO;
    queue_info[0].queueFamilyIndex = s->queue_family;
    queue_info[0].queueCount = 1;
    queue_info[0].pQueuePriorities = queue_priority;
    VkDeviceCreateInfo dev_create_info = {
        .sType = VK_STRUCTURE_TYPE_DEVICE_CREATE_INFO,
        .queueCreateInfoCount = sizeof(queue_info) / sizeof(queue_info[0]),
        .pQueueCreateInfos = queue_info,
        .enabledExtensionCount = device_extension_count,
        .ppEnabledExtensionNames = device_extensions,
    };
    err = vkCreateDevice(s->vk_physical_device, &dev_create_info, s->vk_allocator, &s->vk_device);
    GB_VK_CHECK(err);
    vkGetDeviceQueue(s->vk_device, s->queue_family, 0, &s->vk_queue);

    // Create Descriptor Pool
    VkDescriptorPoolSize pool_sizes[] =
    {
        { VK_DESCRIPTOR_TYPE_SAMPLER, 1000 },
        { VK_DESCRIPTOR_TYPE_COMBINED_IMAGE_SAMPLER, 1000 },
        { VK_DESCRIPTOR_TYPE_SAMPLED_IMAGE, 1000 },
        { VK_DESCRIPTOR_TYPE_STORAGE_IMAGE, 1000 },
        { VK_DESCRIPTOR_TYPE_UNIFORM_TEXEL_BUFFER, 1000 },
        { VK_DESCRIPTOR_TYPE_STORAGE_TEXEL_BUFFER, 1000 },
        { VK_DESCRIPTOR_TYPE_UNIFORM_BUFFER, 1000 },
        { VK_DESCRIPTOR_TYPE_STORAGE_BUFFER, 1000 },
        { VK_DESCRIPTOR_TYPE_UNIFORM_BUFFER_DYNAMIC, 1000 },
        { VK_DESCRIPTOR_TYPE_STORAGE_BUFFER_DYNAMIC, 1000 },
        { VK_DESCRIPTOR_TYPE_INPUT_ATTACHMENT, 1000 }
    };
    VkDescriptorPoolCreateInfo pool_info = {
        .sType = VK_STRUCTURE_TYPE_DESCRIPTOR_POOL_CREATE_INFO,
        .flags = VK_DESCRIPTOR_POOL_CREATE_FREE_DESCRIPTOR_SET_BIT,
        .maxSets = 1000 * GB_ARRAYSIZE(pool_sizes),
        .poolSizeCount = (uint32_t)GB_ARRAYSIZE(pool_sizes),
        .pPoolSizes = pool_sizes,
    };
    err = vkCreateDescriptorPool(s->vk_device, &pool_info, s->vk_allocator, &s->vk_descriptor_pool);
    GB_VK_CHECK(err);

    // Creates texture sampler
    if (!s->vk_sampler) {
        VkSamplerCreateInfo info = {};
        info.sType = VK_STRUCTURE_TYPE_SAMPLER_CREATE_INFO;
        info.magFilter = VK_FILTER_LINEAR;
        info.minFilter = VK_FILTER_LINEAR;
        info.mipmapMode = VK_SAMPLER_MIPMAP_MODE_LINEAR;
        info.addressModeU = VK_SAMPLER_ADDRESS_MODE_REPEAT;
        info.addressModeV = VK_SAMPLER_ADDRESS_MODE_REPEAT;
        info.addressModeW = VK_SAMPLER_ADDRESS_MODE_REPEAT;
        info.minLod = -1000;
        info.maxLod = 1000;
        info.maxAnisotropy = 1.0f;
        err = vkCreateSampler(s->vk_device, &info, s->vk_allocator, &s->vk_sampler);
        GB_VK_CHECK(err);
    }

    // Creates Descriptor Set Layout
    if (!s->vk_descriptor_set_layout) {
        VkSampler sampler[1] = {s->vk_sampler};
        VkDescriptorSetLayoutBinding binding[1] = {};
        binding[0].descriptorType = VK_DESCRIPTOR_TYPE_COMBINED_IMAGE_SAMPLER;
        binding[0].descriptorCount = 1;
        binding[0].stageFlags = VK_SHADER_STAGE_FRAGMENT_BIT;
        binding[0].pImmutableSamplers = sampler;
        VkDescriptorSetLayoutCreateInfo info = {};
        info.sType = VK_STRUCTURE_TYPE_DESCRIPTOR_SET_LAYOUT_CREATE_INFO;
        info.bindingCount = 1;
        info.pBindings = binding;
        err = vkCreateDescriptorSetLayout(s->vk_device, &info, s->vk_allocator, &s->vk_descriptor_set_layout);
        GB_VK_CHECK(err);
    }

    // Creates Pipeline Layout
    if (!s->vk_pipeline_layout) {
        // Using 'vec2 offset' and 'vec2 scale' instead of a full 3d projection matrix
        VkPushConstantRange push_constants[1] = {};
        push_constants[0].stageFlags = VK_SHADER_STAGE_VERTEX_BIT;
        push_constants[0].offset = sizeof(float) * 0;
        push_constants[0].size = sizeof(float) * 4;
        VkDescriptorSetLayout set_layout[1] = { s->vk_descriptor_set_layout };
        VkPipelineLayoutCreateInfo layout_info = {};
        layout_info.sType = VK_STRUCTURE_TYPE_PIPELINE_LAYOUT_CREATE_INFO;
        layout_info.setLayoutCount = 1;
        layout_info.pSetLayouts = set_layout;
        layout_info.pushConstantRangeCount = 1;
        layout_info.pPushConstantRanges = push_constants;
        err = vkCreatePipelineLayout(s->vk_device, &layout_info, s->vk_allocator, &s->vk_pipeline_layout);
        GB_VK_CHECK(err);
    }

    _gb_create_shader_modules(s);
}

static void _gb_setup_vulkan_window(gb_state_t* s, int width, int height) {


    // Check for WSI support
    VkBool32 res;
    vkGetPhysicalDeviceSurfaceSupportKHR(s->vk_physical_device, s->queue_family, s->vk_surface, &res);
    if (res != VK_TRUE) {
        fprintf(stderr, "Error no WSI support on physical device 0\n");
        exit(-1);
    }

    // Select Surface Format
    const VkFormat requestSurfaceImageFormat[] = { VK_FORMAT_B8G8R8A8_UNORM, VK_FORMAT_R8G8B8A8_UNORM, VK_FORMAT_B8G8R8_UNORM, VK_FORMAT_R8G8B8_UNORM };
    const VkColorSpaceKHR requestSurfaceColorSpace = VK_COLORSPACE_SRGB_NONLINEAR_KHR;
    s->vk_surface_format = _gb_select_surface_format(s, requestSurfaceImageFormat,(size_t)GB_ARRAYSIZE(requestSurfaceImageFormat), requestSurfaceColorSpace);

    // Select Present Mode
	if (s->cfg.unlimited_rate) {
    	VkPresentModeKHR present_modes[] = { VK_PRESENT_MODE_MAILBOX_KHR, VK_PRESENT_MODE_IMMEDIATE_KHR, VK_PRESENT_MODE_FIFO_KHR };
    	s->vk_present_mode = _gb_select_present_mode(s, &present_modes[0], GB_ARRAYSIZE(present_modes));
	} else {
    	VkPresentModeKHR present_modes[] = { VK_PRESENT_MODE_FIFO_KHR };
    	s->vk_present_mode = _gb_select_present_mode(s, &present_modes[0], GB_ARRAYSIZE(present_modes));
	}

    // Create SwapChain, RenderPass, Framebuffer, etc.
    GB_ASSERT(s->min_image_count >= 2);
    _gb_create_or_resize_window(s, width, height);
}

static VkSurfaceFormatKHR _gb_select_surface_format(gb_state_t* s, const VkFormat* request_formats, int request_formats_count, VkColorSpaceKHR request_color_space) {

    GB_ASSERT(request_formats != NULL);
    GB_ASSERT(request_formats_count > 0);

    // Per Spec Format and View Format are expected to be the same unless VK_IMAGE_CREATE_MUTABLE_BIT was set at image creation
    // Assuming that the default behavior is without setting this bit, there is no need for separate Swapchain image and image view format
    // Additionally several new color spaces were introduced with Vulkan Spec v1.0.40,
    // hence we must make sure that a format with the mostly available color space, VK_COLOR_SPACE_SRGB_NONLINEAR_KHR, is found and used.
    uint32_t avail_count;
    vkGetPhysicalDeviceSurfaceFormatsKHR(s->vk_physical_device, s->vk_surface, &avail_count, NULL);
    VkSurfaceFormatKHR* avail_format = _gb_alloc(sizeof(VkSurfaceFormatKHR) * avail_count);
    vkGetPhysicalDeviceSurfaceFormatsKHR(s->vk_physical_device, s->vk_surface, &avail_count, avail_format);

    // First check if only one format, VK_FORMAT_UNDEFINED, is available, which would imply that any format is available
    VkSurfaceFormatKHR ret;
    if (avail_count == 1) {
        if (avail_format[0].format == VK_FORMAT_UNDEFINED) {
            ret.format = request_formats[0];
            ret.colorSpace = request_color_space;
            _gb_free(avail_format);
            return ret;
        } else {
            // No point in searching another format
            ret = avail_format[0];
            _gb_free(avail_format);
            return ret;
        }
    } else {
        // Request several formats, the first found will be used
        for (int request_i = 0; request_i < request_formats_count; request_i++) {
            for (uint32_t avail_i = 0; avail_i < avail_count; avail_i++) {
                if (avail_format[avail_i].format == request_formats[request_i] && avail_format[avail_i].colorSpace == request_color_space) {
                    ret = avail_format[avail_i];
                    _gb_free(avail_format);
                    return ret;
                }
            }
        }
        // If none of the requested image formats could be found, use the first available
        ret = avail_format[0];
        _gb_free(avail_format);
        return ret;
    }
}

static VkPresentModeKHR _gb_select_present_mode(gb_state_t* s, const VkPresentModeKHR* request_modes, int request_modes_count) {

    assert(request_modes != NULL);
    assert(request_modes_count > 0);

    // Request a certain mode and confirm that it is available. If not use VK_PRESENT_MODE_FIFO_KHR which is mandatory
    uint32_t avail_count = 0;
    vkGetPhysicalDeviceSurfacePresentModesKHR(s->vk_physical_device, s->vk_surface, &avail_count, NULL);
    VkPresentModeKHR* avail_modes = _gb_alloc(sizeof(VkPresentModeKHR) * avail_count);
    vkGetPhysicalDeviceSurfacePresentModesKHR(s->vk_physical_device, s->vk_surface, &avail_count, avail_modes);

    for (int request_i = 0; request_i < request_modes_count; request_i++) {
        for (uint32_t avail_i = 0; avail_i < avail_count; avail_i++) {
            if (request_modes[request_i] == avail_modes[avail_i]) {
                _gb_free(avail_modes);
                return request_modes[request_i];
            }
        }
    }
    _gb_free(avail_modes);
    return VK_PRESENT_MODE_FIFO_KHR; // Always available
}

static void _gb_create_or_resize_window(gb_state_t* s, int width, int height) {

    _gb_create_window_swap_chain(s, width, height);
    _gb_create_window_command_buffers(s);
}

// Create window swap chain and also destroy old swap chain and in-flight frames data, if any.
static void _gb_create_window_swap_chain(gb_state_t* s, int w, int h) {

    VkResult err;
    VkSwapchainKHR old_swapchain = s->vk_swapchain;
    s->vk_swapchain = VK_NULL_HANDLE;
    err = vkDeviceWaitIdle(s->vk_device);
    GB_VK_CHECK(err);

    // We don't use ImGui_ImplVulkanH_DestroyWindow() because we want to preserve the old swapchain to create the new one.
    // Destroy old Framebuffer
    for (uint32_t i = 0; i < s->image_count; i++) {
        _gb_destroy_frame(s, &s->vk_frames[i]);
    }
    _gb_free(s->vk_frames);
    s->vk_frames = NULL;
    s->image_count  = 0;
    if (s->vk_render_pass) {
        vkDestroyRenderPass(s->vk_device, s->vk_render_pass, s->vk_allocator);
    }
    if (s->vk_pipeline) {
        vkDestroyPipeline(s->vk_device, s->vk_pipeline, s->vk_allocator);
    }

    // If min image count was not specified, request different count of images dependent on selected present mode
    if (s->min_image_count == 0) {
        s->min_image_count = _gb_get_min_image_count_from_present_mode(s->vk_present_mode);
    }

    // Create Swapchain
    {
        VkSwapchainCreateInfoKHR info = {
            .sType = VK_STRUCTURE_TYPE_SWAPCHAIN_CREATE_INFO_KHR,
            .surface = s->vk_surface,
            .minImageCount = s->min_image_count,
            .imageFormat = s->vk_surface_format.format,
            .imageColorSpace = s->vk_surface_format.colorSpace,
            .imageArrayLayers = 1,
            .imageUsage = VK_IMAGE_USAGE_COLOR_ATTACHMENT_BIT,
            .imageSharingMode = VK_SHARING_MODE_EXCLUSIVE,           // Assume that graphics family == present family
            .preTransform = VK_SURFACE_TRANSFORM_IDENTITY_BIT_KHR,
            .compositeAlpha = VK_COMPOSITE_ALPHA_OPAQUE_BIT_KHR,
            .presentMode = s->vk_present_mode,
            .clipped = VK_TRUE,
            .oldSwapchain = old_swapchain,
        };
        VkSurfaceCapabilitiesKHR cap;
        err = vkGetPhysicalDeviceSurfaceCapabilitiesKHR(s->vk_physical_device, s->vk_surface, &cap);
        GB_VK_CHECK(err);
        if (info.minImageCount < cap.minImageCount)
            info.minImageCount = cap.minImageCount;
        else if (cap.maxImageCount != 0 && info.minImageCount > cap.maxImageCount)
            info.minImageCount = cap.maxImageCount;

        if (cap.currentExtent.width == 0xffffffff) {
            info.imageExtent.width = s->width = w;
            info.imageExtent.height = s->height = h;
        }
        else {
            info.imageExtent.width = s->width = cap.currentExtent.width;
            info.imageExtent.height = s->height = cap.currentExtent.height;
        }
        err = vkCreateSwapchainKHR(s->vk_device, &info, s->vk_allocator, &s->vk_swapchain);
        GB_VK_CHECK(err);
        err = vkGetSwapchainImagesKHR(s->vk_device, s->vk_swapchain, &s->image_count, NULL);
        GB_VK_CHECK(err);
        VkImage backbuffers[16] = {};
        assert(s->image_count >= s->min_image_count);
        assert(s->image_count < GB_ARRAYSIZE(backbuffers));
        err = vkGetSwapchainImagesKHR(s->vk_device, s->vk_swapchain, &s->image_count, backbuffers);
        GB_VK_CHECK(err);

        assert(s->vk_frames == NULL);
        s->vk_frames = _gb_alloc(sizeof(struct vulkan_frame) * s->image_count);
        for (uint32_t i = 0; i < s->image_count; i++) {
            s->vk_frames[i].vk_backbuffer = backbuffers[i];
        }
    }
    if (old_swapchain) {
        vkDestroySwapchainKHR(s->vk_device, old_swapchain, s->vk_allocator);
    }

    // Create the Render Pass
    {
        VkAttachmentDescription attachment = {};
        attachment.format = s->vk_surface_format.format;
        attachment.samples = VK_SAMPLE_COUNT_1_BIT;
        attachment.loadOp = VK_ATTACHMENT_LOAD_OP_CLEAR;
        attachment.storeOp = VK_ATTACHMENT_STORE_OP_STORE;
        attachment.stencilLoadOp = VK_ATTACHMENT_LOAD_OP_DONT_CARE;
        attachment.stencilStoreOp = VK_ATTACHMENT_STORE_OP_DONT_CARE;
        attachment.initialLayout = VK_IMAGE_LAYOUT_UNDEFINED;
        attachment.finalLayout = VK_IMAGE_LAYOUT_PRESENT_SRC_KHR;
        VkAttachmentReference color_attachment = {};
        color_attachment.attachment = 0;
        color_attachment.layout = VK_IMAGE_LAYOUT_COLOR_ATTACHMENT_OPTIMAL;
        VkSubpassDescription subpass = {};
        subpass.pipelineBindPoint = VK_PIPELINE_BIND_POINT_GRAPHICS;
        subpass.colorAttachmentCount = 1;
        subpass.pColorAttachments = &color_attachment;
        VkSubpassDependency dependency = {};
        dependency.srcSubpass = VK_SUBPASS_EXTERNAL;
        dependency.dstSubpass = 0;
        dependency.srcStageMask = VK_PIPELINE_STAGE_COLOR_ATTACHMENT_OUTPUT_BIT;
        dependency.dstStageMask = VK_PIPELINE_STAGE_COLOR_ATTACHMENT_OUTPUT_BIT;
        dependency.srcAccessMask = 0;
        dependency.dstAccessMask = VK_ACCESS_COLOR_ATTACHMENT_WRITE_BIT;
        VkRenderPassCreateInfo info = {};
        info.sType = VK_STRUCTURE_TYPE_RENDER_PASS_CREATE_INFO;
        info.attachmentCount = 1;
        info.pAttachments = &attachment;
        info.subpassCount = 1;
        info.pSubpasses = &subpass;
        info.dependencyCount = 1;
        info.pDependencies = &dependency;
        err = vkCreateRenderPass(s->vk_device, &info, s->vk_allocator, &s->vk_render_pass);
        GB_VK_CHECK(err);

        // Creates the window pipeline
        _gb_create_pipeline(s);
    }

    // Create The Image Views
    {
        VkImageViewCreateInfo info = {};
        info.sType = VK_STRUCTURE_TYPE_IMAGE_VIEW_CREATE_INFO;
        info.viewType = VK_IMAGE_VIEW_TYPE_2D;
        info.format = s->vk_surface_format.format;
        info.components.r = VK_COMPONENT_SWIZZLE_R;
        info.components.g = VK_COMPONENT_SWIZZLE_G;
        info.components.b = VK_COMPONENT_SWIZZLE_B;
        info.components.a = VK_COMPONENT_SWIZZLE_A;
        VkImageSubresourceRange image_range = { VK_IMAGE_ASPECT_COLOR_BIT, 0, 1, 0, 1 };
        info.subresourceRange = image_range;
        for (uint32_t i = 0; i < s->image_count; i++) {
            struct vulkan_frame* fd = &s->vk_frames[i];
            info.image = fd->vk_backbuffer;
            err = vkCreateImageView(s->vk_device, &info, s->vk_allocator, &fd->vk_backbuffer_view);
            GB_VK_CHECK(err);
        }
    }

    // Create Framebuffer
    {
        VkImageView attachment[1];
        VkFramebufferCreateInfo info = {};
        info.sType = VK_STRUCTURE_TYPE_FRAMEBUFFER_CREATE_INFO;
        info.renderPass = s->vk_render_pass;
        info.attachmentCount = 1;
        info.pAttachments = attachment;
        info.width = s->width;
        info.height = s->height;
        info.layers = 1;
        for (uint32_t i = 0; i < s->image_count; i++) {
            struct vulkan_frame* fd = &s->vk_frames[i];
            attachment[0] = fd->vk_backbuffer_view;
            err = vkCreateFramebuffer(s->vk_device, &info, s->vk_allocator, &fd->vk_framebuffer);
            GB_VK_CHECK(err);
        }
    }
}

static void _gb_create_window_command_buffers(gb_state_t* s) {

    VkResult err;
    for (uint32_t i = 0; i < s->image_count; i++) {
        struct vulkan_frame* fd = &s->vk_frames[i];
        {
            VkCommandPoolCreateInfo info = {
                .sType = VK_STRUCTURE_TYPE_COMMAND_POOL_CREATE_INFO,
                .flags = VK_COMMAND_POOL_CREATE_RESET_COMMAND_BUFFER_BIT,
                .queueFamilyIndex = s->queue_family,
            };
            err = vkCreateCommandPool(s->vk_device, &info, s->vk_allocator, &fd->vk_command_pool);
            GB_VK_CHECK(err);
        }
        {
            VkCommandBufferAllocateInfo info = {
                .sType = VK_STRUCTURE_TYPE_COMMAND_BUFFER_ALLOCATE_INFO,
                .commandPool = fd->vk_command_pool,
                .level = VK_COMMAND_BUFFER_LEVEL_PRIMARY,
                .commandBufferCount = 1,
            };
            err = vkAllocateCommandBuffers(s->vk_device, &info, &fd->vk_command_buffer);
            GB_VK_CHECK(err);
        }
        {
            VkFenceCreateInfo info = {
                .sType = VK_STRUCTURE_TYPE_FENCE_CREATE_INFO,
                .flags = VK_FENCE_CREATE_SIGNALED_BIT,
            };
            err = vkCreateFence(s->vk_device, &info, s->vk_allocator, &fd->vk_fence);
            GB_VK_CHECK(err);
        }
        {
            VkSemaphoreCreateInfo info = {
                .sType = VK_STRUCTURE_TYPE_SEMAPHORE_CREATE_INFO,
            }; 
            err = vkCreateSemaphore(s->vk_device, &info, s->vk_allocator, &fd->vk_image_acquired_sema);
            GB_VK_CHECK(err);
            err = vkCreateSemaphore(s->vk_device, &info, s->vk_allocator, &fd->vk_render_complete_sema);
            GB_VK_CHECK(err);
        }
    }
}

static int _gb_get_min_image_count_from_present_mode(VkPresentModeKHR present_mode) {

    if (present_mode == VK_PRESENT_MODE_MAILBOX_KHR) {
        return 3;
    }
    if (present_mode == VK_PRESENT_MODE_FIFO_KHR || present_mode == VK_PRESENT_MODE_FIFO_RELAXED_KHR) {
        return 2;
    }
    if (present_mode == VK_PRESENT_MODE_IMMEDIATE_KHR) {
        return 1;
    }
    assert(0);
    return 1;
}

static void _gb_set_min_image_count(gb_state_t* s, uint32_t min_image_count) {

    assert(min_image_count >= 2);
    if (s->min_image_count == min_image_count) {
        return;
    }

    assert(0); // FIXME-VIEWPORT: Unsupported. Need to recreate all swap chains!
    //VkResult err = vkDeviceWaitIdle(s->vk_device);
    //GB_VK_CHECK(err);
    //_gb_destroy_all_viewports_render_buffers(s->vi.Device, s->vi.Allocator);
    //s->min_image_count = min_image_count;
}

static void _gb_create_pipeline(gb_state_t* s) {

    VkPipelineShaderStageCreateInfo stage[2] = {};
    stage[0].sType = VK_STRUCTURE_TYPE_PIPELINE_SHADER_STAGE_CREATE_INFO;
    stage[0].stage = VK_SHADER_STAGE_VERTEX_BIT;
    stage[0].module = s->vk_shader_module_vert;
    stage[0].pName = "main";
    stage[1].sType = VK_STRUCTURE_TYPE_PIPELINE_SHADER_STAGE_CREATE_INFO;
    stage[1].stage = VK_SHADER_STAGE_FRAGMENT_BIT;
    stage[1].module = s->vk_shader_module_frag;
    stage[1].pName = "main";

    VkVertexInputBindingDescription binding_desc[1] = {};
    binding_desc[0].stride = sizeof(gb_vertex_t);
    binding_desc[0].inputRate = VK_VERTEX_INPUT_RATE_VERTEX;

    VkVertexInputAttributeDescription attribute_desc[3] = {};
    attribute_desc[0].location = 0;
    attribute_desc[0].binding = binding_desc[0].binding;
    attribute_desc[0].format = VK_FORMAT_R32G32_SFLOAT;
    attribute_desc[0].offset = offsetof(gb_vertex_t, pos);
    attribute_desc[1].location = 1;
    attribute_desc[1].binding = binding_desc[0].binding;
    attribute_desc[1].format = VK_FORMAT_R32G32_SFLOAT;
    attribute_desc[1].offset = offsetof(gb_vertex_t, uv);
    attribute_desc[2].location = 2;
    attribute_desc[2].binding = binding_desc[0].binding;
    attribute_desc[2].format = VK_FORMAT_R8G8B8A8_UNORM;
    attribute_desc[2].offset = offsetof(gb_vertex_t, col);

    VkPipelineVertexInputStateCreateInfo vertex_info = {};
    vertex_info.sType = VK_STRUCTURE_TYPE_PIPELINE_VERTEX_INPUT_STATE_CREATE_INFO;
    vertex_info.vertexBindingDescriptionCount = 1;
    vertex_info.pVertexBindingDescriptions = binding_desc;
    vertex_info.vertexAttributeDescriptionCount = 3;
    vertex_info.pVertexAttributeDescriptions = attribute_desc;

    VkPipelineInputAssemblyStateCreateInfo ia_info = {};
    ia_info.sType = VK_STRUCTURE_TYPE_PIPELINE_INPUT_ASSEMBLY_STATE_CREATE_INFO;
    ia_info.topology = VK_PRIMITIVE_TOPOLOGY_TRIANGLE_LIST;

    VkPipelineViewportStateCreateInfo viewport_info = {};
    viewport_info.sType = VK_STRUCTURE_TYPE_PIPELINE_VIEWPORT_STATE_CREATE_INFO;
    viewport_info.viewportCount = 1;
    viewport_info.scissorCount = 1;

    VkPipelineRasterizationStateCreateInfo raster_info = {};
    raster_info.sType = VK_STRUCTURE_TYPE_PIPELINE_RASTERIZATION_STATE_CREATE_INFO;
    raster_info.polygonMode = VK_POLYGON_MODE_FILL;
    raster_info.cullMode = VK_CULL_MODE_NONE;
    raster_info.frontFace = VK_FRONT_FACE_COUNTER_CLOCKWISE;
    raster_info.lineWidth = 1.0f;

    VkPipelineMultisampleStateCreateInfo ms_info = {};
    ms_info.sType = VK_STRUCTURE_TYPE_PIPELINE_MULTISAMPLE_STATE_CREATE_INFO;
    ms_info.rasterizationSamples = (s->vk_msaa_samples != 0) ? s->vk_msaa_samples : VK_SAMPLE_COUNT_1_BIT;

    VkPipelineColorBlendAttachmentState color_attachment[1] = {};
    color_attachment[0].blendEnable = VK_TRUE;
    color_attachment[0].srcColorBlendFactor = VK_BLEND_FACTOR_SRC_ALPHA;
    color_attachment[0].dstColorBlendFactor = VK_BLEND_FACTOR_ONE_MINUS_SRC_ALPHA;
    color_attachment[0].colorBlendOp = VK_BLEND_OP_ADD;
    color_attachment[0].srcAlphaBlendFactor = VK_BLEND_FACTOR_ONE;
    color_attachment[0].dstAlphaBlendFactor = VK_BLEND_FACTOR_ONE_MINUS_SRC_ALPHA;
    color_attachment[0].alphaBlendOp = VK_BLEND_OP_ADD;
    color_attachment[0].colorWriteMask = VK_COLOR_COMPONENT_R_BIT | VK_COLOR_COMPONENT_G_BIT | VK_COLOR_COMPONENT_B_BIT | VK_COLOR_COMPONENT_A_BIT;

    VkPipelineDepthStencilStateCreateInfo depth_info = {};
    depth_info.sType = VK_STRUCTURE_TYPE_PIPELINE_DEPTH_STENCIL_STATE_CREATE_INFO;

    VkPipelineColorBlendStateCreateInfo blend_info = {};
    blend_info.sType = VK_STRUCTURE_TYPE_PIPELINE_COLOR_BLEND_STATE_CREATE_INFO;
    blend_info.attachmentCount = 1;
    blend_info.pAttachments = color_attachment;

    VkDynamicState dynamic_states[2] = { VK_DYNAMIC_STATE_VIEWPORT, VK_DYNAMIC_STATE_SCISSOR };
    VkPipelineDynamicStateCreateInfo dynamic_state = {};
    dynamic_state.sType = VK_STRUCTURE_TYPE_PIPELINE_DYNAMIC_STATE_CREATE_INFO;
    dynamic_state.dynamicStateCount = (uint32_t)GB_ARRAYSIZE(dynamic_states);
    dynamic_state.pDynamicStates = dynamic_states;

    VkGraphicsPipelineCreateInfo info = {};
    info.sType = VK_STRUCTURE_TYPE_GRAPHICS_PIPELINE_CREATE_INFO;
    info.flags = s->vk_pipeline_create_flags;
    info.stageCount = 2;
    info.pStages = stage;
    info.pVertexInputState = &vertex_info;
    info.pInputAssemblyState = &ia_info;
    info.pViewportState = &viewport_info;
    info.pRasterizationState = &raster_info;
    info.pMultisampleState = &ms_info;
    info.pDepthStencilState = &depth_info;
    info.pColorBlendState = &blend_info;
    info.pDynamicState = &dynamic_state;
    info.layout = s->vk_pipeline_layout;
    info.renderPass = s->vk_render_pass;
    info.subpass = s->subpass;
    VkResult err = vkCreateGraphicsPipelines(s->vk_device, s->vk_pipeline_cache, 1, &info, s->vk_allocator, &s->vk_pipeline);
    GB_VK_CHECK(err);
}

static gb_texid_t _gb_create_texture(gb_state_t* s, int width, int height, const gb_rgba_t* pixels)  {

    VkResult        err;
    VkDeviceMemory  uploadBufferMemory;
    VkBuffer        uploadBuffer;

    vkQueueWaitIdle(s->vk_queue);

    // Use any command queue
    VkCommandPool command_pool = s->vk_frames[s->frame_index].vk_command_pool;
    VkCommandBuffer command_buffer = s->vk_frames[s->frame_index].vk_command_buffer;
    err = vkResetCommandPool(s->vk_device, command_pool, 0);
    GB_VK_CHECK(err);
    VkCommandBufferBeginInfo begin_info = {};
    begin_info.sType = VK_STRUCTURE_TYPE_COMMAND_BUFFER_BEGIN_INFO;
    begin_info.flags |= VK_COMMAND_BUFFER_USAGE_ONE_TIME_SUBMIT_BIT;
    err = vkBeginCommandBuffer(command_buffer, &begin_info);
    GB_VK_CHECK(err);

    // Allocate texture info
    struct vulkan_texinfo* tex = _gb_alloc(sizeof(struct vulkan_texinfo));
    size_t upload_size = width * height * 4 * sizeof(char);

    // Create the Image:
    {
        VkImageCreateInfo info = {};
        info.sType = VK_STRUCTURE_TYPE_IMAGE_CREATE_INFO;
        info.imageType = VK_IMAGE_TYPE_2D;
        info.format = VK_FORMAT_R8G8B8A8_UNORM;
        info.extent.width = width;
        info.extent.height = height;
        info.extent.depth = 1;
        info.mipLevels = 1;
        info.arrayLayers = 1;
        info.samples = VK_SAMPLE_COUNT_1_BIT;
        info.tiling = VK_IMAGE_TILING_OPTIMAL;
        info.usage = VK_IMAGE_USAGE_SAMPLED_BIT | VK_IMAGE_USAGE_TRANSFER_DST_BIT;
        info.sharingMode = VK_SHARING_MODE_EXCLUSIVE;
        info.initialLayout = VK_IMAGE_LAYOUT_UNDEFINED;
        err = vkCreateImage(s->vk_device, &info, s->vk_allocator, &tex->vk_image);
        GB_VK_CHECK(err);
        VkMemoryRequirements req;
        vkGetImageMemoryRequirements(s->vk_device, tex->vk_image, &req);
        VkMemoryAllocateInfo alloc_info = {};
        alloc_info.sType = VK_STRUCTURE_TYPE_MEMORY_ALLOCATE_INFO;
        alloc_info.allocationSize = req.size;
        alloc_info.memoryTypeIndex = _gb_vulkan_memory_type(s, VK_MEMORY_PROPERTY_DEVICE_LOCAL_BIT, req.memoryTypeBits);
        err = vkAllocateMemory(s->vk_device, &alloc_info, s->vk_allocator, &tex->vk_memory);
        GB_VK_CHECK(err);
        err = vkBindImageMemory(s->vk_device, tex->vk_image, tex->vk_memory, 0);
        GB_VK_CHECK(err);
    }

    // Create the Image View:
    {
        VkImageViewCreateInfo info = {};
        info.sType = VK_STRUCTURE_TYPE_IMAGE_VIEW_CREATE_INFO;
        info.image = tex->vk_image;
        info.viewType = VK_IMAGE_VIEW_TYPE_2D;
        info.format = VK_FORMAT_R8G8B8A8_UNORM;
        info.subresourceRange.aspectMask = VK_IMAGE_ASPECT_COLOR_BIT;
        info.subresourceRange.levelCount = 1;
        info.subresourceRange.layerCount = 1;
        err = vkCreateImageView(s->vk_device, &info, s->vk_allocator, &tex->vk_image_view);
        GB_VK_CHECK(err);
    }

    // Create the Descriptor Set
    tex->vk_descriptor_set = _gb_create_tex_descriptor_set(s, s->vk_sampler, tex->vk_image_view, VK_IMAGE_LAYOUT_SHADER_READ_ONLY_OPTIMAL);

    // Create the Upload Buffer:
    {
        VkBufferCreateInfo buffer_info = {};
        buffer_info.sType = VK_STRUCTURE_TYPE_BUFFER_CREATE_INFO;
        buffer_info.size = upload_size;
        buffer_info.usage = VK_BUFFER_USAGE_TRANSFER_SRC_BIT;
        buffer_info.sharingMode = VK_SHARING_MODE_EXCLUSIVE;
        err = vkCreateBuffer(s->vk_device, &buffer_info, s->vk_allocator, &uploadBuffer);
        GB_VK_CHECK(err);
        VkMemoryRequirements req;
        vkGetBufferMemoryRequirements(s->vk_device, uploadBuffer, &req);
        s->vk_buffer_memory_alignment = (s->vk_buffer_memory_alignment > req.alignment) ? s->vk_buffer_memory_alignment : req.alignment;
        VkMemoryAllocateInfo alloc_info = {};
        alloc_info.sType = VK_STRUCTURE_TYPE_MEMORY_ALLOCATE_INFO;
        alloc_info.allocationSize = req.size;
        alloc_info.memoryTypeIndex = _gb_vulkan_memory_type(s, VK_MEMORY_PROPERTY_HOST_VISIBLE_BIT, req.memoryTypeBits);
        err = vkAllocateMemory(s->vk_device, &alloc_info, s->vk_allocator, &uploadBufferMemory);
        GB_VK_CHECK(err);
        err = vkBindBufferMemory(s->vk_device, uploadBuffer, uploadBufferMemory, 0);
        GB_VK_CHECK(err);
    }

    // Upload to Buffer:
    {
        char* map = NULL;
        err = vkMapMemory(s->vk_device, uploadBufferMemory, 0, upload_size, 0, (void**)(&map));
        GB_VK_CHECK(err);
        memcpy(map, pixels, upload_size);
        VkMappedMemoryRange range[1] = {};
        range[0].sType = VK_STRUCTURE_TYPE_MAPPED_MEMORY_RANGE;
        range[0].memory = uploadBufferMemory;
        range[0].size = upload_size;
        err = vkFlushMappedMemoryRanges(s->vk_device, 1, range);
        GB_VK_CHECK(err);
        vkUnmapMemory(s->vk_device, uploadBufferMemory);
    }

    // Copy to Image:
    {
        VkImageMemoryBarrier copy_barrier[1] = {};
        copy_barrier[0].sType = VK_STRUCTURE_TYPE_IMAGE_MEMORY_BARRIER;
        copy_barrier[0].dstAccessMask = VK_ACCESS_TRANSFER_WRITE_BIT;
        copy_barrier[0].oldLayout = VK_IMAGE_LAYOUT_UNDEFINED;
        copy_barrier[0].newLayout = VK_IMAGE_LAYOUT_TRANSFER_DST_OPTIMAL;
        copy_barrier[0].srcQueueFamilyIndex = VK_QUEUE_FAMILY_IGNORED;
        copy_barrier[0].dstQueueFamilyIndex = VK_QUEUE_FAMILY_IGNORED;
        copy_barrier[0].image = tex->vk_image;
        copy_barrier[0].subresourceRange.aspectMask = VK_IMAGE_ASPECT_COLOR_BIT;
        copy_barrier[0].subresourceRange.levelCount = 1;
        copy_barrier[0].subresourceRange.layerCount = 1;
        vkCmdPipelineBarrier(command_buffer, VK_PIPELINE_STAGE_HOST_BIT, VK_PIPELINE_STAGE_TRANSFER_BIT, 0, 0, NULL, 0, NULL, 1, copy_barrier);

        VkBufferImageCopy region = {};
        region.imageSubresource.aspectMask = VK_IMAGE_ASPECT_COLOR_BIT;
        region.imageSubresource.layerCount = 1;
        region.imageExtent.width = width;
        region.imageExtent.height = height;
        region.imageExtent.depth = 1;
        vkCmdCopyBufferToImage(command_buffer, uploadBuffer, tex->vk_image, VK_IMAGE_LAYOUT_TRANSFER_DST_OPTIMAL, 1, &region);

        VkImageMemoryBarrier use_barrier[1] = {};
        use_barrier[0].sType = VK_STRUCTURE_TYPE_IMAGE_MEMORY_BARRIER;
        use_barrier[0].srcAccessMask = VK_ACCESS_TRANSFER_WRITE_BIT;
        use_barrier[0].dstAccessMask = VK_ACCESS_SHADER_READ_BIT;
        use_barrier[0].oldLayout = VK_IMAGE_LAYOUT_TRANSFER_DST_OPTIMAL;
        use_barrier[0].newLayout = VK_IMAGE_LAYOUT_SHADER_READ_ONLY_OPTIMAL;
        use_barrier[0].srcQueueFamilyIndex = VK_QUEUE_FAMILY_IGNORED;
        use_barrier[0].dstQueueFamilyIndex = VK_QUEUE_FAMILY_IGNORED;
        use_barrier[0].image = tex->vk_image;
        use_barrier[0].subresourceRange.aspectMask = VK_IMAGE_ASPECT_COLOR_BIT;
        use_barrier[0].subresourceRange.levelCount = 1;
        use_barrier[0].subresourceRange.layerCount = 1;
        vkCmdPipelineBarrier(command_buffer, VK_PIPELINE_STAGE_TRANSFER_BIT, VK_PIPELINE_STAGE_FRAGMENT_SHADER_BIT, 0, 0, NULL, 0, NULL, 1, use_barrier);
    }

    VkSubmitInfo end_info = {};
    end_info.sType = VK_STRUCTURE_TYPE_SUBMIT_INFO;
    end_info.commandBufferCount = 1;
    end_info.pCommandBuffers = &command_buffer;
    err = vkEndCommandBuffer(command_buffer);
    GB_VK_CHECK(err);
    err = vkQueueSubmit(s->vk_queue, 1, &end_info, VK_NULL_HANDLE);
    GB_VK_CHECK(err);
    err = vkDeviceWaitIdle(s->vk_device);
    GB_VK_CHECK(err);
    
    if (uploadBuffer) {
        vkDestroyBuffer(s->vk_device, uploadBuffer, s->vk_allocator);
    }
    if (uploadBufferMemory) {
        vkFreeMemory(s->vk_device, uploadBufferMemory, s->vk_allocator);
    }
    return (gb_texid_t)(tex);
}

static void _gb_destroy_texture(gb_state_t* s, struct vulkan_texinfo* tex)  {

    if (tex == NULL) {
        return;
    }
    vkQueueWaitIdle(s->vk_queue);

    if (tex->vk_image_view != VK_NULL_HANDLE) {
        vkDestroyImageView(s->vk_device, tex->vk_image_view, s->vk_allocator);
        tex->vk_image_view = VK_NULL_HANDLE;
    }
    if (tex->vk_image != VK_NULL_HANDLE) {
        vkDestroyImage(s->vk_device, tex->vk_image, s->vk_allocator);
        tex->vk_image = VK_NULL_HANDLE;
    }
    if (tex->vk_memory != VK_NULL_HANDLE) {
        vkFreeMemory(s->vk_device, tex->vk_memory, s->vk_allocator);
        tex->vk_memory = VK_NULL_HANDLE;
    }
    if (tex->vk_descriptor_set != VK_NULL_HANDLE) {
        vkFreeDescriptorSets(s->vk_device, s->vk_descriptor_pool, 1, &tex->vk_descriptor_set);
        tex->vk_descriptor_set = VK_NULL_HANDLE;
    }
    _gb_free(tex);
}

VkDescriptorSet _gb_create_tex_descriptor_set(gb_state_t* s, VkSampler sampler, VkImageView image_view, VkImageLayout image_layout) {

    // Create Descriptor Set:
    VkDescriptorSet descriptor_set;
    {
        VkDescriptorSetAllocateInfo alloc_info = {};
        alloc_info.sType = VK_STRUCTURE_TYPE_DESCRIPTOR_SET_ALLOCATE_INFO;
        alloc_info.descriptorPool = s->vk_descriptor_pool;
        alloc_info.descriptorSetCount = 1;
        alloc_info.pSetLayouts = &s->vk_descriptor_set_layout;
        VkResult err = vkAllocateDescriptorSets(s->vk_device, &alloc_info, &descriptor_set);
        GB_VK_CHECK(err);
    }

    // Update the Descriptor Set:
    {
        VkDescriptorImageInfo desc_image[1] = {};
        desc_image[0].sampler = sampler;
        desc_image[0].imageView = image_view;
        desc_image[0].imageLayout = image_layout;
        VkWriteDescriptorSet write_desc[1] = {};
        write_desc[0].sType = VK_STRUCTURE_TYPE_WRITE_DESCRIPTOR_SET;
        write_desc[0].dstSet = descriptor_set;
        write_desc[0].descriptorCount = 1;
        write_desc[0].descriptorType = VK_DESCRIPTOR_TYPE_COMBINED_IMAGE_SAMPLER;
        write_desc[0].pImageInfo = desc_image;
        vkUpdateDescriptorSets(s->vk_device, 1, write_desc, 0, NULL);
    }
    return descriptor_set;
}


// glsl_shader.vert, compiled with:
// # glslangValidator -V -x -o glsl_shader.vert.u32 glsl_shader.vert
/*
#version 450 core
layout(location = 0) in vec2 aPos;
layout(location = 1) in vec2 aUV;
layout(location = 2) in vec4 aColor;
layout(push_constant) uniform uPushConstant { vec2 uScale; vec2 uTranslate; } pc;

out gl_PerVertex { vec4 gl_Position; };
layout(location = 0) out struct { vec4 Color; vec2 UV; } Out;

void main()
{
    Out.Color = aColor;
    Out.UV = aUV;
    gl_Position = vec4(aPos * pc.uScale + pc.uTranslate, 0, 1);
}
*/
static uint32_t __glsl_shader_vert_spv[] =
{
    0x07230203,0x00010000,0x00080001,0x0000002e,0x00000000,0x00020011,0x00000001,0x0006000b,
    0x00000001,0x4c534c47,0x6474732e,0x3035342e,0x00000000,0x0003000e,0x00000000,0x00000001,
    0x000a000f,0x00000000,0x00000004,0x6e69616d,0x00000000,0x0000000b,0x0000000f,0x00000015,
    0x0000001b,0x0000001c,0x00030003,0x00000002,0x000001c2,0x00040005,0x00000004,0x6e69616d,
    0x00000000,0x00030005,0x00000009,0x00000000,0x00050006,0x00000009,0x00000000,0x6f6c6f43,
    0x00000072,0x00040006,0x00000009,0x00000001,0x00005655,0x00030005,0x0000000b,0x0074754f,
    0x00040005,0x0000000f,0x6c6f4361,0x0000726f,0x00030005,0x00000015,0x00565561,0x00060005,
    0x00000019,0x505f6c67,0x65567265,0x78657472,0x00000000,0x00060006,0x00000019,0x00000000,
    0x505f6c67,0x7469736f,0x006e6f69,0x00030005,0x0000001b,0x00000000,0x00040005,0x0000001c,
    0x736f5061,0x00000000,0x00060005,0x0000001e,0x73755075,0x6e6f4368,0x6e617473,0x00000074,
    0x00050006,0x0000001e,0x00000000,0x61635375,0x0000656c,0x00060006,0x0000001e,0x00000001,
    0x61725475,0x616c736e,0x00006574,0x00030005,0x00000020,0x00006370,0x00040047,0x0000000b,
    0x0000001e,0x00000000,0x00040047,0x0000000f,0x0000001e,0x00000002,0x00040047,0x00000015,
    0x0000001e,0x00000001,0x00050048,0x00000019,0x00000000,0x0000000b,0x00000000,0x00030047,
    0x00000019,0x00000002,0x00040047,0x0000001c,0x0000001e,0x00000000,0x00050048,0x0000001e,
    0x00000000,0x00000023,0x00000000,0x00050048,0x0000001e,0x00000001,0x00000023,0x00000008,
    0x00030047,0x0000001e,0x00000002,0x00020013,0x00000002,0x00030021,0x00000003,0x00000002,
    0x00030016,0x00000006,0x00000020,0x00040017,0x00000007,0x00000006,0x00000004,0x00040017,
    0x00000008,0x00000006,0x00000002,0x0004001e,0x00000009,0x00000007,0x00000008,0x00040020,
    0x0000000a,0x00000003,0x00000009,0x0004003b,0x0000000a,0x0000000b,0x00000003,0x00040015,
    0x0000000c,0x00000020,0x00000001,0x0004002b,0x0000000c,0x0000000d,0x00000000,0x00040020,
    0x0000000e,0x00000001,0x00000007,0x0004003b,0x0000000e,0x0000000f,0x00000001,0x00040020,
    0x00000011,0x00000003,0x00000007,0x0004002b,0x0000000c,0x00000013,0x00000001,0x00040020,
    0x00000014,0x00000001,0x00000008,0x0004003b,0x00000014,0x00000015,0x00000001,0x00040020,
    0x00000017,0x00000003,0x00000008,0x0003001e,0x00000019,0x00000007,0x00040020,0x0000001a,
    0x00000003,0x00000019,0x0004003b,0x0000001a,0x0000001b,0x00000003,0x0004003b,0x00000014,
    0x0000001c,0x00000001,0x0004001e,0x0000001e,0x00000008,0x00000008,0x00040020,0x0000001f,
    0x00000009,0x0000001e,0x0004003b,0x0000001f,0x00000020,0x00000009,0x00040020,0x00000021,
    0x00000009,0x00000008,0x0004002b,0x00000006,0x00000028,0x00000000,0x0004002b,0x00000006,
    0x00000029,0x3f800000,0x00050036,0x00000002,0x00000004,0x00000000,0x00000003,0x000200f8,
    0x00000005,0x0004003d,0x00000007,0x00000010,0x0000000f,0x00050041,0x00000011,0x00000012,
    0x0000000b,0x0000000d,0x0003003e,0x00000012,0x00000010,0x0004003d,0x00000008,0x00000016,
    0x00000015,0x00050041,0x00000017,0x00000018,0x0000000b,0x00000013,0x0003003e,0x00000018,
    0x00000016,0x0004003d,0x00000008,0x0000001d,0x0000001c,0x00050041,0x00000021,0x00000022,
    0x00000020,0x0000000d,0x0004003d,0x00000008,0x00000023,0x00000022,0x00050085,0x00000008,
    0x00000024,0x0000001d,0x00000023,0x00050041,0x00000021,0x00000025,0x00000020,0x00000013,
    0x0004003d,0x00000008,0x00000026,0x00000025,0x00050081,0x00000008,0x00000027,0x00000024,
    0x00000026,0x00050051,0x00000006,0x0000002a,0x00000027,0x00000000,0x00050051,0x00000006,
    0x0000002b,0x00000027,0x00000001,0x00070050,0x00000007,0x0000002c,0x0000002a,0x0000002b,
    0x00000028,0x00000029,0x00050041,0x00000011,0x0000002d,0x0000001b,0x0000000d,0x0003003e,
    0x0000002d,0x0000002c,0x000100fd,0x00010038
};

// glsl_shader.frag, compiled with:
// # glslangValidator -V -x -o glsl_shader.frag.u32 glsl_shader.frag
/*
#version 450 core
layout(location = 0) out vec4 fColor;
layout(set=0, binding=0) uniform sampler2D sTexture;
layout(location = 0) in struct { vec4 Color; vec2 UV; } In;
void main()
{
    fColor = In.Color * texture(sTexture, In.UV.st);
}
*/
static uint32_t __glsl_shader_frag_spv[] =
{
    0x07230203,0x00010000,0x00080001,0x0000001e,0x00000000,0x00020011,0x00000001,0x0006000b,
    0x00000001,0x4c534c47,0x6474732e,0x3035342e,0x00000000,0x0003000e,0x00000000,0x00000001,
    0x0007000f,0x00000004,0x00000004,0x6e69616d,0x00000000,0x00000009,0x0000000d,0x00030010,
    0x00000004,0x00000007,0x00030003,0x00000002,0x000001c2,0x00040005,0x00000004,0x6e69616d,
    0x00000000,0x00040005,0x00000009,0x6c6f4366,0x0000726f,0x00030005,0x0000000b,0x00000000,
    0x00050006,0x0000000b,0x00000000,0x6f6c6f43,0x00000072,0x00040006,0x0000000b,0x00000001,
    0x00005655,0x00030005,0x0000000d,0x00006e49,0x00050005,0x00000016,0x78655473,0x65727574,
    0x00000000,0x00040047,0x00000009,0x0000001e,0x00000000,0x00040047,0x0000000d,0x0000001e,
    0x00000000,0x00040047,0x00000016,0x00000022,0x00000000,0x00040047,0x00000016,0x00000021,
    0x00000000,0x00020013,0x00000002,0x00030021,0x00000003,0x00000002,0x00030016,0x00000006,
    0x00000020,0x00040017,0x00000007,0x00000006,0x00000004,0x00040020,0x00000008,0x00000003,
    0x00000007,0x0004003b,0x00000008,0x00000009,0x00000003,0x00040017,0x0000000a,0x00000006,
    0x00000002,0x0004001e,0x0000000b,0x00000007,0x0000000a,0x00040020,0x0000000c,0x00000001,
    0x0000000b,0x0004003b,0x0000000c,0x0000000d,0x00000001,0x00040015,0x0000000e,0x00000020,
    0x00000001,0x0004002b,0x0000000e,0x0000000f,0x00000000,0x00040020,0x00000010,0x00000001,
    0x00000007,0x00090019,0x00000013,0x00000006,0x00000001,0x00000000,0x00000000,0x00000000,
    0x00000001,0x00000000,0x0003001b,0x00000014,0x00000013,0x00040020,0x00000015,0x00000000,
    0x00000014,0x0004003b,0x00000015,0x00000016,0x00000000,0x0004002b,0x0000000e,0x00000018,
    0x00000001,0x00040020,0x00000019,0x00000001,0x0000000a,0x00050036,0x00000002,0x00000004,
    0x00000000,0x00000003,0x000200f8,0x00000005,0x00050041,0x00000010,0x00000011,0x0000000d,
    0x0000000f,0x0004003d,0x00000007,0x00000012,0x00000011,0x0004003d,0x00000014,0x00000017,
    0x00000016,0x00050041,0x00000019,0x0000001a,0x0000000d,0x00000018,0x0004003d,0x0000000a,
    0x0000001b,0x0000001a,0x00050057,0x00000007,0x0000001c,0x00000017,0x0000001b,0x00050085,
    0x00000007,0x0000001d,0x00000012,0x0000001c,0x0003003e,0x00000009,0x0000001d,0x000100fd,
    0x00010038
};

// Create shader modules if not already created
static void _gb_create_shader_modules(gb_state_t* s) {

    // Create the shader modules
    if (s->vk_shader_module_vert == VK_NULL_HANDLE) {
        VkShaderModuleCreateInfo vert_info = {};
        vert_info.sType = VK_STRUCTURE_TYPE_SHADER_MODULE_CREATE_INFO;
        vert_info.codeSize = sizeof(__glsl_shader_vert_spv);
        vert_info.pCode = (uint32_t*)__glsl_shader_vert_spv;
        VkResult err = vkCreateShaderModule(s->vk_device, &vert_info, s->vk_allocator, &s->vk_shader_module_vert);
        GB_VK_CHECK(err);
    }
    if (s->vk_shader_module_frag == VK_NULL_HANDLE) {
        VkShaderModuleCreateInfo frag_info = {};
        frag_info.sType = VK_STRUCTURE_TYPE_SHADER_MODULE_CREATE_INFO;
        frag_info.codeSize = sizeof(__glsl_shader_frag_spv);
        frag_info.pCode = (uint32_t*)__glsl_shader_frag_spv;
        VkResult err = vkCreateShaderModule(s->vk_device, &frag_info, s->vk_allocator, &s->vk_shader_module_frag);
        GB_VK_CHECK(err);
    }
}

static void _gb_destroy_window(gb_state_t* s) {

    vkQueueWaitIdle(s->vk_queue);

    for (uint32_t i = 0; i < s->image_count; i++) {
        _gb_destroy_frame(s, &s->vk_frames[i]);
    }
    _gb_free(s->vk_frames);
    s->vk_frames = NULL;

    _gb_destroy_window_frame_buffers(s);

    vkDestroyPipeline(s->vk_device, s->vk_pipeline, s->vk_allocator);
    vkDestroyRenderPass(s->vk_device, s->vk_render_pass, s->vk_allocator);
    vkDestroySwapchainKHR(s->vk_device, s->vk_swapchain, s->vk_allocator);
    vkDestroySurfaceKHR(s->vk_instance, s->vk_surface, s->vk_allocator);
}

static void _gb_destroy_frame(gb_state_t* s, struct vulkan_frame* fd) {

    vkFreeCommandBuffers(s->vk_device, fd->vk_command_pool, 1, &fd->vk_command_buffer);
    fd->vk_command_buffer = VK_NULL_HANDLE;

    vkDestroyCommandPool(s->vk_device, fd->vk_command_pool, s->vk_allocator);
    fd->vk_command_pool = VK_NULL_HANDLE;

    vkDestroyFence(s->vk_device, fd->vk_fence, s->vk_allocator);
    fd->vk_fence = VK_NULL_HANDLE;

//    vkDestroyImage(s->vk_device, fd->vk_backbuffer, s->vk_allocator);
//    fd->vk_backbuffer = VK_NULL_HANDLE;

    vkDestroyImageView(s->vk_device, fd->vk_backbuffer_view, s->vk_allocator);
    fd->vk_backbuffer_view = VK_NULL_HANDLE;

    vkDestroyFramebuffer(s->vk_device, fd->vk_framebuffer, s->vk_allocator);
    fd->vk_framebuffer = VK_NULL_HANDLE;

    vkDestroySemaphore(s->vk_device, fd->vk_image_acquired_sema, s->vk_allocator);
    fd->vk_image_acquired_sema = VK_NULL_HANDLE;

    vkDestroySemaphore(s->vk_device, fd->vk_render_complete_sema, s->vk_allocator);
    fd->vk_render_complete_sema = VK_NULL_HANDLE;
}

static void _gb_destroy_frame_buffers(gb_state_t* s, struct vulkan_frame_buffers* buffers) {

    if (buffers->vk_vertex_buffer_memory != VK_NULL_HANDLE) {
        vkFreeMemory(s->vk_device, buffers->vk_vertex_buffer_memory, s->vk_allocator);
        buffers->vk_vertex_buffer_memory = VK_NULL_HANDLE;
    }
    if (buffers->vk_index_buffer_memory != VK_NULL_HANDLE) {
        vkFreeMemory(s->vk_device, buffers->vk_index_buffer_memory, s->vk_allocator);
        buffers->vk_index_buffer_memory = VK_NULL_HANDLE;
    }

    if (buffers->vk_vertex_buffer != VK_NULL_HANDLE) {
        vkDestroyBuffer(s->vk_device, buffers->vk_vertex_buffer, s->vk_allocator);
        buffers->vk_vertex_buffer = VK_NULL_HANDLE;
    }
    if (buffers->vk_index_buffer != VK_NULL_HANDLE) {
        vkDestroyBuffer(s->vk_device, buffers->vk_index_buffer, s->vk_allocator);
        buffers->vk_index_buffer = VK_NULL_HANDLE;
    }
    buffers->vk_vertex_buffer_size = 0;
    buffers->vk_index_buffer_size = 0;
}

// Destroy all Vulkan objects associated with the window
static void _gb_destroy_window_frame_buffers(gb_state_t* s) {

    if (s->vk_buffers == NULL) {
        return;
    }
    for (uint32_t n = 0; n < s->image_count; n++) {
        _gb_destroy_frame_buffers(s, &s->vk_buffers[n]);
    }
    _gb_free(s->vk_buffers);
    s->vk_buffers = NULL;
}

// Destroy all Vulkan common objects
static void _gb_destroy_vulkan(gb_state_t* s) {

    vkDestroyShaderModule(s->vk_device, s->vk_shader_module_frag, s->vk_allocator);
    vkDestroyShaderModule(s->vk_device, s->vk_shader_module_vert, s->vk_allocator);
    vkDestroyPipelineLayout(s->vk_device, s->vk_pipeline_layout, s->vk_allocator);
    vkDestroyDescriptorSetLayout(s->vk_device, s->vk_descriptor_set_layout, s->vk_allocator);
    vkDestroySampler(s->vk_device, s->vk_sampler, s->vk_allocator);
    vkDestroyDescriptorPool(s->vk_device, s->vk_descriptor_pool, s->vk_allocator);

    if (s->cfg.vulkan.validation_layer) {
        // Remove the debug report callback
        PFN_vkDestroyDebugReportCallbackEXT vkDestroyDebugReportCallbackEXT = (PFN_vkDestroyDebugReportCallbackEXT)vkGetInstanceProcAddr(s->vk_instance, "vkDestroyDebugReportCallbackEXT");
        vkDestroyDebugReportCallbackEXT(s->vk_instance, s->vk_debug_report, s->vk_allocator);
    }

    vkDestroyDevice(s->vk_device, s->vk_allocator);
    vkDestroyInstance(s->vk_instance, s->vk_allocator);
}

static void _gb_check_vk_result(VkResult err, int line) {

    if (err == 0) {
        return;
    }
    abort();
    fprintf(stderr, "Vulkan error: VkResult = %d at line:%d\n", err, line);
}

static VKAPI_ATTR VkBool32 VKAPI_CALL _gb_debug_report(VkDebugReportFlagsEXT flags, VkDebugReportObjectTypeEXT objectType,
    uint64_t object, size_t location, int32_t messageCode, const char* pLayerPrefix, const char* pMessage, void* pUserData) {

    (void)flags; (void)object; (void)location; (void)messageCode; (void)pUserData; (void)pLayerPrefix; // Unused arguments
    fprintf(stderr, "[vulkan] Debug report from ObjectType: %i\nMessage: %s\n\n", objectType, pMessage);
    return VK_FALSE;
}

