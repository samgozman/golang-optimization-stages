package stage1

import (
	"context"
	"log"
	"net/http"
	"testing"

	"github.com/samgozman/golang-optimization-stages/utils"
)

// Usage:
//
//	go test -v -cpuprofile=cpu.pprof -memprofile=mem.pprof -benchmem -bench=. -benchtime=1s -count=10
//
// Results:
//
//	BenchmarkServeApp-10                  75          15359607 ns/op         5022236 B/op      63424 allocs/op
//	BenchmarkServeApp-10                  79          15283220 ns/op         5019885 B/op      63420 allocs/op
//	BenchmarkServeApp-10                  80          15376617 ns/op         5019731 B/op      63421 allocs/op
//	BenchmarkServeApp-10                  81          15324918 ns/op         5018524 B/op      63418 allocs/op
//	BenchmarkServeApp-10                  81          15241640 ns/op         5017535 B/op      63417 allocs/op
//	BenchmarkServeApp-10                  80          15320955 ns/op         5017719 B/op      63417 allocs/op
//	BenchmarkServeApp-10                  81          15473488 ns/op         5016617 B/op      63417 allocs/op
//	BenchmarkServeApp-10                  80          15434468 ns/op         5016502 B/op      63416 allocs/op
//	BenchmarkServeApp-10                  74          15740035 ns/op         5015688 B/op      63416 allocs/op
//	BenchmarkServeApp-10                  76          15920205 ns/op         5015834 B/op      63414 allocs/op
func BenchmarkServeApp(b *testing.B) {
	// Start pprof profiling
	if err := utils.StartPprof(); err != nil {
		log.Fatal(err)
	}

	// Create a context
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Run the server
	go ServeApp(ctx)

	// Run the benchmark
	for i := 0; i < b.N; i++ {
		_, err := http.Get("http://localhost:8080/json")
		if err != nil {
			b.Fatal(err)
		}
	}

	// Cancel the context
	cancel()
}
