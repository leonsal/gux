package gb

/*
// Linux Build Tags
#cgo linux,!vulkan CFLAGS: -I gl3w/include
#cgo linux,!wayland CFLAGS: -D_GLFW_X11
#cgo linux,wayland CFLAGS: -D_GLFW_WAYLAND
#cgo linux,gles3 LDFLAGS: -lGLESv3
#cgo linux,!wayland LDFLAGS: -lX11 -lXrandr -lXxf86vm -lXi -lXcursor -lm -lXinerama -ldl -lrt
#cgo linux,wayland LDFLAGS: -lwayland-client -lwayland-cursor -lwayland-egl -lxkbcommon -lm -ldl -lrt

// Windows Build Tags
#cgo windows CFLAGS: -D_GLFW_WIN32 -Iglfw/deps/mingw
#cgo windows LDFLAGS: -lgdi32
#cgo !gles2,windows LDFLAGS: -lopengl32
#cgo gles2,windows LDFLAGS: -lGLESv2

// BSD Build Tags
#cgo freebsd,!wayland netbsd,!wayland openbsd pkg-config: x11 xau xcb xdmcp
#cgo freebsd,wayland netbsd,wayland pkg-config: wayland-client wayland-cursor wayland-egl epoll-shim
#cgo freebsd netbsd openbsd CFLAGS: -D_GLFW_HAS_DLOPEN
#cgo freebsd,!wayland netbsd,!wayland openbsd CFLAGS: -D_GLFW_X11 -D_GLFW_HAS_GLXGETPROCADDRESSARB
#cgo freebsd,wayland netbsd,wayland CFLAGS: -D_GLFW_WAYLAND
#cgo freebsd netbsd openbsd LDFLAGS: -lm
*/
import "C"
