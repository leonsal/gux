# Force link with libs and install
all:
	clear
	cd libs;make
	rm -f /home/leonel/go/bin/gux
	cd ..
	/home/leonel/bin/go/bin/go install .

# Force link with libs, install and run
run: all
	gux


# Force link with libs, install and run showing memory allocations
runalloc: all
	GODEBUG=allocfreetrace=1 gux

