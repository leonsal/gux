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

// Size of a static C-style array. Don't use on pointers!
#define GB_ARRAYSIZE(_ARR)          ((int)(sizeof(_ARR) / sizeof(*(_ARR))))
#define GB_VULKAN_DEBUG_REPORT 1

struct vulkan_frame {
    VkCommandPool       CommandPool;
    VkCommandBuffer     CommandBuffer;
    VkFence             Fence;
    VkImage             Backbuffer;
    VkImageView         BackbufferView;
    VkFramebuffer       Framebuffer;
};

struct vulkan_frame_semaphores {
    VkSemaphore         ImageAcquiredSemaphore;
    VkSemaphore         RenderCompleteSemaphore;
};

struct vulkan_window {
    int                 Width;
    int                 Height;
    VkSwapchainKHR      Swapchain;
    VkSurfaceKHR        Surface;
    VkSurfaceFormatKHR  SurfaceFormat;
    VkPresentModeKHR    PresentMode;
    VkRenderPass        RenderPass;
    VkPipeline          Pipeline;               // The window pipeline may uses a different VkRenderPass than the one passed in ImGui_ImplVulkan_InitInfo
    bool                ClearEnable;
    VkClearValue        ClearValue;
    uint32_t            FrameIndex;             // Current frame being rendered to (0 <= FrameIndex < FrameInFlightCount)
    uint32_t            ImageCount;             // Number of simultaneous in-flight frames (returned by vkGetSwapchainImagesKHR, usually derived from min_image_count)
    uint32_t            SemaphoreIndex;         // Current set of swapchain wait semaphores we're using (needs to be distinct from per frame data)
    struct vulkan_frame*    Frames;
    struct vulkan_frame_semaphores*    FrameSemaphores;
    bool                SwapChainRebuild;
};

struct vulkan_init {
    VkInstance                      Instance;
    VkPhysicalDevice                PhysicalDevice;
    VkDevice                        Device;
    uint32_t                        QueueFamily;
    VkQueue                         Queue;
    VkPipelineCache                 PipelineCache;
    VkDescriptorPool                DescriptorPool;
    uint32_t                        Subpass;
    uint32_t                        MinImageCount;          // >= 2
    uint32_t                        ImageCount;             // >= MinImageCount
    VkSampleCountFlagBits           MSAASamples;            // >= VK_SAMPLE_COUNT_1_BIT (0 -> default to VK_SAMPLE_COUNT_1_BIT)
    const VkAllocationCallbacks*    Allocator;
    //void                            (*CheckVkResultFn)(VkResult err);
};

struct vulkan_data {
    VkRenderPass                RenderPass;
    VkDeviceSize                BufferMemoryAlignment;
    VkPipelineCreateFlags       PipelineCreateFlags;
    VkDescriptorSetLayout       DescriptorSetLayout;
    VkPipelineLayout            PipelineLayout;
    VkPipeline                  Pipeline;
    uint32_t                    Subpass;
    VkShaderModule              ShaderModuleVert;
    VkShaderModule              ShaderModuleFrag;

    // Font data
    VkSampler                   FontSampler;
    VkDeviceMemory              FontMemory;
    VkImage                     FontImage;
    VkImageView                 FontView;
    VkDescriptorSet             FontDescriptorSet;
    VkDeviceMemory              UploadBufferMemory;
    VkBuffer                    UploadBuffer;
    //// Render buffers for main window
    //ImGui_ImplVulkanH_WindowRenderBuffers MainWindowRenderBuffers;

    //ImGui_ImplVulkan_Data()
    //{
    //    memset((void*)this, 0, sizeof(*this));
    //    BufferMemoryAlignment = 256;
    //}
};

// Backend window state
typedef struct {
    GLFWwindow*             w;      // GLFW window pointer
    struct vulkan_init      vi;     // Vulkan initialization info                                    
    struct vulkan_window    vw;     // Vulkan window data
    struct vulkan_data      vd;     // Vulkan data
} gb_state_t;

// Global state
//static VkAllocationCallbacks*   g_Allocator = NULL;
//static VkInstance               g_Instance = VK_NULL_HANDLE;
//static VkPhysicalDevice         g_PhysicalDevice = VK_NULL_HANDLE;
//static VkDevice                 g_Device = VK_NULL_HANDLE;
//static uint32_t                 g_QueueFamily = (uint32_t)-1;
//static VkQueue                  g_Queue = VK_NULL_HANDLE;
static VkDebugReportCallbackEXT g_DebugReport = VK_NULL_HANDLE;
//static VkPipelineCache          g_PipelineCache = VK_NULL_HANDLE;
//static VkDescriptorPool         g_DescriptorPool = VK_NULL_HANDLE;
//static int                      g_MinImageCount = 2;
//static bool                     g_SwapChainRebuild = false;


// Forward declarations of internal functions
static void _gb_glfw_error_callback(int error, const char* description);
static void _gb_check_vk_result(VkResult err);
static void _gb_setup_vulkan(gb_state_t* s, const char** extensions, uint32_t extensions_count);
static void _gb_setup_vulkan_window(gb_state_t* s, struct vulkan_window* wd, VkSurfaceKHR surface, int width, int height);
static VkSurfaceFormatKHR _gb_select_surface_format(VkPhysicalDevice physical_device, VkSurfaceKHR surface,
    const VkFormat* request_formats, int request_formats_count, VkColorSpaceKHR request_color_space);
static VkPresentModeKHR _gb_select_present_mode(VkPhysicalDevice physical_device, VkSurfaceKHR surface,
    const VkPresentModeKHR* request_modes, int request_modes_count);
