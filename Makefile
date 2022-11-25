# Default target
all:
	clear
	cd libs;make
	rm -f /home/leonel/go/bin/gux
	cd ..
	/home/leonel/bin/go/bin/go install .

# Build and run showing memory allocations
runalloc: all
	gux

