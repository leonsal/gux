#
# Makefile to build required libraries
# In windows it requires mingw32-make.exe
#
GCC=-DCMAKE_C_COMPILER=gcc -DCMAKE_CXX_COMPILER=g++
DEBUG=-DCMAKE_BUILD_TYPE=Debug
RELEASE=-DCMAKE_BUILD_TYPE=Release
GEN=-G "Unix Makefiles"
JOBS=-j8
ifeq ($(OS),Windows_NT)
	GEN=-G "MinGW Makefiles"
endif

#
# default target: gcc release
#
all: build
	cmake --build build -- $(JOBS)

build:
	cmake $(GEN) -Bbuild $(GCC) $(RELEASE)


#
# gcc debug
#
debug: build_debug
	cmake --build build_debug -- $(JOBS)

build_debug: 
	cmake $(GEN) -Bbuild_debug $(GCC) $(DEBUG)


#
# clean all targets
#
clean:
	cmake -E remove_directory build build_debug 
	cmake -E remove *.a
	cmake -E remove -f compile_commands.json

