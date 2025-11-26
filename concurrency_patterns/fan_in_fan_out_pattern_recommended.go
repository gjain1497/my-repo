// fan_out_fan_in_recommended.go
// RECOMMENDED APPROACH: 'go' inside, return channel
// Use this for Fan-Out/Fan-In, Pipeline, and functional composition

package main

import (
	"fmt"
	"sync"
	"time"
)

// Generator produces numbers on a channel
func Generator(nums ...int) <-chan int {
	out := make(chan int)
	go func() {
		for _, n := range nums {
			out <- n
		}
		close(out)
	}()
	return out
}

// Square processes numbers in parallel
// Pattern: 'go' inside, return channel
// Benefits:
// - Clean usage (one line)
// - Auto-starts goroutine
// - Returns channel immediately
// - Composable for pipelines
func Square(id int, in <-chan int) <-chan int {
	out := make(chan int)

	go func() { // ⭐ Goroutine inside function
		for n := range in {
			fmt.Printf("Worker %d processing: %d\n", id, n)
			time.Sleep(100 * time.Millisecond)
			out <- n * n
		}
		close(out)
	}()

	return out // ⭐ Returns channel immediately
}

// FanIn merges multiple channels into one
func FanIn(workers []<-chan int) <-chan int {
	out := make(chan int)
	var wg sync.WaitGroup

	// Start a goroutine for each input channel
	for _, ch := range workers {
		wg.Add(1)
		go func(c <-chan int) {
			defer wg.Done()
			for result := range c {
				out <- result
			}
		}(ch)
	}

	// Close output when all inputs done
	go func() {
		wg.Wait()
		close(out)
	}()

	return out
}

func main() {
	start := time.Now()

	// Generate numbers
	input := Generator(1, 2, 3, 4, 5)

	// Fan-Out: Create 3 workers
	numWorkers := 3
	var workers []<-chan int

	for i := 0; i < numWorkers; i++ {
		// ⭐ Clean usage: one line per worker
		workers = append(workers, Square(i, input))
	}

	// Fan-In: Merge results
	results := FanIn(workers)

	// Collect and print results
	fmt.Println("=== Results ===")
	for result := range results {
		fmt.Println("Result:", result)
	}

	elapsed := time.Since(start)
	fmt.Printf("\nTotal time: %v\n", elapsed)
	fmt.Printf("Expected: ~200ms (5 numbers / 3 workers * 100ms)\n")
}
