package main

import (
	"fmt"
	"sync"
	"time"
)

// SafeMap wraps a map with a mutex for thread-safe access
type SafeMap struct {
	mu sync.Mutex
	m  map[int]int
}

func main() {
	fmt.Println("=== Mutex Experiment ===")
	fmt.Println("50 goroutines × 1000 writes = 50,000 unique keys expected")
	fmt.Println()

	runs := 3
	var totalDuration time.Duration

	for r := 1; r <= runs; r++ {
		sm := &SafeMap{m: make(map[int]int)}
		var wg sync.WaitGroup

		start := time.Now()

		for g := range 50 {
			wg.Go(func() {
				for i := range 1000 {
					sm.mu.Lock()
					sm.m[g*1000+i] = i
					sm.mu.Unlock()
				}
			})
		}

		wg.Wait()
		elapsed := time.Since(start)
		totalDuration += elapsed

		fmt.Printf("  Run %d: len=%d, time=%v\n", r, len(sm.m), elapsed)
	}

	avg := totalDuration / time.Duration(runs)
	fmt.Printf("\nAverage time: %v\n", avg)

}