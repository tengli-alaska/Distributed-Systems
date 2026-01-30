package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	fmt.Println("=== sync.Map Experiment ===")
	fmt.Println("50 goroutines × 1000 writes = 50,000 unique keys expected")
	fmt.Println()

	runs := 3
	var totalDuration time.Duration

	for r := 1; r <= runs; r++ {
		var m sync.Map
		var wg sync.WaitGroup

		start := time.Now()

		for g := range 50 {
			wg.Go(func() {
				for i := range 1000 {
					m.Store(g*1000+i, i)
				}
			})
		}

		wg.Wait()
		elapsed := time.Since(start)
		totalDuration += elapsed

		// Count entries using Range
		count := 0
		m.Range(func(key, value any) bool {
			count++
			return true // continue iteration
		})

		fmt.Printf("  Run %d: len=%d, time=%v\n", r, count, elapsed)
	}

	avg := totalDuration / time.Duration(runs)
	fmt.Printf("\nAverage time: %v\n", avg)
}