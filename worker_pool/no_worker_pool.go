package main

import (
	"fmt"
	"runtime"
	"time"
)

func heavyCompute(n int) int {
	time.Sleep(100 * time.Millisecond) // simulate heavy work
	return n * n
}

func main() {
	results := make(chan int, numJobs)

	numJobs := 100_000 // large number to see memory spike
	numbers := make([]int, numJobs)
	for i := 0; i < numJobs; i++ {
		numbers[i] = i + 1
	}

	fmt.Println("Starting jobs without worker pool...")
	start := time.Now()

	fmt.Printf("Memory before jobs: %v MB\n", bToMb(runtimeReadMemStatsHeap()))

	for _, n := range numbers {
		go func(num int) {
			_ = heavyCompute(num)
		}(n)
	}

	// Wait long enough (not ideal, just to observe memory)
	time.Sleep(30 * time.Second)

	// Collect results
	for i := 0; i < numJobs; i++ {
		<-results
	}

	fmt.Printf("Memory after jobs: %v MB\n", bToMb(runtimeReadMemStatsHeap()))
	fmt.Println("All jobs done in", time.Since(start))
}

func bToMb(stats uint64) uint64 {
	return stats / 1024 / 1024
}

func runtimeReadMemStatsHeap() uint64 {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	return m.HeapAlloc
}
