// fan_out_fan_in_alternative.go
// ALTERNATIVE APPROACH: 'go' outside, pass channel
// Use this for Worker Pool pattern where you need lifecycle control

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
// Pattern: 'go' outside, pass channel as parameter
// Benefits:
// - Explicit control over when to start
// - Worker Pool style
// - Can manage lifecycle (start/stop)
// Drawbacks:
// - More verbose (3 lines instead of 1)
// - Easy to forget 'go' keyword
// - Less composable
func Square(id int, in <-chan int, out chan<- int) {
	// ⭐ No goroutine inside
	// ⭐ No return value
	for n := range in {
		fmt.Printf("Worker %d processing: %d\n", id, n)
		time.Sleep(100 * time.Millisecond)
		out <- n * n
	}
	close(out)
}

// FanIn merges multiple channels into one
func FanIn(workers []<-chan int) <-chan int {
	out := make(chan int)
	var wg sync.WaitGroup

	for _, ch := range workers {
		wg.Add(1)
		go func(c <-chan int) {
			defer wg.Done()
			for result := range c {
				out <- result
			}
		}(ch)
	}

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
		// ⭐ Must create channel manually
		out := make(chan int)

		// ⭐ Must store channel before calling
		workers = append(workers, out)

		// ⭐ Must remember 'go' keyword
		go Square(i, input, out)

		// Why this way?
		//the point is we call it using go now
		//so we cannot get anything as a return type
		//from this function.

		//means we cannot write
		//ch := go Square(1, input)

		//so thats why there is no option
		//but to pass channel as paramter in the function

		// - When you call with 'go', you can't capture return value
		// - So we create channel first and pass it in
		// - This gives us the channel reference before starting goroutine
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
