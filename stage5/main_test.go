package stage5

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