static void* _gb_alloc(size_t count);
static void _gb_create_or_resize_window(VkInstance instance, VkPhysicalDevice physical_device, VkDevice device, struct vulkan_window* wd,
    uint32_t queue_family, const VkAllocationCallbacks* allocator, int width, int height, uint32_t min_image_count);
static void _gb_create_window_swap_chain(VkPhysicalDevice physical_device, VkDevice device, struct vulkan_window* wd,
    const VkAllocationCallbacks* allocator, int w, int h, uint32_t min_image_count);
static void _gb_create_window_command_buffers(VkPhysicalDevice physical_device, VkDevice device, struct vulkan_window* wd,
    uint32_t queue_family, const VkAllocationCallbacks* allocator);
static int _gb_get_min_image_count_from_present_mode(VkPresentModeKHR present_mode);
static void _gb_set_min_image_count(gb_state_t* s, uint32_t min_image_count);
static void _gb_destroy_frame(VkDevice device, struct vulkan_frame* fd, const VkAllocationCallbacks* allocator);
static void _gb_destroy_frame_semaphores(VkDevice device, struct vulkan_frame_semaphores* fsd, const VkAllocationCallbacks* allocator);
static void _gb_destroy_all_viewports_render_buffers(VkDevice device, const VkAllocationCallbacks* allocator);
#ifdef GB_VULKAN_DEBUG_REPORT
static VKAPI_ATTR VkBool32 VKAPI_CALL _gb_debug_report(VkDebugReportFlagsEXT flags, VkDebugReportObjectTypeEXT objectType,
    uint64_t object, size_t location, int32_t messageCode, const char* pLayerPrefix, const char* pMessage, void* pUserData);
#endif // IMGUI_VULKAN_DEBUG_REPORT

