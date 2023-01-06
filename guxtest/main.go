package main

import (
	"flag"
	"fmt"
	"os"
	"runtime"

	"github.com/leonsal/gux"
	"github.com/leonsal/gux/gb"
)

var colorList []gb.RGBA

// Command line flags
var (
	oTrace = flag.String("trace", "", "Activate go tool execution tracer writing data to the specified file")
)

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

var (
	mapTests = map[string]testInfo{}
)

type ITest interface {
	draw(*gux.Window)
	destroy(*gux.Window)
}

func main() {

	runtime.LockOSThread()
	flag.Parse()
	args := flag.Args()
	if len(args) == 0 {
		fmt.Println("Test name not supplied")
		os.Exit(1)
	}

	tname := args[0]
	tinfo, ok := mapTests[tname]
	if !ok {
		fmt.Printf("Invalid test name: %s\n", tname)
		os.Exit(1)
	}

	// Create window
	cfg := gb.Config{}
	cfg.DebugPrintCmds = false
	cfg.OpenGL.ES = false
	cfg.Vulkan.ValidationLayer = true
	win, err := gux.NewWindow("title", 2000, 1200, &cfg)
	if err != nil {
		panic(err)
	}

	runTest(win, tinfo, 0)

	//	// Render loop
	//	var cgoCallsStart int64
	//	var statsStart runtime.MemStats
	//	frameCount := 0
	//
	//	// Creates test
	//	test := tinfo.create(win)
	//
	//	for win.StartFrame() {
	//		test.draw(win)
	//		win.Render()
	//		// All the allocations should be done in the first frame
	//		frameCount++
	//		if frameCount == 1 {
	//			cgoCallsStart = runtime.NumCgoCall()
	//			runtime.ReadMemStats(&statsStart)
	//		}
	//	}
	//
	//	// Calculates and shows allocations and cgo calls per frame
	//	cgoCalls := runtime.NumCgoCall() - cgoCallsStart
	//	cgoPerFrame := float64(cgoCalls) / float64(frameCount)
	//	var stats runtime.MemStats
	//	runtime.ReadMemStats(&stats)
	//	allocsPerFrame := float64(stats.Alloc-statsStart.Alloc) / float64(frameCount)
	//	fmt.Println("Frames:", frameCount, "Allocs per frame:", allocsPerFrame, "CGO calls per frame:", cgoPerFrame)
	//
	//	test.destroy(win)
	win.Destroy()
}

func runTest(win *gux.Window, tinfo testInfo, maxFrames int) {

	// Render loop
	var cgoCallsStart int64
	var statsStart runtime.MemStats
	frameCount := 0

	// Creates test
	test := tinfo.create(win)

	for win.StartFrame() {
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

	// Calculates and shows allocations and cgo calls per frame
	cgoCalls := runtime.NumCgoCall() - cgoCallsStart
	cgoPerFrame := float64(cgoCalls) / float64(frameCount)
	var stats runtime.MemStats
	runtime.ReadMemStats(&stats)
	allocsPerFrame := float64(stats.Alloc-statsStart.Alloc) / float64(frameCount)
	fmt.Println("Frames:", frameCount, "Allocs per frame:", allocsPerFrame, "CGO calls per frame:", cgoPerFrame)
	test.destroy(win)
}

// registerTest is used by tests to register themselves
func registerTest(name string, order int, create func(*gux.Window) ITest) {

	mapTests[name] = testInfo{name: name, order: order, create: create}
}

func nextColor(i int) gb.RGBA {

	ci := i % len(colorList)
	return colorList[ci]
}
