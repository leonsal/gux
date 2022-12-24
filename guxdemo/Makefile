# Build/run OpenGL version
gl: libs
	/home/leonel/bin/go/bin/go install .
	guxdemo

# Build/run Vulkan version
vk: libs
	/home/leonel/bin/go/bin/go install -tags vulkan .
	guxdemo

# Build Vulkan executable locally for debugging
vkbuild: libs
	/home/leonel/bin/go/bin/go build -tags vulkan .

# Run showing allocations
runalloc: libs
	GODEBUG=allocfreetrace=1 guxdemo

# Build C libraries and remove executable to force new link
libs:
	cd ../libs;make
	rm -f /home/leonel/go/bin/guxdemo

