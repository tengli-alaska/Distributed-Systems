package main

import (
	"fmt"
	"sync"
)

func main() {
	fmt.Println("=== Collections Experiment ===")
	fmt.Println("50 goroutines × 1000 writes = 50,000 unique keys expected")
	fmt.Println("Attempting concurrent writes to a plain map")

	m := make(map[int]int)
	var wg sync.WaitGroup

	for g := range 50 {
		wg.Go(func() {
			for i := range 1000 {
				// Each goroutine writes to unique keys: g*1000+i
				// g=0: keys 0-999
				// g=1: keys 1000-1999
				m[g*1000+i] = i
			}
		})
	}

	wg.Wait()
	fmt.Printf("Map length: %d\n", len(m))
}