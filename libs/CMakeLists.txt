cmake_minimum_required(VERSION 3.0 FATAL_ERROR)

project(libgux)

set(CMAKE_ARCHIVE_OUTPUT_DIRECTORY ${CMAKE_SOURCE_DIR})
set(CMAKE_EXPORT_COMPILE_COMMANDS ON)

include_directories(./)
include_directories(glfw-3.3.8/include)
include_directories(gl3w/include)

# GLFW dependency subproject
set(GLFW_BUILD_DOCS OFF CACHE BOOL "" FORCE)
set(GLFW_BUILD_TESTS OFF CACHE BOOL "" FORCE)
set(GLFW_BUILD_EXAMPLES OFF CACHE BOOL "" FORCE)
add_subdirectory(glfw-3.3.8 glfw)

set(OPENGL_SOURCES
    gl3w/src/gl3w.c
    src/glfw_opengl.c
)

set(VULKAN_SOURCES
    gl3w/src/gl3w.c
    src/volk.c
    src/glfw_vulkan.c
)

# For windows use local copy of vulkan headers from SDK
if (WIN32) 
    include_directories(vulkan)
endif()

add_library(guxgl STATIC ${OPENGL_SOURCES})
add_library(guxvk STATIC ${VULKAN_SOURCES})


