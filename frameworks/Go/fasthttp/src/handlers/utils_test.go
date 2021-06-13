package handlers

import (
	"encoding/json"
	rand2 "github.com/maxim-kuderko/fast-random"
	"io"
	"math/rand"
	"runtime"
	"sync"
	"testing"
	"time"
)

func BenchmarkJson(b *testing.B) {
	concurrency := runtime.GOMAXPROCS(-1)
	wg := sync.WaitGroup{}
	wg.Add(concurrency)
	w := io.Discard
	world := acquireWorld()
	b.ReportAllocs()
	b.ResetTimer()
	for j := 0; j < concurrency; j++ {
		go func() {
			defer wg.Done()
			for i := 0; i < b.N; i++ {
				b, _ := json.Marshal(world)
				w.Write(b)
			}
		}()
	}
	wg.Wait()

}

func BenchmarkJsonReuse(b *testing.B) {
	concurrency := runtime.GOMAXPROCS(-1)
	wg := sync.WaitGroup{}
	wg.Add(concurrency)
	world := acquireWorld()
	b.ReportAllocs()
	b.ResetTimer()

	for j := 0; j < concurrency; j++ {
		go func() {
			defer wg.Done()
			for i := 0; i < b.N; i++ {
				enc := acquireJsonEncoder(io.Discard)
				enc.Encode(world)
				releaseJsonEncoder(enc)
			}
		}()
	}
	wg.Wait()

}

func BenchmarkMathRand(b *testing.B) {
	concurrency := runtime.GOMAXPROCS(-1)
	wg := sync.WaitGroup{}
	wg.Add(concurrency)
	b.ReportAllocs()
	b.ResetTimer()
	for j := 0; j < concurrency; j++ {
		go func() {
			defer wg.Done()
			for i := 0; i < b.N; i++ {
				rand.Int31n(999)
			}
		}()
	}
	wg.Wait()
}

func BenchmarkConcurrentRand(b *testing.B) {
	concurrency := runtime.GOMAXPROCS(-1)
	wg := sync.WaitGroup{}
	wg.Add(concurrency)
	b.ReportAllocs()
	b.ResetTimer()
	s := rand2.NewSource(concurrency*2, func() int64 {
		return time.Now().UnixNano()
	})
	ra := rand.New(s)
	for j := 0; j < concurrency; j++ {
		go func() {
			defer wg.Done()
			for i := 0; i < b.N; i++ {
				ra.Int31n(999)
			}
		}()
	}
	wg.Wait()
}
