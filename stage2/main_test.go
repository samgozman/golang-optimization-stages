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
//	BenchmarkServeApp-2           74          16027115 ns/op         5014060 B/op      63421 allocs/op
//	BenchmarkServeApp-2           66          15508486 ns/op         5013355 B/op      63418 allocs/op
//	BenchmarkServeApp-2           73          15658756 ns/op         5013058 B/op      63418 allocs/op
//	BenchmarkServeApp-2           70          15703708 ns/op         5012548 B/op      63415 allocs/op
//	BenchmarkServeApp-2           78          15987481 ns/op         5012574 B/op      63414 allocs/op
//	BenchmarkServeApp-2           80          15780026 ns/op         5012611 B/op      63413 allocs/op
//	BenchmarkServeApp-2           79          16164028 ns/op         5012219 B/op      63413 allocs/op
//	BenchmarkServeApp-2           68          15807772 ns/op         5012287 B/op      63410 allocs/op
//	BenchmarkServeApp-2           79          16032428 ns/op         5011579 B/op      63411 allocs/op
//	BenchmarkServeApp-2           72          15638314 ns/op         5011674 B/op      63409 allocs/op
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
