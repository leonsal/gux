package main

import (
	"flag"
	"fmt"
	"os"
	"runtime"
	"sort"

	"github.com/leonsal/gux"
	"github.com/leonsal/gux/gb"
)

var colorList []gb.RGBA

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

	// Parse command line
	flag.Parse()
	args := flag.Args()
	tinfo := testInfo{}
	if len(args) > 0 {
		tname := args[0]
		var ok bool
		tinfo, ok = mapTests[tname]
		if !ok {
			fmt.Printf("Invalid test name: %s\n", tname)
			os.Exit(1)
		}
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

	// Run specified test or run all tests
	if len(tinfo.name) > 0 {
		runTest(win, tinfo, 0)
	} else {
		tests := []testInfo{}
		for _, v := range mapTests {
			tests = append(tests, v)
		}
		sort.Slice(tests, func(i, j int) bool {
			return tests[i].order < tests[j].order
		})
		index := 0
		for {
			abort := runTest(win, tests[index], 200)
			if abort {
				break
			}
			index++
			if index >= len(tests) {
				index = 0
			}
		}
	}

	win.Destroy()
}

func runTest(win *gux.Window, tinfo testInfo, maxFrames int) bool {

	fmt.Printf("Running test: %s \n", tinfo.name)
	var cgoCallsStart int64
	var statsStart runtime.MemStats
	frameCount := 0

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

	// Calculates and shows allocations and cgo calls per frame
	cgoCalls := runtime.NumCgoCall() - cgoCallsStart
	cgoPerFrame := float64(cgoCalls) / float64(frameCount)
	var stats runtime.MemStats
	runtime.ReadMemStats(&stats)
	allocsPerFrame := float64(stats.Alloc-statsStart.Alloc) / float64(frameCount)
	fmt.Println("Frames:", frameCount, "Allocs per frame:", allocsPerFrame, "CGO calls per frame:", cgoPerFrame)
	test.destroy(win)
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
