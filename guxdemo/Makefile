# Build/run OpenGL version
gl: libs
	/home/leonel/bin/go/bin/go install .
	guxdemo

# Build/run Vulkan version
vk: libs
	/home/leonel/bin/go/bin/go install -tags vulkan
	guxdemo

# Build Vulkan executable locally for debugging
vkbuild: libs
	rm -f guxdemo
	/home/leonel/bin/go/bin/go build -a -tags vulkan

# Run showing allocations
runalloc: libs
	/home/leonel/bin/go/bin/go install .
	GODEBUG=allocfreetrace=1 guxdemo

# Build C libraries and remove executable to force new link
libs:
	cd ../libs;make
	rm -f /home/leonel/go/bin/guxdemo

