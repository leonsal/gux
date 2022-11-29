# Force link with libs and install
all:
	clear
	cd ../libs;make
	rm -f /home/leonel/go/bin/guxdemo
	cd ../guxdemo
	/home/leonel/bin/go/bin/go install .

# Force link with libs, install and run
run: all
	guxdemo


# Force link with libs, install and run showing memory allocations
runalloc: all
	GODEBUG=allocfreetrace=1 guxdemo

