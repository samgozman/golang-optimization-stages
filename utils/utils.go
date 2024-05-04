package utils

import (
	"flag"
	"fmt"
	"os"
	"runtime"
	"runtime/pprof"
)

var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to `file`")
var memprofile = flag.String("memprofile", "", "write memory profile to `file`")

// StartPprof starts pprof profiling
func StartPprof() error {
	flag.Parse()
	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			return fmt.Errorf("could not create CPU profile: %w", err)
		}
		if err := pprof.StartCPUProfile(f); err != nil {
			_ = f.Close()
			return fmt.Errorf("could not start CPU profile: %w", err)
		}
		pprof.StopCPUProfile()
		if err := f.Close(); err != nil {
			return fmt.Errorf("could not close CPU profile file: %w", err)
		}
	}

	if *memprofile != "" {
		f, err := os.Create(*memprofile)
		if err != nil {
			return fmt.Errorf("could not create memory profile: %w", err)
		}
		runtime.GC() // get up-to-date statistics
		if err := pprof.WriteHeapProfile(f); err != nil {
			_ = f.Close()
			return fmt.Errorf("could not write memory profile: %w", err)
		}
		if err := f.Close(); err != nil {
			return fmt.Errorf("could not close memory profile file: %w", err)
		}
	}

	return nil
}
