package main

import (
	"fmt"
	"sync"
	"time"
)

// MutexMap uses sync.Mutex (no distinction between read/write)
type MutexMap struct {
	mu sync.Mutex
	m  map[int]int
}

func (sm *MutexMap) Write(key, value int) {
	sm.mu.Lock()
	sm.m[key] = value
	sm.mu.Unlock()
}

func (sm *MutexMap) Read(key int) int {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	return sm.m[key]
}

// RWMutexMap uses sync.RWMutex (allows concurrent reads)
type RWMutexMap struct {
	mu sync.RWMutex
	m  map[int]int
}

func (sm *RWMutexMap) Write(key, value int) {
	sm.mu.Lock()
	sm.m[key] = value
	sm.mu.Unlock()
}

func (sm *RWMutexMap) Read(key int) int {
	sm.mu.RLock()
	defer sm.mu.RUnlock()
	return sm.m[key]
}

func main() {
	fmt.Println("=== RLock Experiment: Read-Heavy Workload ===")
	fmt.Println("90% reads / 10% writes")
	fmt.Println("50 goroutines × 1000 ops = 50,000 total (5,000 writes, 45,000 reads)")
	fmt.Println()

	runs := 3

	// Test Mutex with read-heavy workload
	fmt.Println("--- Mutex (Lock for both read and write) ---")
	var mutexTotal time.Duration
	for r := 1; r <= runs; r++ {
		sm := &MutexMap{m: make(map[int]int)}
		// Pre-populate with data for reads
		for i := range 1000 {
			sm.m[i] = i
		}

		var wg sync.WaitGroup
		var writeCount, readCount int
		var countMu sync.Mutex

		start := time.Now()

		for g := range 50 {
			wg.Go(func() {
				localWrites, localReads := 0, 0
				for i := range 1000 {
					if i%10 == 0 {
						sm.Write(g*1000+i, i)
						localWrites++
					} else {
						_ = sm.Read(i % 1000)
						localReads++
					}
				}
				countMu.Lock()
				writeCount += localWrites
				readCount += localReads
				countMu.Unlock()
			})
		}

		wg.Wait()
		elapsed := time.Since(start)
		mutexTotal += elapsed
		fmt.Printf("  Run %d: len=%d, writes=%d, reads=%d, time=%v\n",
			r, len(sm.m), writeCount, readCount, elapsed)
	}
	fmt.Printf("  Average: %v\n", mutexTotal/time.Duration(runs))

	fmt.Println()

	// Test RWMutex with read-heavy workload
	fmt.Println("--- RWMutex (RLock for reads, Lock for writes) ---")
	var rwmutexTotal time.Duration
	for r := 1; r <= runs; r++ {
		sm := &RWMutexMap{m: make(map[int]int)}
		// Pre-populate with data for reads
		for i := range 1000 {
			sm.m[i] = i
		}

		var wg sync.WaitGroup
		var writeCount, readCount int
		var countMu sync.Mutex

		start := time.Now()

		for g := range 50 {
			wg.Go(func() {
				localWrites, localReads := 0, 0
				for i := range 1000 {
					if i%10 == 0 {
						sm.Write(g*1000+i, i)
						localWrites++
					} else {
						_ = sm.Read(i % 1000)
						localReads++
					}
				}
				countMu.Lock()
				writeCount += localWrites
				readCount += localReads
				countMu.Unlock()
			})
		}

		wg.Wait()
		elapsed := time.Since(start)
		rwmutexTotal += elapsed
		fmt.Printf("  Run %d: len=%d, writes=%d, reads=%d, time=%v\n",
			r, len(sm.m), writeCount, readCount, elapsed)
	}
	fmt.Printf("  Average: %v\n", rwmutexTotal/time.Duration(runs))

}