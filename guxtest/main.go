package main

import (
	"flag"
	"log"
	"os"
	"runtime"
	"runtime/pprof"
	"runtime/trace"
	"sort"

	"github.com/leonsal/gux"
	"github.com/leonsal/gux/gb"
)

// Initializes colors lists
func init() {
	colorList = append(colorList,
		gb.MakeColor(255, 0, 0, 255),
		gb.MakeColor(0, 255, 0, 255),
		gb.MakeColor(0, 0, 255, 255),
		gb.MakeColor(0, 0, 0, 255),
		gb.MakeColor(255, 255, 0, 255),
		gb.MakeColor(0, 255, 255, 255),
		gb.MakeColor(255, 255, 255, 255),
		gb.MakeColor(100, 100, 100, 255),
	)
}

// testInfo describes a test
type testInfo struct {
	name   string                  // Test name
	order  int                     // Show order
	create func(*gux.Window) ITest // Test constructor
}

// ITest is the interface for all tests objects
type ITest interface {
	draw(*gux.Window)
	destroy(*gux.Window)
}

var (
	colorList   = []gb.RGBA{}
	oCpuProfile = flag.String("cpuprofile", "", "Write CPU profile to the specified file")
	oMemProfile = flag.String("memprofile", "", "Write memory profile to the specified file")
	oTrace      = flag.String("trace", "", "Write execution trace to the specified file")
	oFrameCount = flag.Uint("frames", 500, "Number of frames to execute the test")
	mapTests    = map[string]testInfo{}
	traceFile   *os.File
	cpuprofFile *os.File
	memprofFile *os.File
)

func main() {

	runtime.LockOSThread()
	log.SetFlags(log.Lmicroseconds)

	// Parse command line
	flag.Parse()
	args := flag.Args()
	tinfo := testInfo{}
	if len(args) > 0 {
		tname := args[0]
		var ok bool
		tinfo, ok = mapTests[tname]
		if !ok {
			log.Fatalf("Invalid test name: %s\n", tname)
		}
	}

	// Create window
	cfg := gb.Config{}
	cfg.DebugPrintCmds = false
	cfg.OpenGL.ES = false
	cfg.Vulkan.ValidationLayer = true
	win, err := gux.NewWindow("title", 2000, 1200, &cfg)
	if err != nil {
		log.Fatal(err)
	}

	// Starts optional trace/profiling
	traceProfStart()

	// Run specified test or run all tests
	if len(tinfo.name) > 0 {
		runTest(win, tinfo, *oFrameCount)
	} else {
		// Build slice of testInfo structs from map
		tests := []testInfo{}
		for _, v := range mapTests {
			tests = append(tests, v)
		}
		// Sort slice by increasing order field
		sort.Slice(tests, func(i, j int) bool {
			return tests[i].order < tests[j].order
		})
		// Run tests continously unless aborted by closing the window
		index := 0
		for {
			abort := runTest(win, tests[index], *oFrameCount)
			if abort {
				break
			}
			index++
			index %= len(tests)
		}
	}

	// Stops optional trace/profiling
	traceProfStop()

	win.Destroy()
}

// runTest creates the specified test and runs it for the specified number of frames.
// If 'maxFrames' is zero, runs continously till the window is closed.
func runTest(win *gux.Window, tinfo testInfo, maxFrames uint) bool {

	log.Printf("Running test: %s (%d frames) \n", tinfo.name, maxFrames)
	var cgoCallsStart int64
	var statsStart runtime.MemStats
	frameCount := uint(0)

	// Creates test
	test := tinfo.create(win)

	// Render Loop
	abort := false
	for {
		if !win.StartFrame() {
			abort = true
			break
		}
		test.draw(win)
		win.Render()
		// All the allocations should be done in the first frame
		frameCount++
		if frameCount == 1 {
			cgoCallsStart = runtime.NumCgoCall()
			runtime.ReadMemStats(&statsStart)
		}
		if maxFrames > 0 && frameCount > maxFrames {
			break
		}
	}
	test.destroy(win)

	// Calculates and shows allocations and cgo calls per frame
	cgoCalls := runtime.NumCgoCall() - cgoCallsStart
	cgoPerFrame := float64(cgoCalls) / float64(frameCount)
	var stats runtime.MemStats
	runtime.ReadMemStats(&stats)
	allocsPerFrame := float64(stats.Alloc-statsStart.Alloc) / float64(frameCount)
	log.Printf("Frames:%d  Allocs/frame:%f  CGO calls/frame:%f\n\n", frameCount, allocsPerFrame, cgoPerFrame)
	return abort
}

// registerTest is used by tests to register themselves
func registerTest(name string, order int, create func(*gux.Window) ITest) {

	mapTests[name] = testInfo{name: name, order: order, create: create}
}

func nextColor(i int) gb.RGBA {

	ci := i % len(colorList)
	return colorList[ci]
}

// traceProfStart starts optional execution tracing and profiles
func traceProfStart() {

	if len(*oTrace) > 0 {
		var err error
		traceFile, err = os.Create(*oTrace)
		if err != nil {
			log.Fatalf("Error creating execution trace file:%s\n", err)
		}
		err = trace.Start(traceFile)
		if err != nil {
			log.Fatalf("Error starting execution trace:%s\n", err)
		}
	}

	if len(*oCpuProfile) > 0 {
		var err error
		cpuprofFile, err = os.Create(*oCpuProfile)
		if err != nil {
			log.Fatalf("Error creating cpu profile file:%s\n", err)
		}
		err = pprof.StartCPUProfile(cpuprofFile)
		if err != nil {
			log.Fatalf("Error starting cpu profile:%s\n", err)
		}
	}

	if len(*oMemProfile) > 0 {
		var err error
		memprofFile, err = os.Create(*oMemProfile)
		if err != nil {
			log.Fatalf("Error creating memory profile file:%s\n", err)
		}
		err = pprof.WriteHeapProfile(memprofFile)
		runtime.GC()
		if err != nil {
			log.Fatalf("Error writing memory profile:%s\n", err)
		}
	}

}

// traceProfStop stops optional active execution tracing and profiles
func traceProfStop() {

	if len(*oMemProfile) > 0 {
		memprofFile.Close()
		log.Printf("Memory profile saved. To show the profile execute command:\n>go tool pprof -web %s\n", *oMemProfile)
	}

	if len(*oCpuProfile) > 0 {
		pprof.StopCPUProfile()
		cpuprofFile.Close()
		log.Printf("CPU profile saved. To show the profile execute command:\n>go tool pprof -web %s\n", *oCpuProfile)
	}

	if len(*oTrace) > 0 {
		trace.Stop()
		traceFile.Close()
		log.Printf("Execution trace saved. To show the trace execute command:\n>go tool trace %s\n", *oTrace)
	}
}
