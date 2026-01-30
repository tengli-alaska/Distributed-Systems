package main

import (
	"fmt"
	"sync"
	"sync/atomic"
)

func main() {
	fmt.Println("=== Atomicity Experiment ===")
	fmt.Println("50 goroutines × 1000 increments = expected 50,000")
	fmt.Println()

	runs := 5

	fmt.Println("--- Atomic Counter (thread-safe) ---")
	for r := 1; r <= runs; r++ {
		result := testAtomicCounter()
		status := "✓"
		if result != 50000 {
			status = "✗ RACE CONDITION!"
		}
		fmt.Printf("  Run %d: %d %s\n", r, result, status)
	}

	fmt.Println()
	fmt.Println("--- Regular Counter (NOT thread-safe) ---")
	for r := 1; r <= runs; r++ {
		result := testRegularCounter()
		status := "✓"
		if result != 50000 {
			status = "✗ RACE CONDITION!"
		}
		fmt.Printf("  Run %d: %d %s\n", r, result, status)
	}
}

// Thread-safe with atomic operations
func testAtomicCounter() uint64 {
	var ops atomic.Uint64
	var wg sync.WaitGroup

	for range 50 {
		wg.Go(func() {
			for range 1000 {
				ops.Add(1)
			}
		})
	}

	wg.Wait()
	return ops.Load()
}

// NOT thread-safe - will have race conditions
func testRegularCounter() uint64 {
	var ops uint64
	var wg sync.WaitGroup

	for range 50 {
		wg.Go(func() {
			for range 1000 {
				ops++ // Race condition here!
			}
		})
	}

	wg.Wait()
	return ops
}