// Creates Graphics Backend window
gb_window_t gb_create_window(const char* title, int width, int height, gb_config_t* cfg) {

    // Setup error callback and initializes GLFW
    glfwSetErrorCallback(_gb_glfw_error_callback);
    if (glfwInit() == 0) {
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
    memset(s, 0, sizeof(gb_state_t));
    s->w = win;
    s->vi.MinImageCount = 2;
    s->vi.QueueFamily = (uint32_t)-1;
    glfwSetWindowUserPointer(win, s);

    // Get required vulkan extensions from GLFW (WSI)
    uint32_t extensions_count = 0;
    const char** extensions = glfwGetRequiredInstanceExtensions(&extensions_count);
    _gb_setup_vulkan(s, extensions, extensions_count);

    // Create Window Surface
    VkSurfaceKHR surface;
    VkResult err = glfwCreateWindowSurface(s->vi.Instance, win, s->vi.Allocator, &surface);
    _gb_check_vk_result(err);

    // Create Framebuffers
    int w, h;
    glfwGetFramebufferSize(win, &w, &h);
    _gb_setup_vulkan_window(s, &s->vw, surface, w, h);
    return s;
}

void gb_window_destroy(gb_window_t win) {

    gb_state_t* s = (gb_state_t*)(win);
    VkResult err = vkDeviceWaitIdle(s->vi.Device);
    _gb_check_vk_result(err);

}

bool gb_window_start_frame(gb_window_t bw, double timeout) {

    // Checks if user requested window close
    gb_state_t* s = (gb_state_t*)(bw);
    if (glfwWindowShouldClose(s->w)) {
        return false;
    }

    // Poll and handle events, blocking if no events for the specified timeout
    glfwWaitEventsTimeout(timeout);

    // Resize swap chain?
    if (s->vw.SwapChainRebuild) {
        int width, height;
        glfwGetFramebufferSize(s->w, &width, &height);
        if (width > 0 && height > 0) {
            _gb_set_min_image_count(s, s->vi.MinImageCount);
            _gb_create_or_resize_window(s->vi.Instance, s->vi.PhysicalDevice, s->vi.Device, &s->vw, s->vi.QueueFamily,
                    s->vi.Allocator, width, height, s->vi.MinImageCount);
            //g_MainWindowData.FrameIndex = 0;
            s->vw.FrameIndex = 0;
            s->vw.SwapChainRebuild = false;
        }
    }
    return true;
}

void gb_window_render_frame(gb_window_t win, gb_draw_list_t dl) {

}

gb_texid_t gb_create_texture() {

    return 1;
}

void gb_delete_texture(gb_texid_t texid) {

}

void gb_transfer_texture(gb_texid_t texid, int width, int height, const gb_rgba_t* data) {

}

int gb_get_events(gb_window_t win, gb_event_t* events, int ev_count) {

}


//-----------------------------------------------------------------------------
// Internal functions
//-----------------------------------------------------------------------------


static void _gb_setup_vulkan(gb_state_t* s, const char** extensions, uint32_t extensions_count) {
    VkResult err;

    // Create Vulkan Instance
    {
        VkInstanceCreateInfo create_info = {};
        create_info.sType = VK_STRUCTURE_TYPE_INSTANCE_CREATE_INFO;
        create_info.enabledExtensionCount = extensions_count;
        create_info.ppEnabledExtensionNames = extensions;
#ifdef GB_VULKAN_DEBUG_REPORT
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
        err = vkCreateInstance(&create_info, s->vi.Allocator, &s->vi.Instance);
        _gb_check_vk_result(err);
        free(extensions_ext);

        // Load Vulkan functions for the instance
        volkLoadInstance(s->vi.Instance);

        // Get the function pointer (required for any extensions)
        PFN_vkCreateDebugReportCallbackEXT vkCreateDebugReportCallbackEXT = (PFN_vkCreateDebugReportCallbackEXT)vkGetInstanceProcAddr(s->vi.Instance, "vkCreateDebugReportCallbackEXT");
        assert(vkCreateDebugReportCallbackEXT != NULL);

        // Setup the debug report callback
        VkDebugReportCallbackCreateInfoEXT debug_report_ci = {};
        debug_report_ci.sType = VK_STRUCTURE_TYPE_DEBUG_REPORT_CALLBACK_CREATE_INFO_EXT;
        debug_report_ci.flags = VK_DEBUG_REPORT_ERROR_BIT_EXT | VK_DEBUG_REPORT_WARNING_BIT_EXT | VK_DEBUG_REPORT_PERFORMANCE_WARNING_BIT_EXT;
        debug_report_ci.pfnCallback = _gb_debug_report;
        debug_report_ci.pUserData = NULL;
        err = vkCreateDebugReportCallbackEXT(s->vi.Instance, &debug_report_ci, s->vi.Allocator, &g_DebugReport);
        _gb_check_vk_result(err);
#else
        // Create Vulkan Instance without any debug feature
        err = vkCreateInstance(&create_info, g_Allocator, &g_Instance);
        _gb_check_vk_result(err);
        //IM_UNUSED(g_DebugReport);
        
        // Load Vulkan functions for the instance
        volkLoadInstance(g_Instance);
#endif
    }

    // Select GPU
    {
        uint32_t gpu_count;
        err = vkEnumeratePhysicalDevices(s->vi.Instance, &gpu_count, NULL);
        _gb_check_vk_result(err);
        assert(gpu_count > 0);

        VkPhysicalDevice* gpus = (VkPhysicalDevice*)_gb_alloc(sizeof(VkPhysicalDevice) * gpu_count);
        err = vkEnumeratePhysicalDevices(s->vi.Instance, &gpu_count, gpus);
        _gb_check_vk_result(err);

        // If a number >1 of GPUs got reported, find discrete GPU if present, or use first one available. This covers
        // most common cases (multi-gpu/integrated+dedicated graphics). Handling more complicated setups (multiple
        // dedicated GPUs) is out of scope of this sample.
        int use_gpu = 0;
        for (int i = 0; i < (int)gpu_count; i++)
        {
            VkPhysicalDeviceProperties properties;
            vkGetPhysicalDeviceProperties(gpus[i], &properties);
            if (properties.deviceType == VK_PHYSICAL_DEVICE_TYPE_DISCRETE_GPU)
            {
                use_gpu = i;
                break;
            }
        }

        s->vi.PhysicalDevice = gpus[use_gpu];
        free(gpus);
    }

    // Select graphics queue family
    {
        uint32_t count;
        vkGetPhysicalDeviceQueueFamilyProperties(s->vi.PhysicalDevice, &count, NULL);
        VkQueueFamilyProperties* queues = (VkQueueFamilyProperties*)_gb_alloc(sizeof(VkQueueFamilyProperties) * count);
        vkGetPhysicalDeviceQueueFamilyProperties(s->vi.PhysicalDevice, &count, queues);
        for (uint32_t i = 0; i < count; i++) {
            if (queues[i].queueFlags & VK_QUEUE_GRAPHICS_BIT) {
                s->vi.QueueFamily = i;
                break;
            }
        }
        free(queues);
        assert(s->vi.QueueFamily != (uint32_t)-1);
    }

    // Create Logical Device (with 1 queue)
    {
        int device_extension_count = 1;
        const char* device_extensions[] = { "VK_KHR_swapchain" };
        const float queue_priority[] = { 1.0f };
        VkDeviceQueueCreateInfo queue_info[1] = {};
        queue_info[0].sType = VK_STRUCTURE_TYPE_DEVICE_QUEUE_CREATE_INFO;
        queue_info[0].queueFamilyIndex = s->vi.QueueFamily;
        queue_info[0].queueCount = 1;
        queue_info[0].pQueuePriorities = queue_priority;
        VkDeviceCreateInfo create_info = {};
        create_info.sType = VK_STRUCTURE_TYPE_DEVICE_CREATE_INFO;
        create_info.queueCreateInfoCount = sizeof(queue_info) / sizeof(queue_info[0]);
        create_info.pQueueCreateInfos = queue_info;
        create_info.enabledExtensionCount = device_extension_count;
        create_info.ppEnabledExtensionNames = device_extensions;
        err = vkCreateDevice(s->vi.PhysicalDevice, &create_info, s->vi.Allocator, &s->vi.Device);
        _gb_check_vk_result(err);
        vkGetDeviceQueue(s->vi.Device, s->vi.QueueFamily, 0, &s->vi.Queue);
    }

    // Create Descriptor Pool
    {
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
        VkDescriptorPoolCreateInfo pool_info = {};
        pool_info.sType = VK_STRUCTURE_TYPE_DESCRIPTOR_POOL_CREATE_INFO;
        pool_info.flags = VK_DESCRIPTOR_POOL_CREATE_FREE_DESCRIPTOR_SET_BIT;
        pool_info.maxSets = 1000 * GB_ARRAYSIZE(pool_sizes);
        pool_info.poolSizeCount = (uint32_t)GB_ARRAYSIZE(pool_sizes);
        pool_info.pPoolSizes = pool_sizes;
        err = vkCreateDescriptorPool(s->vi.Device, &pool_info, s->vi.Allocator, &s->vi.DescriptorPool);
        _gb_check_vk_result(err);
    }
}

static void _gb_setup_vulkan_window(gb_state_t* s, struct vulkan_window* wd, VkSurfaceKHR surface, int width, int height) {

    wd->Surface = surface;

    // Check for WSI support
    VkBool32 res;
    vkGetPhysicalDeviceSurfaceSupportKHR(s->vi.PhysicalDevice, s->vi.QueueFamily, wd->Surface, &res);
    if (res != VK_TRUE) {
        fprintf(stderr, "Error no WSI support on physical device 0\n");
        exit(-1);
    }

    // Select Surface Format
    const VkFormat requestSurfaceImageFormat[] = { VK_FORMAT_B8G8R8A8_UNORM, VK_FORMAT_R8G8B8A8_UNORM, VK_FORMAT_B8G8R8_UNORM, VK_FORMAT_R8G8B8_UNORM };
    const VkColorSpaceKHR requestSurfaceColorSpace = VK_COLORSPACE_SRGB_NONLINEAR_KHR;
    wd->SurfaceFormat = _gb_select_surface_format(s->vi.PhysicalDevice, wd->Surface, requestSurfaceImageFormat,
        (size_t)GB_ARRAYSIZE(requestSurfaceImageFormat), requestSurfaceColorSpace);

    // Select Present Mode
#ifdef GB_UNLIMITED_FRAME_RATE
    VkPresentModeKHR present_modes[] = { VK_PRESENT_MODE_MAILBOX_KHR, VK_PRESENT_MODE_IMMEDIATE_KHR, VK_PRESENT_MODE_FIFO_KHR };
#else
    VkPresentModeKHR present_modes[] = { VK_PRESENT_MODE_FIFO_KHR };
#endif
    wd->PresentMode = _gb_select_present_mode(s->vi.PhysicalDevice, wd->Surface, &present_modes[0], GB_ARRAYSIZE(present_modes));
    //printf("[vulkan] Selected PresentMode = %d\n", wd->PresentMode);

    // Create SwapChain, RenderPass, Framebuffer, etc.
    assert(s->vi.MinImageCount >= 2);
    _gb_create_or_resize_window(s->vi.Instance, s->vi.PhysicalDevice, s->vi.Device, wd, s->vi.QueueFamily,
        s->vi.Allocator, width, height, s->vi.MinImageCount);
}

static VkSurfaceFormatKHR _gb_select_surface_format(VkPhysicalDevice physical_device, VkSurfaceKHR surface,
    const VkFormat* request_formats, int request_formats_count, VkColorSpaceKHR request_color_space) {

    assert(request_formats != NULL);
    assert(request_formats_count > 0);

    // Per Spec Format and View Format are expected to be the same unless VK_IMAGE_CREATE_MUTABLE_BIT was set at image creation
    // Assuming that the default behavior is without setting this bit, there is no need for separate Swapchain image and image view format
    // Additionally several new color spaces were introduced with Vulkan Spec v1.0.40,
    // hence we must make sure that a format with the mostly available color space, VK_COLOR_SPACE_SRGB_NONLINEAR_KHR, is found and used.
    uint32_t avail_count;
    vkGetPhysicalDeviceSurfaceFormatsKHR(physical_device, surface, &avail_count, NULL);
    VkSurfaceFormatKHR* avail_format = _gb_alloc(sizeof(VkSurfaceFormatKHR) * avail_count);
    vkGetPhysicalDeviceSurfaceFormatsKHR(physical_device, surface, &avail_count, avail_format);

    // First check if only one format, VK_FORMAT_UNDEFINED, is available, which would imply that any format is available
    VkSurfaceFormatKHR ret;
    if (avail_count == 1) {
        if (avail_format[0].format == VK_FORMAT_UNDEFINED) {
            ret.format = request_formats[0];
            ret.colorSpace = request_color_space;
            free(avail_format);
            return ret;
        } else {
            // No point in searching another format
            ret = avail_format[0];
            free(avail_format);
            return ret;
        }
    } else {
        // Request several formats, the first found will be used
        for (int request_i = 0; request_i < request_formats_count; request_i++) {
            for (uint32_t avail_i = 0; avail_i < avail_count; avail_i++) {
                if (avail_format[avail_i].format == request_formats[request_i] && avail_format[avail_i].colorSpace == request_color_space) {
                    ret = avail_format[avail_i];
                    free(avail_format);
                    return ret;
                }
            }
        }
        // If none of the requested image formats could be found, use the first available
        ret = avail_format[0];
        free(avail_format);
        return ret;
    }
}

static VkPresentModeKHR _gb_select_present_mode(VkPhysicalDevice physical_device, VkSurfaceKHR surface,
    const VkPresentModeKHR* request_modes, int request_modes_count) {

    assert(request_modes != NULL);
    assert(request_modes_count > 0);

    // Request a certain mode and confirm that it is available. If not use VK_PRESENT_MODE_FIFO_KHR which is mandatory
    uint32_t avail_count = 0;
    vkGetPhysicalDeviceSurfacePresentModesKHR(physical_device, surface, &avail_count, NULL);
    VkPresentModeKHR* avail_modes = _gb_alloc(sizeof(VkPresentModeKHR) * avail_count);
    vkGetPhysicalDeviceSurfacePresentModesKHR(physical_device, surface, &avail_count, avail_modes);
    //for (uint32_t avail_i = 0; avail_i < avail_count; avail_i++)
    //    printf("[vulkan] avail_modes[%d] = %d\n", avail_i, avail_modes[avail_i]);

    for (int request_i = 0; request_i < request_modes_count; request_i++) {
        for (uint32_t avail_i = 0; avail_i < avail_count; avail_i++) {
            if (request_modes[request_i] == avail_modes[avail_i]) {
                free(avail_modes);
                return request_modes[request_i];
            }
        }
    }
    free(avail_modes);
    return VK_PRESENT_MODE_FIFO_KHR; // Always available
}

static void _gb_create_or_resize_window(VkInstance instance, VkPhysicalDevice physical_device, VkDevice device, struct vulkan_window* wd,
    uint32_t queue_family, const VkAllocationCallbacks* allocator, int width, int height, uint32_t min_image_count) {

    _gb_create_window_swap_chain(physical_device, device, wd, allocator, width, height, min_image_count);
    //ImGui_ImplVulkan_CreatePipeline(device, allocator, VK_NULL_HANDLE, wd->RenderPass, VK_SAMPLE_COUNT_1_BIT, &wd->Pipeline, g_VulkanInitInfo.Subpass);
    _gb_create_window_command_buffers(physical_device, device, wd, queue_family, allocator);
}

// Also destroy old swap chain and in-flight frames data, if any.
static void _gb_create_window_swap_chain(VkPhysicalDevice physical_device, VkDevice device, struct vulkan_window* wd,
    const VkAllocationCallbacks* allocator, int w, int h, uint32_t min_image_count) {

    VkResult err;
    VkSwapchainKHR old_swapchain = wd->Swapchain;
    wd->Swapchain = VK_NULL_HANDLE;
    err = vkDeviceWaitIdle(device);
    _gb_check_vk_result(err);

    // We don't use ImGui_ImplVulkanH_DestroyWindow() because we want to preserve the old swapchain to create the new one.
    // Destroy old Framebuffer
    for (uint32_t i = 0; i < wd->ImageCount; i++) {
        _gb_destroy_frame(device, &wd->Frames[i], allocator);
        _gb_destroy_frame_semaphores(device, &wd->FrameSemaphores[i], allocator);
    }
    free(wd->Frames);
    free(wd->FrameSemaphores);
    wd->Frames = NULL;
    wd->FrameSemaphores = NULL;
    wd->ImageCount = 0;
    if (wd->RenderPass) {
        vkDestroyRenderPass(device, wd->RenderPass, allocator);
    }
    if (wd->Pipeline) {
        vkDestroyPipeline(device, wd->Pipeline, allocator);
    }

    // If min image count was not specified, request different count of images dependent on selected present mode
    if (min_image_count == 0) {
        min_image_count = _gb_get_min_image_count_from_present_mode(wd->PresentMode);
    }

    // Create Swapchain
    {
        VkSwapchainCreateInfoKHR info = {};
        info.sType = VK_STRUCTURE_TYPE_SWAPCHAIN_CREATE_INFO_KHR;
        info.surface = wd->Surface;
        info.minImageCount = min_image_count;
        info.imageFormat = wd->SurfaceFormat.format;
        info.imageColorSpace = wd->SurfaceFormat.colorSpace;
        info.imageArrayLayers = 1;
        info.imageUsage = VK_IMAGE_USAGE_COLOR_ATTACHMENT_BIT;
        info.imageSharingMode = VK_SHARING_MODE_EXCLUSIVE;           // Assume that graphics family == present family
        info.preTransform = VK_SURFACE_TRANSFORM_IDENTITY_BIT_KHR;
        info.compositeAlpha = VK_COMPOSITE_ALPHA_OPAQUE_BIT_KHR;
        info.presentMode = wd->PresentMode;
        info.clipped = VK_TRUE;
        info.oldSwapchain = old_swapchain;
        VkSurfaceCapabilitiesKHR cap;
        err = vkGetPhysicalDeviceSurfaceCapabilitiesKHR(physical_device, wd->Surface, &cap);
        _gb_check_vk_result(err);
        if (info.minImageCount < cap.minImageCount)
            info.minImageCount = cap.minImageCount;
        else if (cap.maxImageCount != 0 && info.minImageCount > cap.maxImageCount)
            info.minImageCount = cap.maxImageCount;

        if (cap.currentExtent.width == 0xffffffff)
        {
            info.imageExtent.width = wd->Width = w;
            info.imageExtent.height = wd->Height = h;
        }
        else
        {
            info.imageExtent.width = wd->Width = cap.currentExtent.width;
            info.imageExtent.height = wd->Height = cap.currentExtent.height;
        }
        err = vkCreateSwapchainKHR(device, &info, allocator, &wd->Swapchain);
        _gb_check_vk_result(err);
        err = vkGetSwapchainImagesKHR(device, wd->Swapchain, &wd->ImageCount, NULL);
        _gb_check_vk_result(err);
        VkImage backbuffers[16] = {};
        assert(wd->ImageCount >= min_image_count);
        assert(wd->ImageCount < GB_ARRAYSIZE(backbuffers));
        err = vkGetSwapchainImagesKHR(device, wd->Swapchain, &wd->ImageCount, backbuffers);
        _gb_check_vk_result(err);

        assert(wd->Frames == NULL);
        wd->Frames = _gb_alloc(sizeof(struct vulkan_frame) * wd->ImageCount);
        wd->FrameSemaphores = _gb_alloc(sizeof(struct vulkan_frame_semaphores) * wd->ImageCount);
        memset(wd->Frames, 0, sizeof(wd->Frames[0]) * wd->ImageCount);
        memset(wd->FrameSemaphores, 0, sizeof(wd->FrameSemaphores[0]) * wd->ImageCount);
        for (uint32_t i = 0; i < wd->ImageCount; i++) {
            wd->Frames[i].Backbuffer = backbuffers[i];
        }
    }
    if (old_swapchain) {
        vkDestroySwapchainKHR(device, old_swapchain, allocator);
    }

    // Create the Render Pass
    {
        VkAttachmentDescription attachment = {};
        attachment.format = wd->SurfaceFormat.format;
        attachment.samples = VK_SAMPLE_COUNT_1_BIT;
        attachment.loadOp = wd->ClearEnable ? VK_ATTACHMENT_LOAD_OP_CLEAR : VK_ATTACHMENT_LOAD_OP_DONT_CARE;
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
        err = vkCreateRenderPass(device, &info, allocator, &wd->RenderPass);
        _gb_check_vk_result(err);

        // We do not create a pipeline by default as this is also used by examples' main.cpp,
        // but secondary viewport in multi-viewport mode may want to create one with:
        //ImGui_ImplVulkan_CreatePipeline(device, allocator, VK_NULL_HANDLE, wd->RenderPass, VK_SAMPLE_COUNT_1_BIT, &wd->Pipeline, bd->Subpass);
    }

    // Create The Image Views
    {
        VkImageViewCreateInfo info = {};
        info.sType = VK_STRUCTURE_TYPE_IMAGE_VIEW_CREATE_INFO;
        info.viewType = VK_IMAGE_VIEW_TYPE_2D;
        info.format = wd->SurfaceFormat.format;
        info.components.r = VK_COMPONENT_SWIZZLE_R;
        info.components.g = VK_COMPONENT_SWIZZLE_G;
        info.components.b = VK_COMPONENT_SWIZZLE_B;
        info.components.a = VK_COMPONENT_SWIZZLE_A;
        VkImageSubresourceRange image_range = { VK_IMAGE_ASPECT_COLOR_BIT, 0, 1, 0, 1 };
        info.subresourceRange = image_range;
        for (uint32_t i = 0; i < wd->ImageCount; i++)
        {
            struct vulkan_frame* fd = &wd->Frames[i];
            info.image = fd->Backbuffer;
            err = vkCreateImageView(device, &info, allocator, &fd->BackbufferView);
            _gb_check_vk_result(err);
        }
    }

    // Create Framebuffer
    {
        VkImageView attachment[1];
        VkFramebufferCreateInfo info = {};
        info.sType = VK_STRUCTURE_TYPE_FRAMEBUFFER_CREATE_INFO;
        info.renderPass = wd->RenderPass;
        info.attachmentCount = 1;
        info.pAttachments = attachment;
        info.width = wd->Width;
        info.height = wd->Height;
        info.layers = 1;
        for (uint32_t i = 0; i < wd->ImageCount; i++) {
            struct vulkan_frame* fd = &wd->Frames[i];
            attachment[0] = fd->BackbufferView;
            err = vkCreateFramebuffer(device, &info, allocator, &fd->Framebuffer);
            _gb_check_vk_result(err);
        }
    }
}

static void _gb_create_window_command_buffers(VkPhysicalDevice physical_device, VkDevice device, struct vulkan_window* wd,
    uint32_t queue_family, const VkAllocationCallbacks* allocator) {

    assert(physical_device != VK_NULL_HANDLE && device != VK_NULL_HANDLE);

    // Create Command Buffers
    VkResult err;
    for (uint32_t i = 0; i < wd->ImageCount; i++) {
        struct vulkan_frame* fd = &wd->Frames[i];
        struct vulkan_frame_semaphores* fsd = &wd->FrameSemaphores[i];
        {
            VkCommandPoolCreateInfo info = {};
            info.sType = VK_STRUCTURE_TYPE_COMMAND_POOL_CREATE_INFO;
            info.flags = VK_COMMAND_POOL_CREATE_RESET_COMMAND_BUFFER_BIT;
            info.queueFamilyIndex = queue_family;
            err = vkCreateCommandPool(device, &info, allocator, &fd->CommandPool);
            _gb_check_vk_result(err);
        }
        {
            VkCommandBufferAllocateInfo info = {};
            info.sType = VK_STRUCTURE_TYPE_COMMAND_BUFFER_ALLOCATE_INFO;
            info.commandPool = fd->CommandPool;
            info.level = VK_COMMAND_BUFFER_LEVEL_PRIMARY;
            info.commandBufferCount = 1;
            err = vkAllocateCommandBuffers(device, &info, &fd->CommandBuffer);
            _gb_check_vk_result(err);
        }
        {
            VkFenceCreateInfo info = {};
            info.sType = VK_STRUCTURE_TYPE_FENCE_CREATE_INFO;
            info.flags = VK_FENCE_CREATE_SIGNALED_BIT;
            err = vkCreateFence(device, &info, allocator, &fd->Fence);
            _gb_check_vk_result(err);
        }
        {
            VkSemaphoreCreateInfo info = {};
            info.sType = VK_STRUCTURE_TYPE_SEMAPHORE_CREATE_INFO;
            err = vkCreateSemaphore(device, &info, allocator, &fsd->ImageAcquiredSemaphore);
            _gb_check_vk_result(err);
            err = vkCreateSemaphore(device, &info, allocator, &fsd->RenderCompleteSemaphore);
            _gb_check_vk_result(err);
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
    if (s->vi.MinImageCount == min_image_count) {
        return;
    }

    assert(0); // FIXME-VIEWPORT: Unsupported. Need to recreate all swap chains!
    VkResult err = vkDeviceWaitIdle(s->vi.Device);
    _gb_check_vk_result(err);
    _gb_destroy_all_viewports_render_buffers(s->vi.Device, s->vi.Allocator);
    s->vi.MinImageCount = min_image_count;
}

static void _gb_create_pipeline(VkDevice device, const VkAllocationCallbacks* allocator, VkPipelineCache pipelineCache,
    VkRenderPass renderPass, VkSampleCountFlagBits MSAASamples, VkPipeline* pipeline, uint32_t subpass) {

//    ImGui_ImplVulkan_Data* bd = ImGui_ImplVulkan_GetBackendData();
//    ImGui_ImplVulkan_CreateShaderModules(device, allocator);
//
//    VkPipelineShaderStageCreateInfo stage[2] = {};
//    stage[0].sType = VK_STRUCTURE_TYPE_PIPELINE_SHADER_STAGE_CREATE_INFO;
//    stage[0].stage = VK_SHADER_STAGE_VERTEX_BIT;
//    stage[0].module = bd->ShaderModuleVert;
//    stage[0].pName = "main";
//    stage[1].sType = VK_STRUCTURE_TYPE_PIPELINE_SHADER_STAGE_CREATE_INFO;
//    stage[1].stage = VK_SHADER_STAGE_FRAGMENT_BIT;
//    stage[1].module = bd->ShaderModuleFrag;
//    stage[1].pName = "main";
//
//    VkVertexInputBindingDescription binding_desc[1] = {};
//    binding_desc[0].stride = sizeof(ImDrawVert);
//    binding_desc[0].inputRate = VK_VERTEX_INPUT_RATE_VERTEX;
//
//    VkVertexInputAttributeDescription attribute_desc[3] = {};
//    attribute_desc[0].location = 0;
//    attribute_desc[0].binding = binding_desc[0].binding;
//    attribute_desc[0].format = VK_FORMAT_R32G32_SFLOAT;
//    attribute_desc[0].offset = IM_OFFSETOF(ImDrawVert, pos);
//    attribute_desc[1].location = 1;
//    attribute_desc[1].binding = binding_desc[0].binding;
//    attribute_desc[1].format = VK_FORMAT_R32G32_SFLOAT;
//    attribute_desc[1].offset = IM_OFFSETOF(ImDrawVert, uv);
//    attribute_desc[2].location = 2;
//    attribute_desc[2].binding = binding_desc[0].binding;
//    attribute_desc[2].format = VK_FORMAT_R8G8B8A8_UNORM;
//    attribute_desc[2].offset = IM_OFFSETOF(ImDrawVert, col);
//
//    VkPipelineVertexInputStateCreateInfo vertex_info = {};
//    vertex_info.sType = VK_STRUCTURE_TYPE_PIPELINE_VERTEX_INPUT_STATE_CREATE_INFO;
//    vertex_info.vertexBindingDescriptionCount = 1;
//    vertex_info.pVertexBindingDescriptions = binding_desc;
//    vertex_info.vertexAttributeDescriptionCount = 3;
//    vertex_info.pVertexAttributeDescriptions = attribute_desc;
//
//    VkPipelineInputAssemblyStateCreateInfo ia_info = {};
//    ia_info.sType = VK_STRUCTURE_TYPE_PIPELINE_INPUT_ASSEMBLY_STATE_CREATE_INFO;
//    ia_info.topology = VK_PRIMITIVE_TOPOLOGY_TRIANGLE_LIST;
//
//    VkPipelineViewportStateCreateInfo viewport_info = {};
//    viewport_info.sType = VK_STRUCTURE_TYPE_PIPELINE_VIEWPORT_STATE_CREATE_INFO;
//    viewport_info.viewportCount = 1;
//    viewport_info.scissorCount = 1;
//
//    VkPipelineRasterizationStateCreateInfo raster_info = {};
//    raster_info.sType = VK_STRUCTURE_TYPE_PIPELINE_RASTERIZATION_STATE_CREATE_INFO;
//    raster_info.polygonMode = VK_POLYGON_MODE_FILL;
//    raster_info.cullMode = VK_CULL_MODE_NONE;
//    raster_info.frontFace = VK_FRONT_FACE_COUNTER_CLOCKWISE;
//    raster_info.lineWidth = 1.0f;
//
//    VkPipelineMultisampleStateCreateInfo ms_info = {};
//    ms_info.sType = VK_STRUCTURE_TYPE_PIPELINE_MULTISAMPLE_STATE_CREATE_INFO;
//    ms_info.rasterizationSamples = (MSAASamples != 0) ? MSAASamples : VK_SAMPLE_COUNT_1_BIT;
//
//    VkPipelineColorBlendAttachmentState color_attachment[1] = {};
//    color_attachment[0].blendEnable = VK_TRUE;
//    color_attachment[0].srcColorBlendFactor = VK_BLEND_FACTOR_SRC_ALPHA;
//    color_attachment[0].dstColorBlendFactor = VK_BLEND_FACTOR_ONE_MINUS_SRC_ALPHA;
//    color_attachment[0].colorBlendOp = VK_BLEND_OP_ADD;
//    color_attachment[0].srcAlphaBlendFactor = VK_BLEND_FACTOR_ONE;
//    color_attachment[0].dstAlphaBlendFactor = VK_BLEND_FACTOR_ONE_MINUS_SRC_ALPHA;
//    color_attachment[0].alphaBlendOp = VK_BLEND_OP_ADD;
//    color_attachment[0].colorWriteMask = VK_COLOR_COMPONENT_R_BIT | VK_COLOR_COMPONENT_G_BIT | VK_COLOR_COMPONENT_B_BIT | VK_COLOR_COMPONENT_A_BIT;
//
//    VkPipelineDepthStencilStateCreateInfo depth_info = {};
//    depth_info.sType = VK_STRUCTURE_TYPE_PIPELINE_DEPTH_STENCIL_STATE_CREATE_INFO;
//
//    VkPipelineColorBlendStateCreateInfo blend_info = {};
//    blend_info.sType = VK_STRUCTURE_TYPE_PIPELINE_COLOR_BLEND_STATE_CREATE_INFO;
//    blend_info.attachmentCount = 1;
//    blend_info.pAttachments = color_attachment;
//
//    VkDynamicState dynamic_states[2] = { VK_DYNAMIC_STATE_VIEWPORT, VK_DYNAMIC_STATE_SCISSOR };
//    VkPipelineDynamicStateCreateInfo dynamic_state = {};
//    dynamic_state.sType = VK_STRUCTURE_TYPE_PIPELINE_DYNAMIC_STATE_CREATE_INFO;
//    dynamic_state.dynamicStateCount = (uint32_t)IM_ARRAYSIZE(dynamic_states);
//    dynamic_state.pDynamicStates = dynamic_states;
//
//    ImGui_ImplVulkan_CreatePipelineLayout(device, allocator);
//
//    VkGraphicsPipelineCreateInfo info = {};
//    info.sType = VK_STRUCTURE_TYPE_GRAPHICS_PIPELINE_CREATE_INFO;
//    info.flags = bd->PipelineCreateFlags;
//    info.stageCount = 2;
//    info.pStages = stage;
//    info.pVertexInputState = &vertex_info;
//    info.pInputAssemblyState = &ia_info;
//    info.pViewportState = &viewport_info;
//    info.pRasterizationState = &raster_info;
//    info.pMultisampleState = &ms_info;
//    info.pDepthStencilState = &depth_info;
//    info.pColorBlendState = &blend_info;
//    info.pDynamicState = &dynamic_state;
//    info.layout = bd->PipelineLayout;
//    info.renderPass = renderPass;
//    info.subpass = subpass;
//    VkResult err = vkCreateGraphicsPipelines(device, pipelineCache, 1, &info, allocator, pipeline);
//    check_vk_result(err);
}

static void _gb_create_shader_modules(VkDevice device, const VkAllocationCallbacks* allocator) {

//    // Create the shader modules
//    ImGui_ImplVulkan_Data* bd = ImGui_ImplVulkan_GetBackendData();
//    if (bd->ShaderModuleVert == VK_NULL_HANDLE)
//    {
//        VkShaderModuleCreateInfo vert_info = {};
//        vert_info.sType = VK_STRUCTURE_TYPE_SHADER_MODULE_CREATE_INFO;
//        vert_info.codeSize = sizeof(__glsl_shader_vert_spv);
//        vert_info.pCode = (uint32_t*)__glsl_shader_vert_spv;
//        VkResult err = vkCreateShaderModule(device, &vert_info, allocator, &bd->ShaderModuleVert);
//        check_vk_result(err);
//    }
//    if (bd->ShaderModuleFrag == VK_NULL_HANDLE)
//    {
//        VkShaderModuleCreateInfo frag_info = {};
//        frag_info.sType = VK_STRUCTURE_TYPE_SHADER_MODULE_CREATE_INFO;
//        frag_info.codeSize = sizeof(__glsl_shader_frag_spv);
//        frag_info.pCode = (uint32_t*)__glsl_shader_frag_spv;
//        VkResult err = vkCreateShaderModule(device, &frag_info, allocator, &bd->ShaderModuleFrag);
//        check_vk_result(err);
//    }
}

static void _gb_destroy_frame(VkDevice device, struct vulkan_frame* fd, const VkAllocationCallbacks* allocator) {

    vkDestroyFence(device, fd->Fence, allocator);
    vkFreeCommandBuffers(device, fd->CommandPool, 1, &fd->CommandBuffer);
    vkDestroyCommandPool(device, fd->CommandPool, allocator);
    fd->Fence = VK_NULL_HANDLE;
    fd->CommandBuffer = VK_NULL_HANDLE;
    fd->CommandPool = VK_NULL_HANDLE;

    vkDestroyImageView(device, fd->BackbufferView, allocator);
    vkDestroyFramebuffer(device, fd->Framebuffer, allocator);
}

static void _gb_destroy_frame_semaphores(VkDevice device, struct vulkan_frame_semaphores* fsd, const VkAllocationCallbacks* allocator) {

    vkDestroySemaphore(device, fsd->ImageAcquiredSemaphore, allocator);
    vkDestroySemaphore(device, fsd->RenderCompleteSemaphore, allocator);
    fsd->ImageAcquiredSemaphore = fsd->RenderCompleteSemaphore = VK_NULL_HANDLE;
}

static void _gb_destroy_all_viewports_render_buffers(VkDevice device, const VkAllocationCallbacks* allocator) {

//  TODO
//    ImGuiPlatformIO& platform_io = ImGui::GetPlatformIO();
//    for (int n = 0; n < platform_io.Viewports.Size; n++)
//        if (ImGui_ImplVulkan_ViewportData* vd = (ImGui_ImplVulkan_ViewportData*)platform_io.Viewports[n]->RendererUserData)
//            ImGui_ImplVulkanH_DestroyWindowRenderBuffers(device, &vd->RenderBuffers, allocator);
}

static void* _gb_alloc(size_t count) {

    void *p = malloc(count);
    if (p == NULL) {
        fprintf(stderr, "NO MEMORY\n");
        abort();
    }
    return p;
}

static void _gb_glfw_error_callback(int error, const char* description) {

    fprintf(stderr, "GLFW Error %d: %s\n", error, description);
}

static void _gb_check_vk_result(VkResult err) {

    if (err == 0) {
        return;
    }
    fprintf(stderr, "[vulkan] Error: VkResult = %d\n", err);
    if (err < 0) {
        abort();
    }
}

#ifdef GB_VULKAN_DEBUG_REPORT
static VKAPI_ATTR VkBool32 VKAPI_CALL _gb_debug_report(VkDebugReportFlagsEXT flags, VkDebugReportObjectTypeEXT objectType,
    uint64_t object, size_t location, int32_t messageCode, const char* pLayerPrefix, const char* pMessage, void* pUserData) {

    (void)flags; (void)object; (void)location; (void)messageCode; (void)pUserData; (void)pLayerPrefix; // Unused arguments
    fprintf(stderr, "[vulkan] Debug report from ObjectType: %i\nMessage: %s\n\n", objectType, pMessage);
    return VK_FALSE;
}
#endif // IMGUI_VULKAN_DEBUG_REPORT

